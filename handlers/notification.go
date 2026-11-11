package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ultrazg/xyz/constant"
	"github.com/ultrazg/xyz/utils"
)

// notificationHeaders 通知类接口共用的 iOS 客户端伪装请求头（模板参考 handlers/notice.go）
func notificationHeaders(token string) map[string]string {
	now := time.Now()
	isoTime := now.Format("2006-01-02T15:04:05Z07:00")

	return map[string]string{
		"Host":                        "api.xiaoyuzhoufm.com",
		"User-Agent":                  "Xiaoyuzhou/2.57.1 (build:1576; iOS 17.4.1)",
		"Market":                      "AppStore",
		"App-BuildNo":                 "1576",
		"OS":                          "ios",
		"x-jike-access-token":         token,
		"Manufacturer":                "Apple",
		"BundleID":                    "app.podcast.cosmos",
		"Connection":                  "keep-alive",
		"Accept-Language":             "zh-Hans-CN;q=1.0, zh-Hant-TW;q=0.9",
		"Model":                       "iPhone14,2",
		"app-permissions":             "4",
		"Accept":                      "*/*",
		"Content-Type":                "application/json",
		"App-Version":                 "2.57.1",
		"WifiConnected":               "true",
		"OS-Version":                  "17.4.1",
		"x-custom-xiaoyuzhou-app-dev": "",
		"Local-Time":                  isoTime,
		"Timezone":                    "Asia/Shanghai",
		"x-jike-device-id":            "",
		"x-jike-device-properties":    "",
	}
}

// bindNotificationBody 把客户端 JSON body 原样读成 map 供上游透传；空 body 返回空 map
func bindNotificationBody(ctx *gin.Context) (map[string]any, bool) {
	p := map[string]any{}

	if ctx.Request.ContentLength == 0 {
		return p, true
	}

	if err := ctx.ShouldBindJSON(&p); err != nil && !errors.Is(err, io.EOF) {
		utils.ReturnBadRequest(ctx, err)

		return nil, false
	}

	if p == nil {
		p = map[string]any{}
	}

	return p, true
}

// proxyNotification 把客户端 body 原样 POST 到上游通知接口，并把上游 JSON 返回给客户端
func proxyNotification(ctx *gin.Context, upstreamPath string) {
	p, ok := bindNotificationBody(ctx)
	if !ok {
		return
	}

	h := ctx.Request.Header
	headers := notificationHeaders(h.Get("x-jike-access-token"))

	url := constant.BaseUrl + upstreamPath
	response, code, err := utils.Request(url, http.MethodPost, p, headers)
	if err != nil {
		ctx.JSON(code, gin.H{
			"code": code,
			"msg":  utils.GetMsg(code),
			"data": err.Error(),
		})

		log.Println(upstreamPath, code, utils.GetMsg(code))

		return
	}

	utils.ReturnJson(response, ctx)
}

// NoticeList 互动通知列表
var NoticeList = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/notice/list")
}

// InboxRemove 删除收件箱消息
var InboxRemove = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/inbox/remove")
}

// SystemNoticeList 系统通知列表
var SystemNoticeList = func(ctx *gin.Context) {
	proxyNotification(ctx, "/api/systemNotice/systemNotice/list")
}

// SystemNoticeUnread 系统通知未读数
var SystemNoticeUnread = func(ctx *gin.Context) {
	proxyNotification(ctx, "/api/systemNotice/unread/get")
}
