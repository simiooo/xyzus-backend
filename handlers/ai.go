package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ultrazg/xyz/utils"
)

// AI 摘要 / 文字稿：只读透传。
// 真实抓包（2026-11-11）：
// - /v1/ai-summary/episode/retrieve-brief body {eid} → {requested, summary:{status,brief,...}|null, confirmation}
// - /v1/episode-transcript/get body {eid, mediaId} → {key, vendorDeclaration, transcriptUrl, highlightUrl}
//
// 文字稿正文放在 CDN 的签名地址上（transcript.xyzcdn.net / transcript-highlight.xyzcdn.net），
// 直接让浏览器 fetch 会被跨域/鉴权拦。这里加一个受限的本地代理：只允许上面两个 host，
// 由后端带上 App 的 UA/Referer 去取，返回解析后的 JSON（通常是一段 [{text,startMs}]）。

var (
	transcriptHTTPClient   = &http.Client{Timeout: 25 * time.Second}
	allowedTranscriptHosts = map[string]bool{
		"transcript.xyzcdn.net":           true,
		"transcript-highlight.xyzcdn.net": true,
	}
)

const maxTranscriptBytes = 8 << 20

// AISummaryEpisodeRetrieveBrief 单集 AI 摘要
var AISummaryEpisodeRetrieveBrief = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/ai-summary/episode/retrieve-brief")
}

// EpisodeTranscriptGet 单集文字稿元数据（返回 CDN 签名地址）
var EpisodeTranscriptGet = func(ctx *gin.Context) {
	proxyNotification(ctx, "/v1/episode-transcript/get")
}

// EpisodeTranscriptContent 取文字稿正文：body {url:<transcriptUrl>} → data 为 CDN 返回的 JSON
var EpisodeTranscriptContent = func(ctx *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ReturnBadRequest(ctx, err)

		return
	}

	parsed, err := url.Parse(strings.TrimSpace(req.URL))
	if err != nil || parsed.Scheme != "https" || !allowedTranscriptHosts[parsed.Hostname()] {
		utils.ReturnBadRequest(ctx, errors.New("transcript url not allowed"))

		return
	}

	httpReq, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
	if err != nil {
		utils.ReturnBadRequest(ctx, err)

		return
	}

	// CDN 会校验 Referer/UA，这里固定成 App 的请求特征（与抓包一致）
	httpReq.Header.Set("User-Agent", "Xiaoyuzhou/3.4.2 (android 34)")
	httpReq.Header.Set("Referer", "https://www.xiaoyuzhoufm.com/")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := transcriptHTTPClient.Do(httpReq)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"code": http.StatusBadGateway, "msg": utils.GetMsg(http.StatusBadGateway)})

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		ctx.JSON(http.StatusBadGateway, gin.H{"code": http.StatusBadGateway, "msg": utils.GetMsg(http.StatusBadGateway)})

		return
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxTranscriptBytes+1))
	if err != nil || len(body) > maxTranscriptBytes {
		ctx.JSON(http.StatusBadGateway, gin.H{"code": http.StatusBadGateway, "msg": utils.GetMsg(http.StatusBadGateway)})

		return
	}

	var data any
	if err = json.Unmarshal(body, &data); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"code": http.StatusBadGateway, "msg": utils.GetMsg(http.StatusBadGateway)})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": utils.GetMsg(http.StatusOK), "data": data})
}
