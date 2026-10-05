package generator

import (
	"os"
	"path/filepath"
)

// BackendDirs are the folders scaffolded for every project.
var BackendDirs = []string{
	"cmd", "db/migrations", "docs/diagrams", "docs/specs",
	"internal/api/authutil", "internal/api/v1", "internal/common",
	"internal/config", "internal/database/dbsqlc",
	"internal/logging/loggedmodule", "internal/schema",
	"internal/services", "internal/validation", "internal/version",
	"pkg/problemdetail",
}

// CreateDir creates a single folder (and any missing parents) under cwd.
func CreateDir(cwd, dir string) error {
	return os.MkdirAll(filepath.Join(cwd, dir), os.ModePerm)
}
