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
	WithHttpCodes(codes map[int]int) map[int]int
	GetText(code int) string
	GetHttpCode(code int) int
	GetError(code int) error
}

type code struct {
	local     string
	codeTexts map[int]string
	httpCodes map[int]int
}

const (
	/* 系统相关 */
	ServerError           = 100001
	/* 应用相关 */
	AuthorizationError    = 200001
	ParamBindError        = 200002
	RequestError          = 200003
	ParamValidateError    = 200004
	UnknownError          = 200005
	DataNotExistError     = 200006
	DataExistError        = 200007
	RequestNotFoundError  = 200008
	DataDeleteError       = 200009
	ResourceNotExistError = 200010
	DataSelectError       = 200011
	DataCreateError       = 200012
	DataUpdateError       = 200013
)

func New(local string) Code {
	var cd = &code{
		local:     ZhCN,
		codeTexts: zhCNText,
		httpCodes: httpCode,
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

func (c *code) WithHttpCodes(codes map[int]int) map[int]int {
	for key, value := range codes {
		c.httpCodes[key] = value
	}
	return c.httpCodes
}

func (c *code) GetText(code int) string {
	return c.codeTexts[code]
}

func (c *code) GetHttpCode(code int) int {
	return c.httpCodes[code]
}

func (c *code) GetError(code int) error {
	return errors.New(c.codeTexts[code])
}
