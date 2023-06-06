package email

import (
	"context"
	"github.com/raylin666/go-utils/mail"
	"go.uber.org/zap"
	"strings"
	"ult/config/autoload"
	"ult/pkg/logger"
	"ult/pkg/proposal"
)

// NotifyHandler 告警通知
func NotifyHandler(ctx context.Context, config autoload.Notify, logger *logger.Logger) func(msg *proposal.AlertMessage) {
	return func(msg *proposal.AlertMessage) {
		go func() {
			if config.Recover.Email.Host == "" || config.Recover.Email.Port == 0 || config.Recover.Email.User == "" || config.Recover.Email.Pass == "" || config.Recover.Email.To == "" {
				logger.UseApp(ctx).Error("email notify config error")
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
				logger.UseApp(ctx).Error("email notify template error", zap.Error(err))
				return
			}

			var m = mail.New(
				mail.WithMailHost(config.Recover.Email.Host),
				mail.WithMailPort(config.Recover.Email.Port),
				mail.WithMailUser(config.Recover.Email.User),
				mail.WithMailPass(config.Recover.Email.Pass))

			if err := m.SendTextHtml(subject, body, strings.Split(config.Recover.Email.To, ",")); err != nil {
				logger.UseApp(ctx).Error("发送告警通知邮件失败", zap.Error(err))
			}

			return
		}()
	}
}
