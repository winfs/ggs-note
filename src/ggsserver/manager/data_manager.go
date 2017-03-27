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

	dataManager = &DataManager{
		agent:    agent,
		DeviceId: deviceId,
	}

	// 初始化操作
	dataManager.UM = newUserManager(dataManager)

	//
	err = dataManager.OnLoad(data, uid)
	if err != nil {
		return
	}

	// 存储uid => DM的映射
	Users.Set(uid, dataManager)

	//
	dataManager.afterLoad(data)

	return
}

// 检查重复登录
func checkMultiLogin(uid string, deviceId uint32) {
	if dm, ok := Users.Get(uid); ok {
		d := dm.(*DataManager)
		if d.DeviceId != deviceId {
			d.SendError(5, "") // 发送重复登录的错误消息
			d.agent.Close()    // 关闭连接
		}

		d.Clear() // 最后的清除工作
	}
}

func (dm *DataManager) Clear() {
	defer func() {
		dm.Destory() //
		dm.agent.SetUserData(nil)
	}()

	dm.L.Lock()
	err := dm.OnSave() //
	if err != nil {
		log.Error("save error: %v", err)
	}
	dm.L.Unlock()
}

// 如果存在则移除
func (d *DataManager) Destory() {
	d2, _ := Users.Get(d.UM.userId())
	if d == d2 {
		Users.Remove(d.UM.userId())
	}
}

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

func (d *DataManager) OnLoad(data map[string]string, uid string) (err error) {
	if err = d.UM.onLoad(data, uid); err != nil {
		log.Error("load user manager error: %v, uid: %v", err, uid)
	}

	return
}

func (d *DataManager) afterLoad(data map[string]string) {
	// ...
}

func (d *DataManager) OnSave() (err error) {
	data := make(map[string]interface{})

	d.UM.onSave(data)

	db.Redis.Cmd("HMSET", GetUserKey(d.UM.userId()), data)

	return
}

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
func (d *DataManager) ResponseDone(desc string, byAction uint32) {
	log.Debug("response success status, desc: %v, action: %v", desc, byAction)
	d.agent.WriteMsg(&msg.Status{
		Code:        proto.Uint32(10000),
		Description: proto.String(desc),
		ByAction:    proto.Uint32(byAction),
	})
}

// 给客户端回复登录响应消息
func (d *DataManager) Response() {
	m := &msg.Login{
		User: d.UM.responseMsg(),
	}

	d.agent.WriteMsg(m)
}

func (d *DataManager) String() (str string) {
	str += fmt.Sprintf("User: %+v\n", d.UM)
	return
}
