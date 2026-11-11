package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ultrazg/xyz/constant"
	"github.com/ultrazg/xyz/utils"
)

// 评论增强相关端点：全部为只读透传，body 与上游一致。
// 真实抓包（2026-11-11，小宇宙 3.4.2）：
// - /v1/comment/list-thread/v2       body {primaryCommentId, order:"SMART"} → {data, totalCount}
// - /v1/comment/list-mine            body {order:"TIME"|"LIKE", limit}       → {data:[Comment+episode]}
// - /v1/comment/list-timestamp-primary body {owner:{id:eid,type:"EPISODE"}, from, to} → {data:[{id,timestamp,type}]}
// - /v1/comment/get                  GET ?commentId=<id>                     → {data:<Comment>}

// CommentListThreadV2 评论回复（v2）
var CommentListThreadV2 = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/comment/list-thread/v2")
}

// CommentListMine 我的评论
var CommentListMine = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/comment/list-mine")
}

// CommentListTimestampPrimary 时点评论（按音频区间 from/to 查询，返回时点标记）
var CommentListTimestampPrimary = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/comment/list-timestamp-primary")
}

// CommentGet 按 id 取单条评论（时点标记 → 评论正文）
var CommentGet = func(ctx *gin.Context) {
	var req struct {
		CommentID string `json:"commentId"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil || req.CommentID == "" {
		utils.ReturnBadRequest(ctx, err)

		return
	}

	headers := notificationHeaders(ctx.Request.Header.Get("x-jike-access-token"))
	upstream := constant.BaseUrl + "/v1/comment/get?commentId=" + url.QueryEscape(req.CommentID)

	response, code, err := utils.Request(upstream, http.MethodGet, map[string]any{}, headers)
	if err != nil {
		ctx.JSON(code, gin.H{
			"code": code,
			"msg":  utils.GetMsg(code),
			"data": err.Error(),
		})

		log.Println("/v1/comment/get", code, utils.GetMsg(code))

		return
	}

	utils.ReturnJson(response, ctx)
}
