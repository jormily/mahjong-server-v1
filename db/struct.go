package db

type TGuests struct {
	GuestAccount string `xorm:"guest_account"`
}

func (TGuests) TableName() string {
	return "t_guests"
}

type TMessage struct {
	Type    string `xorm:"type"`
	Msg     string `xorm:"msg"`
	Version string `xorm:"version"`
}

func (TMessage) TableName() string {
	return "t_message"
}

type TRooms struct {
	Uuid       string `xorm:"uuid"`
	Id         string `xorm:"id"`
	BaseInfo   string `xorm:"base_info"`
	CreateTime int64  `xorm:"create_time"`
	NumOfTurns int64  `xorm:"num_of_turns"`
	NextButton int64  `xorm:"next_button"`
	UserId0    int64  `xorm:"user_id0"`
	UserIcon0  string `xorm:"user_icon0"`
	UserName0  string `xorm:"user_name0"`
	UserScore0 int64  `xorm:"user_score0"`
	UserId1    int64  `xorm:"user_id1"`
	UserIcon1  string `xorm:"user_icon1"`
	UserName1  string `xorm:"user_name1"`
	UserScore1 int64  `xorm:"user_score1"`
	UserId2    int64  `xorm:"user_id2"`
	UserIcon2  string `xorm:"user_icon2"`
	UserName2  string `xorm:"user_name2"`
	UserScore2 int64  `xorm:"user_score2"`
	UserId3    int64  `xorm:"user_id3"`
	UserIcon3  string `xorm:"user_icon3"`
	UserName3  string `xorm:"user_name3"`
	UserScore3 int64  `xorm:"user_score3"`
	Ip         string `xorm:"ip"`
	Port       int64  `xorm:"port"`
	Serverid   int64  `xorm:"serverid"`
}

func (TRooms) TableName() string {
	return "t_rooms"
}

type TUsers struct {
	Userid  int64  `xorm:"userid"`  // 用户ID
	Account string `xorm:"account"` // 账号
	Name    string `xorm:"name"`    // 用户昵称
	Sex     int64  `xorm:"sex"`
	Headimg string `xorm:"headimg"`
	Lv      int64  `xorm:"lv"`    // 用户等级
	Exp     int64  `xorm:"exp"`   // 用户经验
	Coins   int64  `xorm:"coins"` // 用户金币
	Gems    int64  `xorm:"gems"`  // 用户宝石
	Roomid  string `xorm:"roomid"`
	History string `xorm:"history"`
}

func (TUsers) TableName() string {
	return "t_users"
}

type TAccounts struct {
	Account  string `xorm:"account"`
	Password string `xorm:"password"`
}

func (TAccounts) TableName() string {
	return "t_accounts"
}

type TGames struct {
	RoomUuid      string `xorm:"room_uuid"`
	GameIndex     int64  `xorm:"game_index"`
	BaseInfo      string `xorm:"base_info"`
	CreateTime    int64  `xorm:"create_time"`
	Snapshots     string `xorm:"snapshots"`
	ActionRecords string `xorm:"action_records"`
	Result        string `xorm:"result"`
}

func (TGames) TableName() string {
	return "t_games"
}

type TGamesArchive struct {
	RoomUuid      string `xorm:"room_uuid"`
	GameIndex     int64  `xorm:"game_index"`
	BaseInfo      string `xorm:"base_info"`
	CreateTime    int64  `xorm:"create_time"`
	Snapshots     string `xorm:"snapshots"`
	ActionRecords string `xorm:"action_records"`
	Result        string `xorm:"result"`
}

func (TGamesArchive) TableName() string {
	return "t_games_archive"
}
