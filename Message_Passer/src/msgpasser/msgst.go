package msgpasser

import (
	"log"
	"encoding/json"
)

type MessageState interface {
	Start()
	HandleResponse(msg *Message) bool
	Close()
	IsCompleted() bool
	IsSuccessful() bool
}

type MessageSendState struct {
	OriMsg  Message
	msg     *Message
	p				*Passer
	peerNum int

	mutex    chan int
	complete bool
	success  bool

	accepted    []bool
	acceptedNum int
}

func (mss *MessageSendState) Init(msg *Message, p *Passer) {
	mss.OriMsg = *msg
	mss.msg = msg
	mss.p = p

	mss.mutex = make(chan int, 1)
	mss.complete = false
	mss.success = false

	mss.accepted = make([]bool, maxPlayerNum)
	mss.acceptedNum = 0
}

func (mss *MessageSendState) Start() {
	mss.mutex <- 0
	defer func() { <-mss.mutex }()

	mss.p.Broadcast(mss.msg)
	mss.peerNum = mss.p.peerNum

	if mss.p.peerNum == 0 {
		mss.complete = true
		mss.success = true
	}
}

func (mss *MessageSendState) HandleResponse(msg *Message, nodeId int) bool {
	mss.mutex <- 1
	defer func() { <-mss.mutex }()

	if mss.complete {
		return mss.complete
	}

	/* Seq in back message back should be the Seq number of send message */
	if msg.Seq != mss.msg.Seq {
		return mss.complete
	}

	switch msg.Name {
	case ACCPTD_MSG:
		if mss.accepted[nodeId] == false {
			mss.acceptedNum++
			mss.accepted[nodeId] = true
			if mss.acceptedNum == mss.peerNum {
				mss.complete = true
				mss.success = true
			}
		}
	default:
		log.Fatalf("MessageSendState should not handle type %s\n", msg.Name)
	}

	return mss.complete
}

func (mss *MessageSendState) Close() {
	mss.mutex <- 1
	defer func() { <-mss.mutex }()

	mss.complete = true
}

func (mss *MessageSendState) IsCompleted() bool {
	return mss.complete
}

func (mss *MessageSendState) IsSuccessful() bool {
	return mss.success
}

type MessageReceiveState struct {
	OriMsg Message
	msg    *Message
	p			*Passer
	nodeId int

	mutex    chan int
	complete bool
	success  bool

	accepted bool
}

func (mrs *MessageReceiveState) Init(msg *Message, replyMsg *Message, nodeId int, p *Passer) {
	mrs.OriMsg = *msg
	mrs.msg = replyMsg
	mrs.nodeId = nodeId
	mrs.p = p

	mrs.mutex = make(chan int, 1)
	mrs.complete = false
	mrs.success = false

	mrs.accepted = false
}

func (mrs *MessageReceiveState) Start() {
	mrs.mutex <- 1
	defer func() { <-mrs.mutex }()

	mrs.p.Unicast(mrs.msg, mrs.nodeId)
	mrs.accepted = true
}

func (mrs *MessageReceiveState) HandleResponse(msg *Message, nodeId int) bool {
	mrs.mutex <- 1
	defer func() { <-mrs.mutex }()

	if mrs.complete {
		return mrs.complete
	}

	/* Whenever new message comes, the previous one is definitely successful */
	if msg.Seq > mrs.OriMsg.Seq {
		mrs.complete = true
		mrs.success = msg.Result
	}

	return mrs.complete
}

func (mrs *MessageReceiveState) Close() {
	mrs.mutex <- 1
	defer func() { <-mrs.mutex }()

	mrs.complete = true
}

func (mrs *MessageReceiveState) IsCompleted() bool {
	return mrs.complete
}

func (mrs *MessageReceiveState) IsSuccessful() bool {
	return mrs.success
}

type NodeCheckState struct {
	msg       *Message
	latestMsg *Message
	p					*Passer
	peerNum   int

	mutex    chan int
	complete bool
	success  bool

	repliedMsgs []*Message
	repliedNum  int
	seqStat     map[int]int

	maxSeq			int
	maxSeqMsg		*Message

	resultMsg Message
}

func (ncs *NodeCheckState) updateStat(msg *Message, nodeId int) {
	ncs.repliedMsgs[nodeId] = msg
	ncs.repliedNum++
	ncs.seqStat[msg.Seq]++

	if ncs.maxSeq < msg.Seq {
		ncs.maxSeq = msg.Seq
		ncs.resultMsg = *msg
	}
}

func (ncs *NodeCheckState) Init(msg *Message, latestMsg Message, p *Passer) {
	ncs.msg = msg
	ncs.latestMsg = &latestMsg
	ncs.p = p

	ncs.mutex = make(chan int, 1)
	ncs.complete = false
	ncs.success = false

	ncs.repliedMsgs = make([]*Message, maxPlayerNum)
	ncs.repliedNum = 0
	ncs.seqStat = make(map[int]int)
	ncs.maxSeq = -1

	ncs.updateStat(&latestMsg, p.selfNodeId)
}

func (ncs *NodeCheckState) Start() {
	ncs.mutex <- 0
	defer func() { <-ncs.mutex }()

	ncs.p.Broadcast(ncs.msg)
	ncs.peerNum = ncs.p.peerNum
}

func (ncs *NodeCheckState) HandleResponse(msg *Message, nodeId int) bool {
	ncs.mutex <- 1
	defer func() { <-ncs.mutex }()

	if ncs.complete {
		return ncs.complete
	}

	if msg.Seq != ncs.msg.Seq {
		return ncs.complete
	}

	switch msg.Name {
	case CHCKD_MSG:
		if ncs.repliedMsgs[nodeId] == nil {
			ncs.repliedNum++
			ncs.repliedMsgs[nodeId] = new(Message)
			repliedMsg := ncs.repliedMsgs[nodeId]

			json.Unmarshal(msg.Data[0], repliedMsg)
			ncs.updateStat(repliedMsg, nodeId)

			if msg.Result == false {
				ncs.complete = true
				ncs.success = false
				log.Fatalln("Why", msg.PeerName, "think", ncs.latestMsg.PeerName, "is still alive?", msg)
				break
			}
			if ncs.repliedNum  > (ncs.peerNum + 1) / 2 {
				ncs.complete = true
				ncs.success = true

				b := make([][]byte, 1)
				b[0], _ = json.Marshal(ncs.resultMsg)
				ncs.p.Broadcast(NewMessage(CHCK_RSLT_MSG, msg.Seq, true, msg.Time, b))
			}
		}
	case CHCK_RSLT_MSG:
		log.Fatalln("Should not come here")
		ncs.complete = true
		ncs.success = true
		json.Unmarshal(msg.Data[0], &ncs.resultMsg)
	default:
		log.Fatalf("NodeCheckState should not handle type %s\n", msg.Name)
	}

	return ncs.complete
}

func (ncs *NodeCheckState) Close() {
	ncs.mutex <- 1
	defer func() { <-ncs.mutex }()

	ncs.complete = true
}

func (ncs *NodeCheckState) IsCompleted() bool {
	return ncs.complete
}

func (ncs *NodeCheckState) IsSuccessful() bool {
	return ncs.success
}
