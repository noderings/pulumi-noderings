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

//go:build python || all

package examples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVmPy orders an hourly VM with a firewall and a DNS record, then deletes
// it. It is billed: it needs NODERINGS_TEST_BILLED=1 and an SSH public key in
// NODERINGS_TEST_SSH_PUBLIC_KEY (or ~/.ssh/id_ed25519.pub).
func TestVmPy(t *testing.T) {
	skipUnlessBilled(t)

	key := os.Getenv("NODERINGS_TEST_SSH_PUBLIC_KEY")
	if key == "" {
		home, err := os.UserHomeDir()
		require.NoError(t, err)
		b, err := os.ReadFile(filepath.Join(home, ".ssh", "id_ed25519.pub"))
		if err != nil {
			t.Skip("set NODERINGS_TEST_SSH_PUBLIC_KEY or create ~/.ssh/id_ed25519.pub")
		}
		key = string(b)
	}

	opts := getPythonBaseOptions(t).With(integration.ProgramTestOptions{
		Dir: filepath.Join(getCwd(t), "vm-py"),
		Config: map[string]string{
			"zone":         uniqueZone(t, "pulumi-py"),
			"sshPublicKey": key,
		},
		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			assert.Regexp(t, `^\d+\.\d+\.\d+\.\d+$`, stack.Outputs["ipv4"])
			assert.NotEmpty(t, stack.Outputs["provider"])
		},
	})
	integration.ProgramTest(t, &opts)
}
