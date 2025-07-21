package config

type ServerConfig struct {
	LogLevel    string `mapstructure:"SERVER_LOG_LEVEL"`
	Port        string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"SERVER_ENVIRONMENT"`
}

type DatabaseConfig struct {
	User     string `mapstructure:"DATABASE_USER"`
	Password string `mapstructure:"DATABASE_PASSWORD"`
	Host     string `mapstructure:"DATABASE_HOST"`
	Port     string `mapstructure:"DATABASE_PORT"`
	Name     string `mapstructure:"DATABASE_NAME"`
}

type RedisConfig struct {
	User     string `mapstructure:"REDIS_USER"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
}

type GoogleConfig struct {
	ClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	ClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	RedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
}

type MercadoPagoConfig struct {
	Token string `mapstructure:"MERCADO_PAGO_TOKEN"`
}

type JWTConfig struct {
	SecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

type AWSConfig struct {
	DefaultRegion string `mapstructure:"AWS_DEFAULT_REGION"`
}

type FrontendConfig struct {
	AuthRedirect string `mapstructure:"FRONTEND_AUTH_REDIRECT"`
}

type Config struct {
	Server      *ServerConfig
	Database    *DatabaseConfig
	Redis       *RedisConfig
	Google      *GoogleConfig
	MercadoPago *MercadoPagoConfig
	JWT         *JWTConfig
	AWS         *AWSConfig
	Frontend    *FrontendConfig
}
