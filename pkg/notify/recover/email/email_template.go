package email

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// NewHTMLEmail 告警邮件模板
func newHTMLEmail(service, method, host, uri, id string, msg interface{}, t time.Time, stack string) (subject string, body string, err error) {
	mailData := &struct {
		Service   string
		URL       string
		ID        string
		Msg       string
		Stack     string
		Year      int
		Timestamp time.Time
	}{
		Service:   service,
		URL:       fmt.Sprintf("%s %s%s", method, host, uri),
		ID:        id,
		Msg:       fmt.Sprintf("%+v", msg),
		Stack:     stack,
		Year:      time.Now().Year(),
		Timestamp: t,
	}

	// subject 邮件主题
	subject = fmt.Sprintf("%s [系统告警]-%s", service, uri)

	// body 邮件内容
	body, err = getEmailHTMLContent(mailTemplate, mailData)

	return
}

// getEmailHTMLContent 获取邮件模板
func getEmailHTMLContent(mailTpl string, mailData interface{}) (string, error) {
	tpl, err := template.New("email notify tpl").Parse(mailTpl)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, mailData)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

const mailTemplate = `
<!DOCTYPE html>
<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>

    <style type="text/css" rel="stylesheet" media="all">
        /* Media Queries */
        @media only screen and (max-width: 500px) {
            .button {
                width: 100% !important;
            }
        }
    </style>
</head>


<body style="margin: 0; padding: 0; width: 100%;">
<table width="100%" cellpadding="0" cellspacing="0">
    <tr>
        <td style="width: 100%; margin: 0; padding: 0;" align="center">
            <table width="100%" cellpadding="0" cellspacing="0">
                <!-- Logo -->
                <!-- <tr>
                    <td style="padding: 25px 0; text-align: center;">
                        {{.Service}} 系统告警
                    </td>
                </tr> -->

                <!-- Email Body -->
                <tr>
                    <td style="width: 100%; margin: 0; padding: 0; border-top: 1px solid #EDEFF2; border-bottom: 1px solid #EDEFF2; background-color: #FFF;"
                        width="100%">
                        <table style="width: auto; max-width: 750px; margin: 0 auto; padding: 0;" align="center"
                               width="750" cellpadding="0" cellspacing="0">
                            <tr>
                                <td style="font-family: Arial, 'Helvetica Neue', Helvetica, sans-serif; padding: 35px;">
                                    <!-- Greeting -->
                                    <h1 style="margin-top: 0; color: #2F3133; font-size: 22px; font-weight: bold; text-align: left;">
                                        Hello!
                                    </h1>

                                    <!-- Intro -->
                                    <p style="margin-top: 0; color: #74787E; line-height: 2em; font-size: 14px;">
                                        <b><i> 您收到此电子邮件，请紧急安排处理。</i></b>
                                    </p>

                                    <!-- Action Button -->
                                    <table style="width: 100%; margin: 30px auto; padding: 0;"
                                           width="100%" cellpadding="0" cellspacing="0">
                                        <tr style="margin-top: 0; color: #74787E; line-height: 2em;">
                                            <td style="width: 10%;">
                                                请求时间:
                                            </td>
                                            <td style="width: 90%">
                                                {{.Timestamp}}
                                            </td>
                                        </tr>

										<tr style="margin-top: 0; color: #74787E; line-height: 2em;">
                                            <td style="width: 10%;">
                                                请求标识:
                                            </td>
                                            <td style="width: 90%">
                                                {{.ID}}
                                            </td>
                                        </tr>

                                        <tr style="margin-top: 0; color: #74787E; font-size: 16px; line-height: 2em;">
                                            <td style="width: 10%;">
                                                请求地址:
                                            </td>
                                            <td style="width: 90%">
                                                {{.URL}}
                                            </td>
                                        </tr>

                                        <tr style="margin-top: 0; color: #74787E; font-size: 16px; line-height: 2em;">
                                            <td style="width: 10%;">
                                                错误信息:
                                            </td>
                                            <td style="width: 90%">
                                                {{.Msg}}
                                            </td>
                                        </tr>

										<tr style="margin-top: 0; color: #74787E; font-size: 16px; line-height: 2em;">
                                            <td style="width: 10%;"><br /></td>
                                            <td style="width: 90%"><br /></td>
                                        </tr>

										<tr style="margin-top: 0; color: #74787E; font-size: 16px; line-height: 2em;">
                                            <td style="width: 10%;">
                                                错误堆栈:
                                            </td>
                                            <td style="width: 90%">
                                                {{.Stack}}
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>

                <!-- Footer -->
                <tr>
                    <td>
                        <table style="width: auto; max-width: 570px; margin: 0 auto; padding: 0; text-align: center;"
                               align="center" width="750" cellpadding="0" cellspacing="0">
                            <tr>
                                <td style="font-family: Arial, 'Helvetica Neue', Helvetica, sans-serif; color: #AEAEAE; padding: 35px; text-align: center;">
                                    <p style="margin-top: 0; color: #74787E; font-size: 12px; line-height: 1.5em;">
                                        &copy; {{.Year}}
                                        All rights reserved.
                                    </p>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>
`
