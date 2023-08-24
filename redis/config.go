package redis

// Config redis config
type Config struct {
	Client  `mapstructure:"-"`
	Address string       `mapstructure:"address"`
	DB      int          `mapstructure:"db"`
	Options CommonOption `mapstructure:",squash"`
}

// ClusterConfig cluster config
type ClusterConfig struct {
	Client       `mapstructure:"-"`
	StartupNodes []string     `mapstructure:"startup_nodes"`
	Options      CommonOption `mapstructure:",squash"`
}
