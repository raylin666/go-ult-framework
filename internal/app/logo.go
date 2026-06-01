// Package app 提供应用工具包。
package app

import (
	"fmt"
	"runtime"
	"ult/config"

	"github.com/fatih/color"
)

// NewLogo 打印项目服务启动信息。
// 显示项目 Logo、Go 版本、系统信息、服务配置等信息。
//
// 参数:
//   - config: 应用配置
func NewLogo(config *config.Config) {
	// see https://patorjk.com/software/taag/#p=testall&f=Graffiti&t=ult
	var logo = `
	██╗   ██╗	██╗  	     ████████╗
	██║   ██║	██║  	     ╚══██╔══╝
	██║   ██║	██║     	██║   
	██║   ██║	██║     	██║   
	╚██████╔╝	███████╗	██║   
	 ╚═════╝ 	╚══════╝	╚═╝
`
	color.HiYellow(logo)

	contents := fmt.Sprintf(`
GO 版本及路径: %s (%s)
系统类型及架构: %s (%s) - %d 核 CPU
服务名称: %s (%s)
服务版本: %s
服务环境: %s
	`,
		runtime.Version(),
		runtime.GOROOT(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.GOMAXPROCS(0),

		config.App.Name,
		config.App.ID,
		config.App.Version,
		config.Environment,
	)

	color.HiGreen(contents)
}
