package server

// DefaultConfig returns a default config.
func DefaultEnvConfig() *EnvConfig {
	return &EnvConfig{
		Url:              "localhost",
		Port:             8080,
		DBHost:           "",
		DBPort:           1234,
		DBUser:           "",
		DBPass:           "",
		DBName:           "",
		SMTPHost:         "",
		SMTPPort:         1234,
		SMTPUser:         "",
		SMTPPass:         "",
		EmailFromAddress: "",
		RedisAddr:        "",
		RedisPwd:         "",
	}
}

// Config to load the values from environment config.
type EnvConfig struct {
	// Server base URL
	Url string `mapstructure:"URL"`
	// Server port
	Port int `mapstructure:"PORT"`
	// DBHost specifies the database host
	DBHost string `mapstructure:"DB_HOST"`
	// DBPort specifies the database port
	DBPort int `mapstructure:"DB_PORT"`
	// DBUser specifies the database username
	DBUser string `mapstructure:"DB_USER"`
	// DBPass specifies the database password
	DBPass string `mapstructure:"DB_PASS"`
	// DBName specifies the database name
	DBName string `mapstructure:"DB_NAME"`
	// SMTP host URL
	SMTPHost string `mapstructure:"SMTP_HOST"`
	// SMTP port
	SMTPPort int `mapstructure:"SMTP_PORT"`
	// SMTP username
	SMTPUser string `mapstructure:"SMTP_USER"`
	// SMTP password
	SMTPPass string `mapstructure:"SMTP_PASS"`
	// From address for the email to be sent
	EmailFromAddress string `mapstructure:"EMAIL_FROM_ADDRESS"`
	// Redis server address
	RedisAddr string `mapstructure:"REDIS_ADDR"`
	// Redis password
	RedisPwd string `mapstructure:"REDIS_PWD"`
}
