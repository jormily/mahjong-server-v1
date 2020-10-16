package msg

type BaseResponse struct {
	Errcode 	int		`json:"errcode"`
	Errmsg 		string	`json:"errmsg"`
}

type Token struct {
	UserId 		int64	`json:"userId"`
	Time 		int64	`json:"time"`
	LifeTime	int64	`json:"lifeTime"`
}

type RoomConf struct {
	Type 			string	`json:"type"`
	BaseScore		int		`json:"difen"`
	ZiMo			int		`json:"zimo"`
	JiangDui		bool	`json:"jiangdui"`
	HuanSanZhang 	bool	`json:"huansanzhang"`
	DianGangHua		string	`json:"dianganghua"`
	DianGanH		int64	`json:"diangangh"`
	MengQing		bool 	`json:"menqing"`
	TianDiHu		bool	`json:"tiandihu"`
	MaxFan			int 	`json:"zuidafanshu"`
	MaxGames		int		`json:"jushuxuanze"`
	Creator 		int		`json:"creator"`
}

//////////////////////////////////
// account server msg
//////////////////////////////////
type ServerInfo  struct{
	Version int 	`json:"version"`
	Hall  	string	`json:"hall"`
	Appweb 	string	`json:"appweb"`
}

type ServerVersion  struct{
	Version int	 	`json:"version"`
}

type GuestResponse struct {
	BaseResponse
	Account		string	`json:"account"`
	Halladdr	string	`json:"halladdr"`
	Sign		string	`json:"sign"`
}

type BaseInfoResponse struct {
	BaseResponse
	Name 		string	`json:"name"`
	Sex 		int 	`json:"sex"`
	HeadImgUrl	string	`json:"headimgurl"`
}

//////////////////////////////////
// hall server msg
//////////////////////////////////
type LoginResponse struct {
	BaseResponse
	Account 	string 	`json:"account"`
	UserId 		int64 	`json:"userid"`
	Name 		string 	`json:"name"`
	Level  		int64 	`json:"lv"`
	Exp   		int64 	`json:"exp"`
	Coins  		int64 	`json:"coins"`
	Gems 		int64 	`json:"gems"`
	Ip 			string 	`json:"ip"`
	Sex 		int64 	`json:"sex"`
	RoomId		string 	`json:"roomid"`
}

type UserStateResponse struct {
	BaseResponse
	Gems		int 	`json:"gems"`
}

type MessageResponse struct {
	BaseResponse
	Msg 		string 	`json:"msg"`
	Version 	string 	`json:"version"`
}

type GameServerResponse struct {
	BaseResponse
	Ip 			string 	`json:"ip"`

}

type CreateRoomResponse struct {
	BaseResponse
	RoomId	string `json:"roomid"`
}

type EnterRoomResponse struct {
	BaseResponse
	RoomId		string 	`json:"roomid"`
	Ip 			string  `json:"ip"`
	Port		string  `json:"port"`
	Token		string  `json:"token"`
	Time		int64  	`json:"time"`
	Sign 		string 	`json:"sign"`
}

type EnterRoomInfoResponse struct {
	BaseResponse
	RoomId		string 	`json:"roomid"`
	Ip 			string  `json:"ip"`
	Port		string  `json:"port"`
	Token		Token  `json:"token"`
	Time		int64  	`json:"time"`
	Sign 		string 	`json:"sign"`
}

type PingResponse struct {
	BaseResponse
	Runing		bool	`json:"runing"`
}

//////////////////////////////////
// game server msg
//////////////////////////////////
type EnterTokenResponse struct {
	BaseResponse
	Tok 	Token	`json:"token"`
}

type EmptyMsg struct {

}

type LoginMsg struct {
	Token 	string 	`json:"token"`
	RoomId	string 	`json:"roomid"`
	Time	int64	`json:"time"`
	Sign 	string 	`json:"sign"`
}

type PlayerInfo struct {
	UserId 	int64	`json:"userid"`
	Ip 		string	`json:"ip"`
	Name 	string	`json:"name"`
	Score 	int		`json:"score"`
	Online	bool	`json:"online"`
	Ready	bool	`json:"ready"`
	SeateIndex	int `json:"seatindex"`
}

type LoginResult struct {
	RoomId 		string 			`json:"roomid"`
	Conf 		RoomConf		`json:"conf"`
	GameCount	int				`json:"numofgames"`
	Seats		[]PlayerInfo 	`json:"seats"`
}

type LoginResultResponse struct {
	BaseResponse
	Result 	LoginResult 	`json:"data"`
}

type ReadyResponse struct {
	UserId 	int64	`json:"userid"`
	Ready 	bool 	`json:"ready"`
}

type ChangeThreeResponse struct {
	UserId 	int64	`json:"userid"`
	ChangeCards	[]int	`json:"huanpais"`
}

type ChangeMethodNotify struct {
	Method 	int		`json:"method"`
}

type MakeLackNotify struct {

}

type PlayCardNotify struct {
	Card 	int		`json:"pai"`
	UserId 	int64	`json:"userId"`
}

type ActionNotify struct {
	Id 			int 	`json:"si"`
	Card 		int 	`json:"pai"`
	CanDraw		bool	`json:"hu"`
	CanPeng		bool	`json:"peng"`
	CanGang		bool	`json:"gang"`
	CardGang	[]int 	`json:"gangpai"`
}

type PassNotify struct {
	UserId 	int64	`json:"userId"`
	Card 	int 	`json:"pai"`
}

type GangNotify struct {
	UserId 	int64	`json:"userId"`
	Card 	int 	`json:"pai"`
	Type 	string 	`json:"gangtype"`
}

type HuNotify struct {
	Id 		int		`json:"seatindex"`
	ZiMo	bool	`json:"iszimo"`
	Card 	int 	`json:"hupai"`
}

type PlayerSettle struct {
	UserId 		int64 	`json:"userId"`
	Peng		[]int	`json:"pengs"`
	Actions 	[]int 	`json:"actions"`
	WangGang	[]int	`json:"wangangs"`
	DianGang	[]int	`json:"diangangs"`
	AnGang		[]int 	`json:"angangs"`
	NumOfGen	int 	`json:"numofgen"`
	Holds		[]int	`json:"holds"`
	Fan 		int 	`json:"fan"`
	Score 		int 	`json:"score"`
	TotalScore	int		`json:"totalscore"`
	QingYiSe	bool 	`json:"qingyise"`
	Pattern 	string 	`json:"pattern"`
	IsGangHu	bool	`json:"isganghu"`
	MenQing		bool	`json:"menqing"`
	ZhongZhang 	bool	`json:"zhongzhang"`

	JingGouHu	bool	`json:"jingouhu"`
	HaiDiHu		bool	`json:"haidihu"`
	TianHu		bool	`json:"tianhu"`
	DiHu		bool	`json:"dihu"`
	Huorder		int		`json:"huorder"`
	HuInfo 		[]int 	`json:"huinfo"`  //todo:完善
}

type TotalSettle struct {
	NumZiMo 	int 	`json:"numzimo"`
	NumJiaoPai	int		`json:"numjiaopai"`
	NumDianPao	int 	`json:"numdianpao"`
	NumAnGang 	int 	`json:"numangang"`
	NumMingGang	int		`json:"numminggang"`
	NumChaDaJiao int 	`json:"numchadajiao"`
}

type GameOverNotify struct {
	PlayerSettleList 	[]PlayerSettle 		`json:"results"`
	TotalSettleList		[]TotalSettle		`json:"endinfo"`
	IsOver 				bool 				`json:"isOver"`
}