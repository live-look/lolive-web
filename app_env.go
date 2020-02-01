package camforchat

const (
	// AppEnvDevelopment for development
	AppEnvDevelopment AppEnv = "development"
	// AppEnvTest for tests
	AppEnvTest AppEnv = "test"
	// AppEnvProduction for production
	AppEnvProduction AppEnv = "production"
)

// AppEnv describes application environment
type AppEnv string
