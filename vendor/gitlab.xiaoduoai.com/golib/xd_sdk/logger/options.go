package logger

type Config struct {
	Logger  `mapstructure:"-"`
	Options Options `mapstructure:",squash"`
}

type Options struct {
	Level     string `mapstructure:"level" json:"level" toml:"level"`
	File      string `mapstructure:"file" json:"file" toml:"file"`
	ErrFile   string `mapstructure:"err_file" json:"err_file" toml:"err_file"`
	CrashFile string `mapstructure:"crash_file" json:"crash_file" toml:"crash_file"`
	AppName   string `mapstructure:"app_name" json:"app_name" toml:"app_name"`
	Format    string `mapstructure:"format" json:"format" toml:"format"`
	WithStack bool   `mapstructure:"with_stack" json:"with_stack" toml:"with_stack"`
	WithTrace bool   `mapstructure:"with_trace" json:"with_trace" toml:"with_trace"`
}

func newOptions(opts ...Option) Options {
	options := Options{
		Level:     "",
		File:      "",
		ErrFile:   "",
		WithTrace: true,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

type Option func(*Options)

func WithLevel(level string) Option {
	return func(options *Options) {
		options.Level = level
	}
}

func WithFile(file string) Option {
	return func(options *Options) {
		options.File = file
	}
}

func WithErrFile(errFile string) Option {
	return func(options *Options) {
		options.ErrFile = errFile
	}
}

// WithStack 配置日志中是否添加错误堆栈信息
func WithStack(enable bool) Option {
	return func(options *Options) {
		options.WithStack = enable
	}
}

// WithTrace 配置日志中是否添加 trace 字段
func WithTrace(enable bool) Option {
	return func(options *Options) {
		options.WithTrace = true
	}
}

//func WithShadowFile(shadowFile string) Option {
//	return func(options *Options) {
//		options.ShadowFile = shadowFile
//	}
//}
