package manager

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

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
	GameSkeleton = *service.Skeleton
)

type DataManager struct {
	UM *UserManager

	agent    *gate.Agent
	L        sync.Mutex
	DeviceId int32
}

/***************************************************/

func init() {
	rand.Seed(time.Now().UnixNano()) // 以当前纳秒数作为随机数种子
}

func NewDataManager(uid string, agent gate.Agent, deviceId uint32) (dataManager *DataManager, err error) {
	checkMultiLogin(uid, deviceId)

	data, err := db.Redis.Cmd("HGETALL", getUserKey(uid)).Map()
	if err != nil {
		return
	}

	dataManager = &DataManager{
		agent:    agent,
		DeviceId: deviceId,
	}
	dataManager.UM = newUserManager(dataManager)

	err = dataManager.OnLoad(data, uid)
	if err != nil {
		return
	}

	Users.Set(uid, dataManager)
	dataManager.afterLoad(data)

	return
}

// 检查重复登录
func checkMultiLogin(uid string, deviceId uint32) {
	if dm, ok := Users.Get(uid); ok {
		d := dm.(*DataManager)
		if d.DeviceId != deviceId {
			d.SendError(5, "") // 发送重复登录的错误消息
			d.agent.Close()
		}

		d.Clear()
	}
}

func (d *DataManager) Clear() {
	defer func() {
		dm.Destroy() //
		dm.agent.SetUserData(nil)
	}()

	dm.L.Lock()
	err := dm.OnSave() //
	if err != nil {
		log.Error("save error: %v", err)
	}
	dm.L.Unlock()
}

func (d *DataManager) Destory() {
	d2, _ := Users.Get(d.UM.userId())
	if d == d2 {
		Users.Remove(d.UM.userId)
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

	db.Redis.Cmd("HMSET", getUserKey(d.UM.userId()), data)

	return
}

func (d *DataManager) WriteMsg(msg interface{}) {
	d.agent.WriteMsg(msg)
}

// 响应失败
func (d *DataManager) ResponseErr(code uint32, desc string, byAction uint32) {
	log.Debug("response error status, code: %v, desc: %v, action: %v", code, desc, byAction)
	d.agent.WriteMsg(&msg.Status{
		Code:        proto.Uint32(code),
		Description: proto.String(desc),
		ByAction:    proto.Uint32(byAction),
	})
}

// 响应成功
func (d *DataManager) ResponseDone(desc string, byAction uint32) {
	log.Debug("response success status, desc: %v, action: %v", desc, byAction)
	d.agent.WriteMsg(&msg.Status{
		Code:        proto.Uint32(10000),
		Description: proto.String(desc),
		ByAction:    proto.Uint32(byAction),
	})
}

// 登录响应
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
