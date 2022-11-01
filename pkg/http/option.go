package http

import (
	"github.com/raylin666/go-utils/middleware"
	"time"
	"ult/pkg/proposal"
)

type Option func(opt *option)

type option struct {
	cors struct {
		domains []string
	}
	pprof       bool
	rate        bool
	openBrowser string
	alertNotify proposal.NotifyHandler
	timeout     time.Duration
	middlewares []middleware.HTTPHandler
}

func EnableCors(domains []string) Option {
	return func(opt *option) {
		opt.cors.domains = domains
	}
}

func EnablePProf() Option {
	return func(opt *option) {
		opt.pprof = true
	}
}

func EnableRate() Option {
	return func(opt *option) {
		opt.rate = true
	}
}

func EnableOpenBrowser(uri string) Option {
	return func(opt *option) {
		opt.openBrowser = uri
	}
}

func EnableAlertNotify(handler proposal.NotifyHandler) Option {
	return func(opt *option) {
		opt.alertNotify = handler
	}
}

func WithTimeout(ts time.Duration) Option {
	return func(opt *option) {
		opt.timeout = ts
	}
}

func WithMiddleware(m ...middleware.HTTPHandler) Option {
	return func(opt *option) {
		opt.middlewares = append(opt.middlewares, m...)
	}
}
