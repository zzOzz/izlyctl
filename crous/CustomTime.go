package crous

import (
	"github.com/Sirupsen/logrus"
	"strings"
	"time"
	"fmt"
)
//const ctLayout = "2006-01-02T15:04:05"
const inputCtLayout = "2006-01-02T15:04:05.000-07"
//const outputCtLayout = "2006-01-02 15:04:05.000000000 -0700 MST"
const outputCtLayout = "2006-01-02T15:04:05.000-07"

var nilTime = (time.Time{}).UnixNano()

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	logrus.Debugf(">>> I'm UnmarshalJSON  %s", s)
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(inputCtLayout, s)
	return
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	logrus.Debugf(">>> I'm Marshalling  %s", ct)
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	//return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(outputCtLayout))), nil
}