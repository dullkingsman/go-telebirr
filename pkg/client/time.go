package client

import (
	"strconv"
	"time"
)

func GetCurrentUnixTimestampString(utc ...bool) string {
	var t = time.Now()

	if len(utc) > 0 && utc[0] == true {
		t = t.UTC()
	}

	return strconv.FormatInt(t.Unix(), 10)
}

func GetCurrentUnixMilliTimestampString(utc ...bool) string {
	var t = time.Now()

	if len(utc) > 0 && utc[0] == true {
		t = t.UTC()
	}

	return strconv.FormatInt(t.UnixMilli(), 10)
}
