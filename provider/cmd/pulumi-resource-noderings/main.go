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

// Command pulumi-resource-noderings is the NodeRings provider plugin for Pulumi.
package main

import (
	"context"
	_ "embed"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"

	noderings "github.com/noderings/pulumi-noderings/provider"
)

//go:embed schema-embed.json
var pulumiSchema []byte

func main() {
	tfbridge.Main(context.Background(), "noderings", noderings.Provider(),
		tfbridge.ProviderMetadata{PackageSchema: pulumiSchema})
}
