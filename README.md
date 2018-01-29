# Liver

一个基于 golang 的简单直播状态收集工具。

## 安装

liver 基于 [jsoniter](https://github.com/json-iterator/go) 、 [gorequest](https://github.com/parnurzeal/gorequest) 和 [color](https://github.com/fatih/color)。

`go get github.com/eleflea/liver`安装 liver 。

`go build`编译即可。

## 用法

默认读取运行目录下`default.json`，注意B站地址不支持 **short room id**。
使用相对地址`liver [path/to/config/file]`。
使用绝对地址`liver -f [path/to/config/file]`。

### 配置文件

- `show_time`: `bool`是否显示抓取消耗的时间。
- `on_color`: `string`正在直播的该行文字颜色，为以下值之一。black, blue, cyan, green, hiBlack, hiBlue, hiCyan, higreen, hiMagenta, hiRed, hiWhite, hiYellow, magenta, red, white, yellow.
- `off_color`: `string`不在直播的该行文字颜色，同`on_color`。

## 支持的直播网站

目前支持的直播网站：

- 战旗
- Bilibili
- 熊猫
- 斗鱼
- 虎牙
- 全民
- 龙珠
- 火猫
