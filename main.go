package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/fatih/color"
	"github.com/json-iterator/go"
	"github.com/parnurzeal/gorequest"
)

const (
	defaultConfig = "default.json"

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

type settings struct {
	ShowTime             bool   `json:"show_time"`
	OnColor              string `json:"on_color"`
	OffColor             string `json:"off_color"`
	nameMax, platformMax int
}

// Ups is todo list about up
type Ups struct {
	Up       []*up         `json:"ups"`
	Len      int           `json:"len"`
	Time     time.Duration `json:"time"`
	Code     int           `json:"code"`
	Msg      string        `json:"msg"`
	Settings settings      `json:"settings"`
}

// replace std json pkg with json-iter
var json = jsoniter.ConfigDefault

// colorful output map
var colorMap = map[string]func(string, ...interface{}){
	"black":     color.Black,
	"blue":      color.Blue,
	"cyan":      color.Cyan,
	"green":     color.Green,
	"hiBlack":   color.HiBlack,
	"hiBlue":    color.HiBlue,
	"hiCyan":    color.HiCyan,
	"higreen":   color.HiGreen,
	"hiMagenta": color.HiMagenta,
	"hiRed":     color.HiRed,
	"hiWhite":   color.HiWhite,
	"hiYellow":  color.HiYellow,
	"magenta":   color.Magenta,
	"red":       color.Red,
	"white":     color.White,
	"yellow":    color.Yellow,
}

func handleFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

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

func strWidth(str string) int {
	return (utf8.RuneCountInString(str) + len(str)) / 2
}

func pad(str string, length int) string {
	num := length - strWidth(str)
	if num < 0 {
		return str
	}
	return str + strings.Repeat(" ", num)
}

func load(fileName string, upSet *Ups) error {
	var nameMax, platformMax int
	config, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrap(err, "read config fail")
	}
	err = json.Unmarshal(config, upSet)
	if err != nil {
		return errors.Wrap(err, "config format wrong")
	}
	upSet.Len = len(upSet.Up)
	for _, v := range upSet.Up {
		v.Platform = domain(v.URL)
		nameLen := strWidth(v.Name)
		platformLen := strWidth(v.Platform)
		if nameLen > nameMax {
			nameMax = nameLen
		}
		if platformLen > platformMax {
			platformMax = platformLen
		}
	}
	upSet.Settings.nameMax = nameMax
	upSet.Settings.platformMax = platformMax
	return nil
}

func errorMark(code int) rune {
	if code != 0 {
		return '*'
	}
	return ' '
}

func onOff(status bool) string {
	if status == true {
		return "ON"
	}
	return "OFF"
}

func displayWithColor(str, color string) {
	display, ok := colorMap[color]
	if ok == true {
		display("%s", str)
	} else {
		fmt.Println(str)
	}
}

func show(fileName string) {
	start := time.Now()
	// load json
	var upSet Ups
	err := load(fileName, &upSet)
	handleFatal(err)
	signal := make(chan int, upSet.Len)
	request := gorequest.New().Timeout(time.Second * 3)
	// run each goroutine of query
	for _, v := range upSet.Up {
		go mux(v, request, signal)
	}
	// wait all of goroutine end
	for i := upSet.Len; i > 0; i-- {
		<-signal
	}
	upSet.Time = time.Now().Sub(start)
	// sort and colorful print the result
	sort.Slice(upSet.Up, func(i, j int) bool {
		return upSet.Up[i].Islive
	})
	set := upSet.Settings
	for _, v := range upSet.Up {
		line := fmt.Sprintf("%s | %s | %s%c", pad(v.Name, set.nameMax),
			pad(v.Platform, set.platformMax), onOff(v.Islive), errorMark(v.Code))
		if v.Islive == true {
			displayWithColor(line, set.OnColor)
		} else {
			displayWithColor(line, set.OffColor)
		}
	}
	if set.ShowTime != false {
		fmt.Println(upSet.Time)
	}
}

func exPath() string {
	ex, err := os.Executable()
	handleFatal(err)
	return filepath.Dir(ex)
}

func fileName(args []string) string {
	switch length := len(args); length {
	case 0:
		return filepath.Join(exPath(), defaultConfig)
	case 1:
		return filepath.Join(exPath(), args[0])
	case 2:
		if args[0] != "-f" {
			handleFatal(errors.New(`unknown args "` + args[0] + `"`))
		}
		return args[1]
	default:
		return filepath.Join(exPath(), defaultConfig)
	}
}

func main() {
	fn := fileName(os.Args[1:])
	// run and show result
	show(fn)
	// press enter to exit
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
