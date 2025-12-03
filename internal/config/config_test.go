package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_WriteRead(t *testing.T) {
	t.Parallel()

	tbool := true
	hash := "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	tempDir := t.TempDir()
	testConfigFile := filepath.Join(tempDir, "test_config.toml")

	expectedConfig := &Config{
		Version: 1,
		Sandboxes: map[string]Sandbox{
			"daytona": {
				Enable: &tbool,
				Source: "github.com/daytona/daytona",
				Hash:   hash,
				Config: &map[string]any{
					"foo": "bar",
				},
			},
		},

		Interceptors: map[string]Interceptor{
			"logger": {
				Enable:   &tbool,
				Source:   "github.com/acme/interceptor",
				Hash:     hash,
				Priority: 10,
				Config: &map[string]any{
					"baz": int64(123),
				},
			},
		},

		Services: map[string]Service{
			"com.spotify/music": {
				Enable: &tbool,
				Hash:   hash,
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
	t.Parallel() // Run tests in parallel

	validHash := "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	t.Run("it sets Enable to true when it is missing", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigFile := filepath.Join(tempDir, "test_config.toml")
		tomlData := `
		version = 1

		[services."com.spotify/music"]
		hash = "` + validHash + `"
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

	t.Run("it does not override Enable when it is set to false", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigFile := filepath.Join(tempDir, "test_config.toml")
		tomlData := `
		version = 1

		[services."com.spotify/music"]
		enable = false
		hash = "` + validHash + `"
		`

		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)
		cfg, err := LoadFile(testConfigFile)
		require.NoError(t, err)
		service := cfg.Services["com.spotify/music"]
		require.NotNil(t, service)
		require.NotNil(t, service.Enable)
		assert.False(t, *service.Enable)
	})

	t.Run("it does not override Enable when it is set to true", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigFile := filepath.Join(tempDir, "test_config.toml")
		tomlData := `
		version = 1

		[services."com.spotify/music"]
		enable = true
		hash = "` + validHash + `"
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
}

func TestConfig_Validation(t *testing.T) {
	t.Parallel() // Run tests in parallel

	validHash := "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	t.Run("valid hash format should pass validation", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigFile := filepath.Join(tempDir, "test_config.toml")
		tomlData := `
		version = 1

		[sandboxes.my-sandbox]
		hash = "` + validHash + `"
		`
		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		_, err = LoadFile(testConfigFile)
		assert.NoError(t, err)
	})

	t.Run("invalid hash format - wrong prefix", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigFile := filepath.Join(tempDir, "test_config.toml")
		tomlData := `
		version = 1

		[sandboxes.my-sandbox]
		hash = "sha1-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		`
		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		_, err = LoadFile(testConfigFile)
		assert.ErrorContains(t, err, "invalid hash format")
	})

	t.Run("invalid hash format - wrong length", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigFile := filepath.Join(tempDir, "test_config.toml")
		tomlData := `
		version = 1

		[sandboxes.my-sandbox]
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85" # 63 chars
		`
		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		_, err = LoadFile(testConfigFile)
		assert.ErrorContains(t, err, "invalid hash format")
	})

	t.Run("invalid hash format - non-hex characters", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigFile := filepath.Join(tempDir, "test_config.toml")
		tomlData := `
		version = 1

		[sandboxes.my-sandbox]
		hash = "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85G" # 'G' is not hex
		`
		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		_, err = LoadFile(testConfigFile)
		assert.ErrorContains(t, err, "invalid hash format")
	})

	t.Run("nil hash should fail validation", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigFile := filepath.Join(tempDir, "test_config.toml")
		tomlData := `
		version = 1

		[sandboxes.my-sandbox]
		source = "some-source"
		`
		err := os.WriteFile(testConfigFile, []byte(tomlData), 0o644)
		require.NoError(t, err)

		_, err = LoadFile(testConfigFile)
		assert.ErrorContains(t, err, "requires a hash")
	})
}

