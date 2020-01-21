package cmd

import (
	"testing"

	"os"
)

func TestAppGetenv(t *testing.T) {
	app := &Application{}

	os.Setenv("APP_ENV", "")

	if app.getenv() != AppEnv("development") {
		t.Errorf("Expected equals: %v and %v", app.getenv(), AppEnv("development"))
	}

	app.env = AppEnv("")
	os.Setenv("APP_ENV", "production")

	if app.getenv() != AppEnvProduction {
		t.Errorf("Expected equals: %v and %v", app.getenv(), AppEnvProduction)
	}

	app.env = AppEnv("")
	os.Setenv("APP_ENV", "test")

	if app.getenv() != AppEnv("test") {
		t.Errorf("Expected equals: %v and %v", app.getenv(), AppEnv("development"))
	}

	app.env = AppEnv("")
	os.Setenv("APP_ENV", "foo")

	if app.getenv() != AppEnv("development") {
		t.Errorf("Expected equals: %v and %v", app.getenv(), AppEnv("development"))
	}
}
