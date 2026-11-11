package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ultrazg/xyz/constant"
	"github.com/ultrazg/xyz/utils"
)

type FollowingBody struct {
	Uid string `form:"uid"`
}

// FollowingList 查询「我」关注的人
var FollowingList = func(ctx *gin.Context) {
	var params *FollowingBody

	err := ctx.ShouldBind(&params)
	if err != nil {
		utils.ReturnBadRequest(ctx, err)

		return
	}

	if params.Uid == "" {
		utils.ReturnBadRequest(ctx, nil)

		return
	}

	p := map[string]any{
		"uid": params.Uid,
	}
	h := ctx.Request.Header
	XJikeAccessToken := h.Get("x-jike-access-token")
	url := constant.BaseUrl + "/v1/user-relation/list-following"
	headers := map[string]string{
		"x-jike-device-id":    utils.DeviceId,
		"x-jike-access-token": XJikeAccessToken,
		"Content-Type":        "application/json",
	}

	response, code, err := utils.Request(url, http.MethodPost, p, headers)
	if err != nil {
		ctx.JSON(code, gin.H{
			"code": code,
			"msg":  utils.GetMsg(code),
			"data": err.Error(),
		})

		log.Println("/v1/user-relation/list-following", code, utils.GetMsg(code))

		return
	}

	utils.ReturnJson(response, ctx)
}

// FollowerList 查询关注「我」的人
var FollowerList = func(ctx *gin.Context) {
	var params *FollowingBody

	err := ctx.ShouldBind(&params)
	if err != nil {
		utils.ReturnBadRequest(ctx, err)

		return
	}

	if params.Uid == "" {
		utils.ReturnBadRequest(ctx, nil)

		return
	}

	p := map[string]any{
		"uid": params.Uid,
	}
	h := ctx.Request.Header
	XJikeAccessToken := h.Get("x-jike-access-token")
	url := constant.BaseUrl + "/v1/user-relation/list-follower"
	headers := map[string]string{
		"x-jike-device-id":    utils.DeviceId,
		"x-jike-access-token": XJikeAccessToken,
		"Content-Type":        "application/json",
	}

	response, code, err := utils.Request(url, http.MethodPost, p, headers)
	if err != nil {
		ctx.JSON(code, gin.H{
			"code": code,
			"msg":  utils.GetMsg(code),
			"data": err.Error(),
		})

		log.Println("/v1/user-relation/list-follower", code, utils.GetMsg(code))

		return
	}

	utils.ReturnJson(response, ctx)
}

// FollowingUpdatesList 关注动态（谁推荐了哪一期）：真实抓包是 POST body {}（不是 GET）。
var FollowingUpdatesList = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/personal-update/list-following-updates")
}

// SubscriptionOftenPlayed 常听节目
var SubscriptionOftenPlayed = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/subscription/list-often-played")
}
