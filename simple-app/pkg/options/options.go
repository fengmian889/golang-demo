package options

import (
	"fmt"
	"simple-app/pkg/flag"
	"simple-app/pkg/log"

	"github.com/spf13/pflag"
)

type CliOptions interface {
	Flags() (fss flag.NamedFlagSets)
	Validate() []error
}

// SimpleOptions 是简化版的配置选项，只包含最基本的配置
// 用于演示命令行参数传递流程
type SimpleOptions struct {
	ServerRunOptions *ServerRunOptions       `json:"server"   mapstructure:"server"`
	InsecureServing  *InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	Log              *log.Options            `json:"log"      mapstructure:"log"`
}

type ServerRunOptions struct {
	Mode    string `json:"mode"    mapstructure:"mode"`
	Healthz bool   `json:"healthz" mapstructure:"healthz"`
}

type InsecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

func NewSimpleOptions() *SimpleOptions {
	return &SimpleOptions{
		ServerRunOptions: &ServerRunOptions{
			Mode:    "release", // 默认发布模式
			Healthz: true,      // 默认启用健康检查
		},
		InsecureServing: &InsecureServingOptions{
			BindAddress: "127.0.0.1", // 默认本地地址
			BindPort:    8080,        // 默认 8080 端口
		},
		//Log: log.NewOptions(), // 使用日志包的默认配置
	}
}

func (o *SimpleOptions) Flags() (fss flag.NamedFlagSets) {
	o.ServerRunOptions.AddFlags(fss.FlagSet("server"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	//o.Log.AddFlags(fss.FlagSet("logs"))
	return fss
}

func (o *SimpleOptions) Validate() []error {
	var errs []error
	errs = append(errs, o.ServerRunOptions.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	return errs
}

func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Mode, "server.mode", s.Mode,
		"服务器运行模式 (debug, release, test)")
	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz,
		"是否启用健康检查端点")
}

// Validate 验证服务器运行配置
func (s *ServerRunOptions) Validate() []error {
	var errors []error

	// 验证模式是否有效
	validModes := []string{"debug", "release", "test"}
	isValidMode := false
	for _, mode := range validModes {
		if s.Mode == mode {
			isValidMode = true
			break
		}
	}
	if !isValidMode {
		errors = append(errors,
			fmt.Errorf("invalid server mode: %s, must be one of: %v", s.Mode, validModes))
	}

	return errors
}

// AddFlags 为 HTTP 服务配置添加命令行标志
func (s *InsecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "insecure.bind-address", s.BindAddress,
		"HTTP 服务器绑定的 IP 地址")
	fs.IntVar(&s.BindPort, "insecure.bind-port", s.BindPort,
		"HTTP 服务器监听的端口号")
}

// Validate 验证 HTTP 服务配置
func (s *InsecureServingOptions) Validate() []error {
	var errors []error

	// 验证端口范围
	if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(errors,
			fmt.Errorf("bind port %d must be between 0 and 65535", s.BindPort))
	}

	return errors
}
