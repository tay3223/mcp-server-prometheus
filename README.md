# Prometheus MCP 服务器 🚀

大家好！这是一个超级简单易用的工具，它可以帮助你把 Prometheus（一个监控系统）的数据展示在 Cursor（一个超智能的代码编辑器）中。通过这个工具，你可以直接在 Cursor 中查看各种监控数据，就像跟 AI 助手聊天一样简单！

## 🌟 这个工具能做什么？

* 📊 **查看实时数据**：随时查看监控指标
* 📈 **查看历史趋势**：了解一段时间内的数据变化
* 🏷️ **查看标签信息**：浏览所有监控标签
* 📝 **查看详细信息**：获取监控指标的详细说明
* 🎯 **查看监控目标**：了解正在监控哪些服务
* ⚠️ **查看告警信息**：及时发现系统问题
* 📋 **查看监控规则**：了解监控的具体设置
* 🔐 **安全可靠**：支持API密钥保护

## 🚀 三步搞定部署

### 第一步：下载和安装

1. 确保你的电脑已经安装了 [Go语言](https://golang.org/dl/)（版本 1.21 或更高）

2. 打开终端，输入以下命令：
   ```bash
   # 下载代码
   git clone https://github.com/tay3223/mcp-server-prometheus.git
   
   # 进入项目目录
   cd mcp-server-prometheus
   
   # 编译程序（Windows用户使用 mcp-server-prometheus.exe）
   go build -o mcp-server-prometheus cmd/server/main.go
   ```

### 第二步：启动服务

1. 找到你的 Prometheus 服务地址（比如：http://localhost:9090）

2. 启动服务器（选择下面任意一种方式）：

   方式一：直接运行（会自动生成API密钥）
   ```bash
   ./mcp-server-prometheus --prometheus=http://localhost:9090
   ```

   方式二：使用自定义API密钥
   ```bash
   ./mcp-server-prometheus --prometheus=http://localhost:9090 --api-key=你的密钥
   ```

   方式三：使用环境变量
   ```bash
   export PROMETHEUS_HOST=http://localhost:9090
   export MCP_API_KEY=你的密钥
   ./mcp-server-prometheus
   ```

3. 看到类似下面的输出就说明启动成功了：
   ```
   生成的API Key: xxxxxxxx（记住这个密钥！）
   连接到Prometheus服务器: http://localhost:9090
   启动Prometheus MCP服务器...
   ```

### 第三步：配置 Cursor

1. 打开 Cursor 编辑器

2. 按下 `Command + Shift + P`（Windows用户按 `Ctrl + Shift + P`）

3. 输入 `Open Settings (JSON)`，回车

4. 在设置文件中添加下面的配置（记得替换其中的值）：
   ```json
   {
     "mcpServers": {
       "prometheus": {
         "url": "http://localhost:8080",  // 改成你的MCP服务器地址
         "defaultArguments": {
           "X-API-Key": "你的密钥"  // 改成你的API密钥
         }
       }
     }
   }
   ```

5. 保存设置文件

## 🎮 如何使用

1. 打开 Cursor，创建或打开任意文件

2. 确保 Cursor 的 AI 助手设置为 "Auto" 模式：
   - 点击右下角的 AI 图标
   - 选择 "Auto" 模式

3. 在编辑器中输入类似下面的问题：
   ```
   查询一下现在系统的运行状态
   ```

4. AI 助手会自动调用 Prometheus 查询，并给你友好的回答！

### 📝 常用查询示例

```python
# 查看某个服务是否在线
"查询 nginx 服务的在线状态"

# 查看 CPU 使用率
"最近一小时的 CPU 使用率是多少？"

# 查看内存使用情况
"当前系统内存使用情况如何？"

# 查看磁盘空间
"哪些磁盘的使用率超过了 80%？"

# 查看网络流量
"过去 30 分钟的网络流入流量是多少？"
```

## 🔧 常见问题

1. **问**：启动时报错 "connection refused"？
   **答**：检查你的 Prometheus 地址是否正确，并确保可以访问

2. **问**：Cursor 中无法连接到服务器？
   **答**：检查服务器地址和 API 密钥是否配置正确

3. **问**：查询时提示 "unauthorized"？
   **答**：检查 API 密钥是否正确

## 🤝 需要帮助？

- 遇到问题？欢迎[提交 Issue](https://github.com/tay3223/mcp-server-prometheus/issues)
- 有改进建议？欢迎[提交 PR](https://github.com/tay3223/mcp-server-prometheus/pulls)
- 想学习更多？查看[Prometheus 官方文档](https://prometheus.io/docs/introduction/overview/)

## 📜 开源协议

本项目采用 [MIT 协议](LICENSE) 开源，欢迎自由使用！ 