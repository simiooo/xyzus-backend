> [!WARNING]
> 通过第三方登录会使小宇宙账号遭到**风控**或**封禁**，请谨慎使用。

<h1 align="center">
  <img src="./logo.png" alt="xyzus-backend" width="100">
  <br>xyzus-backend<br>
</h1>
<h4 align="center">宇宙电台 · 小宇宙播客 Web 后端（xyz 补丁版）</h4>

## 免责声明

**⚠️ 本项目仅供学习、研究使用，请遵守国家法律，严禁用于任何非法用途。**

## 这是什么

[ultrazg/xyz](https://github.com/ultrazg/xyz)（小宇宙 FM API 代理，Go 单进程）的**补丁版**，
为 [xyzus-web](https://github.com/simiooo/xyzus-web)（宇宙电台 Web 客户端）补齐上游没有的代理端点。
钉住上游 commit `22cfe78`（v1.10.3），在其上新增 **22 条代理路由** + 一处 `local-time` 修复：

| 域 | 端点 |
|---|---|
| 通知 / 收件箱 | `/notice_list`、`/inbox_list_v2`、`/inbox_remove`、`/system_notice_list`、`/system_notice_unread` |
| 首页内容流 | `/tab_featured_list`、`/tab_discovery_list`、`/tab_category_list` |
| 评论增强 | `/comment_list_thread_v2`、`/comment_list_mine`、`/comment_list_timestamp_primary`、`/comment_get` |
| AI / 文字稿 | `/ai_summary_episode_retrieve_brief`、`/episode_transcript_get`、`/episode_transcript_content` |
| 关注 / 订阅 | `/following_updates_list`、`/subscription_often_played` |
| 播放同步 / 续听 | `/playback_progress_list`、`/playback_progress_update`、`/episode_played_create`、`/episode_played_list_for_cache`（最近收听增量同步） |
| 相关单集 | `/related_episode_list` |

**local-time 修复**（`utils/http.go`）：对每个上游请求无条件动态生成 `Local-Time`（对齐 App 格式：毫秒 + 无冒号时区，
如 `2026-11-11T21:57:37.943+0800`），修复默认头 `...Z07:00` 格式不符导致播放进度写入被上游「200 但静默丢弃」。

完整端点契约（请求体 / 响应 / 实测）见 [xyzus-web README §4](https://github.com/simiooo/xyzus-web)。

## 部署：用户入口在 xyzus-web

**面向用户的部署入口是 [xyzus-web](https://github.com/simiooo/xyzus-web) 的 Docker Compose**——
它一次性构建前端 + 本补丁后端，无需单独起后端：

```bash
git clone https://github.com/simiooo/xyzus-web.git
cd xyzus-web
docker compose up -d --build
# 本机 http://localhost:8080/ ；局域网 http://<主机IP>:8080/（手机 / 其它电脑可访问）
```

> xyzus-web 的 `backend.Dockerfile` 会克隆上游 `22cfe78` 并 apply `backend-patches/notification-endpoints.patch`，
> 产物与本仓库源码等价。因此 **compose 部署不依赖本仓库**；本仓库是后端开发的主仓库。

## 独立运行（后端开发）

环境：Go 1.22+。

```bash
git clone https://github.com/simiooo/xyzus-backend.git
cd xyzus-backend
go mod tidy
go run .            # 默认监听 23020；-p 3000 改端口；-d 启动时打开文档
```

- 接口地址：`http://localhost:23020/login`
- 文档地址：`http://localhost:23020/docs`

或用 Docker 单独跑后端（本仓库 compose 从补丁源码构建，映射 `23020`）：

```bash
docker compose up -d --build
```

前端 dev 联调时把 Vite 代理指向它（在 xyzus-web 里）：`XYZ_API_TARGET=http://localhost:23020 pnpm dev`。

## 补丁维护

本仓库把补丁**直接提交进源码**（`handlers/`、`router/router.go`、`utils/http.go`）；
xyzus-web 另维护一份等价的 `backend-patches/notification-endpoints.patch`，用于在干净上游 pin 上可复现构建。
**两处需保持同步**：改后端后同步重新生成 xyzus-web 的补丁文件，例如：

```bash
git diff 22cfe78 -- handlers router utils > ../xyzus-web/backend-patches/notification-endpoints.patch
```

## License

[The MIT License](./LICENSE)（上游 ultrazg/xyz 同为 MIT）。