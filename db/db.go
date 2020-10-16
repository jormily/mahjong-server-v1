package db

import (
	"encoding/base64"
	"errors"
	"strconv"
)

var (
	ErrNotFound = errors.New("data not found")
)

func UserExist(account string) bool {
	model := &TUsers{Account: account}
	has, err := database.Get(model)
	if err != nil {
		return false
	}
	return has
}

func CreateUser(account string,name string,coins int64,gems int64,sex int64,headimg string) error {
	model := &TUsers{
		Account: account,
		Name: base64.StdEncoding.EncodeToString([]byte(name)),
		Coins: coins,
		Gems: gems,
		Sex: sex,
		Headimg: headimg,
	}

	_,err := database.Insert(model)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(account string) (*TUsers, error){
	model := &TUsers{Account: account}
	has, err := database.Get(model)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil,ErrNotFound
	}
	name,_ := base64.StdEncoding.DecodeString(model.Name)
	model.Name = string(name)

	return model,nil
}

func GetUserById(userId int64) (*TUsers,error){
	model := &TUsers{Userid: userId}
	has, err := database.Get(model)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil,ErrNotFound
	}
	name,_ := base64.StdEncoding.DecodeString(model.Name)
	model.Name = string(name)

	return model,nil
}

func GetMessage(typ string,version string) (*TMessage, error){
	model := &TMessage{}
	model.Type = typ
	if version != "null" {
		model.Version = version
	}

	var has,err = database.Get(model)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil,ErrNotFound
	}
	return model,nil
}

func RoomExist(roomid string) bool {
	model := &TRooms{Id: roomid}
	has,err := database.Get(model)
	if err != nil {
		return false
	}
	return has
}

func GetRoomModel(roomid string) (*TRooms, error){
	model := new(TRooms)
	has,err := database.Where("id=?",roomid).Get(model)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil,ErrNotFound
	}
	return model,nil
}

func GetRoomModelByServerId(serverId int64) ([]*TRooms){
	models := make([]*TRooms,0)
	database.Where("serverid = ?",serverId).Find(&models)
	return models
}

func SetUserRoomId(userid int64,roomid string) error {
	_,err := database.Exec("update `t_users` set `roomid` = ? where `userid` = ?",roomid,userid)
	if err != nil {
		return err
	}

	return nil
}

func IsRoomExist(roomid string) bool {
	roomModel := &TRooms{Id: roomid}
	has,_ := database.Get(roomModel)
	return has
}

func InsertRoom(roomid string,config string,ip string,port int,createTime int64) (int64, bool) {
	model := &TRooms{}
	model.Uuid = strconv.FormatInt(createTime,10) + roomid
	model.Id = roomid
	model.Ip = ip
	model.Port = int64(port)
	model.BaseInfo = config
	model.CreateTime = createTime
	affected, err := database.Insert(model)
	if err != nil {
		return affected,false
	}
	return affected,true
}