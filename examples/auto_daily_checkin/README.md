# auto_daily_checkin

基于当前 Go SDK 的自动签到示例。

这个示例的目标不是再维护一套独立登录状态，而是直接复用本项目 SDK 的 `session` 文件。自动签到配置只记录“账号与 session 文件的关联关系、最近签到状态、守护进程信息”，不会额外再抄一份 token。

## 功能概览

- `add-account`
  - 交互式登录并添加自动签到账号
  - 支持二维码登录和手机验证码登录
- `del-account`
  - 从自动签到配置中移除账号
  - 只移除配置，不删除原 session 文件
- `start`
  - 启动自动签到主循环
  - 不带子命令时默认等价于 `start`
- `daemon`
  - 后台启动 `start`
- `stop`
  - 停止当前后台实例
- `status`
  - 查看 daemon 状态
  - 查看账号昵称、最后一次签到时间、会员过期时间、今日是否已签到

## 运行逻辑

1. 启动时读取自动签到配置文件。
2. 如果当前没有自动签到账号，会尝试把项目默认 SDK session 自动迁移进来。
3. 每轮同步会先读取账号对应的 session 文件，并刷新用户信息。
4. 通过 `youth_month_vip_record` 判断今天是否已经签到。
5. 如果今天未签到且当前时间已过北京时间 `01:00`，执行签到流程：
   - 先调用 `YouthListenSong`
   - 再调用 `YouthVip`
6. 如果今天已签到，则只更新状态，不重复执行。

## 项目架构

目录内文件按职责拆分，避免单文件过大：

- [`main.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/main.go)
  - 命令入口
  - 解析 `start / daemon / stop / status / add-account / del-account`
- [`checkin.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/checkin.go)
  - 自动签到主循环
  - 前台运行逻辑
  - 具体签到执行流程
- [`add_account.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/add_account.go)
  - 交互式登录
  - 账号添加与删除
  - 调用现有登录 example 写入 session
- [`status_command.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/status_command.go)
  - `status` 输出整理
- [`status_helpers.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/status_helpers.go)
  - 读取账号状态
  - 规范化 `userid`
  - 提取最后签到日期和会员到期时间
- [`config.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/config.go)
  - 自动签到配置模型
  - 默认 session 迁移
- [`daemon.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/daemon.go)
  - pid 文件
  - 后台启动
  - 停止与状态判断
- [`paths.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/paths.go)
  - 默认路径
  - 运行目录
- [`logger.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/logger.go)
  - 控制台和日志双写
- [`status_helpers_test.go`](/root/KuGouMusicApi/kugou-api-go/examples/auto_daily_checkin/status_helpers_test.go)
  - 关键状态提取逻辑单测

## 默认文件位置

默认运行目录：

```text
~/.config/kugou-music-api/auto_daily_checkin/
```

默认会生成这些文件：

- `config.json`
  - 自动签到账号配置
- `auto_daily_checkin.log`
  - daemon 和 start 日志
- `auto_daily_checkin.pid`
  - 后台进程 pid 文件
- `sessions/`
  - `add-account` 新增账号时写入的独立 session 文件

默认共享 SDK session：

```text
~/.kugou_music_api_session.json
```

如果自动签到配置为空，启动时会尝试把这个默认 session 自动迁移进自动签到配置。

## 使用方法

### 查看状态

```bash
go run ./examples/auto_daily_checkin status
```

### 添加账号

交互式选择登录方式：

```bash
go run ./examples/auto_daily_checkin add-account
```

指定二维码登录：

```bash
go run ./examples/auto_daily_checkin add-account -method qr
```

指定手机验证码登录：

```bash
go run ./examples/auto_daily_checkin add-account -method cellphone
```

指定新账号 session 文件位置：

```bash
go run ./examples/auto_daily_checkin add-account \
  -method qr \
  -session-file /path/to/account-session.json
```

### 删除账号

交互式删除：

```bash
go run ./examples/auto_daily_checkin del-account
```

按 `userid` 删除：

```bash
go run ./examples/auto_daily_checkin del-account -user 51815528
```

### 前台启动

```bash
go run ./examples/auto_daily_checkin
```

或：

```bash
go run ./examples/auto_daily_checkin start
```

调整轮询间隔：

```bash
go run ./examples/auto_daily_checkin start -interval 2m
```

### 后台启动

```bash
go run ./examples/auto_daily_checkin daemon
```

停止后台实例：

```bash
go run ./examples/auto_daily_checkin stop
```

## 编译方法

编译到当前目录：

```bash
go build -o auto_daily_checkin ./examples/auto_daily_checkin
```

编译后运行：

```bash
./auto_daily_checkin status
./auto_daily_checkin daemon
./auto_daily_checkin stop
```

## 可选参数

- `-config`
  - 自动签到配置文件路径
- `-source-session-file`
  - 默认迁移来源的 SDK session 文件路径
- `-log`
  - 日志文件路径
- `-pid`
  - pid 文件路径
- `-interval`
  - 主循环检查间隔，默认 `1m`

## 环境变量

- `KUGOU_AUTO_CHECKIN_DIR`
  - 覆盖默认运行目录
- `KUGOU_SESSION_FILE`
  - 影响共享 SDK 默认 session 路径
- `KUGOU_AUTO_CHECKIN_VIP_INTERVAL`
  - `YouthVip` 多次领取之间的等待间隔
  - 默认 `30s`

## 说明

- `status` 中“最后一次签到时间”优先显示本工具自己记录的完整时间。
- 如果账号是在本工具之外已经签到过，接口通常只能返回日期，因此会显示：
  - `2026-04-26 (接口仅返回日期)`
- 当前逻辑会在北京时间 `01:00` 之后才尝试执行当天签到。
- 当前示例已经对齐本项目的 `lite` 平台默认配置。
