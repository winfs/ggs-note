package manager

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"ggs/gate"
	"ggs/log"
	"ggs/service"

	"github.com/golang/protobuf/proto"
	"github.com/streamrail/concurrent-map"

	"ggsserver/db"
	"ggsserver/msg"
)

/***************************************************/

var (
	Users        = cmap.New()
	GameSkeleton *service.Skeleton
)

type DataManager struct {
	UM *UserManager

	agent    gate.Agent
	L        sync.Mutex
	DeviceId uint32
}

/***************************************************/

func init() {
	rand.Seed(time.Now().UnixNano()) // 以当前纳秒数作为随机数种子
}

func NewDataManager(uid string, agent gate.Agent, deviceId uint32) (dataManager *DataManager, err error) {
	// 检查用户是否重复登录，如果是，则回复重复登录错误信息，并断开连接
	checkMultiLogin(uid, deviceId)

	// 获取用户的所有信息
	data, err := db.Redis.Cmd("HGETALL", GetUserKey(uid)).Map()
	if err != nil {
		return
	}

	// 初始化DataManager
	dataManager = &DataManager{
		agent:    agent,
		DeviceId: deviceId,
	}
	dataManager.UM = newUserManager(dataManager)

	// 初始化所有Manager数据
	err = dataManager.OnLoad(data, uid)
	if err != nil {
		return
	}

	// 存储uid => DM的映射
	Users.Set(uid, dataManager)

	// 后续的一些数据初始化
	dataManager.afterLoad(data)

	return
}

// 检查重复登录
func checkMultiLogin(uid string, deviceId uint32) {
	// 如果已有该用户id
	if dm, ok := Users.Get(uid); ok {
		d := dm.(*DataManager)
		// 如果在不同的设备上登录(id不一致)
		if d.DeviceId != deviceId {
			d.SendError(5, "") // 发送重复登录的错误消息
			d.agent.Close()    // 断开连接
		}

		d.Clear() // 最后的保存工作
	}
}

// 保存数据，如果失败则移除该用户id
func (dm *DataManager) Clear() {
	defer func() {
		dm.Destory() // 从Users中移除该用户id
		dm.agent.SetUserData(nil)
	}()

	dm.L.Lock()
	err := dm.OnSave() // 保存工作
	if err != nil {
		log.Error("save error: %v", err)
	}
	dm.L.Unlock()
}

func (d *DataManager) Destory() {
	d2, _ := Users.Get(d.UM.userId())
	if d == d2 {
		Users.Remove(d.UM.userId()) // 移除该元素
	}
}

// 初始化所有Manager数据
func (d *DataManager) OnLoad(data map[string]string, uid string) (err error) {
	if err = d.UM.onLoad(data, uid); err != nil {
		log.Error("load user manager error: %v, uid: %v", err, uid)
	}

	return
}

// onload里的数据加载后才可以进行初始化的数据可以在此操作
func (d *DataManager) afterLoad(data map[string]string) {
	// ...
}

// 保存所有Manager数据
func (d *DataManager) OnSave() (err error) {
	data := make(map[string]interface{})

	d.UM.onSave(data)

	db.Redis.Cmd("HMSET", GetUserKey(d.UM.userId()), data) // 保存到redis中

	return
}

/***************************************************/

// 给客户端回复消息
func (d *DataManager) WriteMsg(msg interface{}) {
	d.agent.WriteMsg(msg)
}

// 给客户端回复响应失败消息
func (d *DataManager) ResponseErr(code uint32, desc string, byAction uint32) {
	log.Debug("response error status, code: %v, desc: %v, action: %v", code, desc, byAction)
	d.agent.WriteMsg(&msg.Status{
		Code:        proto.Uint32(code),
		Description: proto.String(desc),
		ByAction:    proto.Uint32(byAction),
	})
}

// 给客户端回复响应成功消息
func (d *DataManager) ResponseDone(desc string, byAction uint32, params ...string) {
	log.Debug("response success status, desc: %v, action: %v, param: %v", desc, byAction, params)
	m := &msg.Status{
		Code:        proto.Uint32(10000),
		Description: proto.String(desc),
		ByAction:    proto.Uint32(byAction),
	}

	if len(params) > 0 && params[0] != "" {
		m.Param = proto.String(params[0])
	}

	d.agent.WriteMsg(m)
}

// 给客户端回复登录响应消息
func (d *DataManager) Response() {
	m := &msg.Login{
		User: d.UM.responseMsg(),
	}

	d.agent.WriteMsg(m)
}

// 给客户端回复错误消息
// code: 5, 重复登录
// code: 7, 屏蔽
// code: 8, 下线
func (d *DataManager) SendError(code uint32, desc string) {
	log.Debug("send error, code: %v, desc: %v", code, desc)
	msg := &msg.Error{
		Code: proto.Uint32(code),
		Desc: proto.String(desc),
	}
	d.agent.WriteMsg(msg)
}

func (d *DataManager) String() (str string) {
	str += fmt.Sprintf("User: %+v\n", d.UM)
	return
}
