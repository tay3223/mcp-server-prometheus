package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/tay3223/mcp-server-prometheus/pkg/prometheus"
)

// PrometheusServer 封装了Prometheus MCP服务器
type PrometheusServer struct {
	server   *server.MCPServer
	client   *prometheus.Client
	address  string
	apiKey   string
	httpAddr string
}

// NewPrometheusServer 创建一个新的Prometheus MCP服务器
func NewPrometheusServer(address string, apiKey string, httpAddr string) (*PrometheusServer, error) {
	client, err := prometheus.NewClient(address)
	if err != nil {
		return nil, err
	}

	s := server.NewMCPServer(
		"Prometheus MCP",
		"1.0.0",
	)

	ps := &PrometheusServer{
		server:   s,
		client:   client,
		address:  address,
		apiKey:   apiKey,
		httpAddr: httpAddr,
	}

	// 注册工具
	ps.registerTools()

	return ps, nil
}

// registerTools 注册所有MCP工具
func (ps *PrometheusServer) registerTools() {
	// 即时查询工具
	instantQueryTool := mcp.NewTool("query",
		mcp.WithDescription("执行Prometheus即时查询"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("PromQL查询表达式"),
		),
		mcp.WithString("time",
			mcp.Description("查询时间(RFC3339格式)，默认为当前时间"),
		),
		mcp.WithString("X-API-Key",
			mcp.Required(),
			mcp.Description("API Key用于认证"),
		),
	)
	ps.server.AddTool(instantQueryTool, ps.withAuth(ps.handleInstantQuery))

	// 范围查询工具
	rangeQueryTool := mcp.NewTool("query_range",
		mcp.WithDescription("执行Prometheus范围查询"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("PromQL查询表达式"),
		),
		mcp.WithString("start",
			mcp.Required(),
			mcp.Description("开始时间(RFC3339格式)"),
		),
		mcp.WithString("end",
			mcp.Required(),
			mcp.Description("结束时间(RFC3339格式)"),
		),
		mcp.WithString("step",
			mcp.Required(),
			mcp.Description("步长(例如：15s, 1m, 1h)"),
		),
		mcp.WithString("X-API-Key",
			mcp.Required(),
			mcp.Description("API Key用于认证"),
		),
	)
	ps.server.AddTool(rangeQueryTool, ps.withAuth(ps.handleRangeQuery))

	// 标签值工具
	labelValuesTool := mcp.NewTool("label_values",
		mcp.WithDescription("获取Prometheus标签值"),
		mcp.WithString("label",
			mcp.Required(),
			mcp.Description("标签名称"),
		),
		mcp.WithString("start",
			mcp.Description("开始时间(RFC3339格式)"),
		),
		mcp.WithString("end",
			mcp.Description("结束时间(RFC3339格式)"),
		),
		mcp.WithString("X-API-Key",
			mcp.Required(),
			mcp.Description("API Key用于认证"),
		),
	)
	ps.server.AddTool(labelValuesTool, ps.withAuth(ps.handleLabelValues))

	// 指标元数据工具
	metadataTool := mcp.NewTool("metadata",
		mcp.WithDescription("获取Prometheus指标元数据"),
		mcp.WithString("metric",
			mcp.Description("指标名称，如果为空则返回所有指标"),
		),
		mcp.WithString("X-API-Key",
			mcp.Required(),
			mcp.Description("API Key用于认证"),
		),
	)
	ps.server.AddTool(metadataTool, ps.withAuth(ps.handleMetadata))

	// 目标工具
	targetsTool := mcp.NewTool("targets",
		mcp.WithDescription("获取Prometheus抓取目标信息"),
		mcp.WithString("X-API-Key",
			mcp.Required(),
			mcp.Description("API Key用于认证"),
		),
	)
	ps.server.AddTool(targetsTool, ps.withAuth(ps.handleTargets))

	// 告警工具
	alertsTool := mcp.NewTool("alerts",
		mcp.WithDescription("获取Prometheus告警信息"),
		mcp.WithString("X-API-Key",
			mcp.Required(),
			mcp.Description("API Key用于认证"),
		),
	)
	ps.server.AddTool(alertsTool, ps.withAuth(ps.handleAlerts))

	// 规则工具
	rulesTool := mcp.NewTool("rules",
		mcp.WithDescription("获取Prometheus规则信息"),
		mcp.WithString("X-API-Key",
			mcp.Required(),
			mcp.Description("API Key用于认证"),
		),
	)
	ps.server.AddTool(rulesTool, ps.withAuth(ps.handleRules))
}

