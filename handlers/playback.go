package handlers

import "github.com/gin-gonic/gin"

// 播放状态同步（真实抓包 2026-11-11）：
// - POST /v1/playback-progress/list   body {eids:[...]} → {data:[{eid,pid,progress,playedAt}]}
// - POST /v1/playback-progress/update body {data:[{eid,pid,playedAt,progress}]} → {}
// - POST /v1/episode-played/create    body {eid, mediaSource?} → {}
// - POST /v1/episode-played/list-for-cache body {version?, loadMoreKey?} → {data, version, needRefresh}
// 上游按 playedAt 的先后合并多端进度（最后写入者优先），客户端只负责带上时间戳。

// PlaybackProgressList 批量取服务端播放进度（跨端续听）
var PlaybackProgressList = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/playback-progress/list")
}

// PlaybackProgressUpdate 上报播放进度（带 playedAt，服务端做多端合并）
var PlaybackProgressUpdate = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/playback-progress/update")
}

// EpisodePlayedCreate 记录一次「开始收听」（更新最近播放/收听历史）
var EpisodePlayedCreate = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/episode-played/create")
}

// EpisodePlayedListForCache 最近收听本地缓存增量同步（带版本游标，对齐 App 3.4.2）
var EpisodePlayedListForCache = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/episode-played/list-for-cache")
}
