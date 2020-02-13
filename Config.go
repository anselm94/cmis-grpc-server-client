package docserverclient

// Config holds the configuration for both the server and client
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

// NewDefaultConfig creates a Default `Config`.
// Seek attention in updating the values accordingly
func NewDefaultConfig() *Config {
	return &Config{
		AppHost:    "localhost", // Server host address
		AppPort:    ":9998",     // Server port
		DBHost:     "localhost", // DB host address
		DBPort:     "5432",      // DB port
		DBUser:     "firstuser", // DB username
		DBPassword: "password",  // DB password
		DBName:     "firstdb",   // DB name
		DBSSLMode:  "disable",   // Disable SSL connection to DB
	}
}
