package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Client 封装了Prometheus API客户端
type Client struct {
	api v1.API
}

// NewClient 创建一个新的Prometheus客户端
func NewClient(address string) (*Client, error) {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		api: v1.NewAPI(client),
	}, nil
}

// Query 执行即时查询
func (c *Client) Query(ctx context.Context, query string, ts time.Time) (model.Value, error) {
	val, _, err := c.api.Query(ctx, query, ts)
	return val, err
}

// QueryRange 执行范围查询
func (c *Client) QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, error) {
	val, _, err := c.api.QueryRange(ctx, query, r)
	return val, err
}

// Series 查找与标签匹配器匹配的序列
func (c *Client) Series(ctx context.Context, matches []string, startTime, endTime time.Time) ([]model.LabelSet, error) {
	val, _, err := c.api.Series(ctx, matches, startTime, endTime)
	return val, err
}

// LabelNames 获取标签名称
func (c *Client) LabelNames(ctx context.Context, matches []string, startTime, endTime time.Time) ([]string, error) {
	val, _, err := c.api.LabelNames(ctx, matches, startTime, endTime)
	return val, err
}

// LabelValues 获取标签值
func (c *Client) LabelValues(ctx context.Context, label string, matches []string, startTime, endTime time.Time) (model.LabelValues, error) {
	val, _, err := c.api.LabelValues(ctx, label, matches, startTime, endTime)
	return val, err
}

// Targets 获取目标信息
func (c *Client) Targets(ctx context.Context) (v1.TargetsResult, error) {
	return c.api.Targets(ctx)
}

// Alerts 获取告警信息
func (c *Client) Alerts(ctx context.Context) (v1.AlertsResult, error) {
	return c.api.Alerts(ctx)
}

// Rules 获取规则信息
func (c *Client) Rules(ctx context.Context) (v1.RulesResult, error) {
	return c.api.Rules(ctx)
}

// Metadata 获取指标元数据
func (c *Client) Metadata(ctx context.Context, metric string, limit string) (map[string][]v1.Metadata, error) {
	return c.api.Metadata(ctx, metric, limit)
} 