package main

import "fmt"

// golang 选项设计模式
type Options struct {
	StrOption1 string
	StrOption2 string
	StrOption3 string
	IntOption1 int
	IntOption2 int
	IntOption3 int
}

type Option func(opts *Options)

func InitOption(opts ...Option) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	fmt.Printf("init option:%#v\n", options)
}

func WithStringOption1(str string) Option {
	return func(opts *Options) {
		opts.StrOption1 = str
	}
}

func WithStringOption2(str string) Option {
	return func(opts *Options) {
		opts.StrOption2 = str
	}
}

func WithStringOption3(str string) Option {
	return func(opts *Options) {
		opts.StrOption3 = str
	}
}

func WithIntOption1(val int) Option {
	return func(opts *Options) {
		opts.IntOption1 = val
	}
}

func WithIntOption2(val int) Option {
	return func(opts *Options) {
		opts.IntOption2 = val
	}
}

func WithIntOption3(val int) Option {
	return func(opts *Options) {
		opts.IntOption3 = val
	}
}

func main() {
	InitOption(
		WithStringOption1("str1"),
		WithStringOption3("str3"),
		WithIntOption1(1),
		WithIntOption2(2),
		WithIntOption3(3),
		WithStringOption2("str2"),
	)
}
