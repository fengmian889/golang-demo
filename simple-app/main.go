package main

import (
	"fmt"
	"simple-app/pkg/app"
	"simple-app/pkg/log"
	"simple-app/pkg/options"
)

func NewApp(basename string) *app.App {
	opts := options.NewSimpleOptions()

	application := app.NewApp("Simple IAM Server", // 应用名称
		basename,                   // 二进制文件名
		app.WithOptions(opts),      // 配置选项（必需）
		app.WithRunFunc(run(opts)), // 运行函数（必需）
	)

	return application
}

func run(opts *options.SimpleOptions) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()
		fmt.Printf("🚀 启动简化版 IAM 服务器...\n")
		fmt.Printf("📋 配置信息:\n")
		fmt.Printf("   服务器模式: %s\n", opts.ServerRunOptions.Mode)
		fmt.Printf("   健康检查: %v\n", opts.ServerRunOptions.Healthz)
		fmt.Printf("   HTTP 端口: %d\n", opts.InsecureServing.BindPort)
		//fmt.Printf("   日志级别: %s\n", opts.Log.Level)
		//fmt.Printf("   日志格式: %s\n", opts.Log.Format)

		fmt.Printf("✅ 服务器已启动！\n")
		fmt.Printf("💡 在真实的 IAM 项目中，这里会:\n")
		fmt.Printf("   - 连接数据库 (MySQL)\n")
		fmt.Printf("   - 连接缓存 (Redis)\n")
		fmt.Printf("   - 启动 HTTP 服务器 (端口 %d)\n", opts.InsecureServing.BindPort)
		fmt.Printf("   - 启动 HTTPS 服务器\n")
		fmt.Printf("   - 启动 gRPC 服务器\n")
		fmt.Printf("\n🔧 尝试修改参数重新运行:\n")
		fmt.Printf("   %s --server.mode=debug\n", basename)
		fmt.Printf("   %s --insecure.bind-port=9090 --log.level=debug\n", basename)

		return nil
	}
}

func main() {
	// 步骤6: 创建并运行应用程序
	app := NewApp("simple-app")
	app.Run()
}
