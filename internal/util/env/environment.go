package env

import (
	"github.com/joho/godotenv"
	"os"
)

type IEnvironment interface {
	Init()
	Get(key string) string
	Set(key string, value string) error
	GetHostname() (string, error)
}

type Environment struct {}

// New
// Returns new Environment.
func New() IEnvironment {
	return &Environment{}
}

func (e *Environment) Init() {
	var err error

	// Load default env file.
	err = godotenv.Load(".env")
	if err != nil {
		panic("Panicked while loading environment.")
	}

	appEnv := os.Getenv(AppEnvironment)
	appName := os.Getenv(AppName)
	if len(appEnv) < 1 {
		panic("APP_ENVIRONMENT variable is not set.")
	}
	if len(appName) < 1 {
		panic("APP_NAME variable is not set.")
	}

	// Overload specific env file.
	fileName := ".env." + appEnv
	if _, err := os.Stat(fileName); err == nil {
		err = godotenv.Overload(fileName)
		if err != nil {
			panic("Panicked while loading environment.")
		}
	}
}

func (e *Environment) Get(key string) string {
	return os.Getenv(key)
}

func (e *Environment) Set(key string, value string) error {
	return os.Setenv(key, value)
}

func (e *Environment) GetHostname() (string, error) {
	return os.Hostname()
}
