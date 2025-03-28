package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/promcp/pkg/mcp"
)

// generateAPIKey 生成一个随机的API Key
func generateAPIKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func main() {
	// 解析命令行参数
	prometheusAddr := flag.String("prometheus", "http://localhost:9090", "Prometheus服务器地址")
	apiKey := flag.String("api-key", "", "API Key用于认证(如果不提供将自动生成)")
	flag.Parse()

	// 从环境变量中获取Prometheus地址(如果有)
	if envAddr := os.Getenv("PROMETHEUS_HOST"); envAddr != "" {
		*prometheusAddr = envAddr
	}

	// 从环境变量中获取API Key(如果有)
	if envKey := os.Getenv("MCP_API_KEY"); envKey != "" {
		*apiKey = envKey
	}

	// 如果没有提供API Key，则生成一个
	if *apiKey == "" {
		*apiKey = generateAPIKey()
		fmt.Printf("生成的API Key: %s\n", *apiKey)
	}

	fmt.Printf("连接到Prometheus服务器: %s\n", *prometheusAddr)

	// 创建MCP服务器
	server, err := mcp.NewPrometheusServer(*prometheusAddr, *apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建服务器失败: %v\n", err)
		os.Exit(1)
	}

	// 启动服务器
	fmt.Println("启动Prometheus MCP服务器...")
	if err := server.StartServer(); err != nil {
		fmt.Fprintf(os.Stderr, "服务器启动失败: %v\n", err)
		os.Exit(1)
	}
}
