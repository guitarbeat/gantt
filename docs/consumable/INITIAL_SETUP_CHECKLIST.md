# Initial Setup Checklist

**Purpose:** One-time setup tasks for new developers/users. Delete this file after completion.

---

## ✅ Prerequisites Installation

- [ ] Install Go 1.21+
  ```bash
  go version  # Verify installation
  ```

- [ ] Install LaTeX distribution
  - Windows: MiKTeX or TeX Live
  - Mac: MacTeX
  - Linux: TeX Live
  ```bash
  pdflatex --version  # Verify installation
  ```

- [ ] Install Git
  ```bash
  git --version  # Verify installation
  ```

- [ ] (Optional) Install Python 3.8+ for preview images
  ```bash
  python --version  # Verify installation
  ```

---

## ✅ Repository Setup

- [ ] Clone the repository
  ```bash
  git clone https://github.com/guitarbeat/gantt.git
  cd gantt
  ```

- [ ] Install Go dependencies
  ```bash
  go mod download
  go mod tidy
  ```

- [ ] Build the application
  ```bash
  go build -o plannergen.exe ./cmd/planner
  # or
  make build
  ```

- [ ] Verify build
  ```bash
  ./plannergen.exe --help
  ```

---

## ✅ Development Tools (Optional)

- [ ] Install pre-commit hooks
  ```bash
  pip install pre-commit
  pre-commit install
  ```

- [ ] Install golangci-lint
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

- [ ] (Optional) Install Python dependencies for preview
  ```bash
  pip install pdf2image Pillow
  ```

- [ ] (Optional) Install poppler for pdf2image
  - Windows: Download from GitHub releases
  - Mac: `brew install poppler`
  - Linux: `sudo apt-get install poppler-utils`

---

## ✅ First Build Test

- [ ] Create or use sample CSV file
  ```bash
  # Use the included sample
  set PLANNER_CSV_FILE=input_data/research_timeline_v5_comprehensive.csv
  ```

- [ ] Generate your first planner
  ```bash
  ./plannergen.exe
  ```

- [ ] Verify output
  ```bash
  # Check generated files
  ls generated/
  # Open the PDF
  start generated/planner.pdf  # Windows
  open generated/planner.pdf   # Mac
  xdg-open generated/planner.pdf  # Linux
  ```

---

## ✅ Run Tests

- [ ] Run all tests
  ```bash
  go test ./...
  # or
  make test
  ```

- [ ] Run with coverage
  ```bash
  go test -v -race -coverprofile=coverage.txt ./...
  # or
  make test-coverage
  ```

---

## ✅ Documentation Review

- [ ] Read [User Guide](../tasks/USER_GUIDE.md)
- [ ] Read [Developer Guide](../tasks/DEVELOPER_GUIDE.md)
- [ ] Bookmark [Troubleshooting](../tasks/TROUBLESHOOTING.md)
- [ ] Review [Documentation Index](../README.md)

---

## ✅ Configuration

- [ ] Choose a preset (academic, compact, presentation)
- [ ] (Optional) Create custom configuration
- [ ] Set environment variables if needed
  ```bash
  # Windows PowerShell
  $env:PLANNER_CSV_FILE = "input_data/my_timeline.csv"
  $env:PLANNER_OUTPUT_DIR = "output"
  
  # Mac/Linux Bash
  export PLANNER_CSV_FILE="input_data/my_timeline.csv"
  export PLANNER_OUTPUT_DIR="output"
  ```

---

## ✅ IDE Setup (Optional)

- [ ] Install VS Code Go extension
- [ ] Configure code formatting on save
- [ ] Set up debugging configuration
- [ ] Install recommended extensions

---

## 🎉 Setup Complete!

Once all items are checked:
1. You're ready to start using the planner
2. Delete this file: `docs/consumable/INITIAL_SETUP_CHECKLIST.md`
3. Start with the [User Guide](../tasks/USER_GUIDE.md)

---

**Setup Date:** _____________  
**Completed By:** _____________  
**Notes:** _____________
