// Package email 提供邮件告警通知功能。
// 当系统发生错误时，通过邮件发送告警通知。
package email

import (
	"context"
	"strings"
	"ult/config/autoload"
	"ult/pkg/logger"
	"ult/pkg/proposal"

	"github.com/raylin666/go-utils/v2/mail"
	"go.uber.org/zap"
)

// NotifyHandler 创建邮件告警通知处理函数。
// 根据配置发送 HTML 格式的告警邮件。
//
// 参数:
//   - ctx: 上下文
//   - config: 告警通知配置
//   - logger: 日志记录器
//
// 返回:
//   - func(msg *proposal.AlertMessage): 告警通知处理函数
func NotifyHandler(ctx context.Context, config autoload.Notify, logger *logger.Logger) func(msg *proposal.AlertMessage) {
	return func(msg *proposal.AlertMessage) {
		go func() {
			if config.Recover.Email.Host == "" || config.Recover.Email.Port == 0 || config.Recover.Email.User == "" || config.Recover.Email.Pass == "" || config.Recover.Email.To == "" {
				logger.UseApp(ctx).Error("发送告警邮件通知邮件配置错误")
				return
			}

			subject, body, err := newHTMLEmail(
				msg.ProjectName,
				msg.Method,
				msg.HOST,
				msg.URI,
				msg.TraceID,
				msg.ErrorMessage,
				msg.Timestamp,
				msg.ErrorStack)
			if err != nil {
				logger.UseApp(ctx).Error("发送告警邮件通知邮件模板错误", zap.Error(err))
				return
			}

			m, err := mail.New(
				mail.WithMailHost(config.Recover.Email.Host),
				mail.WithMailPort(config.Recover.Email.Port),
				mail.WithMailUser(config.Recover.Email.User),
				mail.WithMailPass(config.Recover.Email.Pass))

			if err := m.SendTextHtml(subject, body, strings.Split(config.Recover.Email.To, ",")); err != nil {
				logger.UseApp(ctx).Error("发送告警通知邮件失败", zap.Error(err))
			}
		}()
	}
}
