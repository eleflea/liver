package main

import (
	"strings"

	"github.com/parnurzeal/gorequest"
)

func mux(u *up, request *gorequest.SuperAgent, signal chan int) {
	switch u.Platform {
	case zhanqi:
		request.Get(u.URL).End(func(resp gorequest.Response, body string, errs []error) {
			getZhanqi(body, errs, u)
		})
	case bili:
		id := tail(u.URL)
		request.Get(biliRoomInfoURL + id).EndBytes(func(resp gorequest.Response, body []byte, errs []error) {
			getBili(body, errs, u)
		})
	case panda:
		id := tail(u.URL)
		request.Get(pandaRoomInfoURL + id).EndBytes(func(resp gorequest.Response, body []byte, errs []error) {
			getPanda(body, errs, u)
		})
	case douyu:
		request.Get(u.URL).End(func(resp gorequest.Response, body string, errs []error) {
			getDouyu(body, errs, u)
		})
	case huya:
		request.Get(u.URL).End(func(resp gorequest.Response, body string, errs []error) {
			getHuya(body, errs, u)
		})
	case quanmin:
		request.Get(u.URL).End(func(resp gorequest.Response, body string, errs []error) {
			getQuanmin(body, errs, u)
		})
	case longzhu:
		id := tail(u.URL)
		request.Get(longzhuRoomInfoURL + id).EndBytes(func(resp gorequest.Response, body []byte, errs []error) {
			getLongzhu(body, errs, u)
		})
	case huomao:
		request.Get(u.URL).End(func(resp gorequest.Response, body string, errs []error) {
			getHuomao(body, errs, u)
		})
	default:
		u.Islive = false
		u.Code = 3
		u.Msg = "unsupport site error"
	}
	signal <- 0
	return
}

func getZhanqi(body string, errs []error, u *up) {
	if len(errs) != 0 {
		u.Islive = false
		u.Code = 1
		u.Msg = "get page error"
		return
	}
	start := strings.Index(body, `","status":"`)
	if start == -1 {
		u.Islive = false
		u.Code = 2
		u.Msg = "search room status error"
		return
	}
	if body[start+12] == '4' {
		u.Islive = true
		return
	}
	u.Islive = false
}

func getBili(body []byte, errs []error, u *up) {
	if len(errs) != 0 {
		u.Islive = false
		u.Code = 1
		u.Msg = "get page error"
		return
	}
	if json.Get(body, "data", "live_status").ToInt() == 1 {
		u.Islive = true
		return
	}
	u.Islive = false
}

func getPanda(body []byte, errs []error, u *up) {
	if len(errs) != 0 {
		u.Islive = false
		u.Code = 1
		u.Msg = "get page error"
		return
	}
	if json.Get(body, "data", "videoinfo", "status").ToString() == "2" {
		u.Islive = true
		return
	}
	u.Islive = false
}

func getDouyu(body string, errs []error, u *up) {
	if len(errs) != 0 {
		u.Islive = false
		u.Code = 1
		u.Msg = "get page error"
		return
	}
	start := strings.Index(body, `,"show_status":`)
	if start == -1 {
		u.Islive = false
		u.Code = 2
		u.Msg = "search room status error"
		return
	}
	if body[start+15] == '1' {
		u.Islive = true
		return
	}
	u.Islive = false
}

func getHuya(body string, errs []error, u *up) {
	if len(errs) != 0 {
		u.Islive = false
		u.Code = 1
		u.Msg = "get page error"
		return
	}
	start := strings.Index(body, `","state":"`)
	if start == -1 {
		u.Islive = false
		u.Code = 2
		u.Msg = "search room status error"
		return
	}
	if body[start+11:start+13] == "ON" {
		u.Islive = true
		return
	}
	u.Islive = false
}

func getQuanmin(body string, errs []error, u *up) {
	if len(errs) != 0 {
		u.Islive = false
		u.Code = 1
		u.Msg = "get page error"
		return
	}
	start := strings.Index(body, `","status":"`)
	if start == -1 {
		u.Islive = false
		u.Code = 2
		u.Msg = "search room status error"
		return
	}
	if body[start+12] == '2' {
		u.Islive = true
		return
	}
	u.Islive = false
}

func getLongzhu(body []byte, errs []error, u *up) {
	if len(errs) != 0 {
		u.Islive = false
		u.Code = 1
		u.Msg = "get page error"
		return
	}
	if json.Get(body, "data", "live", "isLive").ToBool() == true {
		u.Islive = true
		return
	}
	u.Islive = false
}

func getHuomao(body string, errs []error, u *up) {
	if len(errs) != 0 {
		u.Islive = false
		u.Code = 1
		u.Msg = "get page error"
		return
	}
	start := strings.Index(body, `,"is_live":`)
	if start == -1 {
		u.Islive = false
		u.Code = 2
		u.Msg = "search room status error"
		return
	}
	if body[start+11] == '1' {
		u.Islive = true
		return
	}
	u.Islive = false
}
