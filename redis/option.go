package redis

const (
	defaultTimeout = 2 // 2 second
	defaultIdle    = 200
)

// CommonOption CommonOption.
type CommonOption struct {
	Username       string `mapstructure:"username,omitempty"`
	Password       string `mapstructure:"password,omitempty"`
	ConnectTimeout int    `mapstructure:"connect_timeout,omitempty"`
	ReadTimeout    int    `mapstructure:"read_timeout,omitempty"`
	WriteTimeout   int    `mapstructure:"write_timeout,omitempty"`
	MaxIdle        int    `mapstructure:"max_idle,omitempty"`
}

// Convert convert option to options func.
func (c *CommonOption) Convert() []Options {
	var result []Options
	if c.Username != "" {
		result = append(result, WithUsername(c.Username))
	}
	if c.Password != "" {
		result = append(result, WithPassword(c.Password))
	}
	if c.ConnectTimeout != 0 {
		result = append(result, WithConnectTimeout(c.ConnectTimeout))
	}
	if c.ReadTimeout != 0 {
		result = append(result, WithReadTimeout(c.ReadTimeout))
	}
	if c.WriteTimeout != 0 {
		result = append(result, WithWriteTimeout(c.WriteTimeout))
	}
	if c.MaxIdle != 0 {
		result = append(result, WithMaxIdle(c.MaxIdle))
	}
	return result
}

// Options Options
type Options func(o *CommonOption)

// WithUsername WithUsername
func WithUsername(username string) Options {
	return func(o *CommonOption) {
		o.Username = username
	}
}

// WithPassword WithPassword
func WithPassword(password string) Options {
	return func(o *CommonOption) {
		o.Password = password
	}
}

// WithConnectTimeout WithConnectTimeout
func WithConnectTimeout(timeout int) Options {
	return func(o *CommonOption) {
		o.ConnectTimeout = timeout
	}
}

// WithReadTimeout WithReadTimeout
func WithReadTimeout(timeout int) Options {
	return func(o *CommonOption) {
		o.ReadTimeout = timeout
	}
}

// WithWriteTimeout WithWriteTimeout
func WithWriteTimeout(timeout int) Options {
	return func(o *CommonOption) {
		o.WriteTimeout = timeout
	}
}

// WithMaxIdle WithMaxIdle
func WithMaxIdle(max int) Options {
	return func(o *CommonOption) {
		o.MaxIdle = max
	}
}

func defaultOption() *CommonOption {
	return &CommonOption{
		ConnectTimeout: defaultTimeout,
		ReadTimeout:    defaultTimeout,
		WriteTimeout:   defaultTimeout,
		MaxIdle:        defaultIdle,
	}
}
