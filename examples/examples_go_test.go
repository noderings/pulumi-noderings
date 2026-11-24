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

//go:build go || all

package examples

import (
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

// TestDnsGo creates a zone with two records using the Go SDK.
func TestDnsGo(t *testing.T) {
	skipWithoutToken(t)
	opts := getGoBaseOptions(t).With(integration.ProgramTestOptions{
		Dir:    filepath.Join(getCwd(t), "dns-go"),
		Config: map[string]string{"zone": uniqueZone(t, "pulumi-go")},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			assert.NotEmpty(t, stack.Outputs["nameservers"])
		},
	})
	integration.ProgramTest(t, &opts)
}
