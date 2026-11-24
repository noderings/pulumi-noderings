// Copyright 2026, NodeRings.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package examples runs the example programs against the live NodeRings API.
//
// Each test deploys a program, checks that a second preview shows no changes,
// refreshes it and destroys it. The tests need NODERINGS_TOKEN and are skipped
// without it. Tests that order a VM cost money and also need
// NODERINGS_TEST_BILLED=1. See CONTRIBUTING.md.
package examples

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/require"
)

// testDomain is the parent of the throwaway zones the tests create.
func testDomain() string {
	if d := os.Getenv("NODERINGS_TEST_DOMAIN"); d != "" {
		return d
	}
	return "example.com"
}

// uniqueZone returns a zone name that no other test run uses.
func uniqueZone(t *testing.T, prefix string) string {
	t.Helper()
	b := make([]byte, 4)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return fmt.Sprintf("%s-%s.%s", prefix, hex.EncodeToString(b), testDomain())
}

func skipWithoutToken(t *testing.T) {
	t.Helper()
	if os.Getenv("NODERINGS_TOKEN") == "" {
		t.Skip("NODERINGS_TOKEN is not set; skipping a test that calls the NodeRings API")
	}
}

func skipUnlessBilled(t *testing.T) {
	t.Helper()
	skipWithoutToken(t)
	if os.Getenv("NODERINGS_TEST_BILLED") != "1" {
		t.Skip("NODERINGS_TEST_BILLED is not 1; skipping a test that orders a billed VM")
	}
}

func getCwd(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	return cwd
}

// getBaseOptions runs every program against the provider binary in ../bin.
func getBaseOptions(t *testing.T) integration.ProgramTestOptions {
	t.Helper()
	binPath, err := filepath.Abs("../bin")
	require.NoError(t, err)
	return integration.ProgramTestOptions{
		LocalProviders: []integration.LocalDependency{
			{Package: "noderings", Path: binPath},
		},
		// The API rejects some calls during the first seconds after a create,
		// and the provider retries them, so allow a slow run.
		RetryFailedSteps: true,
	}
}

func getJSBaseOptions(t *testing.T) integration.ProgramTestOptions {
	t.Helper()
	return getBaseOptions(t).With(integration.ProgramTestOptions{
		Dependencies: []string{"@noderings/pulumi"},
	})
}

func getPythonBaseOptions(t *testing.T) integration.ProgramTestOptions {
	t.Helper()
	sdk, err := filepath.Abs(filepath.Join("..", "sdk", "python", "bin"))
	require.NoError(t, err)
	return getBaseOptions(t).With(integration.ProgramTestOptions{
		Dependencies: []string{sdk},
	})
}

func getGoBaseOptions(t *testing.T) integration.ProgramTestOptions {
	t.Helper()
	sdk, err := filepath.Abs("../sdk")
	require.NoError(t, err)
	return getBaseOptions(t).With(integration.ProgramTestOptions{
		Dependencies: []string{fmt.Sprintf("github.com/noderings/pulumi-noderings/sdk=%s", sdk)},
	})
}
