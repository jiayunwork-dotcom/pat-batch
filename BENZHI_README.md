制药 PAT 批次对比的 Go 命令行：两份 CSV（测量 vs 规格）算出各参数 CPK，并按批次标 OOS/OOT，退出码区分在控和超标；pat-batch serve 用 /api/analyze 吃 JSON 做同一套统计。

# pat-batch 是制药工艺过程分析（PAT）批次对比与过程能力检测工具

读取批次测量 CSV 与参数规格 CSV，计算各参数 CPK 并检测超标（OOS）与趋势偏移（OOT）。

## 构建 / 运行 / 测试

```text
go build ./...
go run . -measurements example/measurements.csv -specs example/specs.csv
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
