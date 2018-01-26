# Liver

一个 golang 的简单直播状态收集工具。

## 安装

liver 基于 [jsoniter](https://github.com/json-iterator/go) 和 [gorequest](https://github.com/parnurzeal/gorequest) 。

```bash
go get github.com/json-iterator/go
go get github.com/parnurzeal/gorequest
```

安装 liver 。
`go get github.com/eleflea/liver`

## 用法

改写`settings.go`，注意B站地址不支持 **short room id** ，`go build`编译运行即可。

目前支持的直播网站：

- 战旗
- Bilibili
- 熊猫
- 斗鱼
- 虎牙
- 全民
- 龙珠
- 火猫
