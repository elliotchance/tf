package tf_test

import (
	"github.com/elliotchance/tf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSetEnv(t *testing.T) {
	t.Run("Replace", func(t *testing.T) {
		require.NoError(t, os.Setenv("FOO", "bar"))
		require.Equal(t, os.Getenv("FOO"), "bar")

		resetEnv := tf.SetEnv(t, "FOO", "baz")
		assert.Equal(t, os.Getenv("FOO"), "baz")

		resetEnv()
		assert.Equal(t, os.Getenv("FOO"), "bar")
	})

	t.Run("New", func(t *testing.T) {
		_, exists := os.LookupEnv("BAR")
		require.False(t, exists)

		resetEnv := tf.SetEnv(t, "BAR", "baz")
		assert.Equal(t, os.Getenv("BAR"), "baz")

		resetEnv()

		_, exists = os.LookupEnv("BAR")
		require.False(t, exists)
	})
}

func TestSetEnvs(t *testing.T) {
	t.Run("Multiple", func(t *testing.T) {
		require.NoError(t, os.Setenv("FOO", "bar"))
		require.Equal(t, os.Getenv("FOO"), "bar")

		_, exists := os.LookupEnv("BAR")
		require.False(t, exists)

		resetEnv := tf.SetEnvs(t, map[string]string{
			"FOO": "baz",
			"BAR": "baz",
		})

		assert.Equal(t, os.Getenv("FOO"), "baz")
		assert.Equal(t, os.Getenv("BAR"), "baz")

		resetEnv()

		assert.Equal(t, os.Getenv("FOO"), "bar")

		_, exists = os.LookupEnv("BAR")
		require.False(t, exists)
	})
}
