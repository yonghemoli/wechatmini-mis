package main

import (
	"os"
	"strings"
	"testing"
)

func TestSQLFilesCanBeSplit(t *testing.T) {
	files := []string{"../../sql/init-schema.sql", "../../sql/init-seed.sql", "../../sql/cleanup-legacy-tables.sql", "../../sql/mini-business-test-data.sql"}
	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		statements := splitSQLStatements(string(raw))
		if len(statements) == 0 {
			t.Errorf("%s has no SQL statements", path)
		}
		for i, statement := range statements {
			if strings.TrimSpace(statement) == "" {
				t.Errorf("%s statement %d is empty", path, i)
			}
		}
	}
}

func TestCleanupContainsKnownLegacyTables(t *testing.T) {
	raw, err := os.ReadFile("../../sql/cleanup-legacy-tables.sql")
	if err != nil {
		t.Fatal(err)
	}
	content := string(raw)
	for _, table := range []string{"orders", "service_targets", "addresses", "meal_packages", "dishes", "shops", "services", "service_types"} {
		if !strings.Contains(content, "DROP TABLE IF EXISTS "+table) {
			t.Errorf("cleanup does not drop %s", table)
		}
	}
}
