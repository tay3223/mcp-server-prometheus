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
* **安全认证**：支持API Key认证机制

## 安装

从源代码构建：

```bash
git clone https://github.com/tay3223/mcp-server-prometheus.git
cd mcp-server-prometheus
go build -o mcp-server-prometheus cmd/server/main.go
```

## 使用方法

### 命令行参数

```bash
./mcp-server-prometheus --prometheus=http://your-prometheus-instance:9090 [--api-key=your-api-key]
```

### 环境变量

也可以通过环境变量配置：

```bash
export PROMETHEUS_HOST=http://your-prometheus-instance:9090
export MCP_API_KEY=your-api-key
./mcp-server-prometheus
```

## 在Cursor中使用

在Cursor中配置MCP服务器：

```json
{
  "mcpServers": {
    "prometheus": {
      "command": "/path/to/mcp-server-prometheus",
      "env": {
        "PROMETHEUS_HOST": "http://your-prometheus-instance:9090"
      },
      "defaultArguments": {
        "X-API-Key": "your-api-key"
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
query(
  query: "up",
  X-API-Key: "your-api-key"
)

# 执行范围查询
query_range(
  query: "rate(node_cpu_seconds_total{mode=\"system\"}[5m])",
  start: "2023-04-01T00:00:00Z",
  end: "2023-04-01T01:00:00Z",
  step: "5m",
  X-API-Key: "your-api-key"
)

# 获取标签值
label_values(
  label: "job",
  X-API-Key: "your-api-key"
)

# 获取指标元数据
metadata(
  metric: "up",
  X-API-Key: "your-api-key"
)
```

## 安全性

服务器使用API Key进行认证。每个请求都需要提供有效的API Key。API Key可以通过以下方式设置：

1. 命令行参数：`--api-key`
2. 环境变量：`MCP_API_KEY`
3. 如果未提供，服务器会自动生成一个随机的API Key

建议：
* 使用足够长的随机API Key（建议32字节以上）
* 通过安全的方式传输API Key
* 定期轮换API Key
* 在生产环境中使用HTTPS

## 许可证

MIT

## 贡献

欢迎提交问题和PR！ 