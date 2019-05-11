Drone.io Badge Server
===

Droneの各ステージにおける結果に基づいてバッジを返すサーバです．

## API
### POST /generate
```json
{
  "build_number": 12, // ビルド番号
  "repo_namespace": "atpons", // レポジトリの名前空間
  "repo_name": "drone-badge-server" // レポジトリ名
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