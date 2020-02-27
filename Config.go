package docserverclient

// Config holds the configuration for both the server and client
type Config struct {
	GrpcAppHost string
	GrpcAppPort string
	CmisAppHost string
	CmisAppPort string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
}

// NewDefaultConfig creates a Default `Config`.
// Seek attention in updating the values accordingly
func NewDefaultConfig() *Config {
	return &Config{
		GrpcAppHost: "localhost", // gRPC Server host address
		GrpcAppPort: ":9998",     // gRPC Server port
		CmisAppHost: "localhost", // CMIS Server host address
		CmisAppPort: ":8000",     // CMIS Server port
		DBHost:      "localhost", // DB host address
		DBPort:      "5432",      // DB port
		DBUser:      "firstuser", // DB username
		DBPassword:  "password",  // DB password
		DBName:      "firstdb",   // DB name
		DBSSLMode:   "disable",   // Disable SSL connection to DB
	}
}
