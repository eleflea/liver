package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
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

func pad(str string, length int) string {
	width := (utf8.RuneCountInString(str) + len(str)) / 2
	if length-width < 0 {
		return str
	}
	return str + strings.Repeat(" ", length-width)
}

func load(upSet *Ups) (length, nameMax, platformMax int) {
	err := json.UnmarshalFromString(config, upSet)
	if err != nil {
		log.Fatalln(err)
	}
	length = len(upSet.Up)
	for _, v := range upSet.Up {
		v.Platform = domain(v.URL)
		nameLen := len(v.Name)
		platformLen := len(v.Platform)
		if nameLen > nameMax {
			nameMax = nameLen
		}
		if platformLen > platformMax {
			platformMax = platformLen
		}
	}
	return
}

func errorMark(code int) rune {
	if code != 0 {
		return '*'
	}
	return ' '
}

func main() {
	start := time.Now()
	// load json
	var upSet Ups
	length, nameMax, platformMax := load(&upSet)
	upSet.Len = length
	signal := make(chan int, length)
	request := gorequest.New().Timeout(time.Second * 3)
	// run each goroutine of query
	for _, v := range upSet.Up {
		go mux(v, request, signal)
	}
	// wait all of goroutine end
	for i := length; i > 0; i-- {
		<-signal
	}
	upSet.Time = time.Now().Sub(start)
	// sort and colorful print the result
	sort.Slice(upSet.Up, func(i, j int) bool {
		return upSet.Up[i].Islive
	})
	for _, v := range upSet.Up {
		line := fmt.Sprintf("%s | %s | %t%c", pad(v.Name, nameMax),
			pad(v.Platform, platformMax), v.Islive, errorMark(v.Code))
		if v.Islive == true {
			color.Yellow("%s", line)
		} else {
			fmt.Println(line)
		}
	}
	fmt.Println(upSet.Time)
	// press enter to exit
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
