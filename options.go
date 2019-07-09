package acm

import (
	"fmt"
)

// apiPath API相对路径
type apiPath string

func (ap apiPath) URL(host string, port uint16) string {
	return fmt.Sprintf("http://%s:%d/%s", host, port, ap)
}

const (
	getServer apiPath = "diamond-server/diamond"
	getConfig apiPath = "diamond-server/config.co"
	setConfig apiPath = "diamond-server/basestone.do?method=syncUpdateAll"
	delConfig apiPath = "diamond-server/datum.do?method=deleteAllDatums"
)

// authCreds 鉴权凭证
type authCreds struct {
	accessKey string
	secretKey string
}

func (ac *authCreds) check() {
	if ac.accessKey == "" {
		panic("aliacm: access key is required")
	} else if ac.secretKey == "" {
		panic("aliacm: secret key is required")
	}
}

// Options ACM客户端配置
type Options struct {
	endpoint  string
	namespace string
	groupName string
	authCreds authCreds
}

type Option func(options *Options)

func newOptions(fns []Option) Options {
	opts := Options{
		groupName: "DEFAULT_GROUP",
	}
	for _, fn := range fns {
		fn(&opts)
	}
	if opts.endpoint == "" {
		panic("aliacm: endpoint is required")
	} else if opts.namespace == "" {
		panic("aliacm: namespace is required")
	}
	opts.authCreds.check()
	return opts
}

func Endpoint(ep string) Option {
	return func(options *Options) {
		options.endpoint = ep
	}
}

func Namespace(ns string) Option {
	return func(options *Options) {
		options.namespace = ns
	}
}

func GroupName(gn string) Option {
	return func(options *Options) {
		options.groupName = gn
	}
}

func AuthCreds(ak, sk string) Option {
	return func(options *Options) {
		options.authCreds = authCreds{
			accessKey: ak,
			secretKey: sk,
		}
	}
}
