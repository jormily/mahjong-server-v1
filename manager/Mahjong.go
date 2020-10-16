package manager

import (
	"fmt"
	"github.com/jormily/mahjong-server/msg"
	"github.com/jormily/mahjong-server/tool"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// 牌值
var (
	Card_Value_Start = 0
	Card_Value_End = 26
)
// 牌类型
var (
	Card_Type_Tong 	= 0
	Card_Type_Tiao 	= 1
	Card_Type_Wan	= 2
)
// 玩家操作
var (
	ACTION_CHUPAI = 1
	ACTION_MOPAI = 2
	ACTION_PENG = 3
	ACTION_GANG = 4
	ACTION_HU = 5
	ACTION_ZIMO = 6
)

var (
	Act_State_HU 	= 1
	Act_State_Peng 	= 2
	Act_State_Gang 	= 4
	Act_State_AnG	= 8
	Act_State_DianG = 16
	Act_State_WanG 	= 32
	Act_State_Play	= 64
)

func getCardType(card int) int {
	return card/9
}

func getCardValue(card int) int {
	card = (card + 1)%9
	if card == 0 {
		card = 9
	}
	return card-1
}

func getCardSlice(cards *[]int) *[27]int {
	s := [27]int{}
	for _,card := range *cards {
		s[card]++
	}
	return &s
}

func checkLack(cards *[]int,lack int) bool {
	for i:=len(*cards)-1;i>=0;i-- {
		if getCardType((*cards)[i]) == lack {
			return false
		}
	}
	return true
}

// 巧七对
func isSevenPair(cards *[]int,array *[27]int) bool {
	if array == nil {
		array = getCardSlice(cards)
	}

	if len(*cards) != 14 {
		return false
	}

	// 要保证ar数据正常
	for _, cnt := range array {
		if !(cnt == 2 || cnt == 0) {
			return false
		}
		}
	return true
}

// 龙七对
func isBigSevenPair(cards *[]int,array *[27]int) bool {
	if array == nil {
		array = getCardSlice(cards)
	}

	if len(*cards) != 14 {
		return false
	}

	// 要保证ar数据正常
	big := false
	for _, cnt := range array {
		if !(cnt == 2 || cnt == 4 || cnt == 0) {
			return false
		}
		if cnt == 4 {
			big = true
		}
	}
	return big
}

// 清一色
func isSameColor(cards *[]int) bool {
	if len(*cards) == 0 {
		return false
	}
	color := getCardValue((*cards)[0])
	for _,card := range *cards {
		if getCardType(card) != color {
			return false
		}
	}
	return true
}

// 大单吊
func isSinglePair(cards *[]int) bool {
	if len(*cards) == 2 && (*cards)[0] == (*cards)[1] {
		return true
	}
	return false
}

// 大对子
func isBigPair(cards *[]int,array *[27]int) bool {
	if array == nil {
		array = getCardSlice(cards)
	}
	if len(*cards) != 5 {
		return false
	}

	// 要保证ar数据正常
	for _, cnt := range array {
		if !(cnt == 3 || cnt == 2 || cnt == 0) {
			return false
		}
	}
	return true
}

func isSomeGroup(markList tool.IntList,array [27]int) bool {
	checkGroup := func(array [27]int) bool {
		var index int
		for i:=0;i<3;i++{
			for j:=0;j<7;j++{
				index = i*9+j
				if array[index] > 0 {
					array[index+1] -= array[index]
					array[index+2] -= array[index]
					array[index] = 0

					if array[index+1] < 0 || array[index+2] < 0 {
						return false
					}
				}
			}
			if array[i*9+7] > 0 || array[i*9+8] > 0 {
				return false
			}
		}
		return true
	}

	getMark := func(list tool.IntList,array *[27]int) int {
		for card,cnt := range array {
			if cnt >= 3 && list.IndexOf(card) == 0 {
				return card
			}
		}
		return -1
	}

	if mark := getMark(markList,&array); mark >= 0 {
		markList.Insert(mark)
		if isSomeGroup(markList.Clone(),array) {
			return true
		}else{
			array[mark] = array[mark] - 3
			return  isSomeGroup(markList,array)
		}
	}else{
		return checkGroup(array)
	}

	return false	
}

func isNormalH(cards *[]int,array [27]int) bool {
	//一个对子+N坎牌（三个相同的牌、三个数值连续牌)
	if (len(*cards)-2)%3 != 0 {
		//打印错误，牌数错误
		return false
	}

	for i:=0;i<len(array);i++ {
		if array[i] >= 2 {
			array[i] -= 2
			if isSomeGroup(tool.IntList{},array) {
				log.Info("true-2")
				return true
			}
			array[i] += 2
		}
	}

	return false
}

func checkHu(cards []int,lack int,args...int) bool {
	log.Info(cards)
	log.Info(lack)

	if !checkLack(&cards,lack) {
		return false
	}

	array := getCardSlice(&cards)
	if isBigSevenPair(&cards,array) || isSevenPair(&cards,array) {
		log.Info("true-1")
		return true
	}

	return isNormalH(&cards,*array)
}

func arrangeCards(cards *[]int) map[int][]int {
	array := make(map[int][]int)
	for _,c := range *cards {
		typ := getCardType(c)
		val := getCardValue(c)
		if array[typ] == nil {
			array[typ] = make([]int,9)
		}
		array[typ][val]++
	}
	return array
}


type PlayData struct {
	seat 		int

	holds 		[]int
	folds 		[]int
	agCards 	[]int
	dgCards 	[]int
	wgCards 	[]int
	peCards		[]int
	lack		int
	changeCards []int
	gangCards 	[]int
	huCard 		int
	huList 		[]int

	cardMap 	map[int]int
	tingMap 	map[int]int

	canGang		bool
	canPen		bool
	canHu		bool
	isZiMo		bool
	canPlay		bool
	state 		int


	hufan	 	int
	drawed		bool
	drawS 		bool
	drawG		bool

	action 		[]int
	fan 		int
	score  		int

	lastfg 		int		// 上一次放杠玩家

	numZiMo 	int
	numJiePao 	int
	numDianPao 	int
	numAnGang 	int
	numMingGang int
	numChaJiao 	int
}

func newPlayData(seat int) *PlayData {
	this := new(PlayData)
	this.seat = seat
	//this.init()
	return this
}

func (this *PlayData) setState(idx int,state int){
	if state == 0 {
		this.state = this.state & (1<<idx)
	}
}

func (this *PlayData) getHoldCardType() [3]int {
	tl := [3]int{}
	for _,card := range this.holds {
		ct := getCardType(card)
		tl[ct]++
	}
	return tl
}

func (this *PlayData) checkCards(c interface{}) bool {
	switch c.(type) {
	case int:
		if this.cardMap[c.(int)] > 0 {
			return true
		}else{
			return false
		}
	case []int:
		cardMap := tool.Slice2Map(c.([]int))
		for card,cnt := range cardMap {
			if this.cardMap[card]  < cnt {
				return false
			}
		}
		return true
	default:
		return false
		
	}
}

func (this *PlayData) checkCardLack(card int) bool {
	return this.lack == getCardType(card)
}

/**
是否打缺
 */
func (this *PlayData) checkLack() bool {
	for _,card := range this.holds {
		if getCardType(card) == this.lack {
			return false
		}
	}
	return true
}

func (this *PlayData) checkTing() bool {
	if !this.checkLack() {
		return false
	}

	return true
}

func (this *PlayData)init(){
	this.holds 	= []int{}
	this.folds 	= []int{}
	this.agCards = []int{}
	this.dgCards = []int{}
	this.wgCards = []int{}
	this.peCards	= []int{}
	this.huList = []int{}

	this.lack = -1
	this.hufan = -1

	this.changeCards = nil
	this.action =	[]int{}
	this.cardMap =	map[int]int{}
	this.tingMap =	map[int]int{}

	this.canGang =	false
	this.canPen =	false
	this.canHu =	false
	this.canPlay =	false
	this.drawed =	false
	this.drawS = 	false
	this.drawG =	false
}

type MahjongLogic struct {
	round 		int
	cardList 	[]int
	cardIndex 	int
	banker 		int
	turn 		int
	state 		string
	actionList 	[]int
	drawCardList []int
	playList 	[4]*PlayData
	config 		msg.RoomConf

	playCnt 	int
	playCard 	int

	room		*Room
}

/**
new logic class
 */
func NewMahjongLogic(room *Room) *MahjongLogic {
	this := new(MahjongLogic)
	this.round = 0
	this.cardList = make([]int,0)
	this.banker = -1
	this.turn = -1
	this.playList = [4]*PlayData{}
	this.room = room
	this.init()

	return this
}


/**
init logic
 */
func (this *MahjongLogic)init(){
	this.drawCardList = []int{}
	this.state = "idle"
	this.playCnt = 0
	this.playCard = -1
	this.round ++

	if this.banker == -1 {
		this.banker = tool.Rand(0,3)
	}else{
		this.banker = (this.banker+1)%4
	}
	this.turn = this.banker

	for i := 0;i < 4; i ++ {
		if this.playList[i] == nil {
			this.playList[i] = newPlayData(i)
		}
		this.playList[i].init()
	}

	if len(this.cardList) == 0 {
		for i:=Card_Value_Start;i<=Card_Value_End;i++ {
			for j:=0;j<4;j++ {
				this.cardList = append(this.cardList,i)
			}
		}
	}
	this.cardIndex = 0
}

func (this *MahjongLogic) getPlayData(id int) *PlayData {
	return this.playList[id]
}

func (this *MahjongLogic) checkPeng(id int,card int) {
	for k,v := range this.playList {
		if k == id || v.drawed {
			continue
		}

		if v.checkCardLack(card) {
			continue
		}

		if v.cardMap[card] >= 2 {
			v.canPen = true
		}
	}

}

func (this *MahjongLogic) checkDianGang(id int,card int) {
	//如果没有牌了，则不能再杠
	if len(this.cardList) <= this.cardIndex {
		return
	}

	for k,v := range this.playList {
		if v.drawed || k == id || v.checkCardLack(card) {
			continue
		}

		if v.cardMap[card] >= 3 {
			v.gangCards = append(v.gangCards, card)
			v.canGang = true
		}
	}
}

func (this *MahjongLogic) checkHu(id int,card int) {
	if card >= 0 {
		for sid,playData := range this.playList {
			if sid != id {
				//log.Info("id = %d",sid)
				if playData.drawed {
					continue
				}

				cards := append([]int{}, playData.holds...)
				cards = append(cards, card)

				if checkHu(cards,playData.lack) {
					playData.canHu = true
					playData.huCard = card
					playData.isZiMo = false
				}
			}
		}
	}else {
		//log.Info("id = %d",id)
		playData := this.getPlayData(id)
		if checkHu(playData.holds,playData.lack) {
			playData.canHu = true
			playData.huCard = playData.holds[len(playData.holds)-1]
			playData.isZiMo = true
		}
	}
}

func (this *MahjongLogic) checkAnGang(id int) {
	if len(this.cardList) <= this.cardIndex {
		return
	}
	playData := this.getPlayData(id)
	for card,cnt := range playData.cardMap {
		if !playData.checkCardLack(card) && cnt == 4 {
			playData.gangCards = append(playData.gangCards,card)
			playData.canGang = true
		}
	}
}

func (this *MahjongLogic) checkWanGang(id int) {
	if len(this.cardList) <= this.cardIndex {
		return
	}
	playData := this.getPlayData(id)
	for _,card := range playData.peCards {
		if playData.cardMap[card] > 0 {
			playData.gangCards = append(playData.gangCards,card)
			playData.canGang = true
		}
	}
}

func (this *MahjongLogic) clearOpt(id int) {
	if id >= 0 && id < len(this.playList) {
		playData := this.getPlayData(id)
		playData.canHu = false
		playData.canGang = false
		playData.canPen = false
		if len(playData.gangCards) > 0 {
			playData.gangCards = []int{}
		}
	}else{
		for _, v := range this.playList {
			v.canHu = false
			v.canGang = false
			v.canPen = false
			if len(v.gangCards) > 0 {
				v.gangCards = []int{}
			}
		}
	}
}

func (this *MahjongLogic) hasOpt(id int) bool {
	if id >= 0 && id < len(this.playList) {
		playData := this.getPlayData(id)
		return playData.canGang || playData.canHu || playData.canPen
	}else{
		for _, v := range this.playList {
			if  v.canGang || v.canHu || v.canPen {
				return true
			}
		}
		return false
	}
}

func (this *MahjongLogic) initOptions(id int,card int,state int) {
	if state & Act_State_HU > 0 {
		this.checkHu(id,card)
	}

	if state & Act_State_AnG > 0 {
		this.checkAnGang(id)
	}

	if state & Act_State_WanG > 0 {
		this.checkWanGang(id)
	}

	if state & Act_State_DianG > 0 {
		this.checkDianGang(id,card)
	}

	if state & Act_State_Peng > 0 {
		this.checkPeng(id,card)
	}
}

func (this *MahjongLogic) dealOptions(card int){
	for k,v := range this.playList {
		if v.canGang || v.canHu || v.canPen {
			v.canPen = false
			v.canGang = false
			v.gangCards = []int{}

			this.room.Send(k,"game_action_push",msg.ActionNotify{
				Id: k,
				Card: card,
				CanPeng: v.canPen,
				CanGang: v.canGang,
				CanDraw: v.canHu,
				CardGang: v.gangCards,
			})
		}
	}
}

func (this *MahjongLogic) notifyOptions(id int,card int) bool {
	flag := false
	for k,v := range this.playList {
		if v.canGang || v.canHu || v.canPen {
			flag = true
			this.room.Send(k,"game_action_push",msg.ActionNotify{
				Id: k,
				Card: card,
				CanPeng: v.canPen,
				CanGang: v.canGang,
				CanDraw: v.canHu,
				CardGang: v.gangCards,
			})
		}
	}

	return flag
}

func (this *MahjongLogic)Shuffle(){
	rand.Seed(time.Now().UnixNano())
	for i:=this.cardIndex;i<len(this.cardList);i++{
		j := tool.Rand(this.cardIndex,len(this.cardList)-1)
		this.cardList[i],this.cardList[j] = this.cardList[j],this.cardList[i]
	}

	// 配牌逻辑
	if CnfDebug {
		list :=	tool.IntList(this.cardList)
		for _,v := range CardCnf {
			list.RemoveByValue(v)
		}
		for k,v := range CardCnf {
			list.Insert(v,k+1)
		}
	}
}

func (this *MahjongLogic)GetCardCount() int {
	return len(this.cardList) - this.cardIndex
}

func (this *MahjongLogic)Deal() int {
	if this.cardIndex < len(this.cardList) {
		card := this.cardList[this.cardIndex]
		this.cardIndex++
		//return card
		playData := this.playList[this.turn]
		playData.holds = append(playData.holds,card)
		playData.cardMap[card]++
		return card
	}
	return -1
}

func (this *MahjongLogic)StartDeal(){
	var playData *PlayData
	var card int
	for i:=0;i<4;i++{
		for j:=0;j<13;j++{
			playData = this.playList[i]
			card = this.cardList[this.cardIndex]
			playData.holds = append(playData.holds,card)
			playData.cardMap[card]++
			this.cardIndex ++
		}
	}
	playData = this.playList[this.banker]
	card = this.cardList[this.cardIndex]
	playData.holds = append(playData.holds,card)
	playData.cardMap[card]++
	this.cardIndex ++
}

func (this *MahjongLogic) StartGameMessage() {
	for i := 0; i< len(this.playList); i++ {
		this.room.Send(i,"game_holds_push",this.playList[i].holds)
		this.room.Send(i,"mj_count_push",this.GetCardCount())
		this.room.Send(i,"game_num_push",1)
		this.room.Send(i,"game_begin_push",this.banker)
		if this.config.HuanSanZhang {
			this.state = "huanpai"
			this.room.Send(i,"game_huanpai_push",nil)
		}else{
			this.state = "dingque"
			this.room.Send(i,"game_dingque_push",nil)
		}
	}
}

func (this *MahjongLogic)StartGame(){
	// 初始化数据
	this.init()
	this.Shuffle()
	this.StartDeal()
	this.StartGameMessage()
}

func (this *MahjongLogic) OnReady(id int){
	playerData := this.room.GetPlayerData(id)
	playerData.Ready = true
	this.room.BroadcastEx("user_ready_push",msg.ReadyResponse{
		UserId: this.room.GetUserId(id),
		Ready: true,
	},id);

	for _,v := range this.room.GetAllPlayer() {
		if !v.Ready {
			return
		}
	}

	this.StartGame()
}

func (this *MahjongLogic) MoveNext(id int) *MahjongLogic {
	if id >= 0 && id < 4 {
		this.turn = id
		return this
	}

	for {
		this.turn = (this.turn+1)%4
		if !this.getPlayData(this.turn).drawed {
			return this
		}
	}
}

func (this *MahjongLogic) GainCard(){
	this.playCard = -1
	card := this.Deal()
	if card == -1 {
		//todo:游戏结束
		this.DoGameOver()
		return
	}

	this.room.Broadcast("mj_count_push",this.GetCardCount())
	this.room.Send(this.turn,"game_mopai_push",card)
	playData := this.getPlayData(this.turn)
	playData.canPlay = true
	this.room.Broadcast("game_chupai_push",this.room.GetUserId(this.turn))

	this.initOptions(this.turn,-1,Act_State_HU | Act_State_WanG | Act_State_AnG)
	this.notifyOptions(this.turn,card)


}

func (this *MahjongLogic) checkChanged() bool {
	for i:=0;i<len(this.playList);i++ {
		if this.playList[i].changeCards == nil {
			return false
		}
	}
	return true
}

func (this *MahjongLogic)OnChangeThree(id int,c1,c2,c3 int){
	if this.state != "huanpai" {
		log.Info("state err huanpai")
		return
	}

	playData := this.playList[id]
	if playData.changeCards != nil {
		log.Info("has done change three")
		return
	}

	cards := []int{c1,c2,c3}
	if !playData.checkCards(cards) {
		return
	}

	for _,c := range cards {
		playData.holds = tool.SliceRemoveByValue(playData.holds,c)
		playData.cardMap[c] = playData.cardMap[c] - 1
	}
	playData.changeCards = cards

	this.room.Send(id,"game_holds_push",msg.ChangeThreeResponse{
		UserId: this.room.GetUserId(id),
		ChangeCards: cards,
	})
	this.room.BroadcastEx("game_holds_push",msg.ChangeThreeResponse{
		UserId: this.room.GetUserId(id),
		ChangeCards: []int{},
	},id)

	if this.checkChanged() {
		this.ChangeThreeCard()
	}
}

func (this *MahjongLogic)ChangeThreeCard(){
	var changeThreeFunc = func(playData *PlayData,cards []int){
		for _,card := range cards {
			playData.holds = append(playData.holds, card)
			playData.cardMap[card]++
		}
	}

	for i:=0;i<4;i++ {
		changeThreeFunc(this.playList[i], this.playList[(i+1)%4].changeCards)
	}
	this.state = "dingque"

	changeMethod := msg.ChangeMethodNotify{Method: 1}
	for i:=0;i<len(this.playList);i++ {
		this.room.Send(i,"game_huanpai_over_push",changeMethod)
		this.room.Send(i,"game_holds_push",this.playList[i].holds)
	}
	this.room.Broadcast("game_dingque_push",nil)
}

func (this *MahjongLogic)checkLack() bool {
	for _,playData := range this.playList {
		if playData.lack < 0 {
			return false
		}
	}
	return true
}

func (this *MahjongLogic) OnMakeLack(id int,lack int) {
	if lack < 0 || lack > 3 {
		return
	}

	if this.state != "dingque" {
		return
	}

	if this.playList[id].lack > 0 {
		return
	}
	this.playList[id].lack = lack

	if !this.checkLack() {
		this.room.Broadcast("game_dingque_notify_push",this.room.GetUserId(id))
		return
	}

	array := [4]int{}
	for i,v := range this.playList {
		array[i] = v.lack
	}
	this.room.Broadcast("game_dingque_finish_push",array)
	this.room.Broadcast("game_playing_push",nil)
	//todo:检查所有玩家是否有可听牌的

	this.state = "playing"
	playData := this.playList[this.turn]
	playData.canPlay = true
	this.room.Broadcast("game_chupai_push",this.room.GetUserId(this.turn))
	//todo:检查是否可杠、可胡
}

func (this *MahjongLogic) OnPlayCard (id int,card int) {
	log.Infof("OnPlayCard card:%d",card)
	playData := this.getPlayData(id)
	if this.turn != id || playData.drawed || !playData.canPlay {
		return
	}

	if !playData.checkCards(card) {
		return
	}

	playData.holds = tool.SliceRemoveByValue(playData.holds,card)
	playData.canPlay = false
	playData.cardMap[card]--

	this.playCard = card
	this.playCnt++

	this.room.Broadcast("game_chupai_notify_push",msg.PlayCardNotify{
		UserId: this.room.GetUserId(id),
		Card: card,
	})

	this.initOptions(id,card,Act_State_Peng | Act_State_DianG | Act_State_HU)
	if !this.notifyOptions(id,card) {
		this.room.Broadcast("guo_notify_push",msg.PassNotify{
			UserId: this.room.GetUserId(id),
			Card: this.playCard,
		})
		this.MoveNext(-1).GainCard()
	}
}

func (this *MahjongLogic) othersCanHu(id int) bool {
	index := this.turn
	for {
		index = (index+1)%4
		if index == this.turn {
			return false
		}else{
			playData := this.getPlayData(index)
			if playData.canHu && index != id {
				return true
			}
		}
	}
}

func (this *MahjongLogic) OnPeng(id int) {
	playData := this.getPlayData(id)
	if !playData.canPen || playData.drawed {
		return
	}

	if this.othersCanHu(id) {
		return
	}

	//if playData.cardMap[this.playCard] < 2 {
	//	return
	//}
	playData.holds = tool.SliceRemoveByValue(playData.holds,this.playCard)
	playData.holds = tool.SliceRemoveByValue(playData.holds,this.playCard)
	playData.cardMap[this.playCard] -= 2
	playData.peCards = append(playData.peCards,this.playCard)

	this.room.Broadcast("peng_notify_push",msg.PassNotify{
		UserId: this.room.GetUserId(id),
		Card: this.playCard,
	})
	this.playCard = -1
	this.clearOpt(id)

	this.MoveNext(id)
	playData.canPlay = true

	this.room.Broadcast("game_chupai_push",this.room.GetUserId(id))

}

func (this *MahjongLogic) DoGang(id int,gtype string,card int){
	playData := this.getPlayData(id)
	switch gtype {
	case "angang":
		playData.agCards = append(playData.agCards,card)
		playData.holds = tool.SliceRemoveByValue(playData.holds,card,4)
		playData.cardMap[card] -= 4
	case "diangang":
		playData.dgCards = append(playData.dgCards,card)
		playData.holds = tool.SliceRemoveByValue(playData.holds,card,3)
		playData.cardMap[card] -= 3
	case "wangang":
		playData.wgCards = append(playData.wgCards,card)
		playData.peCards = tool.SliceRemoveByValue(playData.peCards,card)
		playData.holds = tool.SliceRemoveByValue(playData.holds,card,1)
		playData.cardMap[card] -= 1
	}

	this.room.Broadcast("gang_notify_push",msg.GangNotify{
		UserId: this.room.GetUserId(id),
		Card: card,
		Type: gtype,
	})

	this.MoveNext(id).GainCard()
}

func (this *MahjongLogic) OnGang (id int,card int) {
	playData := this.getPlayData(id)
	if !playData.canGang { //|| playData.drawed {
		return
	}

	if this.othersCanHu(id) {
		return
	}

	var gtype string
	switch playData.cardMap[card] {
	case 1:
		gtype = "wangang"
	case 3:
		gtype = "diangang"
	case 4:
		gtype = "angang"
	default:
		return
	}
	fmt.Println(gtype)
	this.clearOpt(-1)

	this.room.Broadcast("hangang_notify_push",id)

	//todo:是否可以抢杠

	this.DoGang(id,gtype,card)
}

func (this *MahjongLogic) OnPass(id int){
	playData := this.getPlayData(id)
	if !(playData.canGang || playData.canPen || playData.canHu) {
		return
	}

	this.room.Send(id,"guo_result",nil)
	this.clearOpt(id)

	if this.hasOpt(-1) {
		return
	}

	this.MoveNext(-1).GainCard()

}

func (this *MahjongLogic)OnHu(id int){
	playData := this.getPlayData(id)
	if playData.drawed {
		return
	}
	playData.drawed = true
	playData.huList = append(playData.huList,playData.huCard)
	// todo:抢杠

	this.room.Broadcast("hu_push",msg.HuNotify{
		Id: id,
		ZiMo: playData.isZiMo,
		Card: playData.huCard,
	})
	this.clearOpt(id)
	this.dealOptions(this.playCard)

	if this.othersCanHu(id) {
		return
	}

	if this.checkGameOver() {
		this.DoGameOver()
		return
	}

	this.clearOpt(-1)
	this.turn = id
	this.MoveNext(-1).GainCard()
	
}

func (this *MahjongLogic)checkGameOver() bool {
	cnt := 0
	for _,playData := range this.playList {
		if playData.drawed {
			cnt++
		}
	}

	return cnt >= 3
}

func (this *MahjongLogic)DoGameOver(){
	players := make([]msg.PlayerSettle,4)
	settles := make([]msg.TotalSettle,4)

	this.room.ClearReady()

	for i:=0;i<4;i++ {
		playData := this.playList[i]
		player := &players[i]
		player.UserId = this.room.GetUserId(i)
		player.Peng = playData.peCards
		player.Actions = []int{}
		player.WangGang = playData.wgCards
		player.DianGang = playData.dgCards
		player.AnGang = playData.agCards
		player.NumOfGen = 0
		player.Holds = playData.holds
		player.Fan = 0
		player.Score = 0
		player.TotalScore = 0
		player.HuInfo = []int{}

		player.Huorder = i

		settle := settles[i]
		settle.NumZiMo = playData.numZiMo
		settle.NumJiaoPai = playData.numJiePao
		settle.NumDianPao = playData.numDianPao
		settle.NumAnGang = playData.numAnGang
		settle.NumMingGang = playData.numMingGang
		settle.NumChaDaJiao = playData.numChaJiao
	}

	this.room.Broadcast("game_over_push",msg.GameOverNotify{
		PlayerSettleList: players,
		TotalSettleList: settles,
	})
}