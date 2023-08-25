package mkdocs

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func writeConfig(configFile string, config *Config) error {
	err := os.MkdirAll(filepath.Dir(configFile), 0o755)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, overrideData(data), 0o644)
}

// simple solution to replace the empty strings, as mapping those fields
// into the MkDocs config is not supported yet
func overrideData(data []byte) []byte {
	content := "# This file is autogenerated by the 'modulegen' tool.\n" + string(data)
	content = setEmoji(content, "generator", "to_svg")
	content = setEmoji(content, "index", "twemoji")
	return []byte(content)
}

func setEmoji(content string, key string, value string) string {
	old := "emoji_" + key + `: ""`
	new := "emoji_" + key + ": !!python/name:materialx.emoji." + value
	return strings.ReplaceAll(content, old, new)
}
