# go-probeorch

健康探测编排库 + `probed` 管理服务：注册 HTTP/TCP/自定义探针，按间隔调度，环形结果与成功率聚合。

```bash
go build ./...
go test ./... -count=1
go run ./cmd/probed -addr :8100 -web web
```
