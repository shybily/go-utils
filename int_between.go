package utils

import (
	"math"
	"math/rand"
	"time"
)

type IntBetween []int

func (t IntBetween) Rand() int {
	if t.Max() == t.Min() {
		return t.Max()
	}
	return rand.Intn(t.Max()-t.Min()) + t.Min()
}

func (t IntBetween) RandTimeDuration() time.Duration {
	r := t.Rand()
	return time.Duration(r) * time.Second
}

func (t IntBetween) Speed(count float64) time.Duration {
	r := float64(t.Rand())
	s := int(math.Floor(r / 60 * count))
	return time.Second * time.Duration(s)
}

func (t IntBetween) Max() int {
	if t[0] >= t[1] {
		return t[0]
	} else {
		return t[1]
	}
}

func (t IntBetween) Min() int {
	if t[0] <= t[1] {
		return t[0]
	} else {
		return t[1]
	}
}