// withAuth 验证API Key的中间件
func (ps *PrometheusServer) withAuth(handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		apiKey, ok := request.Params.Arguments["X-API-Key"].(string)
		if !ok || apiKey != ps.apiKey {
			return nil, fmt.Errorf("unauthorized: invalid API key")
		}
		return handler(ctx, request)
	}
}

// handleInstantQuery 处理即时查询请求
func (ps *PrometheusServer) handleInstantQuery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query must be a string")
	}

	var ts time.Time
	if timeStr, ok := request.Params.Arguments["time"].(string); ok && timeStr != "" {
		var err error
		ts, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid time format: %v", err)
		}
	} else {
		ts = time.Now()
	}

	result, err := ps.client.Query(ctx, query, ts)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// handleRangeQuery 处理范围查询请求
func (ps *PrometheusServer) handleRangeQuery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query must be a string")
	}

	startStr, ok := request.Params.Arguments["start"].(string)
	if !ok {
		return nil, fmt.Errorf("start must be a string")
	}
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %v", err)
	}

	endStr, ok := request.Params.Arguments["end"].(string)
	if !ok {
		return nil, fmt.Errorf("end must be a string")
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end time: %v", err)
	}

	stepStr, ok := request.Params.Arguments["step"].(string)
	if !ok {
		return nil, fmt.Errorf("step must be a string")
	}
	step, err := parseDuration(stepStr)
	if err != nil {
		return nil, fmt.Errorf("invalid step: %v", err)
	}

	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	result, err := ps.client.QueryRange(ctx, query, r)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// handleLabelValues 处理标签值请求
func (ps *PrometheusServer) handleLabelValues(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	label, ok := request.Params.Arguments["label"].(string)
	if !ok {
		return nil, fmt.Errorf("label must be a string")
	}

	var start, end time.Time
	if startStr, ok := request.Params.Arguments["start"].(string); ok && startStr != "" {
		var err error
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %v", err)
		}
	} else {
		start = time.Now().Add(-1 * time.Hour)
	}

	if endStr, ok := request.Params.Arguments["end"].(string); ok && endStr != "" {
		var err error
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %v", err)
		}
	} else {
		end = time.Now()
	}

	result, err := ps.client.LabelValues(ctx, label, nil, start, end)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// handleMetadata 处理元数据请求
func (ps *PrometheusServer) handleMetadata(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	metric := ""
	if metricVal, ok := request.Params.Arguments["metric"].(string); ok {
		metric = metricVal
	}

	result, err := ps.client.Metadata(ctx, metric, "")
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// handleTargets 处理目标请求
func (ps *PrometheusServer) handleTargets(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := ps.client.Targets(ctx)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// handleAlerts 处理告警请求
func (ps *PrometheusServer) handleAlerts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := ps.client.Alerts(ctx)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// handleRules 处理规则请求
func (ps *PrometheusServer) handleRules(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := ps.client.Rules(ctx)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// StartServer 启动MCP服务器
func (ps *PrometheusServer) StartServer() error {
	// 如果指定了HTTP地址，则启动HTTP服务器
	if ps.httpAddr != "" {
		fmt.Printf("启动HTTP服务器在: %s\n", ps.httpAddr)
		mux := http.NewServeMux()

		// 添加健康检查端点
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})

		// 添加状态端点
		mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"server":"Prometheus MCP Server","version":"1.0.0","status":"running"}`))
		})

		// 启动HTTP服务器
		return http.ListenAndServe(ps.httpAddr, mux)
	}

	// 否则使用标准IO模式
	return server.ServeStdio(ps.server)
}

// 辅助函数，解析持续时间字符串
func parseDuration(s string) (time.Duration, error) {
	// 尝试直接解析
	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}

	// 尝试解析数字（秒）
	if seconds, err := strconv.ParseFloat(s, 64); err == nil {
		return time.Duration(seconds * float64(time.Second)), nil
	}

	return 0, fmt.Errorf("无法解析持续时间: %s", s)
}
