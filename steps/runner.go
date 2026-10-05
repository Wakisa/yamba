package steps

import (
	"github.com/wakisa/yamba/generator"
)

// Step names, shown as the checklist on the progress screen.
const (
	StepBackendFolders = "Create backend folders"
	StepUIFolders      = "Create UI folders"
	StepCobraRoot      = "Generate Cobra root.go"
	StepStubFiles      = "Generate stub files"
)

// Task is one small unit of scaffolding work (a single folder or file),
// so the UI can report progress item by item.
type Task struct {
	Step  string // which checklist step this task belongs to
	Label string // what was created, shown in the activity log
	Run   func() error
}

// Names returns the ordered, de-duplicated step names for a list of tasks.
func Names(tasks []Task) []string {
	var names []string
	for _, t := range tasks {
		if len(names) == 0 || names[len(names)-1] != t.Step {
			names = append(names, t.Step)
		}
	}
	return names
}

// BuildTasks returns every scaffold task in the order it should run.
func BuildTasks(cwd, projectName string, withUI bool) []Task {
	var tasks []Task

	for _, d := range generator.BackendDirs {
		tasks = append(tasks, Task{
			Step:  StepBackendFolders,
			Label: "Created folder " + d,
			Run:   func() error { return generator.CreateDir(cwd, d) },
		})
	}

	if withUI {
		for _, d := range generator.UIDirs {
			tasks = append(tasks, Task{
				Step:  StepUIFolders,
				Label: "Created UI folder " + d,
				Run:   func() error { return generator.CreateDir(cwd, d) },
			})
		}
	}

	tasks = append(tasks, Task{
		Step:  StepCobraRoot,
		Label: "Added cmd/root.go",
		Run:   func() error { return generator.GenerateCobraRoot(cwd, projectName) },
	})

	for _, f := range generator.StubFiles {
		tasks = append(tasks, Task{
			Step:  StepStubFiles,
			Label: "Wrote " + f.Path,
			Run:   func() error { return generator.WriteStub(cwd, f) },
		})
	}

	return tasks
}
