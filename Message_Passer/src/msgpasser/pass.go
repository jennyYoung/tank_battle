package msgpasser

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const maxMessageLength = 50000
const timeSlot = 200 * time.Millisecond
const tryDisconnectSlotNum = 3
const disconnectSlotNum = 75
const maxPlayerNum = 16

const UIAddr = "127.0.0.1:8888"
const serverAddr = "127.0.0.1:51425"
const defaultPort = 9999
const debug = true
const info = false

type Node struct {
	name  string
	addrS string
	id    int

	addr *net.UDPAddr
	conn net.Conn

	firstConnect bool
	lastSeq      int
	mrs          MessageReceiveState
	ncs          NodeCheckState

	udpSourceSeq int
	udpDestSeq   int
	udpHistory   MessageQueue
	udpStatus    bool
	sendMutex    chan int

	connecting    bool
	tryDisconnect int
	disconnect    bool

	receiveQueue chan Message
	sendQueue    chan Message
	quit         chan int

	crypto CryptoTool
}

func (n *Node) init(name string, addr string, id int) {
	n.name = name
	n.addrS = addr
	n.id = id

	var err error
	n.addr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	n.conn, err = net.DialUDP("udp", nil, n.addr)
	if err != nil {
		log.Fatalln(err)
	}

	n.firstConnect = false
	n.lastSeq = 0
	n.mrs.complete = true
	n.mrs.OriMsg.PeerName = n.name
	n.ncs.complete = true

	n.connecting = false
	n.tryDisconnect = 0
	n.disconnect = false

	n.sendQueue = make(chan Message, 20)
	n.receiveQueue = make(chan Message, 20)
	n.quit = make(chan int, 20)

	n.udpSourceSeq = 1
	n.udpDestSeq = 0
	n.udpStatus = true
	n.udpHistory.Init(1500)
	n.sendMutex = make(chan int, 1)

	n.crypto.Init()
}

type TimeOnNode struct {
	time   time.Time
	nodeId int
}

type Passer struct {
	name string
	port int

	selfNodeId int
	peerNum    int
	nodes      []Node
	name2Node  map[string]*Node

	seq         int
	result      bool
	mss         MessageSendState
	waitingData DataQueue

	init         bool
	readyPeerNum int
	ready        bool
	work         bool
	mutex        chan int

	timeDistributor  bool
	timeDistributed  bool
	timeDeviation    time.Duration
	timeDeviationMap map[string]string

	readyPullQueue MessageQueue
	pullQueue      chan Message

	lastValidTimeCh chan TimeOnNode
	lastValidTimes  []time.Time
	lastValidTime   time.Time

	crypto CryptoTool

	quit chan int

	log *log.Logger
}

func (p *Passer) timeNow() time.Time {
	return time.Now().Add(p.timeDeviation)
}

func (p *Passer) validTime() time.Time {
	if p.mss.IsCompleted() == false {
		return p.mss.OriMsg.Time
	}
	return p.timeNow()
}

func (p *Passer) Broadcast(msg *Message) {
	for i := range p.nodes {
		if p.nodes[i].disconnect == false && p.nodes[i].connecting == false {
			p.nodes[i].sendQueue <- *msg
		}
	}
}

func (p *Passer) Unicast(msg *Message, nodeId int) {
	if p.nodes[nodeId].disconnect == false {
		p.nodes[nodeId].sendQueue <- *msg
	}
}

