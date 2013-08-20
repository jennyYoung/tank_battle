package msgpasser

import (
	"encoding/json"
	"log"
	"time"
)

const ACCPT_MSG = "accept"
const ACCPTD_MSG = "accepted"
const NLL_MSG = "null"
const RESEND_MSG = "resend"

const CHCK_MSG = "check"
const CHCKD_MSG = "checked"
const CHCK_RSLT_MSG = "check_result"

const TIME_MSG = "time"

const GAME_ROOM_DATA = "game room info"
const GAME_START_DATA = "game start"

const CNNCT_DATA = "connection success"
const FRZ_DATA = "freeze"
const UNFRZ_DATA = "unfreeze"
const DSCNNCT_DATA = "disconnect"

type Jsonable interface {
	Byte() []byte
}

type Data struct {
	Name      string
	Forwarded bool
}

type Message struct {
	/* Message Type */
	Name     string
	PeerName string
	Seq      int
	/* result of last message */
	Result bool
	/* time this message represents */
	Time time.Time

	Data [][]byte
	Sign []byte

	UDPSourceSeq int
	UDPDestSeq   int
	UDPStatus    bool
}

func NewMessage(name string, seq int, result bool, time time.Time, data [][]byte) *Message {
	return &Message{name, "", seq, result, time, data, nil, 0, 0, false}
}

func (msg *Message) WithoutData() *Message {
	msg.Data = nil
	return msg
}

func (msg Message) Byte() (b []byte) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln("Json UIData fail!")
	}
	return
}

func (msg Message) String() (s string) {
	s = string(msg.Byte())
	return
}

type UIData struct {
	Name       string
	Time       int64
	TimeFrame  int64
	Tank_id    int
	Key_code   int
	X          int
	Y          int
	Forwarded  bool
	Tank_life  int
	Tank_score int
}

func (d *UIData) Byte() (b []byte) {
	b, err := json.Marshal(d)
	if err != nil {
		log.Fatalln("Json UIData fail!")
	}
	return
}

type PlayerInfo struct {
	Ip   string
	Name string
}

type GameRoomData struct {
	Name      string
	Id        int
	Players   []PlayerInfo
	Count     int
	LocalName string
	State     string
	NewPlayer string
}

type ConnectData struct {
	Name string
}

func (d ConnectData) Byte() (b []byte) {
	b, err := json.Marshal(d)
	if err != nil {
		log.Fatalln("Json UIData fail!")
	}
	return
}

type FreezeData struct {
	Name string
	Peer string
}

func (d FreezeData) Byte() (b []byte) {
	b, err := json.Marshal(d)
	if err != nil {
		log.Fatalln("Json UIData fail!")
	}
	return
}

type UnfreezeData struct {
	Name      string
	Peer      string
	Forwarded bool
}

func (d UnfreezeData) Byte() (b []byte) {
	b, err := json.Marshal(d)
	if err != nil {
		log.Fatalln("Json UIData fail!")
	}
	return
}

type DisconnectData struct {
	Name string
	Peer string
}

func (d DisconnectData) Byte() (b []byte) {
	b, err := json.Marshal(d)
	if err != nil {
		log.Fatalln("Json UIData fail!")
	}
	return
}
