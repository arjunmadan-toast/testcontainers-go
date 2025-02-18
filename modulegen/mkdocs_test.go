package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go/modulegen/internal/mkdocs"
)

func TestGetMkDocsConfigFile(t *testing.T) {
	tmpCtx := NewContext(filepath.Join(t.TempDir(), "testcontainers-go"))
	cfgFile := tmpCtx.MkdocsConfigFile()
	err := os.MkdirAll(tmpCtx.RootDir, 0o777)
	require.NoError(t, err)

	err = os.WriteFile(cfgFile, []byte{}, 0o777)
	require.NoError(t, err)

	file := tmpCtx.MkdocsConfigFile()
	require.NotNil(t, file)

	assert.True(t, strings.HasSuffix(file, filepath.Join("testcontainers-go", "mkdocs.yml")))
}

func TestReadMkDocsConfig(t *testing.T) {
	tmpCtx := NewContext(filepath.Join(t.TempDir(), "testcontainers-go"))
	err := os.MkdirAll(tmpCtx.RootDir, 0o777)
	require.NoError(t, err)

	err = copyInitialMkdocsConfig(t, tmpCtx)
	require.NoError(t, err)

	config, err := mkdocs.ReadConfig(tmpCtx.MkdocsConfigFile())
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "Testcontainers for Go", config.SiteName)
	assert.Equal(t, "https://github.com/testcontainers/testcontainers-go", config.RepoURL)
	assert.Equal(t, "edit/main/docs/", config.EditURI)

	// theme
	theme := config.Theme
	assert.Equal(t, "material", theme.Name)

	// nav bar
	nav := config.Nav
	assert.Equal(t, "index.md", nav[0].Home)
	assert.Greater(t, len(nav[2].Features), 0)
	assert.Greater(t, len(nav[3].Modules), 0)
	assert.Greater(t, len(nav[4].Examples), 0)
}

func TestExamples(t *testing.T) {
	ctx := getTestRootContext(t)
	examples, err := ctx.GetExamples()
	require.NoError(t, err)
	examplesDocs, err := ctx.GetExamplesDocs()
	require.NoError(t, err)

	// we have to remove the index.md file from the examples docs
	assert.Equal(t, len(examplesDocs)-1, len(examples))

	// all example modules exist in the documentation
	for _, example := range examples {
		found := false
		for _, exampleDoc := range examplesDocs {
			markdownName := example + ".md"

			if markdownName == exampleDoc {
				found = true
				continue
			}
		}
		assert.True(t, found, "example %s is not present in the docs", example)
	}
}

func copyInitialMkdocsConfig(t *testing.T, tmpCtx *Context) error {
	ctx := getTestRootContext(t)
	return mkdocs.CopyConfig(ctx.MkdocsConfigFile(), tmpCtx.MkdocsConfigFile())
}
