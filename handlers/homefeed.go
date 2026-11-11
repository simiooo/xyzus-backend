package handlers

import (
	"github.com/gin-gonic/gin"
)

// 首页 feed / 收件箱 v2 / 相关单集等上游代理端点。
// 与通知类接口共用 handlers/notification.go 里的 bindNotificationBody + proxyNotification：
// 客户端 JSON body 原样透传给上游，上游 JSON 原样返回，不额外改写字段。

// InboxListV2 收件箱列表（v2，body 透传 limit / loadMoreKey）
var InboxListV2 = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v2/inbox/list")
}

// TabFeaturedList 首页「精选」tab 列表（body 透传 loadMoreKey / returnAll / forceRefresh）
var TabFeaturedList = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/tab/featured/list")
}

// TabDiscoveryList 首页「发现」tab 列表（body 透传 returnAll / forceRefresh）
var TabDiscoveryList = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/tab/discovery/list")
}

// TabCategoryList 首页「分类」tab 列表（body 透传 loadMoreKey）
var TabCategoryList = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/tab/category/list")
}

// RelatedEpisodeList 相关单集列表（body 透传 eid）
var RelatedEpisodeList = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/related-episode/list")
}