func (p *Passer) listen(port int) {
	service := fmt.Sprintf(":%d", port)
	addr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		p.log.Fatalln(err)
	}

	ln, err := net.ListenUDP("udp", addr)
	if err != nil {
		p.log.Fatalln(err)
	}

	for {
		b := make([]byte, maxMessageLength)
		n, _, err := ln.ReadFromUDP(b)
		if err != nil {
			p.log.Fatalln(err)
			continue
		}

		var data Data
		err = json.Unmarshal(b[:n], &data)
		if err != nil {
			p.log.Fatalln(err)
		}

		if data.Name != GAME_ROOM_DATA && p.init == false {
			continue
		}

		if data.Forwarded {
			p.Push(b[:n])
			continue
		}

		if data.Name != NLL_MSG && debug || info {
			p.log.Println("RECEIVE", string(b[:n]))
		}

		switch data.Name {
		case ACCPT_MSG, ACCPTD_MSG, NLL_MSG, CHCK_MSG, CHCKD_MSG, CHCK_RSLT_MSG, TIME_MSG:
			var msg Message
			err = json.Unmarshal(b[:n], &msg)
			if err != nil {
				p.log.Fatalln(err)
			}

			node, ok := p.name2Node[msg.PeerName]
			if ok == false {
				p.log.Println("Cannot find specified peer!", msg.PeerName)
				continue
			}

			if p.init == false {
				continue
			}

			if msg.Seq > 0 && node.connecting == false && p.work && msg.Sign != nil {
				h := md5.New()
				for _, data := range msg.Data {
					io.WriteString(h, string(data))
				}
				if bytes.Compare(p.crypto.Dec(msg.Sign), h.Sum(nil)) != 0 {
					p.log.Fatalln("Check inconsistency!!")
				}
			}

			node.receiveQueue <- msg

		case GAME_ROOM_DATA:
			var data GameRoomData
			err = json.Unmarshal(b[:n], &data)
			if err != nil {
				p.log.Fatalln(err)
			}

			if data.State == "close" {
				p.Uninit()

				/* time.Sleep(5 * time.Second)
				cmd := exec.Command("java", "-jar", "ttt.jar")
				err := cmd.Run()
				if err != nil {
					log.Println(err)
				}*/

			} else if data.State == "start" || ((data.State == "join" || data.State == "rejoin") && data.NewPlayer == data.LocalName) {
				p.mutex <- 1
				if p.init {
					p.log.Println("receive old start msg")
					if p.work {
						p.sendToLocal(ConnectData{CNNCT_DATA}, serverAddr)
					}
					<-p.mutex
					continue
				}
				count := data.Count
				names := make([]string, count)
				addrs := make([]string, count)
				for i := 0; i < count; i++ {
					player := data.Players[i]
					names[i] = player.Name
					addrs[i] = fmt.Sprintf("%s:%d", player.Ip, defaultPort)
				}
				p.log.Println(names, addrs, data.LocalName)
				p.Init2(names, addrs, data.LocalName)
				if data.State == "join" || data.State == "rejoin" {
					p.timeDistributor = false
				}
				<-p.mutex
			} else {
				p.mutex <- 1
				var i int
				if data.NewPlayer != data.LocalName {
					for i = 0; i < maxPlayerNum; i++ {
						if p.nodes[i].disconnect == false && p.nodes[i].name == data.NewPlayer {
							break
						}
					}
					if i < maxPlayerNum {
						p.log.Println("The requested join one", data.NewPlayer, "has been in peers")
						if p.nodes[i].connecting == false {
							p.sendToLocal(ConnectData{CNNCT_DATA}, serverAddr)
						}
						<-p.mutex
						continue
					}
					for i = 0; i < maxPlayerNum; i++ {
						if p.nodes[i].disconnect {
							break
						}
					}
					for _, player := range data.Players {
						if player.Name == data.NewPlayer {
							p.nodes[i].init(player.Name, fmt.Sprintf("%s:%d", player.Ip, defaultPort), i)
							p.nodes[i].connecting = true
							p.name2Node[p.nodes[i].name] = &p.nodes[i]
							go p.receive(&p.nodes[i])
							go p.send(&p.nodes[i])
							break
						}
					}
				}
				<-p.mutex
			}

		default:
			p.log.Println("Unhandled type", data.Name)
		}
	}
}

func (p *Passer) writeMessage(msg *Message, n *Node) {
	var extendData []byte
	var err error

	if msg.Seq > 0 && n.connecting == false {
		h := md5.New()
		for _, data := range msg.Data {
			io.WriteString(h, string(data))
		}
		msg.Sign = n.crypto.Enc(h.Sum(nil))
	}

	extendData, err = json.Marshal(msg)
	if err != nil {
		p.log.Fatalln(err)
	}
	_, err = n.conn.Write(extendData)
	if err != nil {
		// p.log.Println(err, string(extendData))
	}

	if msg.Name != NLL_MSG && debug || info {
		p.log.Println("SEND", string(extendData), "to", n.name)
	}
}

