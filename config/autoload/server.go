// Package autoload 提供配置模块自动加载。
package autoload

// Server 服务器配置结构体。
type Server struct {
	Http Http // HTTP 服务器配置
}

// Http HTTP 服务器配置结构体。
// 定义 HTTP 服务器的网络、地址、端口和 CORS 配置。
type Http struct {
	Network string `yaml:"network"` // 网络协议 (tcp)
	Host    string `yaml:"host"`    // 服务主机地址
	Port    int    `yaml:"port"`    // 服务端口
	Cors    struct {
		Domains string `yaml:"domains"` // CORS 允许的域名 (all 或逗号分隔的域名列表)
	} `yaml:"cors"`
}
