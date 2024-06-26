package models

import (
	"dy/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Message 消息
type Message struct {
	gorm.Model
	FromId   uint   `json:"fromId"`   // 发送者
	TargetId uint   `json:"targetId"` // 接收者
	Type     int    `json:"type"`     // 发送类型  群聊、广播、私聊
	Media    int    `json:"media"`    // 消息类型  文字、图片、音频
	Content  string `json:"content"`  // 消息内容
	Pic      string `json:"pic"`
	Url      string `json:"url"`
	Desc     string `json:"desc"`
	Amount   int    `json:"amount"` // 其他数字统计
	State    int    `json:"state"`  // 已读0，未读1
}
type VideoOffer struct {
	Sdp  string `json:"sdp"`
	Type string `json:"type"`
}

func GetMessagesByFromIdAndTargetIdAndType(fromId, targetId, typ uint) []Message {
	msgs := make([]Message, 0)
	utils.DB.Where("from_id = ? and target_id = ? and type = ?", fromId, targetId, typ).Find(&msgs)
	return msgs
}

type MessageWithName struct {
	Message
	Name string `json:"name"`
}

func GetGroupMessagesByTargetIdAndType(targetId, typ uint) []Message {
	tmsgs := make([]MessageWithName, 0)
	//utils.DB.Where("target_id = ? and type = ?", targetId, typ).Find(&msgs)
	utils.DB.Table("messages").
		Select("messages.*, user_basic.name as `name`").
		Joins("JOIN user_basic ON messages.from_id = user_basic.id").
		Where("messages.target_id = ? AND messages.type = ?", targetId, typ).
		Find(&tmsgs)
	var msgs = make([]Message, len(tmsgs))
	for i, tmsg := range tmsgs {
		msgs[i] = tmsg.Message
		msgs[i].Desc = tmsg.Name
	}
	return msgs
}

func GetLastMessageByUserIdAndType(fromId, targetId, typ uint) Message {
	var msg Message
	utils.DB.Where("((target_id = ? and from_id = ?) or (target_id = ? and from_id = ?)) and type = ?", fromId, targetId, targetId, fromId, typ).Order("created_at desc").First(&msg)
	return msg
}
func GetGroupLastMessageByTargetIdAndType(targetId, typ uint) Message {
	var msg Message
	utils.DB.Where("target_id = ? and type = ?", targetId, typ).Order("created_at desc").First(&msg)
	return msg
}

func GetNumOfUnreadMessageByUserId(fromId, targetId, typ uint) int64 {
	var messageCount int64
	utils.DB.Model(&Message{}).Where("state = ? AND from_id = ? AND target_id = ? AND type = ?", 1, fromId, targetId, typ).Count(&messageCount)
	return messageCount
}

// Node 链接结点
type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	// 检验token
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isvalid := true // check token
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 用户关系
	// userId 和 node绑定 加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//发送
	go sendProc(node)
	//接收
	go recvProc(node)
	sendMsg(uint(userId), []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws] >>>>> ", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}

		}
	}
}

func recvProc(node *Node) {
	defer func() {
		// 当recvProc退出时，清理掉相应的连接节点
		fmt.Println("有用户掉线，清除结点")
		rwLocker.Lock()
		delete(clientMap, getUserIdByNode(node))
		rwLocker.Unlock()
		node.Conn.Close()
	}()
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data)
		broadMsg(data)
		fmt.Println("[ws] <<<<< ", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func InitProc() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成upd数据发送协程
func udpSendProc() {
	println(utils.CONF.App.Ipv4)
	addrstr := strings.Split(utils.CONF.App.Ipv4, ".")
	addrs := make([]int, len(addrstr))
	for i, addr := range addrstr {
		addrs[i], _ = strconv.Atoi(addr)
	}
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		//IP:   net.IPv4(192, 168, 137, 255),
		IP:   net.IPv4(byte(addrs[0]), byte(addrs[1]), byte(addrs[2]), byte(addrs[3])),
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}

		}
	}
}

// 完成upd数据接收协程
func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		var buf [512]byte
		_, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
		}
		dispatch(buf[0:])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		//fmt.Println(err)
		return
	}
	if msg.FromId == 0 && msg.Content == "ping" {
		return
	}

	switch msg.Type {
	case 1:
		sendMsg(msg.TargetId, data) //私信
	case 2:
		sendGroupMsg(msg, data) //群发
		//case 3:sendAllMsg()//广播
		//case 4:send()

	}
	if msg.Media == 5 {
		fmt.Println("视频数据", string(data))
		return
	}
	//保存一条消息到数据库
	saveMsg(msg)
}

// 私信
func sendMsg(userId uint, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[int64(userId)]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

// 群聊
func sendGroupMsg(msg Message, data []byte) {
	//根据msg.TargetId,也就是群ID,找到这个群的成员
	memberIds := GetGroupMembersByGroupID(msg.TargetId)
	//根据ID找到在线的成员，给在线的成员推送一条消息
	for _, id := range memberIds {
		sendMsg(uint(id), data)
	}
}

// 通过Node获取UserId的辅助函数
func getUserIdByNode(node *Node) int64 {
	for userId, n := range clientMap {
		if n == node {
			return userId
		}
	}
	return 0
}
func IsOnline(userId uint) bool {
	rwLocker.RLock()
	_, ok := clientMap[int64(userId)]
	rwLocker.RUnlock()
	return ok
}

func saveMsg(msg Message) {
	msg.State = 1 //未读
	utils.DB.Create(&msg)
}
