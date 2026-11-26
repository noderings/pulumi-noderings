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

// Package noderings maps the NodeRings Terraform provider into a Pulumi provider.
package noderings

import (
	"bytes"
	"path"
	"regexp"

	// Embed bridge-metadata.json in the provider.
	_ "embed"

	pfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"

	upstream "github.com/noderings/terraform-provider-noderings/shim"

	"github.com/noderings/pulumi-noderings/provider/pkg/version"
)

const (
	// mainPkg is the package name in every language registry.
	mainPkg = "noderings"
	// mainMod puts every resource in the top-level module, so the names read
	// noderings.VirtualMachine, noderings.DnsRecord and so on.
	mainMod = "index"
)

//go:embed cmd/pulumi-resource-noderings/bridge-metadata.json
var metadata []byte

// Provider returns the Pulumi mapping of the NodeRings Terraform provider.
func Provider() tfbridge.ProviderInfo {
	prov := tfbridge.ProviderInfo{
		P:       pfbridge.ShimProvider(upstream.NewProviderWithUserAgent(version.Version, userAgent())),
		Name:    "noderings",
		Version: version.Version,

		DisplayName: "NodeRings",
		Publisher:   "NodeRings",
		LogoURL:     "https://raw.githubusercontent.com/noderings/pulumi-noderings/main/docs/logo.svg",
		// Pulumi downloads the provider plugin from this repository's GitHub releases.
		PluginDownloadURL: "github://api.github.com/noderings",
		Description:       "A Pulumi package for the NodeRings infrastructure marketplace: order and manage compute, networking, DNS, access control and more.",
		Keywords:          []string{"noderings", "vps", "dns", "category/cloud"},
		License:           "Apache-2.0",
		Homepage:          "https://noderings.com",
		Repository:        "https://github.com/noderings/pulumi-noderings",
		// The upstream provider's GitHub organization, used for doc links.
		GitHubOrg:    "noderings",
		MetadataInfo: tfbridge.NewProviderMetadata(metadata),
		DocRules:     &tfbridge.DocRuleInfo{EditRules: docEditRules},

		Config: map[string]*tfbridge.SchemaInfo{
			"token": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"NODERINGS_TOKEN", "NODERINGS_API_TOKEN", "NR_API_TOKEN"}},
				Secret:  tfbridge.True(),
			},
			"endpoint": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"NODERINGS_API_URL"}},
			},
			"organization_id": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"NODERINGS_ORGANIZATION_ID"}},
			},
		},

		JavaScript: &tfbridge.JavaScriptInfo{
			PackageName:          "@noderings/pulumi",
			RespectSchemaVersion: true,
		},
		Python: &tfbridge.PythonInfo{
			PackageName:          "pulumi_noderings",
			RespectSchemaVersion: true,
			PyProject:            struct{ Enabled bool }{true},
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: path.Join(
				"github.com/noderings/pulumi-noderings/sdk/",
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			GenerateResourceContainerTypes: true,
			GenerateExtraInputTypes:        true,
			RespectSchemaVersion:           true,
		},
		CSharp: &tfbridge.CSharpInfo{
			RootNamespace:        "NodeRings",
			Namespaces:           map[string]string{mainPkg: "Pulumi"},
			RespectSchemaVersion: true,
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
	}

	// Map every upstream resource and data source into the top-level module:
	// noderings_virtual_machine becomes noderings:index/virtualMachine:VirtualMachine,
	// data.noderings_marketplace_offer becomes noderings:index/getMarketplaceOffer:getMarketplaceOffer.
	prov.MustComputeTokens(tokens.SingleModule("noderings_", mainMod, tokens.MakeStandard(mainPkg)))
	prov.MustApplyAutoAliases()

	return prov
}

// terraformPhrases matches upstream wording that describes what Terraform does,
// so the Pulumi docs can say Pulumi instead. Links to the Terraform provider and
// HCL code blocks are left alone.
var terraformPhrases = regexp.MustCompile(
	`\bTerraform (refuses|starts|stops|unlocks|waits|does not|then|creates|deletes|removes|keeps|reads|writes|cannot)\b` +
		`|\b(outside|by|in|from|through) Terraform\b` +
		`|\bthe Terraform state\b|\bTerraform state\b`)

func docEditRules(defaults []tfbridge.DocsEdit) []tfbridge.DocsEdit {
	return append(defaults, tfbridge.DocsEdit{
		Path: "*",
		Edit: func(_ string, content []byte) ([]byte, error) {
			return terraformPhrases.ReplaceAllFunc(content, func(m []byte) []byte {
				return bytes.ReplaceAll(m, []byte("Terraform"), []byte("Pulumi"))
			}), nil
		},
	})
}

// userAgent identifies this provider to the NodeRings API, so traffic from
// Pulumi is not counted as Terraform.
func userAgent() string {
	return "pulumi-noderings/" + version.Version + " (+https://www.pulumi.com/registry/packages/noderings/)"
}
