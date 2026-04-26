# Subsonic Server (KuGou)

基于 `kugou-api-go` 的 Subsonic 兼容服务示例，支持常用接口与音频反向代理。

## 已实现接口

- `ping`
- `getLicense`
- `getMusicFolders`
- `getIndexes`
- `getAlbumList2`
- `getAlbum`
- `getSong`
- `search2`
- `getPlaylists`
- `getPlaylist`
- `getCoverArt`
- `stream`（反向代理到酷狗音频地址）

## 启动

```bash
cd kugou-api-go
SUBSONIC_ADDR=:8089 \
SUBSONIC_USER=admin \
SUBSONIC_PASSWORD=admin \
SUBSONIC_DATA_FILE=./examples/subsonic_server/data/state.json \
go run ./examples/subsonic_server
```

## 认证

采用 Subsonic token 认证：

- `u`: 用户名
- `s`: 盐
- `t`: `md5(password + salt)`

示例：

```bash
SALT=c19b2d
TOKEN=$(printf 'admin%s' "$SALT" | md5sum | awk '{print $1}')

curl "http://127.0.0.1:8089/rest/ping.view?u=admin&t=${TOKEN}&s=${SALT}&v=1.16.1&c=test&f=json"
```

## 说明

- `stream` 为服务端反向代理，客户端不需要直接访问酷狗 URL。
- `search2` 优先走 `Search`，若返回空则自动回退 `TopSong` 本地过滤。
- 服务会把歌曲缓存、封面缓存、歌单缓存持久化到 `SUBSONIC_DATA_FILE`。
- 服务启动时会自动读取会话文件（`KUGOU_SESSION_FILE` 或默认 `~/.kugou_music_api_session.json`），用于返回“我的歌单”。
