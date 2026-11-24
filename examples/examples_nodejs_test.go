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

//go:build nodejs || all

package examples

import (
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

// TestDnsTs creates a zone with A, MX, TXT, CAA and geo records and DNSSEC,
// then changes a TTL in place.
func TestDnsTs(t *testing.T) {
	skipWithoutToken(t)
	opts := getJSBaseOptions(t).With(integration.ProgramTestOptions{
		Dir:    filepath.Join(getCwd(t), "dns-ts"),
		Config: map[string]string{"zone": uniqueZone(t, "pulumi-ts")},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			assert.NotEmpty(t, stack.Outputs["nameservers"])
			assert.Equal(t, "pending_ns", stack.Outputs["zoneState"])
			assert.Regexp(t, `^\d+ 13 2 [0-9a-f]{64}$`, stack.Outputs["dsRecord"])
		},
		EditDirs: []integration.EditDir{{
			Dir:      filepath.Join(getCwd(t), "dns-ts", "step2"),
			Additive: true,
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				assert.Equal(t, []interface{}{"192.0.2.20"}, stack.Outputs["apexRecords"])
			},
		}},
	})
	integration.ProgramTest(t, &opts)
}
