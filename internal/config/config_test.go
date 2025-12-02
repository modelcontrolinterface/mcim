package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	tempDir := t.TempDir()
	testConfigFile := filepath.Join(tempDir, "test_config.toml")

	expectedConfig := &Config{
		Version: 1,
		Sandboxes: map[string]Sandbox{
			"daytona": {
				Source:  "github.com/daytona/daytona",
				Version: "latest",
				Config: map[string]any{
					"foo": "bar",
				},
			},
		},
		Interceptors: map[string]Interceptor{
			"acme": {
				Source:   "github.com/acme/interceptor",
				Version:  "v1.0.0",
				Priority: 10,
				Config: map[string]any{
					"baz": int64(123),
				},
			},
		},
		Services: map[string]Service{
			"com.spotify/music": {
				Version: "v1.2.3",
				Config: map[string]any{
					"enabled": true,
				},
			},
		},
	}

	err := expectedConfig.ToFile(testConfigFile)
	require.NoError(t, err)

	actualConfig, err := FromFile(testConfigFile)
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
