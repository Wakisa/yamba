package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateUIFolders(cwd string) {
	dirs := []string{
		"ui/src/components/form",
		"ui/src/components/ui",
		"ui/src/context",
		"ui/src/data/forms",
		"ui/src/data/queries",
		"ui/src/hooks",
		"ui/src/lib/auth",
		"ui/src/routes/app",
		"ui/src/services",
	}
	for _, d := range dirs {
		os.MkdirAll(filepath.Join(cwd, d), os.ModePerm)
		fmt.Println("[✔] Created UI folder:", d)
	}
}
