# pat-batch

制药工艺过程分析（PAT）批次对比工具。

读取批次测量 CSV 与参数规格 CSV，计算各参数的均值 / 标准差 / 过程能力指数 CPK，
并按批次检测 **超标（OOS）** 与 **趋势偏移（OOT）**，输出文本或 JSON 报告。

- 退出码 `0`：全部在控
- 退出码 `1`：检测到 OOS/OOT（或运行时错误）
- 退出码 `2`：参数缺失（用法错误）

## 用法

```
pat-batch -measurements <csv> -specs <csv> [-format text|json]
```

## CSV 格式

measurements.csv:
```
batch,parameter,value
B1,Temperature,121.3
```

specs.csv:
```
parameter,target,low,high
Temperature,120,118,122
```

## 实现说明

纯标准库实现，离线可构建。