func (p *Passer) resendMessage(seq int, n *Node) {
	n.sendMutex <- 1
	defer func() { <-n.sendMutex }()

	msgs := n.udpHistory.GetAllUDPAfter(seq + 1)
	if msgs == nil {
		p.log.Fatalln("resend", n.name, "Cannot find anything in the history?!", seq, n.udpSourceSeq)
	} else {
		/* TODO: can do optimization here */
		for _, msg := range msgs {
			msg.UDPDestSeq = n.udpDestSeq
			msg.UDPStatus = n.udpStatus
			p.writeMessage(&msg, n)
		}
	}
}

func (p *Passer) sendMessage(msg *Message, n *Node) {
	n.sendMutex <- 1
	defer func() { <-n.sendMutex }()

	msg.PeerName = p.name

	msg.UDPSourceSeq = n.udpSourceSeq
	n.udpSourceSeq++
	msg.UDPDestSeq = n.udpDestSeq
	msg.UDPStatus = n.udpStatus
	n.udpHistory.Push(msg)

	p.writeMessage(msg, n)
}

func (p *Passer) send(n *Node) {

	for {
		select {
		case msg := <-n.sendQueue:
			p.sendMessage(&msg, n)
		case <-time.After(timeSlot):
			p.mutex <- 1
			msg := NewMessage(NLL_MSG, p.seq, p.result, p.validTime(), nil)
			if p.seq == 0 || n.connecting || n.udpDestSeq == 0 {
				msg.Data = [][]byte{p.crypto.MarshalPublicKey()}
			}
			<-p.mutex
			p.sendMessage(msg, n)
		case <-n.quit:
			return
		}
	}
}

func (p *Passer) startNextMessage() bool {
	if p.mss.IsCompleted() == false || p.ready == false || p.work == false {
		return false
	}

	started := false
	for {
		data, time := p.waitingData.PullAll()
		if data == nil {
			break
		}
		started = true
		msg := NewMessage(ACCPT_MSG, p.seq, p.result, time, data)
		msg.PeerName = p.name
		p.mss.Init(msg, p)
		p.mss.Start()

		if p.mss.IsCompleted() {
			p.readyPullQueue.TryPush(&p.mss.OriMsg)
			p.lastValidTimeCh <- TimeOnNode{time, p.selfNodeId}
			continue
		}
		break
	}

	return started
}

func (p *Passer) disconnectNode(n *Node, msg *Message) {
	p.log.Println("Disconnect", n.name)
	if n.disconnect {
		return
	}

	if msg != nil {
		if n.mrs.IsCompleted() == false || n.lastSeq < msg.Seq {
			n.mrs.complete = true
			n.mrs.OriMsg = *msg
			p.readyPullQueue.TryPush(msg)
		}
	}

	if n.ncs.IsCompleted() == false {
		n.ncs.Close()
	}

	n.disconnect = true
	n.quit <- 1
	n.quit <- 1
	p.lastValidTimeCh <- TimeOnNode{p.timeNow(), n.id}
	delete(p.name2Node, n.name)
	p.peerNum--
	p.sendToLocal(UnfreezeData{UNFRZ_DATA, n.name, false}, UIAddr)
	p.sendToLocal(DisconnectData{DSCNNCT_DATA, n.name}, serverAddr)

	if p.mss.IsCompleted() == false {
		p.mss.Close()
		p.seq++
		p.result = p.mss.IsSuccessful()
	}
	p.startNextMessage()

}

