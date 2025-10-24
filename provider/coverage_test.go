package provider

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DavidKrau/terraform-provider-simplemdm/internal/apicatalog"
)

func TestAPICatalogCoverage(t *testing.T) {
	resourceDefinitions := ResourceDefinitionMap()
	dataSourceDefinitions := DataSourceDefinitionMap()

	for _, endpoint := range apicatalog.Catalog {
		if endpoint.ResourceType != "" {
			if _, ok := resourceDefinitions[endpoint.ResourceType]; !ok {
				t.Fatalf("no resource registered for %s (%s)", endpoint.ResourceType, endpoint.Endpoint)
			}
		}

		if endpoint.DataSourceType != "" {
			if _, ok := dataSourceDefinitions[endpoint.DataSourceType]; !ok {
				t.Fatalf("no data source registered for %s (%s)", endpoint.DataSourceType, endpoint.Endpoint)
			}
		}
	}
}

func TestResourceDocumentationCoverage(t *testing.T) {
	for _, definition := range ResourceDefinitions() {
		if definition.DocsPath != "" {
			assertFileExists(t, definition.DocsPath)
		}

		for _, exampleDir := range definition.ExampleDirs {
			assertExampleExists(t, exampleDir)
		}

		if len(definition.TestFiles) == 0 {
			if !definition.TestsOptional {
				t.Fatalf("resource %s is missing acceptance tests", definition.TypeName)
			}
		} else {
			for _, testFile := range definition.TestFiles {
				assertFileExists(t, testFile)
			}
		}
	}
}

func TestDataSourceDocumentationCoverage(t *testing.T) {
	for _, definition := range DataSourceDefinitions() {
		if definition.DocsPath != "" {
			assertFileExists(t, definition.DocsPath)
		}

		for _, exampleDir := range definition.ExampleDirs {
			assertExampleExists(t, exampleDir)
		}

		if len(definition.TestFiles) == 0 {
			t.Fatalf("data source %s is missing acceptance tests", definition.TypeName)
		}

		for _, testFile := range definition.TestFiles {
			assertFileExists(t, testFile)
		}
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	fullPath := filepath.Join("..", path)

	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("expected %s to exist: %v", fullPath, err)
	}
}

func assertExampleExists(t *testing.T, dir string) {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join("..", dir, "*.tf"))
	if err != nil {
		t.Fatalf("failed to glob %s: %v", dir, err)
	}

	if len(matches) == 0 {
		t.Fatalf("expected at least one example in %s", dir)
	}
}
