package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go/modulegen/internal/dependabot"
)

func TestGetDependabotConfigFile(t *testing.T) {
	ctx := NewContext(filepath.Join(t.TempDir(), "testcontainers-go"))

	githubDir := ctx.GithubDir()
	cfgFile := ctx.DependabotConfigFile()
	err := os.MkdirAll(githubDir, 0o777)
	require.NoError(t, err)

	err = os.WriteFile(cfgFile, []byte{}, 0o777)
	require.NoError(t, err)

	file := ctx.DependabotConfigFile()
	require.NotNil(t, file)

	assert.True(t, strings.HasSuffix(file, filepath.Join("testcontainers-go", ".github", "dependabot.yml")))
}

func TestExamplesHasDependabotEntry(t *testing.T) {
	ctx := getTestRootContext(t)
	examples, err := ctx.GetExamples()
	require.NoError(t, err)
	dependabotUpdates, err := dependabot.GetUpdates(ctx.DependabotConfigFile())
	require.NoError(t, err)

	exampleUpdates := []dependabot.Update{}
	// exclude the Go modules from the examples updates
	for _, update := range dependabotUpdates {
		if strings.HasPrefix(update.Directory, "/examples/") {
			exampleUpdates = append(exampleUpdates, update)
		}
	}

	assert.Equal(t, len(exampleUpdates), len(examples))

	// all example modules exist in the dependabot updates
	for _, example := range examples {
		found := false
		for _, exampleUpdate := range exampleUpdates {
			dependabotDir := "/examples/" + example

			assert.Equal(t, exampleUpdate.Schedule.Interval, "monthly")

			if dependabotDir == exampleUpdate.Directory {
				found = true
				continue
			}
		}
		assert.True(t, found, "example %s is not present in the dependabot updates", example)
	}
}

func TestModulesHasDependabotEntry(t *testing.T) {
	ctx := getTestRootContext(t)
	modules, err := ctx.GetModules()
	require.NoError(t, err)
	dependabotUpdates, err := dependabot.GetUpdates(ctx.DependabotConfigFile())
	require.NoError(t, err)

	moduleUpdates := []dependabot.Update{}
	// exclude the Go modules from the examples updates
	for _, update := range dependabotUpdates {
		if strings.HasPrefix(update.Directory, "/modules/") {
			moduleUpdates = append(moduleUpdates, update)
		}
	}
	assert.Equal(t, len(moduleUpdates), len(modules))

	// all module modules exist in the dependabot updates
	for _, module := range modules {
		found := false
		for _, moduleUpdate := range moduleUpdates {
			dependabotDir := "/modules/" + module

			assert.Equal(t, moduleUpdate.Schedule.Interval, "monthly")

			if dependabotDir == moduleUpdate.Directory {
				found = true
				continue
			}
		}
		assert.True(t, found, "module %s is not present in the dependabot updates", module)
	}
}

func getTestRootContext(t *testing.T) *Context {
	current, err := os.Getwd()
	require.NoError(t, err)
	return NewContext(filepath.Dir(current))
}