func (p *Passer) receive(n *Node) {
	var extraSlotNum time.Duration = tryDisconnectSlotNum

	for {
		select {
		case msg := <-n.receiveQueue:

			if msg.UDPStatus == false {
				p.resendMessage(msg.UDPDestSeq, n)
			}
			if msg.UDPSourceSeq != n.udpDestSeq+1 && (msg.Seq != 0 || n.udpDestSeq != 0) {
				if msg.UDPSourceSeq > n.udpDestSeq+1 {
					n.udpStatus = false
				}
				continue
			}
			if n.udpDestSeq == 0 {
				p.log.Println(msg)
				n.crypto.ParsePublicKey(msg.Data[0])
				p.log.Println("Parse public key from", n.name, n.addrS, n.id, n.firstConnect, n.lastSeq, n.udpSourceSeq, n.udpDestSeq, n.crypto)
			}
			n.udpDestSeq = msg.UDPSourceSeq
			n.udpStatus = true

			if n.firstConnect == false {
				if msg.UDPDestSeq == 0 {
					continue
				}
				if p.ready == false {
					p.mutex <- 1
					n.firstConnect = true
					p.readyPeerNum++
					p.log.Println("Ready number inc with seq", msg.UDPDestSeq)
					if p.timeDistributor {
						p.timeDeviationMap[n.name] = time.Now().Sub(msg.Time).String()
					}
					if p.readyPeerNum >= p.peerNum {
						p.ready = true
						p.log.Println("Passer has been ready!")
						if p.timeDistributor {
							b := make([][]byte, 1)
							b[0], _ = json.Marshal(p.timeDeviationMap)
							p.log.Println(p.timeDeviationMap, string(b[0]))
							p.Broadcast(NewMessage(TIME_MSG, 0, true, p.timeNow(), b))

							time.Sleep(500 * time.Millisecond)
							p.work = true
							p.seq = 1
							p.startNextMessage()
							p.sendToLocal(ConnectData{CNNCT_DATA}, serverAddr)
						} else if p.timeDistributed {
							p.work = true
							p.seq = 1
							p.startNextMessage()
							p.sendToLocal(ConnectData{CNNCT_DATA}, serverAddr)
						}
					} else {
						p.ready = false
					}
					<-p.mutex
				} else {
					p.mutex <- 1
					n.firstConnect = true
					p.log.Println("Add", n.name, "to the game!")

					p.timeDeviationMap[n.name] = p.timeNow().Sub(msg.Time).String()
					b := make([][]byte, 1)
					b[0], _ = json.Marshal(p.timeDeviationMap)
					p.log.Println(p.timeDeviationMap, string(b[0]))
					p.Unicast(NewMessage(TIME_MSG, 0, true, p.validTime(), b), n.id)

					time.Sleep(500 * time.Millisecond)
					n.connecting = false
					p.peerNum++
					p.sendToLocal(ConnectData{CNNCT_DATA}, serverAddr)
					<-p.mutex
				}
			}

			if n.tryDisconnect != 0 {
				n.tryDisconnect = 0
				p.sendToLocal(UnfreezeData{UNFRZ_DATA, n.name, true}, UIAddr)
				extraSlotNum = tryDisconnectSlotNum
				if n.ncs.IsCompleted() == false {
					n.ncs.Close()
				}
			}

			switch msg.Name {
			case ACCPTD_MSG:
				p.mutex <- 1
				if p.mss.IsCompleted() == false && p.mss.HandleResponse(&msg, n.id) {
					if p.mss.IsSuccessful() {
						p.readyPullQueue.TryPush(&p.mss.OriMsg)
					}
					p.seq++
					p.result = p.mss.IsSuccessful()
					/* for latency */
					if p.startNextMessage() == false {
						p.Broadcast(NewMessage(NLL_MSG, p.seq, p.result, p.validTime(), nil))
					}
					p.lastValidTimeCh <- TimeOnNode{p.timeNow(), p.selfNodeId}
				}
				p.lastValidTimeCh <- TimeOnNode{msg.Time, n.id}
				<-p.mutex
			case ACCPT_MSG:
				if n.mrs.IsCompleted() == false && n.mrs.HandleResponse(&msg, n.id) {
					if n.mrs.IsSuccessful() {
						p.readyPullQueue.TryPush(&n.mrs.OriMsg)
					}
					/* for latency */
					p.Broadcast(NewMessage(NLL_MSG, p.seq, p.result, p.validTime(), nil))
				}
				if n.mrs.IsCompleted() {
					if msg.Seq != n.lastSeq+1 && n.lastSeq != 0 {
						p.log.Fatalf("It seems that you have to handle package(%d, %d) loss\n", n.lastSeq+1, msg.Seq)
					}
					n.mrs.Init(&msg, NewMessage(ACCPTD_MSG, msg.Seq, true, p.validTime(), nil), n.id, p)
					n.mrs.Start()
					n.lastSeq = msg.Seq
				}
				p.lastValidTimeCh <- TimeOnNode{msg.Time, n.id}
			case NLL_MSG:
				if n.mrs.IsCompleted() == false && n.mrs.HandleResponse(&msg, n.id) {
					if n.mrs.IsSuccessful() {
						p.readyPullQueue.TryPush(&n.mrs.OriMsg)
					}
					/* for latency */
					p.Broadcast(NewMessage(NLL_MSG, p.seq, p.result, p.validTime(), nil))
				}
				if p.work {
					p.lastValidTimeCh <- TimeOnNode{msg.Time, n.id}
				}
			case CHCK_RSLT_MSG:
				var lmsg Message
				json.Unmarshal(msg.Data[0], &lmsg)
				ln, ok := p.name2Node[lmsg.PeerName]
				if ok == false {
					p.log.Println("CHCK_RSLT_MSG Cannot find specified peer!", lmsg)
					continue
				}
				p.mutex <- 1
				p.disconnectNode(ln, &lmsg)
				<-p.mutex
			case CHCKD_MSG:
				var lmsg Message
				json.Unmarshal(msg.Data[0], &lmsg)
				ln, ok := p.name2Node[lmsg.PeerName]
				if ok == false {
					p.log.Println("CHCKD_MSG Cannot find specified peer!", lmsg)
					continue
				}

				if ln.ncs.IsCompleted() == false && ln.ncs.HandleResponse(&msg, n.id) {
					if ln.ncs.IsSuccessful() {
						p.mutex <- 1
						p.disconnectNode(ln, &ln.ncs.resultMsg)
						<-p.mutex
					}
					/*else the check process is unsuccessful?!?*/
				}
			case CHCK_MSG:
				var lmsg Message
				json.Unmarshal(msg.Data[0], &lmsg)
				if lmsg.PeerName == p.name {
					b := make([][]byte, 1)
					b[0], _ = json.Marshal(&p.mss.OriMsg)
					n.sendQueue <- *NewMessage(CHCKD_MSG, msg.Seq, false, msg.Time, b)
				}
				ln, ok := p.name2Node[lmsg.PeerName]
				if ok == false {
					p.log.Println("CHCKD_MSG Cannot find specified peer!", lmsg)
					continue
				}

				b := make([][]byte, 1)
				b[0], _ = json.Marshal(&ln.mrs.OriMsg)
				n.sendQueue <- *NewMessage(CHCKD_MSG, msg.Seq, true, msg.Time, b)
			case TIME_MSG:
				p.mutex <- 1
				if p.work == false {
					var inter interface{}
					err := json.Unmarshal(msg.Data[0], &inter)
					if err != nil {
						p.log.Fatalln("Unmarshal error")
					}
					timeMap, ok := inter.(map[string]interface{})
					if ok == false {
						p.log.Fatalln("Time deviation map gettttt")
					}
					p.log.Println(timeMap)
					p.timeDeviation, err = time.ParseDuration(timeMap[p.name].(string))
					if err != nil {
						p.log.Fatalln("Time deviation gettttt")
					}
					time.Sleep(500 * time.Millisecond)
					p.timeDistributed = true
					if p.ready {
						p.work = true
						p.seq = 1
						p.startNextMessage()
						p.sendToLocal(ConnectData{CNNCT_DATA}, serverAddr)
					}
				}
				<-p.mutex
			}
		case <-time.After(timeSlot * extraSlotNum):
			if p.work == false || n.connecting {
				continue
			}

			if n.tryDisconnect == 0 {
				n.tryDisconnect = 1
				p.log.Println("FREEZE", n.name, "NOW!!!")
				p.sendToLocal(FreezeData{FRZ_DATA, n.name}, UIAddr)
				extraSlotNum = disconnectSlotNum
				continue
			}

			if n.tryDisconnect == 1 {
				n.tryDisconnect = 2
				p.log.Println("Checking", n.name, "connectivity!")
				b := make([][]byte, 1)
				b[0], _ = json.Marshal(&n.mrs.OriMsg)
				n.ncs.Init(NewMessage(CHCK_MSG, 0, true, p.timeNow(), b), n.mrs.OriMsg, p)
				n.ncs.Start()
				extraSlotNum = tryDisconnectSlotNum
				continue
			}

			p.log.Println("Check", n.name, "connectivity but failed!")
			p.log.Println("AM I OUT OF NETWORK?!")
			p.sendToLocal(UnfreezeData{UNFRZ_DATA, p.name, false}, UIAddr)
			p.sendToLocal(DisconnectData{DSCNNCT_DATA, p.name}, serverAddr)
			if p.init {
				p.Uninit()
			}

		case <-n.quit:
			return
		}

	}
}

