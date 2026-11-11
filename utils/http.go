package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 15,
}

var (
	extraHeadersOnce sync.Once
	extraHeaders     map[string]string
)

// 本地联调覆盖：上游对「token + 设备头」是绑定校验的，而 xyz 默认的 iOS 2.57.1
// 请求头模板与新版 App token 不匹配。设置 XYZ_EXTRA_HEADERS（JSON map）后，
// 这里会把其中的请求头覆盖到所有上游请求上，便于用真实账号 token 联调。
// 未设置该环境变量时完全不生效，生产 Docker 构建不受影响。
func loadExtraHeaders() map[string]string {
	extraHeadersOnce.Do(func() {
		raw := os.Getenv("XYZ_EXTRA_HEADERS")
		if raw == "" {
			return
		}
		var parsed map[string]string
		if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
			extraHeaders = parsed
		}
	})
	return extraHeaders
}

func Request(url, method string, body map[string]any, headers map[string]string) (*http.Response, int, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request body: %v", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %v", err)
	}

	if headers != nil {
		// 请求头
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	for key, value := range loadExtraHeaders() {
		req.Header.Set(key, value)
	}

	// local-time 每次请求都用当前时间，格式对齐 App（毫秒 + 无冒号时区，如 2026-11-11T21:57:37.943+0800）。
	// 上游用该头校验/排序播放进度写入；notificationHeaders 默认的 ...Z07:00（无毫秒、时区带冒号）
	// 会让 playback-progress/update 返回 200 但静默丢弃。统一在此覆盖，生产无需任何环境变量。
	req.Header.Set("Local-Time", time.Now().Format("2006-01-02T15:04:05.000-0700"))

	var resp *http.Response
	var retryErr error
	//for i := 0; i < 3; i++ {
	resp, retryErr = client.Do(req)
	//if retryErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
	//	break
	//}

	//time.Sleep(time.Second * 2)
	//}

	//defer resp.Body.Close()

	if retryErr != nil {
		return nil, 0, fmt.Errorf("failed to send request: %v", retryErr)
	}

	// 未登录
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, resp.StatusCode, fmt.Errorf("received 401 status code: %d", resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	return resp, resp.StatusCode, nil
}
