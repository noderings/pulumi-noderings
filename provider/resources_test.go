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

package noderings

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/noderings/pulumi-noderings/provider/pkg/version"
)

func TestMain(m *testing.M) {
	// The release build sets the version with -ldflags; tests need one too.
	if version.Version == "" {
		version.Version = "0.0.1-test"
	}
	os.Exit(m.Run())
}

// Every Terraform resource and data source must have a Pulumi token, so a new
// upstream resource cannot silently go missing from the SDKs.
func TestEveryUpstreamItemIsMapped(t *testing.T) {
	prov := Provider()

	prov.P.ResourcesMap().Range(func(name string, _ shim.Resource) bool {
		info, ok := prov.Resources[name]
		if assert.True(t, ok, "resource %s has no mapping", name) {
			assert.True(t, strings.HasPrefix(string(info.Tok), "noderings:index/"), "resource %s token %s", name, info.Tok)
		}
		return true
	})
	prov.P.DataSourcesMap().Range(func(name string, _ shim.Resource) bool {
		info, ok := prov.DataSources[name]
		if assert.True(t, ok, "data source %s has no mapping", name) {
			assert.True(t, strings.HasPrefix(string(info.Tok), "noderings:index/get"), "data source %s token %s", name, info.Tok)
		}
		return true
	})
	assert.NotEmpty(t, prov.Resources)
	assert.NotEmpty(t, prov.DataSources)
}

func TestTokenIsSecretAndReadsTheEnvironment(t *testing.T) {
	token := Provider().Config["token"]
	require.NotNil(t, token)
	require.NotNil(t, token.Secret)
	assert.True(t, *token.Secret)
	require.NotNil(t, token.Default)
	assert.Contains(t, token.Default.EnvVars, "NODERINGS_TOKEN")
}

// The committed schema must cover everything the provider maps. If this fails,
// run `make tfgen` and commit the result.
func TestCommittedSchemaIsCurrent(t *testing.T) {
	raw, err := os.ReadFile("cmd/pulumi-resource-noderings/schema.json")
	require.NoError(t, err)
	var schema struct {
		Name      string                     `json:"name"`
		Resources map[string]json.RawMessage `json:"resources"`
		Functions map[string]json.RawMessage `json:"functions"`
	}
	require.NoError(t, json.Unmarshal(raw, &schema))
	assert.Equal(t, "noderings", schema.Name)

	prov := Provider()
	for name, info := range prov.Resources {
		assert.Contains(t, schema.Resources, string(info.Tok), "schema is missing %s; run make tfgen", name)
	}
	for name, info := range prov.DataSources {
		assert.Contains(t, schema.Functions, string(info.Tok), "schema is missing %s; run make tfgen", name)
	}
}

// The API should see Pulumi traffic as Pulumi, not as Terraform.
func TestUserAgentNamesPulumi(t *testing.T) {
	assert.Regexp(t, `^pulumi-noderings/\S+ \(\+https://www\.pulumi\.com/registry/packages/noderings/\)$`, userAgent())
}
