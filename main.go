package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/json-iterator/go"
	"github.com/parnurzeal/gorequest"
)

const (
	zhanqi  = "zhanqi"
	bili    = "bilibili"
	panda   = "panda"
	douyu   = "douyu"
	huya    = "huya"
	quanmin = "quanmin"
	longzhu = "longzhu"
	huomao  = "huomao"
	unknown = "unknown"

	biliRoomInfoURL    = "https://api.live.bilibili.com/room/v1/Room/get_info?from=room&room_id="
	pandaRoomInfoURL   = "http://www.panda.tv/api_room?roomid="
	longzhuRoomInfoURL = "http://yoyo-api.longzhu.com/api/room/init?domain="
)

// up represents a up
type up struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Islive   bool   `json:"islive"`
	Platform string `json:"platform"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
}

// Ups is todo list about up
type Ups struct {
	Up   []*up         `json:"ups"`
	Len  int           `json:"len"`
	Time time.Duration `json:"time"`
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
}

// replace std json pkg with json-iter
var json = jsoniter.ConfigDefault

func domain(url string) string {
	start := strings.Index(url, "://")
	if start == -1 {
		return unknown
	}
	rest := url[start+3:]
	end := strings.Index(rest, "/")
	if end == -1 {
		end = len(rest)
	}
	// loop reversely find domain between two dot
	dotCount := 0
	domainEnd := 0
	for i := len(rest) - 1; i >= 0; i-- {
		if rest[i] == '.' {
			if dotCount == 0 {
				domainEnd = i
			}
			if dotCount == 1 {
				return rest[i+1 : domainEnd]
			}
			dotCount++
		}
	}
	// if there is no second dot
	return rest[:domainEnd]
}

// get url path between first '/' and next '?'
func tail(url string) string {
	start := strings.LastIndex(url, "/")
	if start == -1 {
		return ""
	}
	end := strings.Index(url, "?")
	if end == -1 {
		return url[start+1:]
	}
	return url[start+1 : end]
}

func main() {
	start := time.Now()
	// load json
	var upSet Ups
	err := json.UnmarshalFromString(config, &upSet)
	if err != nil {
		fmt.Println(err)
	}

	upSet.Len = len(upSet.Up)
	signal := make(chan int, upSet.Len)
	request := gorequest.New()
	// run each goroutine of query
	for _, v := range upSet.Up {
		go mux(v, request, signal)
	}
	// wait all of goroutine end
	for i := upSet.Len; i > 0; i-- {
		<-signal
	}
	upSet.Time = time.Now().Sub(start)

	for _, v := range upSet.Up {
		if v.Islive == true {
			fmt.Printf("%s | %s | %t\n", v.Name, v.Platform, v.Islive)
		}
	}
	for _, v := range upSet.Up {
		if v.Islive == false {
			fmt.Printf("%s | %s | %t\n", v.Name, v.Platform, v.Islive)
		}
	}
	fmt.Println(upSet.Time)
	// press enter to exit
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
