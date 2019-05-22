package tf

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// SetEnv sets an environment variable and returns a reset function to ensure
// the environment is always returned to it's previous state:
//
//   resetEnv := tf.SetEnv(t, "HOME", "/somewhere/else")
//   defer resetEnv()
//
// If you would like to set multiple environment variables, see SetEnvs().
func SetEnv(t *testing.T, name, value string) (resetEnv func()) {
	originalValue, exists := os.LookupEnv(name)
	assert.NoError(t, os.Setenv(name, value))

	return func() {
		if exists {
			assert.NoError(t, os.Setenv(name, originalValue))
		} else {
			assert.NoError(t, os.Unsetenv(name))
		}
	}
}

// SetEnvs works the same way as SetEnv, but on multiple environment variables:
//
//   resetEnv := tf.SetEnvs(t, map[string]string{
//       "HOME":  "/somewhere/else",
//       "DEBUG": "on",
//   })
//   defer resetEnv()
//
func SetEnvs(t *testing.T, env map[string]string) (resetEnv func()) {
	type original struct {
		value  string
		exists bool
	}

	originalValues := map[string]original{}
	for name, value := range env {
		originalValue, exists := os.LookupEnv(name)
		originalValues[name] = original{ originalValue, exists }
		assert.NoError(t, os.Setenv(name, value))
	}

	return func() {
		for name, originalValue := range originalValues {
			if originalValue.exists {
				assert.NoError(t, os.Setenv(name, originalValue.value))
			} else {
				assert.NoError(t, os.Unsetenv(name))
			}
		}
	}
}
