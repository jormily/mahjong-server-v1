package manager

import (
	"encoding/json"
	"errors"
	//"github.com/kudoochui/kudos/component/timers"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"time"

	"github.com/jormily/mahjong-server/db"
	"github.com/jormily/mahjong-server/msg"
)

var RoomMgr = NewRoomManager()

func generateRoomId() string {
	roomId := 0
	for i := 0;i < 6;i++ {
		rand.Seed(time.Now().UnixNano())
		roomId = roomId*10 + (rand.Intn(9) + 1)
	}
	return strconv.FormatInt(int64(roomId),10)
}

type Location struct {
	roomid 	string
	seate	int
}

type PlayerData struct {
	// 基础数据
	UserId 		int64
	Score 		int
	Name 		string

	// 状态数据
	State 		int
	Ready		bool
	Seate		int

	// 游戏数据
	NumZiMo		int 	//自摸
	NumJiePao	int
	NumDianPao	int
	NumAnGang	int
	NumMingGang	int
	NumChaJiao	int

}


type Room struct {
	uuid 		string
	id 			string
	turns		int
	createTime	int64
	nextButton 	int
	players		[4]*PlayerData
	conf 		msg.RoomConf
	mahjong		*MahjongLogic
	//timer := timers.NewTimer()
	//timer 		timers.Timers
}

func NewRoom(roomModel *db.TRooms) *Room {
	this := new(Room)
	this.uuid = roomModel.Uuid
	this.id = roomModel.Id
	this.turns = 0
	this.createTime = roomModel.CreateTime
	this.nextButton = int(roomModel.NextButton)
	this.players = [4]*PlayerData{}
	json.Unmarshal([]byte(roomModel.BaseInfo),&this.conf)
	this.mahjong = NewMahjongLogic(this)

	for i:=0;i<4;i++ {
		this.players[i] = new(PlayerData)
		this.players[i].State = 0
		this.players[i].Seate = i
	}

	return this
}

func (this *Room)GetUserId(seat int) int64 {
	if this.players[seat] == nil {
		return 0
	}
	return this.players[seat].UserId
}

func (this *Room)GetPlayerId(userid int64) int{
	for _,v := range this.players {
		if v.UserId == userid {
			return v.Seate
		}
	}
	return  -1
}

func (this *Room)ClearReady(){
	for _,v := range this.players {
		v.Ready = false
	}
}

func (this *Room)GetPlayerData(id interface{}) *PlayerData {
	switch id.(type) {
		case int:
			if id.(int) < 0 || id.(int) >= 4{
				return nil
			}
			return this.players[id.(int)]
		case int64:
			for _,v := range this.players {
				if v.UserId == id.(int64) {
					return v
				}
			}
			return nil
		default:
			return nil
	}
}

func (this *Room)GetAllPlayer() []*PlayerData {
	return this.players[:]
}

func (this *Room)GetLogic() *MahjongLogic {
	return this.mahjong
}

func (this *Room)PlayerEnter(userid int64,name string,score int) (int,error) {
	for _,v := range this.players {
		if v.UserId == 0 {
			v.UserId = userid
			v.Name = name
			v.Score = score
			return v.Seate,nil
		}
	}
	return -1,errors.New("room has no seate")
}

func (this *Room)Send(index int,event string,msg interface{}){
	if index < 0 || index >= len(this.players) {
		return
	}
	player := this.players[index]
	PlayerMgr.Send(player.UserId,event,msg)
}

func (this *Room)Broadcast(event string,msg interface{}){
	for _,v := range this.players {
		PlayerMgr.Send(v.UserId,event,msg)
	}
}

func (this *Room)BroadcastEx(event string,msg interface{},id int){
	for _,v := range this.players {
		if  v.Seate != id {
			PlayerMgr.Send(v.UserId,event,msg)
		}
	}
}

type RoomManager struct {
	location 	map[int64]*Location
	rooms		map[string]*Room
	creatingRooms map[string]bool
	totalRooms	int
}

func NewRoomManager() *RoomManager {
	this := new(RoomManager)
	this.location = map[int64]*Location{}
	this.rooms = map[string]*Room{}
	this.creatingRooms = map[string]bool{}
	this.totalRooms = 0
	return this
}

func checkConfig(conf *msg.RoomConf) bool {
	return true
}

func (this *RoomManager) GenerateRoomId() string {
	var roomid string
	for {
		roomid = generateRoomId()
		if db.RoomExist(roomid) {
			continue
		}else{
			break
		}
	}
	return roomid
}

func (this *RoomManager) CreateAllRoom(){
	roomModels := db.GetRoomModelByServerId(viper.GetInt64("game_svr.id"))
	for _,roomModel := range roomModels {
		room := NewRoom(roomModel)
		this.rooms[roomModel.Id] = room
	}
}

func (this *RoomManager) CreateRoom(roomid string,userid int64) error {
	if _,ok := this.rooms[roomid];ok {
		return nil
	}

	roomModel,err := db.GetRoomModel(roomid)
	if err != nil {
		return errors.New("room not find")
	}

	room := NewRoom(roomModel)
	this.rooms[roomid] = room

	return nil
}

func (this *RoomManager) EnterRoom(roomid string,userid int64,name string,score int) error {
	l,ok := this.GetUserLocation(userid)
	if ok && l.roomid == roomid {
		return errors.New("room is exist")
	}

	if room,ok := this.rooms[roomid];ok {
		if seate,err := room.PlayerEnter(userid,name,score);err == nil {
			this.location[userid] = &Location{
				roomid: roomid,
				seate: seate,
			}
			db.SetUserRoomId(userid,roomid)
			return nil
		}else{
			return err
		}
	}else{
		return errors.New("room is not run")
	}
}

func (this *RoomManager) GetUserLocation(userid int64) (*Location,bool) {
	if l,ok := this.location[userid];ok {
		return l,true
	}
	return nil,false
}

func (this *RoomManager) GetRoomByUserId(userid int64) *Room {
	if l,ok := this.GetUserLocation(userid);ok {
		if r,ok := this.rooms[l.roomid];ok {
			return r
		}
	}
	return nil
}

func (this *RoomManager) GetRoom(roomid string) *Room {
	if r,ok := this.rooms[roomid];ok {
		return r
	}
	return nil
}


