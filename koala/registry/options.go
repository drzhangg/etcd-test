package registry

import "time"

type Options struct {
	Address      []string //注册中心地址
	Timeout      time.Duration
	RegistryPath string //注册路径
	HeartBeat    int64  //心跳时间
}

type Option func(opts *Options)

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

func WithAddress(adds []string) Option {
	return func(opts *Options) {
		opts.Address = adds
	}
}

func WithRegistryPath(path string) Option {
	return func(opts *Options) {
		opts.RegistryPath = path
	}
}

func WithHearBeat(hearBeat int64) Option {
	return func(opts *Options) {
		opts.HeartBeat = hearBeat
	}
}
