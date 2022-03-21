package relayer

import "time"

type block struct {
	origin   BlockOrigin
	height   uint64
	time     time.Time
	payloads []Payload
}

func NewBlock(origin BlockOrigin, height uint64, time time.Time, payloads []Payload) block {
	return block{
		origin:   origin,
		height:   height,
		time:     time,
		payloads: payloads,
	}
}

func (b block) Origin() BlockOrigin {
	return b.origin
}

func (b block) Height() uint64 {
	return b.height
}

func (b block) Time() time.Time {
	return b.time
}

func (b block) Payloads() []Payload {
	return b.payloads
}
