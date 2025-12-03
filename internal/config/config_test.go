package config

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_WriteRead(t *testing.T) {
	t.Parallel()

	isEnabled := true
	testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")

	expectedConfig := &Config{
		Version: 1,

		Sandboxes: map[string]Sandbox{
			"daytona": {
				Enable: &isEnabled,
				Source: "github.com/daytona/daytona",
				Hash:   "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				Config: &map[string]any{
					"foo": "bar",
				},
			},
		},

		Interceptors: map[string]Interceptor{
			"logger": {
				Enable:   &isEnabled,
				Source:   "github.com/acme/interceptor",
				Hash:     "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				Priority: 10,
				Config: &map[string]any{
					"baz": int64(123),
				},
			},
		},

		Services: map[string]Service{
			"com.spotify/music": {
				Enable: &isEnabled,
				Hash:   "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				Config: &map[string]any{
					"enabled": true,
				},
			},
		},
	}

	err := expectedConfig.WriteFile(testConfigFile)
	require.NoError(t, err)

	actualConfig, err := LoadFile(testConfigFile)
	require.NoError(t, err)

	assert.Equal(t, expectedConfig, actualConfig)

	content, err := os.ReadFile(testConfigFile)
	require.NoError(t, err)

	var decodedContent map[string]any
	_, err = toml.Decode(string(content), &decodedContent)
	require.NoError(t, err)

	assert.Equal(t, int64(1), decodedContent["version"])
	assert.Contains(t, decodedContent, "sandboxes")
	assert.Contains(t, decodedContent, "interceptors")
	assert.Contains(t, decodedContent, "services")
}

func TestConfig_Defaults(t *testing.T) {
	t.Parallel()

	t.Run("should set `enable` to true if unset", func(t *testing.T) {
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")
		tomlData := `
		[services."com.spotify/music"]
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		cfg, err := LoadFile(testConfigFile)
		require.NoError(t, err)

		service := cfg.Services["com.spotify/music"]
		require.NotNil(t, service)
		require.NotNil(t, service.Enable)

		assert.True(t, *service.Enable)
	})

	t.Run("should not override `enable` if set", func(t *testing.T) {
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")

		randomBool := rand.Intn(2) == 1
		tomlData := fmt.Sprintf(`
		[services."com.spotify/music"]
		enable = %v
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`, randomBool)

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		cfg, err := LoadFile(testConfigFile)
		require.NoError(t, err)

		service := cfg.Services["com.spotify/music"]
		require.NotNil(t, service)
		require.NotNil(t, service.Enable == &randomBool)

		assert.False(t, *service.Enable)
	})
}

func TestConfig_Validation(t *testing.T) {
	t.Parallel()

	t.Run("should pass valid hash format", func(t *testing.T) {
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")
		tomlData := `
		[sandboxes.my-sandbox]
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		_, err = LoadFile(testConfigFile)
		assert.NoError(t, err)
	})

	t.Run("should fail invalid hash format with unsupported prefix", func(t *testing.T) {
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")
		tomlData := `
		[sandboxes.my-sandbox]
		hash = "sha1-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		_, err = LoadFile(testConfigFile)
		assert.ErrorContains(t, err, "unsupported hash prefix")
	})

	t.Run("should fail nil hash", func(t *testing.T) {
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")
		tomlData := `
		[sandboxes.my-sandbox]
		source = "some-source"
		`

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		_, err = LoadFile(testConfigFile)
		assert.ErrorContains(t, err, "unsupported hash")
	})
}
