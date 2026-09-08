package steps

import (
	"github.com/wakisa/yamba/generator"
)

// RunStep executes a single scaffold step.
// Accepts the step name, current working directory, project name, and whether to scaffold UI.
func RunStep(step, cwd, projectName string, withUI bool) {
	switch step {
	case "Create backend folders":
		generator.CreateBackendFolders(cwd)
		if withUI {
			generator.CreateUIFolders(cwd)
		}

	case "Generate Cobra root.go":
		generator.GenerateCobraRoot(cwd, projectName)

	case "Generate stub files":
		generator.GenerateStubFiles(cwd)
	}
}
