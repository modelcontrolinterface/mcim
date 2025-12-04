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
	defaultConnectTimeout := "10s"
	testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")

	expectedConfig := &Config{
		Version: 1,

		Sandboxes: map[string]Sandbox{
			"my-sandbox": {
				Enable:         &isEnabled,
				AutoReconnect:  &isEnabled,
				ConnectTimeout: &defaultConnectTimeout,
				Source:         "static.my-sandbox.io/mci.json",
				Hash:           "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				Config: &map[string]any{
					"foo": "bar",
				},
			},
		},

		Interceptors: map[string]Interceptor{
			"logger": {
				Enable:         &isEnabled,
				AutoReconnect:  &isEnabled,
				ConnectTimeout: &defaultConnectTimeout,
				Source:         "static.mcilog.io/mci.json",
				Hash:           "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				Priority:       10,
				Config: &map[string]any{
					"baz": int64(123),
				},
			},
		},

		Services: map[string]Service{
			"my-service": {
				Enable:         &isEnabled,
				AutoReconnect:  &isEnabled,
				ConnectTimeout: &defaultConnectTimeout,
				Hash:           "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
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

	t.Run("should set `enable`, `auto_reconnect`, and `connect_timeout` to defaults if unset", func(t *testing.T) {
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")
		tomlData := `
		[services."my-service"]
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		cfg, err := LoadFile(testConfigFile)
		require.NoError(t, err)

		service := cfg.Services["my-service"]
		require.NotNil(t, service)

		require.NotNil(t, service.Enable)
		assert.True(t, *service.Enable)

		require.NotNil(t, service.AutoReconnect)
		assert.True(t, *service.AutoReconnect)

		require.NotNil(t, service.ConnectTimeout)
		assert.Equal(t, "10s", *service.ConnectTimeout)
	})

	t.Run("should not override `enable` if set", func(t *testing.T) {
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")

		randomBool := rand.Intn(2) == 1
		tomlData := fmt.Sprintf(`
		[services."my-service"]
		enable = %v
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`, randomBool)

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		cfg, err := LoadFile(testConfigFile)
		require.NoError(t, err)

		service := cfg.Services["my-service"]
		require.NotNil(t, service)
		require.NotNil(t, service.Enable)

		assert.True(t, *service.Enable == randomBool)
	})

	t.Run("should not override `auto_reconnect` if set", func(t *testing.T) {
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")

		randomBool := rand.Intn(2) == 1
		tomlData := fmt.Sprintf(`
		[services."my-service"]
		auto_reconnect = %v
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`, randomBool)

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		cfg, err := LoadFile(testConfigFile)
		require.NoError(t, err)

		service := cfg.Services["my-service"]
		require.NotNil(t, service)
		require.NotNil(t, service.AutoReconnect)

		assert.True(t, *service.AutoReconnect == randomBool)
	})

	t.Run("should not override `connect_timeout` if set", func(t *testing.T) {
		customTimeout := "30s"
		testConfigFile := filepath.Join(t.TempDir(), "test_config.toml")
		tomlData := fmt.Sprintf(`
		[services."my-service"]
		connect_timeout = "%s"
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`, customTimeout)

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		cfg, err := LoadFile(testConfigFile)
		require.NoError(t, err)

		service := cfg.Services["my-service"]
		require.NotNil(t, service)
		require.NotNil(t, service.ConnectTimeout)

		assert.Equal(t, customTimeout, *service.ConnectTimeout)
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
		[sandboxes.my-service]
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
