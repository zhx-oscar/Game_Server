package Const

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func UTF8Width(str string) int {
	if len(str) > 0 {
		width := 0
		for i := 0; i < len(str); {
			if str[i] > 239 {
				i += 4
				width += 2
			} else if str[i] > 223 {
				i += 3
				width += 2
			} else if str[i] > 128 {
				i += 2
				width += 2
			} else {
				i += 1
				width += 1
			}
		}
		return width
	} else {
		return 0
	}
}

func GetZoneOpenTime() time.Time {
	s := viper.GetString("ZoneOpenTime")
	time,err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(fmt.Sprintf("GetZoneOpenTime err=%s", err))
	}
	return time
}