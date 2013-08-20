package msgpasser

import (
	"log"
	"time"
)

type MessageQueue struct {
	q          []Message
	capa       int
	start, end int

	mutex   chan int
	lastSeq int
}

func (msq *MessageQueue) Init(capa int) {
	if capa < 5 {
		log.Fatalln("Too small capacity")
	}
	msq.capa = capa
	msq.q = make([]Message, capa)
	msq.start = 0
	msq.end = 0

	msq.mutex = make(chan int, 1)
	msq.lastSeq = 0
}

func (msq *MessageQueue) inc(v int) int {
	v = (v + 1) % msq.capa
	return v
}

func (msq *MessageQueue) size() int {
	return (msq.end + msq.capa - msq.start) % msq.capa
}

func (msq *MessageQueue) PushInOrder(msg *Message) {
	msq.mutex <- 1
	defer func() { <-msq.mutex }()

	if msq.lastSeq != 0 && msg.Seq != msq.lastSeq+1 {
		log.Fatalln("Wrong message to push!")
	}
	msq.lastSeq = msg.Seq

	if msq.size() == msq.capa-1 {
		msq.start = msq.inc(msq.start)
	}
	msq.q[msq.end] = *msg
	msq.end = msq.inc(msq.end)
}

func (msq *MessageQueue) PullAllAfter(seq int) []Message {
	msq.mutex <- 1
	defer func() { <-msq.mutex }()

	var nq []Message
	record := false
	off := 0
	for i := msq.start; i != msq.end; i = msq.inc(i) {
		if msq.q[i].Seq == seq {
			record = true
			nq = make([]Message, (msq.end+msq.capa-i)%msq.capa)
		}
		if record {
			nq[off] = msq.q[i]
			off++
		}
	}

	return nq
}

func (msq *MessageQueue) TryPush(msg *Message) {
	msq.mutex <- 1
	defer func() { <-msq.mutex }()

	if msq.size() == msq.capa-1 {
		log.Println("Cannot push anymore! Is this an error?")
		return
	}
	msq.q[msq.end] = *msg
	msq.end = msq.inc(msq.end)
}

func (msq *MessageQueue) Pull() *Message {
	msq.mutex <- 1
	defer func() { <-msq.mutex }()

	if msq.size() == 0 {
		return nil
	}

	msg := msq.q[msq.start]
	msq.start = msq.inc(msq.start)

	return &msg
}

func (msq *MessageQueue) PullMinTime() *Message {
	msq.mutex <- 1
	defer func() { <-msq.mutex }()

	if msq.size() == 0 {
		return nil
	}

	mi := msq.start
	mtime := msq.q[mi].Time
	for i := msq.start; i != msq.end; i = msq.inc(i) {
		if mtime.After(msq.q[i].Time) ||
			mtime.Equal(msq.q[i].Time) && msq.q[mi].Name > msq.q[i].Name ||
			mtime.Equal(msq.q[i].Time) && msq.q[mi].Name == msq.q[i].Name && msq.q[mi].Seq > msq.q[i].Seq {
			mi = i
			mtime = msq.q[i].Time
		}
	}

	msg := msq.q[msq.start]
	for i := msq.start; i != mi; i = msq.inc(i) {
		tmsg := msq.q[msq.inc(i)]
		msq.q[msq.inc(i)] = msg
		msg = tmsg
	}
	msq.start = msq.inc(msq.start)

	return &msg
}

func (msq *MessageQueue) Push(msg *Message) {
	msq.mutex <- 1
	defer func() { <-msq.mutex }()

	if msq.size() == msq.capa-1 {
		msq.start = msq.inc(msq.start)
	}
	msq.q[msq.end] = *msg
	msq.end = msq.inc(msq.end)
}

func (msq *MessageQueue) GetAllUDPAfter(seq int) []Message {
	msq.mutex <- 1
	defer func() { <-msq.mutex }()

	var nq []Message
	record := false
	off := 0
	for i := msq.start; i != msq.end; i = msq.inc(i) {
		if msq.q[i].UDPSourceSeq == seq {
			record = true
			nq = make([]Message, (msq.end+msq.capa-i)%msq.capa)
		}
		if record {
			nq[off] = msq.q[i]
			off++
		}
	}

	return nq
}

type DataQueue struct {
	q          [][]byte
	capa       int
	start, end int
	p          *Passer

	mutex chan int
	time  time.Time
}

func (dq *DataQueue) Init(capa int, p *Passer) {
	if capa < 5 {
		log.Fatalln("Too small capacity")
	}
	dq.capa = capa
	dq.q = make([][]byte, capa)
	dq.start = 0
	dq.end = 0
	dq.p = p

	dq.mutex = make(chan int, 1)
}

func (dq *DataQueue) inc(v int) int {
	v = (v + 1) % dq.capa
	return v
}

func (dq *DataQueue) size() int {
	return (dq.end + dq.capa - dq.start) % dq.capa
}

func (dq *DataQueue) TryPush(data []byte) {
	dq.mutex <- 1
	defer func() { <-dq.mutex }()

	if dq.size() == dq.capa-1 {
		log.Println("Cannot push anymore! Is this an error?")
		return
	}
	if dq.size() == 0 {
		dq.time = dq.p.timeNow()
	}
	dq.q[dq.end] = data
	dq.end = dq.inc(dq.end)
}

func (dq *DataQueue) PullAll() ([][]byte, time.Time) {
	dq.mutex <- 1
	defer func() { <-dq.mutex }()

	if dq.size() == 0 {
		return nil, dq.time
	}

	nq := make([][]byte, dq.size())
	off := 0
	for i := dq.start; i != dq.end; i = dq.inc(i) {
		nq[off] = dq.q[i]
		off++
	}
	dq.start = 0
	dq.end = 0

	return nq, dq.time
}
