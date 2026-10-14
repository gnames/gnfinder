package config

// LogConfig provides typical settings for application logs.
type LogConfig struct {
	// Format can be 'json', 'text' or 'tint' (user-facing and colored).
	Format string `mapstructure:"format"      yaml:"format"`
	// Level of logging -- 'error', 'warn', 'info', 'debug'
	Level string `mapstructure:"level"       yaml:"level"`
	// Destination can be a log file (to default place), STDERR or STDOUT
	Destination string `mapstructure:"destination" yaml:"destination"`
}
