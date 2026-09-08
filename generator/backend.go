package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateBackendFolders(cwd string) {
	dirs := []string{
		"cmd", "db/migrations", "docs/diagrams", "docs/specs",
		"internal/api/authutil", "internal/api/v1", "internal/common",
		"internal/config", "internal/database/dbsqlc",
		"internal/logging/loggedmodule", "internal/schema",
		"internal/services", "internal/validation", "internal/version",
		"pkg/problemdetail",
	}
	for _, d := range dirs {
		os.MkdirAll(filepath.Join(cwd, d), os.ModePerm)
		fmt.Println("[✔] Created:", d)
	}
}
