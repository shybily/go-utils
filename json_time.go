package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"
const DefaultTimeZone = "Asia/Shanghai"

func ParseWithLocation(name string, timeStr string) (time.Time, error) {
	locationName := name
	if l, err := time.LoadLocation(locationName); err != nil {
		return time.Time{}, err
	} else {
		lt, _ := time.ParseInLocation(TimeLayout, timeStr, l)
		return lt, nil
	}
}

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

func NewJSONTimeFromString(t string) JSONTime {
	lt, _ := ParseWithLocation(DefaultTimeZone, t)
	return NewJSONTimeFromTime(lt)
}

func NewJSONTimeFromTime(time time.Time) JSONTime {
	t := JSONTime{time}
	return t
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeLayout))
	return []byte(formatted), nil
}

func (t JSONTime) FormatTime() string {
	return fmt.Sprintf("%s", t.Format(TimeLayout))
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