func (p *Passer) findMinTime(times []time.Time) time.Time {
	mtime := p.timeNow()

	for i, time := range times {
		if mtime.After(time) && p.nodes[i].disconnect == false && p.nodes[i].connecting == false {
			mtime = time
		}
	}

	return mtime
}

func (p *Passer) update() {
	for {
		select {
		case ton := <-p.lastValidTimeCh:
			if ton.nodeId < maxPlayerNum && p.lastValidTimes[ton.nodeId].Before(ton.time) {
				p.lastValidTimes[ton.nodeId] = ton.time
			}
			if p.work == false {
				continue
			}
			lastValidTime := p.findMinTime(p.lastValidTimes)
			selfValidTime := p.validTime()
			if lastValidTime.After(selfValidTime) {
				lastValidTime = selfValidTime
			}
			p.lastValidTime = lastValidTime

			// p.log.Println("Last valid time is", p.lastValidTime, "update", ton.nodeId, p.nodes[ton.nodeId].name, ton.time)

			for {
				msg := p.readyPullQueue.PullMinTime()
				if msg == nil {
					break
				}

				// p.log.Println("Pull new message from ready PULL QUEUE for check")

				if msg.Time.Before(p.lastValidTime) {
					p.pullQueue <- *msg
				} else {
					p.readyPullQueue.TryPush(msg)
					break
				}
			}
		case <-p.quit:
			return
		}
	}
}

