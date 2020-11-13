package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"
const DefaultTimeZone = "Asia/Shanghai"

func ParseWithLocation(locationName string, timeStr string) (time.Time, error) {
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

func NewJSONTimeFromString(t string) (JSONTime, error) {
	if lt, err := ParseWithLocation(DefaultTimeZone, t); err != nil {
		return JSONTime{}, err
	} else {
		return NewJSONTimeFromTime(lt), nil
	}

}

func NewJSONTimeFromTime(time time.Time) JSONTime {
	t := JSONTime{time}
	return t
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.FormatTime())
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if t1, err := NewJSONTimeFromString(s); err != nil {
		return err
	} else {
		*t = t1
	}
	return nil
}

func (t JSONTime) FormatTime() string {
	return fmt.Sprintf("%s", t.Format(TimeLayout))
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan value of time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
