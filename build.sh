#!/bin/bash

# Build script for terraform-provider-simplemdm
# This script builds the provider binary, generates documentation, and runs validations

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Terraform Provider SimpleMDM Build Script ===${NC}"
echo ""

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

if ! command -v terraform &> /dev/null; then
    echo -e "${YELLOW}Warning: Terraform is not installed. Skipping terraform fmt...${NC}"
    SKIP_TF_FMT=1
else
    SKIP_TF_FMT=0
fi

echo -e "${GREEN}✓ Prerequisites checked${NC}"
echo ""

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -f terraform-provider-simplemdm
rm -rf ~/.terraform.d/plugins/github.com/DavidKrau/simplemdm/dev/
echo -e "${GREEN}✓ Cleaned${NC}"
echo ""

# Download dependencies
echo -e "${YELLOW}Downloading Go dependencies...${NC}"
go mod download
echo -e "${GREEN}✓ Dependencies downloaded${NC}"
echo ""

# Run go fmt
echo -e "${YELLOW}Running go fmt...${NC}"
go fmt ./...
echo -e "${GREEN}✓ Code formatted${NC}"
echo ""

# Run go vet
echo -e "${YELLOW}Running go vet...${NC}"
go vet ./...
echo -e "${GREEN}✓ Code vetted${NC}"
echo ""

# Build the provider first
echo -e "${YELLOW}Building provider binary...${NC}"
go build -o terraform-provider-simplemdm
echo -e "${GREEN}✓ Provider binary built: terraform-provider-simplemdm${NC}"
echo ""

# Install provider locally for documentation generation
echo -e "${YELLOW}Installing provider locally for documentation generation...${NC}"
PROVIDER_DIR="${HOME}/.terraform.d/plugins/github.com/DavidKrau/simplemdm/dev/$(go env GOOS)_$(go env GOARCH)"
mkdir -p "${PROVIDER_DIR}"
cp terraform-provider-simplemdm "${PROVIDER_DIR}/terraform-provider-simplemdm_vdev"
echo -e "${GREEN}✓ Provider installed locally${NC}"
echo ""

# Format terraform examples if terraform is available
if [ $SKIP_TF_FMT -eq 0 ]; then
    echo -e "${YELLOW}Formatting Terraform examples...${NC}"
    terraform fmt -recursive ./examples/ || echo -e "${YELLOW}Warning: Some files could not be formatted${NC}"
    echo -e "${GREEN}✓ Examples formatted${NC}"
    echo ""
fi

# Generate documentation using the locally installed provider (optional)
if [ "${SKIP_DOCS:-0}" -eq 0 ]; then
    echo -e "${YELLOW}Generating documentation...${NC}"
    export TF_CLI_CONFIG_FILE="${PWD}/.terraformrc"
    if go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name simplemdm 2>/dev/null; then
        echo -e "${GREEN}✓ Documentation generated${NC}"
    else
        echo -e "${YELLOW}Warning: Documentation generation failed. Using existing documentation.${NC}"
        echo -e "${YELLOW}This is expected if the provider isn't published to the registry yet.${NC}"
    fi
    unset TF_CLI_CONFIG_FILE
else
    echo -e "${YELLOW}Skipping documentation generation (SKIP_DOCS=1)${NC}"
fi
echo ""

# Verify the binary
if [ -f "terraform-provider-simplemdm" ]; then
    echo -e "${GREEN}✓ Binary verification successful${NC}"
    ls -lh terraform-provider-simplemdm
else
    echo -e "${RED}Error: Binary not found after build${NC}"
    exit 1
fi
echo ""

# Verify documentation was generated
echo -e "${YELLOW}Verifying documentation...${NC}"
if [ -d "docs/resources" ] && [ -d "docs/data-sources" ]; then
    RESOURCE_COUNT=$(find docs/resources -name "*.md" | wc -l)
    DATASOURCE_COUNT=$(find docs/data-sources -name "*.md" | wc -l)
    echo -e "${GREEN}✓ Documentation verified:${NC}"
    echo "  - Resources: $RESOURCE_COUNT files"
    echo "  - Data Sources: $DATASOURCE_COUNT files"
else
    echo -e "${RED}Error: Documentation directories not found${NC}"
    exit 1
fi
echo ""

echo -e "${GREEN}=== Build Complete ===${NC}"
echo -e "Binary: ${GREEN}terraform-provider-simplemdm${NC}"
echo -e "Documentation: ${GREEN}docs/${NC}"