func (p *Passer) sendToLocal(j Jsonable, addrS string) {
	fmt.Println(fmt.Sprintf("[%s]", p.name), "***TRIGGER***", string(j.Byte()), "to", addrS)

	addr, err := net.ResolveUDPAddr("udp", addrS)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Write(j.Byte())
	/* if err != nil {
	  p.log.Fatalln(err)
	}*/

	time.Sleep(500 * time.Millisecond)
}

func (p *Passer) sendToUI() {
	var ltime time.Time
	// fout, _ := os.Create(fmt.Sprintf("%s_passerd.txt", p.name))
	// defer fout.Close()

	addr, err := net.ResolveUDPAddr("udp", UIAddr)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case msg := <-p.pullQueue:
			if ltime.After(msg.Time) {
				p.log.Println("Pulled Message time order error", ltime, msg.Time, msg.PeerName)
			}
			ltime = msg.Time

			for _, data := range msg.Data {

				fmt.Println(fmt.Sprintf("[%s]", p.name), "***PULL***", string(data), "of", msg.Name, "from", msg.PeerName, "at", msg.Seq, "at", msg.Time, "after", p.timeNow().Sub(msg.Time))
				// fout.Write(data)
				// fout.Write([]byte("\n"))

				var d UIData
				json.Unmarshal(data, &d)
				d.Time = msg.Time.UnixNano() / 1000000
				d.TimeFrame = p.timeNow().UnixNano() / 1000000
				data, _ = json.Marshal(d)
				fmt.Println(string(data))

				_, err = conn.Write(data)
				/* if err != nil {
					  p.log.Fatalln(err)
				  }*/
			}
		case <-time.After(50 * time.Millisecond):
			var d UIData
			d.Name = "heartbeat"
			d.Time = p.lastValidTime.UnixNano() / 1000000
			d.TimeFrame = p.timeNow().UnixNano() / 1000000
			b, _ := json.Marshal(d)
			_, err = conn.Write(b)

		case <-p.quit:
			return
		}
	}
}

