# subsonic_playlists example

读取 Subsonic 服务中的歌单列表，并继续读取第一个歌单的歌曲条目。

## 环境变量

- `SUBSONIC_BASE`：Subsonic 服务地址，默认 `http://127.0.0.1:8089`
- `SUBSONIC_USER`：Subsonic 用户名，默认 `admin`
- `SUBSONIC_PASSWORD`：Subsonic 密码，默认 `admin`
- `SUBSONIC_CLIENT`：Subsonic 客户端名，默认 `subsonic-playlists-example`

## 运行

```bash
cd kugou-api-go
SUBSONIC_BASE='http://127.0.0.1:18089' \
SUBSONIC_USER='joe' \
SUBSONIC_PASSWORD='sesame' \
go run ./examples/subsonic_playlists
```

输出会先显示歌单数量和前几个歌单，再显示第一个歌单的前 10 首歌曲。
