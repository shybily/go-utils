package utils

import (
	"encoding/json"
	"testing"
	"time"
)

func TestJSONTime_UnmarshalJSON(t *testing.T) {
	tm := NewJSONTimeFromTime(time.Now())
	b, e := json.Marshal(tm)
	if e != nil {
		t.Fatal(e)
		return
	}
	t1 := JSONTime{}
	e = json.Unmarshal(b, &t1)
	if e != nil {
		t.Fatal(e)
		return
	}
}