func (p *Passer) Init(port int) {
	p.name2Node = make(map[string]*Node)

	p.log = log.New(os.Stdout, fmt.Sprintf("[default]\t"), log.LstdFlags|log.Lshortfile)
	p.init = false
	p.mutex = make(chan int, 1)

	if p.port != port {
		go p.listen(port)
	}
	p.port = port
	// TODO: sendToUI start, and send ack back
}

func (p *Passer) Init2(names []string, addrs []string, selfName string) {
	p.log = log.New(os.Stdout, fmt.Sprintf("[%s]\t", selfName), log.Lshortfile)

	p.selfNodeId = maxPlayerNum - 1
	p.peerNum = len(names) - 1
	p.name = selfName
	p.nodes = make([]Node, maxPlayerNum)

	p.seq = 0
	p.result = true
	p.mss.complete = true
	p.waitingData.Init(1000, p)

	p.readyPeerNum = 0
	if p.readyPeerNum >= p.peerNum {
		p.ready = true
	} else {
		p.ready = false
	}
	p.work = false
	p.readyPullQueue.Init(500)
	p.pullQueue = make(chan Message, 200)

	p.lastValidTimeCh = make(chan TimeOnNode, 20)
	p.lastValidTimes = make([]time.Time, maxPlayerNum)

	p.timeDistributor = true
	for _, name := range names {
		if name < selfName {
			p.timeDistributor = false
		}
	}
	p.timeDistributed = false
	p.timeDeviationMap = make(map[string]string)
	if p.timeDistributor {
		p.timeDeviation = time.Duration(0)
	}

	p.crypto.Init()
	p.crypto.GenerateKey()

	p.quit = make(chan int, 20)

	j := 0
	for i := range names {
		if names[i] == selfName {
			continue
		}
		p.nodes[j].init(names[i], addrs[i], j)
		p.name2Node[p.nodes[j].name] = &p.nodes[j]
		go p.receive(&p.nodes[j])
		go p.send(&p.nodes[j])
		j += 1
	}
	for ; j < maxPlayerNum; j++ {
		p.nodes[j].disconnect = true
	}
	go p.update()
	go p.sendToUI()
	p.init = true
}

func (p *Passer) Uninit() {
	/*for _, node := range p.nodes {
		if node.disconnect == false {
			node.quit <- 1
			node.quit <- 1
			node.disconnect = true
		}
	}*/
	p.log.Println("Passer uninit now")
	p.init = false
	for i := range p.nodes {
		if p.nodes[i].disconnect == false {
			p.nodes[i].quit <- 1
			p.nodes[i].quit <- 1
			p.nodes[i].disconnect = true
		}
	}
	p.quit <- 1
	p.quit <- 1

	/* back to local time */
	p.timeDeviation = time.Duration(0)

	p.Init(p.port)
}

func (p *Passer) Push(data []byte) {
	p.log.Println("Push", string(data))
	p.mutex <- 1
	time.Sleep(time.Microsecond)
	p.waitingData.TryPush(data)
	p.startNextMessage()
	<-p.mutex
}

/*func (p *Passer) Pull() Message {
	msg := <-p.pullQueue
	prefix := fmt.Sprintf("[%s]", p.name)
	fmt.Println(prefix, "Pull", string(msg.Data), "of", msg.Name, "from", msg.PeerName, "at", msg.Time, "at", time.Now(), "at", msg.Seq)
	return msg
}*/
