package code

import (
	"errors"
	"strings"
)

const (
	// ZhCN 简体中文 - 中国
	ZhCN = "zh-cn"
	// EnUS 英文 - 美国
	EnUS = "en-us"
)

var cInterface Code

func init() { New(ZhCN) }

var (
	inLocal = []string{ZhCN, EnUS}
)

type Code interface {
	WithTexts(local string, texts map[int]string) map[int]string
	GetText(code int) string
	GetError(code int) error
}

type code struct {
	local     string
	codeTexts map[int]string
}

const (
	/* 系统相关 */
	ServerError           = 100001
	AuthorizationError    = 100002
	ParamBindError        = 100003
	RequestError          = 100004
	ParamValidateError    = 100005
	UnknownError          = 100006
	DataNotExistError     = 100007
	DataExistError        = 100008
	RequestNotFoundError  = 100009
	DataDeleteError       = 100010
	ResourceNotExistError = 100011
	DataSelectError       = 100012
	DataCreateError       = 100013
	DataUpdateError       = 100014
)

func New(local string) Code {
	var cd = &code{
		local:     ZhCN,
		codeTexts: zhCNText,
	}

	local = strings.ToLower(local)
	for _, value := range inLocal {
		if value == local {
			cd.local = value
			switch cd.local {
			case ZhCN:
				cd.codeTexts = zhCNText
			case EnUS:
				cd.codeTexts = enUSText
			}
			break
		}
	}

	cInterface = cd
	return cd
}

func Get() Code { return cInterface }

func (c *code) WithTexts(local string, texts map[int]string) map[int]string {
	local = strings.ToLower(local)
	for key, value := range texts {
		c.codeTexts[key] = value
	}
	return c.codeTexts
}

func (c *code) GetText(code int) string {
	return c.codeTexts[code]
}

func (c *code) GetError(code int) error {
	return errors.New(c.codeTexts[code])
}
