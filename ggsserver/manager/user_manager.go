package manager

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"ggs/log"

	"github.com/Jarvis-Li/golibs/api"
	jhttp "github.com/Jarvis-Li/golibs/http"
	jtime "github.com/Jarvis-Li/golibs/time"
	"github.com/golang/protobuf/proto"

	"ggsserver/conf"
	"ggsserver/db"
	"ggsserver/msg"
)

/************************************************************/

type UserManager struct {
	User *User
	DM   *DataManager
}

type User struct {
	Id                  string
	ServerID            string
	Name                string
	OfflineTime         time.Time
	Avatar              string
	CountryName         string
	Platform            string
	RegisterTime        int
	HasUsedGift         int
	HasGotSubscribeGift int
	/*
		MoneyChargedToday   float64
		Progress            map[string]bool
		FirstCharges        []float64
		MoneyCharged        float64
		MoneyChargedSuiMo   float64
		ChargeDays          int
		AttendedWars        map[int]bool
		LastLogoutToday     bool
	*/
	IsNew bool
}

/************************************************************/

func newUserManager(DM *DataManager) *UserManager {
	user := new(UserManager)
	user.User = new(User)
	//user.User.Progress = make(map[string]bool)
	//user.User.AttendedWars = make(map[string]bool)
	user.DM = DM
	return user
}

// 初始化用户数据
func (um *UserManager) onLoad(data map[string]string, uid string) error {
	var err error
	um.User.Id = uid
	um.User.ServerID = getServerID(uid)

	// 查找当前用户离线时间
	if str, ok := data["offline_time"]; ok {
		unixTime, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		um.User.OfflineTime = time.Unix(unixTime, 0)
		//um.setLastLogoutToday()
	} else { // 没有，则为新用户
		um.User.IsNew = true
	}

	if str, ok := data["avatar"]; ok {
		avatar, _ := url.QueryUnescape(str) // url解码
		um.User.Avatar = avatar
	}

	if str, ok := data["name"]; ok {
		um.User.Name = str
	}

	if str, ok := data["platform"]; ok {
		um.User.Platform = str
	}

	if str, ok := data["has_use_gift"]; ok {
		num, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		um.User.HasUsedGift = num
	}

	if str, ok := data["has_got_subscribe_gift"]; ok {
		num, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		um.User.HasGotSubscribeGift = num
	}

	if str, ok := data["country_name"]; ok {
		um.User.CountryName = str
	} else {
		if um.User.Name != "" {
			um.User.CountryName = string([]rune(um.User.Name)[0])
		}
	}

	if str, ok := data["register_time"]; ok {
		registerTime, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		um.User.RegisterTime = registerTime
	}

	// ...

	//um.loadProgress(data)

	return nil
}

// 更新用户的数据
func (um *UserManager) OnSave(data map[string]interface{}) error {
	data["offline_time"] = time.Now().Unix()
	data["name"] = um.User.Name
	data["country_name"] = um.User.CountryName
	data["has_use_gift"] = um.User.HasUsedGift

	//um.saveProgress(data)
	return nil
}

// 返回用户响应消息
func (um *UserManager) responseMsg() *msg.User {
	return &msg.User{
		Id:          proto.String(um.User.Id),
		Name:        proto.String(um.UserName()),
		Avatar:      proto.String(um.User.Avatar),
		OfflineTime: proto.Uint32(uint32(time.Now().Sub(um.User.OfflineTime).Seconds())),
		CountryName: proto.String(um.countryName()),
		//KingLisence: um.DM.KLM.msg(),
		//Charges:     um.User.FirstCharges,
		IsNew: proto.Bool(um.User.IsNew),
	}
}

func (um *UserManager) UserName() string {
	if conf.Server.IsMultiple() { // 如果有合服
		return um.UserNameWithServer()
	}
	return um.User.Name
}

func (um *UserManager) UserNameWithServer() string {
	return getServerName(um.User.ServerID) + um.User.Name
}

func (um *UserManager) countryName() string {
	return um.User.CountryName
}

// 给客户端发送用户响应消息
func (um *UserManager) response() {
	msg := um.responseMsg() // 获取用户响应消息
	um.DM.agent.WriteMsg(msg)
}

func (um *UserManager) responseName() {
	um.DM.agent.WriteMsg(&msg.User{
		Name: proto.String(um.UserName()),
	})
}

func (um *UserManager) responseCountryName() {
	um.DM.agent.WriteMsg(&msg.User{
		CountryName: proto.String(um.countryName()),
	})
}

func (um *UserManager) String() string {
	return fmt.Sprintf("%+v\n", um.User)
}

/************************************************************/

func getServerID(uid string) string {
	idInfo := strings.Split(id, ":")
	return idInfo[0]
}

func getServerName(serverID string) string {
	return conf.Server.Names[serverID] + "."
}
