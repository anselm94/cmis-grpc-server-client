package docserverclient

type Config struct {
	AppHost    string
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func NewDefaultConfig() *Config {
	return &Config{
		AppHost:    "localhost",
		AppPort:    ":9999",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "firstuser",
		DBPassword: "password",
		DBName:     "firstdb",
		DBSSLMode:  "disable",
	}
}
