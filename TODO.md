# TODO

## P0 - 签名/加密接口全量修复
- [x] 将 `sdk/API_SIGN_CHECK.md` 中 44 个 `Need=YES` 接口全部改为“原函数手写修复”（不新增平替函数）。
- [x] 清零 `待校对` 的 4 个接口：
  - [x] `/audio/related`
  - [x] `/brush`
  - [x] `/register/dev`
  - [x] `/user/video/collect`
- [x] 清零 `部分修复` 的 30 个接口（逐个补齐动态参数、签名、加解密、返回后处理）。
- [x] 每修复 1 个接口，同步更新：
  - [x] `sdk/API_FIX_STATUS.md` -> `已校对修复`
  - [x] `sdk/API_SIGN_CHECK.md` -> 对应状态更新

## P0 - 接口待实测收尾（当前剩余 49 个）
- [ ] 将 `sdk/API_FIX_STATUS.md` 中 `待实测` 接口从 49 降到 0（按模块分批联调验收：`rank/*`、`scene/*`、`youth/*`、`longaudio/*` 等）。
- [ ] 备注：你尚未实测，当前这批先按“代码已接入，待你验收”推进，实测后再改状态。
- [x] 代码侧已逐个补齐默认参数/登录校验/后处理（先不改 `API_FIX_STATUS`，待你统一验收后再切状态）。

## P0 - 登录链路稳定性
- [ ] 验证并修复手机验证码登录全流程（发送验证码 -> 登录 -> 用户信息）在概念版默认配置下稳定通过。
- [ ] 验证并修复 `UserDetail` 在以下场景的稳定性：
  - [ ] 新登录后立即查询
  - [ ] 旧会话恢复后查询
  - [ ] userid 为科学计数法历史脏数据时查询
- [ ] 保证登录失败或用户信息失败时不落盘无效会话。

## P1 - 示例与回归
- [ ] 回归并确认以下示例都可运行：
  - [ ] `examples/login_cellphone`
  - [ ] `examples/login_session`
  - [ ] `examples/token_refresh`
  - [ ] `examples/user_info_refresh`
  - [ ] `examples/daily_recommend`
- [ ] 为签名接口补充最小可复用示例（优先高频接口：`search`、`playlist_detail`、`user_detail`）。

## P1 - 文档收尾
- [ ] README 增加“已手写签名修复接口列表”章节（从 `API_FIX_STATUS.md` 自动生成或手工维护）。
- [ ] README 明确默认平台策略：
  - [ ] 默认 `lite`
  - [ ] 如何切普通版（`platform=normal` 或 `sdk.WithLite(false)`）
- [ ] README 增加排障条目：
  - [ ] `20006`
  - [ ] `20018`
  - [ ] userid 科学计数法问题

## P2 - 工具化
- [ ] 强化 `tools/audit`：区分“签名风险”“参数组装风险”“后处理风险”并单独计数。
- [ ] 新增“状态一致性检查”工具：确保 `API_FIX_STATUS.md` 与 `API_SIGN_CHECK.md` 对同一接口状态一致。
- [ ] 在 CI（或本地脚本）加入：
  - [ ] `go test ./...`
  - [ ] 生成文件一致性检查（`tools/gen` 后无差异）
