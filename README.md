# Prometheus MCP服务器

这是一个基于Go实现的Prometheus MCP（Model Context Protocol）服务器，可以将Prometheus数据暴露给支持MCP协议的大型语言模型应用（如Cursor）。

## 功能特性

* **即时查询**：在特定时间执行PromQL查询
* **范围查询**：在一段时间内执行PromQL查询
* **标签值查询**：获取特定标签的值
* **元数据查询**：获取指标的元数据信息
* **目标查询**：获取Prometheus抓取目标信息
* **告警查询**：获取当前活动的告警
* **规则查询**：获取记录规则和告警规则

## 安装

从源代码构建：

```bash
git clone https://github.com/yourusername/promcp.git
cd promcp
go build -o promcp-server cmd/server/main.go
```

## 使用方法

### 命令行参数

```bash
./promcp-server -prometheus=http://your-prometheus-instance:9090
```

### 环境变量

也可以通过环境变量配置：

```bash
PROMETHEUS_HOST=http://your-prometheus-instance:9090 ./promcp-server
```

## 在Cursor中使用

在Cursor中配置MCP服务器：

```json
{
  "mcpServers": {
    "prometheus": {
      "command": "/path/to/promcp-server",
      "env": {
        "PROMETHEUS_HOST": "http://your-prometheus-instance:9090"
      }
    }
  }
}
```

## 可用工具

本服务器提供以下MCP工具：

* `query`：执行即时PromQL查询
* `query_range`：执行时间范围内的PromQL查询
* `label_values`：获取特定标签的值
* `metadata`：获取指标元数据
* `targets`：获取Prometheus抓取目标信息
* `alerts`：获取当前告警
* `rules`：获取记录和告警规则

## 示例

以下是在Cursor中使用MCP工具的示例：

```
# 执行即时查询
query(query: "up")

# 执行范围查询
query_range(
  query: "rate(node_cpu_seconds_total{mode=\"system\"}[5m])",
  start: "2023-04-01T00:00:00Z",
  end: "2023-04-01T01:00:00Z",
  step: "5m"
)

# 获取标签值
label_values(label: "job")

# 获取指标元数据
metadata(metric: "up")
```

## 许可证

MIT

## 贡献

欢迎提交问题和PR！ 