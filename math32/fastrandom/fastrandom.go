package fastrandom

import (
	"math/rand"
	"time"
)

type XorShift struct {
	s [2]uint64
}

func New(seed uint32) *XorShift {
	xs := &XorShift{}
	if seed == 0 {
		t := uint64(time.Now().UnixNano())
		xs.s[0] = t
		xs.s[1] = t ^ (t >> 32) ^ 0xbadcaffe
	} else {
		xs.s[0] = uint64(seed)
		xs.s[1] = uint64(seed) ^ 0xdeadbeef
	}

	if xs.s[0] == 0 && xs.s[1] == 0 {
		xs.s[0] = 1
	}

	return xs
}

func NewWithDefaults() *XorShift {
	return New(rand.Uint32())
}

func (xs *XorShift) nextUint64() uint64 {
	s1 := xs.s[0]
	s0 := xs.s[1]

	result := s0 + s1

	xs.s[0] = s0
	s1 ^= s1 << 23
	s1 ^= s1 >> 17
	s1 ^= s0
	s1 ^= s0 >> 26
	xs.s[1] = s1

	return result
}

// Float32 generates the next pseudo-random number in the sequence
// and returns it as a float32 within the range [0, 1).
func (xs *XorShift) Float32() float32 {
	val := xs.nextUint64()

	return float32(val>>40) / 16777216.0
}
