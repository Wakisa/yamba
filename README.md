# Yamba

Yamba is a project generator for Go applications.  
It scaffolds backend and optional UI folders, sets up Cobra for CLI commands, and provides a Bubble Tea interface for interactive project initialization.

## ✨ Features
- Interactive CLI powered by [Cobra](https://github.com/spf13/cobra) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- Confirmation prompt to avoid generating files in the wrong directory.
- Choice between **backend only** or **backend + UI** scaffolding.
- Automatic creation of common project folders (`cmd`, `internal`, `pkg`, `docs`, etc.).
- Built-in `version` command (`yamba version`, `yamba --version`, `yamba -v`).

## 🚀 Installation
Clone the repository and build the binary:

```bash
git clone https://github.com/wakisa/yamba.git
cd yamba
go build -o yamba .

Optionally move it into your PATH (e.g., C:\\Program Files\\yamba on Windows):
`move yamba.exe "C:\Program Files\yamba\yamba.exe"`
