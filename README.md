Drone.io Badge Server
===

Droneの各ステージにおける結果に基づいてバッジを返すサーバです．

## 注意
- Droneが複数台のときのことを考慮していません．
- 取得した結果をインメモリDBに保存しています．

## API
### POST /generate
```json
{
  "build_number": 12,
  "repo_namespace": "atpons",
  "repo_name": "drone-badge-server"
}
```

### GET /<repo_id>/<stage_number>
レポジトリIDとステージ番号に対応したMapを書いてください．
```go
var BadgeRoute = badge.Repo{
	100: { // レポジトリID
		1: {"build_pass.png", "build_fail.png"}, // ステージ番号
		2: {"webhook_pass.png", "webhook_fail.png"}, // ステージ番号
	},
}
```
