package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type TimeStringBetween struct {
	Src   []string
	start time.Time
	end   time.Time
}

func NewTimeStringBetween(data []byte) (*TimeStringBetween, error) {
	t := &TimeStringBetween{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TimeStringBetween) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Src)
}

func (t *TimeStringBetween) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &t.Src)
	if err != nil {
		return err
	}
	if len(t.Src) != 2 {
		return fmt.Errorf("invalid length")
	}
	var (
		timezone = "Asia/Shanghai"
		tmp      = make([]time.Time, 2)
		date     = time.Now().Format("2006-01-02")
	)
	for i, s := range t.Src {
		tmp[i], err = ParseWithLocation(timezone, fmt.Sprintf("%s %s:00", date, s))
		if err != nil {
			return err
		}
	}
	if tmp[0].Before(tmp[1]) {
		t.start = tmp[0]
		t.end = tmp[1]
	} else {
		t.start = tmp[1]
		t.end = tmp[0]
	}
	return nil
}

func (t *TimeStringBetween) Between(target time.Time) bool {
	return target.After(t.start) && target.Before(t.end)
}

func (t *TimeStringBetween) End() time.Time {
	return t.end
}

func (t *TimeStringBetween) Start() time.Time {
	return t.start
}
