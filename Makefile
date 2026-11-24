# Build the NodeRings Pulumi provider and its SDKs.
#
#   make                  schema, provider binary and all SDKs
#   make tfgen            regenerate the Pulumi schema from the Terraform provider
#   make generate_sdks    regenerate the committed SDK sources from the schema
#   make build_sdks       build every SDK for publishing (needs Node, Python, .NET)
#   make install_plugin   install the provider binary into the local Pulumi plugin cache
#   make test_provider    unit tests, no account needed
#   make test             example programs against the live API (needs NODERINGS_TOKEN)
#
# VERSION defaults to the latest v* tag; set it when building a release.

PACK         := noderings
PROJECT      := github.com/noderings/pulumi-noderings
PROVIDER     := pulumi-resource-$(PACK)
TFGEN        := pulumi-tfgen-$(PACK)
WORKING_DIR  := $(shell pwd)
SCHEMA_DIR   := provider/cmd/$(PROVIDER)
SCHEMA       := $(SCHEMA_DIR)/schema.json
LANGUAGES    := nodejs python go dotnet

VERSION      ?= $(shell git describe --tags --match 'v*' --abbrev=0 2>/dev/null | sed 's/^v//')
ifeq ($(VERSION),)
VERSION      := 0.0.1-dev
endif
LDFLAGS      := -X $(PROJECT)/provider/pkg/version.Version=$(VERSION)

# Keep test runs fast; long-running steps set their own timeouts.
TESTPARALLELISM ?= 4

.PHONY: all build tfgen schema provider generate_sdks build_sdks \
	build_nodejs build_python build_go build_dotnet \
	install_plugin install_nodejs_sdk test test_provider lint tidy clean help

all: build

build: provider build_sdks

help:
	@sed -n '1,13p' Makefile | sed 's/^# \{0,1\}//'

# --- schema --------------------------------------------------------------

tfgen: schema

schema:
	cd provider && go build -ldflags "$(LDFLAGS)" -o $(WORKING_DIR)/bin/$(TFGEN) ./cmd/$(TFGEN)
	PULUMI_SKIP_UPDATE_CHECK=1 bin/$(TFGEN) schema --out $(SCHEMA_DIR)
	python3 -c "import json; json.dump(json.load(open('$(SCHEMA)')), open('$(SCHEMA_DIR)/schema-embed.json', 'w'), separators=(',', ':'))"

$(SCHEMA_DIR)/schema-embed.json: $(SCHEMA)
	python3 -c "import json; json.dump(json.load(open('$(SCHEMA)')), open('$(SCHEMA_DIR)/schema-embed.json', 'w'), separators=(',', ':'))"

# --- provider ------------------------------------------------------------

provider: $(SCHEMA_DIR)/schema-embed.json
	cd provider && go build -ldflags "$(LDFLAGS)" -o $(WORKING_DIR)/bin/$(PROVIDER) ./cmd/$(PROVIDER)

install_plugin: provider
	pulumi plugin rm resource $(PACK) $(VERSION) --yes 2>/dev/null || true
	pulumi plugin install resource $(PACK) $(VERSION) --file bin/$(PROVIDER)

# --- SDKs ----------------------------------------------------------------

# The committed SDK sources carry no version; build_* stamps it in.
generate_sdks:
	for lang in $(LANGUAGES); do \
		rm -rf sdk/$$lang; \
		pulumi package gen-sdk $(SCHEMA) --language $$lang --out sdk/ || exit 1; \
	done
	cd sdk && go mod tidy

build_sdks: build_nodejs build_python build_go build_dotnet

build_nodejs:
	cd sdk/nodejs && \
		yarn install && \
		yarn run tsc && \
		cp package.json README.md ../../LICENSE bin/ && \
		sed -i.bak 's/$${VERSION}/$(VERSION)/g' bin/package.json && rm bin/package.json.bak

build_python:
	cd sdk/python && \
		rm -rf bin && mkdir bin && \
		cp -r pulumi_$(PACK) pyproject.toml bin/ && \
		cp ../../README.md bin/README.md && \
		sed -i.bak 's/^  version = .*/  version = "$(VERSION)"/' bin/pyproject.toml && rm bin/pyproject.toml.bak && \
		python3 -c "import json; p='bin/pulumi_$(PACK)/pulumi-plugin.json'; d=json.load(open(p)); d['version']='$(VERSION)'; json.dump(d, open(p, 'w'), indent=2)" && \
		cd bin && python3 -m build .

build_go:
	cd sdk && go build ./...

build_dotnet:
	cd sdk/dotnet && \
		echo "$(VERSION)" > version.txt && \
		dotnet build -c Release /p:Version=$(VERSION)

install_nodejs_sdk: build_nodejs
	cd sdk/nodejs/bin && yarn link

# --- tests ---------------------------------------------------------------

test_provider:
	cd provider && go test -v -count=1 -parallel $(TESTPARALLELISM) ./...

# Free tests need NODERINGS_TOKEN. Add NODERINGS_TEST_BILLED=1 to also order a VM.
test: provider build_python install_nodejs_sdk
	cd examples && go test -v -count=1 -tags=all -timeout 2h -parallel $(TESTPARALLELISM) ./...

lint:
	cd provider && go vet ./...
	cd examples && go vet -tags=all ./...

tidy:
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd examples && go mod tidy

clean:
	rm -rf bin sdk/nodejs/bin sdk/nodejs/node_modules sdk/python/bin sdk/dotnet/bin sdk/dotnet/obj
