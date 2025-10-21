# PhD Dissertation Planner - Technical Documentation

Complete technical documentation for the PhD Dissertation Planner project.

## Table of Contents

1. [Project Structure](#project-structure)
2. [Task Index System](#task-index-system)
3. [CSV File Organization](#csv-file-organization)
4. [Color Coding System](#color-coding-system)
5. [Phase Ordering](#phase-ordering)
6. [Visual Features](#visual-features)
7. [Statistics & Metrics](#statistics--metrics)
8. [LaTeX Implementation](#latex-implementation)
9. [Code Organization](#code-organization)
10. [Design Principles](#design-principles)

---

## Project Structure

### Overview

Clean, consolidated structure with clear separation of concerns:

```
phd-dissertation-planner/
├── input_data/              # All inputs
├── output_data/             # All outputs  
├── internal/                # All code
├── main.go                  # Entry point
├── plannergen               # Binary
└── Makefile                 # Build automation
```

### Detailed Structure

#### Input Directory (input_data/)

All input files in one place:

```
input_data/
├── config.yaml                      # Configuration
├── proposal_and_setup.csv           # Setup & proposal tasks
├── research_and_experiments.csv     # Research aims tasks
├── publications.csv                 # Publication tasks
└── dissertation_and_defense.csv     # Writing & defense tasks
```

#### Output Directory (output_data/)

All generated files organized by type:

```
output_data/
├── latex/          # LaTeX source files (.tex)
├── pdfs/           # Compiled PDFs
├── auxiliary/      # LaTeX auxiliary files (.aux, .log, .fls, .synctex.gz)
└── binaries/       # Binary outputs
```

#### Code Directory (internal/)

All application code organized by functionality:

```
internal/
├── app/            # Application logic
│   ├── cli.go                  # Command-line interface
│   ├── generator.go            # Generation orchestration
│   └── template_funcs.go       # Template functions
│
├── calendar/       # Calendar generation
│   ├── calendar.go             # Calendar structures
│   ├── cell_builder.go         # Cell rendering
│   └── task_stacker.go         # Task layout
│
├── core/           # Core functionality
│   ├── colors.go               # Color generation
│   ├── config.go               # Configuration
│   ├── config_manager.go       # Config management
│   ├── defaults.go             # Default values
│   ├── errors.go               # Error handling
│   ├── logger.go               # Logging
│   ├── reader.go               # CSV reading
│   ├── task.go                 # Task structures
│   └── validation.go           # Data validation
│
└── templates/      # LaTeX templates
    ├── monthly/                # Monthly calendar templates
    │   ├── body.tpl            # Calendar body
    │   ├── calendar.tpl        # Calendar grid
    │   ├── document.tpl        # Document structure
    │   ├── header.tpl          # Page headers
    │   ├── macros.tpl          # LaTeX macros
    │   ├── page.tpl            # Page layout
    │   └── toc.tpl             # Table of contents
    └── rendering.go            # Template rendering
```

#### Root Files

```
├── main.go              # Application entry point
├── plannergen           # Compiled binary (gitignored)
├── Makefile             # Build automation
├── README.md            # User documentation
├── DOCUMENTATION.md     # This file
├── go.mod               # Go module definition
├── go.sum               # Go dependencies
├── .gitignore           # Git ignore rules
└── .gitattributes       # Git attributes
```

### Design Principles

#### 1. Clear Separation
- **Input** - Everything you edit (`input_data/`)
- **Output** - Everything generated (`output_data/`)
- **Code** - Everything that runs (`internal/`)

#### 2. Minimal Hierarchy
- No deep nesting
- Maximum 3 levels deep
- Easy to navigate

#### 3. Logical Grouping
- Related files together
- Clear naming conventions
- Obvious purpose from structure

#### 4. No Redundancy
- Single source of truth
- No duplicate configs
- No scattered files

### Consolidation History

**Before:**
```
├── cmd/planner/main.go
├── configs/config.yaml
├── generated/plannergen
├── input_data/*.csv
├── internal/...
└── pkg/templates/...
```

**After:**
```
├── main.go
├── input_data/config.yaml
├── input_data/*.csv
├── plannergen
└── internal/templates/...
```

**Changes Made:**
1. Moved config: `configs/` → `input_data/`
2. Moved binary: `generated/` → root
3. Moved main: `cmd/planner/` → root
4. Moved templates: `pkg/` → `internal/`
5. Removed directories: `cmd/`, `configs/`, `generated/`, `pkg/`
6. Created `output_data/` for centralized output

---

## Task Index System

### Overview

The Task Index is a comprehensive, hierarchical table of contents that organizes all 108 dissertation tasks across 17 phases and 4 major sections.

### 3-Level Hierarchy

#### Level 1: Sections (Major Groupings)

Four major sections organize the PhD timeline:

```
═══════════════════════════════════════════════════
  SETUP & PROPOSAL
═══════════════════════════════════════════════════
```

**Visual Design:**
- Extra large font (`\LARGE`)
- Bold text
- Horizontal rule separator (0.8pt thick)
- Generous spacing (0.5cm above, 0.3cm below)

**The Four Sections:**

1. **Setup & Proposal** (5 phases, 40 tasks, 10 milestones)
   - Initial preparation and proposal defense
   
2. **Research Aims** (4 phases, 29 tasks, 4 milestones)
   - Core experimental work
   
3. **Publications & Tools** (5 phases, 13 tasks, 5 milestones)
   - Dissemination and software development
   
4. **Dissertation & Defense** (3 phases, 26 tasks, 13 milestones)
   - Writing and completion

#### Level 2: Phases (Colored Headers)

Each phase has a colored header box:

```
┌─────────────────────────────────────────────────┐
│ PhD Proposal          19 tasks, 6 milestones    │
└─────────────────────────────────────────────────┘
```

**Features:**
- Colored background matching calendar legend
- Phase name in large bold text
- Right-aligned statistics (tasks, milestones, completion %)
- Full-width colored bar (98% of line width)
- Consistent padding (2pt top/bottom)

#### Level 3: Tasks (Table Rows)

Individual tasks organized in clean tables:

```
┌───┬──────────────────────────────────┬─────────┬─────────┐
│ # │ Task                             │ Start   │ End     │
├───┼──────────────────────────────────┼─────────┼─────────┤
│ 1 │ Draft Timeline v1                │ Aug 29  │ Sep 02  │
│ 2 │ Develop Specific Aims & Outline  │ Sep 02  │ Sep 08  │
│ 3 │ Submit Outline to Advisor ★      │ Sep 05  │ Sep 09  │
└───┴──────────────────────────────────┴─────────┴─────────┘
```

### Table Layout

#### 4-Column Structure

Each task table has four columns:

1. **# (Number)**
   - Sequential task number within phase
   - Centered alignment
   - Easy reference for discussions
   - Column width: `c` (centered, auto-width)

2. **Task (Name)**
   - Full task name with hyperlink to calendar
   - Flexible width (uses remaining space)
   - Visual indicators:
     - **Bold + ★** for milestones
     - **Gray + ✓** for completed tasks
   - Clickable link to jump to calendar entry
   - Column width: `X` (flexible, ragged right)

3. **Start (Date)**
   - Task start date in short format
   - Format: "Jan 01", "Feb 15", etc.
   - Footnotesize for compact display
   - Column width: `l` (left-aligned, auto-width)

4. **End (Date)**
   - Task end date in same format
   - Shows task duration at a glance
   - Footnotesize for compact display
   - Column width: `l` (left-aligned, auto-width)

#### LaTeX Implementation

```latex
\begin{tabularx}{\linewidth}{
  @{\hspace{0.5em}}c                    % # column (centered)
  @{\hspace{0.8em}}>{\RaggedRight}X     % Task column (flexible)
  @{\hspace{0.8em}}l                    % Start column (left)
  @{\hspace{0.8em}}l                    % End column (left)
  @{\hspace{0.5em}}
}
\hline
\textbf{\#} & \textbf{Task} & \textbf{Start} & \textbf{End} \\
\hline
% Task rows...
\hline
\end{tabularx}
```

**Key Features:**
- `tabularx` package for flexible width
- `X` column type for task names (auto-adjusts to available space)
- `c` for centered number column
- `l` for left-aligned date columns
- `@{\hspace{...}}` for consistent spacing between columns
- `\hline` for horizontal borders (top, header separator, bottom)
- `>{\RaggedRight}` for better text wrapping in task column

---

## CSV File Organization

### File Structure

The planner uses **4 independent CSV files** that are merged at runtime:

#### 1. proposal_and_setup.csv (40 tasks)

**Phases:**
- Project Metadata (10 tasks, 2 milestones, 100% complete)
- PhD Proposal (19 tasks, 6 milestones)
- Committee Management (3 tasks, 1 milestone)
- Microscope Setup (5 tasks)
- Laser System (3 tasks)

**Focus:** Initial setup, proposal development, and equipment preparation

#### 2. research_and_experiments.csv (29 tasks)

**Phases:**
- Aim 1 - AAV-based Vascular Imaging (8 tasks, 1 milestone)
- Aim 2 - Dual-channel Imaging Platform (7 tasks)
- Aim 3 - Stroke Study & Analysis (11 tasks, 1 milestone)
- Data Management & Analysis (3 tasks, 2 milestones)

**Focus:** Core experimental work and research aims

#### 3. publications.csv (13 tasks)

**Phases:**
- SLAVV-T Development (3 tasks, 1 milestone)
- AR Platform Development (4 tasks)
- Research Paper (3 tasks, 1 milestone)
- Methodology Paper (2 tasks, 1 milestone)
- Manuscript Submissions (1 task, 1 milestone)

**Focus:** Publications, tools, and dissemination

#### 4. dissertation_and_defense.csv (26 tasks)

**Phases:**
- Dissertation Writing (11 tasks, 4 milestones)
- Committee Review & Defense (6 tasks, 2 milestones)
- Final Submission & Graduation (9 tasks, 9 milestones)

**Focus:** Writing, defense, and completion

### File Independence

Each CSV file is **completely independent**:
- ✅ No cross-file task dependencies
- ✅ Can be edited separately
- ✅ Merged automatically at build time
- ✅ Alphabetically sorted for consistent ordering
- ✅ Self-contained phases

### Data Sources Display

The Task Index shows merged file information at the top:

```
┌─────────────────┬──────────────────────────────────────────┐
│ Data Sources:   │ 4 CSV file(s) merged                     │
│ Files:          │ dissertation_and_defense.csv,            │
│                 │ proposal_and_setup.csv,                  │
│                 │ publications.csv,                        │
│                 │ research_and_experiments.csv             │
│ Total Tasks:    │ 108 tasks (32 milestones) | 10 completed│
└─────────────────┴──────────────────────────────────────────┘
```

---

## Color Coding System

### Color Generation Algorithm

Colors are generated using the `GenerateCategoryColor()` function in `internal/core/colors.go`:

```go
func GenerateCategoryColor(category string) string {
    // Normalize category name
    normalizedCategory := strings.ToUpper(strings.TrimSpace(category))
    
    // Create hash for consistent color assignment
    hash := 0
    for i, char := range normalizedCategory {
        hash = hash*31 + int(char) + i*7
    }
    
    // Use golden angle (137.5°) for optimal color distribution
    hue := float64(hash%360) * 137.5
    hue = hue - float64(int(hue/360.0)*360)
    
    // Optimized saturation and lightness
    saturation := 0.75  // High saturation for vibrancy
    lightness := 0.65   // Balanced lightness for contrast
    
    // Convert HSL to RGB and return hex
    r, g, b := hslToRgb(hue, saturation, lightness)
    return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
```

**Key Features:**
- **Consistent** - Same phase name always produces same color
- **Distributed** - Golden angle ensures visual distinction
- **Accessible** - High saturation and balanced lightness
- **Automatic** - No manual color assignment needed

### Color Consistency

Phase colors match exactly between Task Index and calendar pages:

| Phase | RGB Color | Hex Color | Usage |
|-------|-----------|-----------|-------|
| Project Metadata | 199,232,98 | #C7E862 | TOC header & calendar legend |
| PhD Proposal | 232,232,98 | #E8E862 | TOC header & calendar legend |
| Committee Management | 210,232,98 | #D2E862 | TOC header & calendar legend |
| Microscope Setup | 204,232,98 | #CCE862 | TOC header & calendar legend |
| Laser System | 215,232,98 | #D7E862 | TOC header & calendar legend |
| Aim 1 - AAV Imaging | 98,232,215 | #62E8D7 | TOC header & calendar legend |
| Aim 2 - Dual-channel | 232,160,98 | #E8A062 | TOC header & calendar legend |
| Aim 3 - Stroke Study | 115,232,98 | #73E862 | TOC header & calendar legend |
| Data Management | 193,232,98 | #C1E862 | TOC header & calendar legend |
| SLAVV-T Development | 98,143,232 | #628FE8 | TOC header & calendar legend |
| AR Platform | 232,98,154 | #E8629A | TOC header & calendar legend |
| Research Paper | 232,98,121 | #E86279 | TOC header & calendar legend |
| Methodology Paper | 232,98,187 | #E862BB | TOC header & calendar legend |
| Manuscript Submissions | 220,232,98 | #DCE862 | TOC header & calendar legend |
| Dissertation Writing | 160,98,232 | #A062E8 | TOC header & calendar legend |
| Committee Review | 126,98,232 | #7E62E8 | TOC header & calendar legend |
| Final Submission | 98,109,232 | #626DE8 | TOC header & calendar legend |

### Benefits

- **Visual Continuity** - Same colors throughout document
- **Easy Identification** - Quickly recognize phases
- **Professional Appearance** - Consistent color scheme
- **Automatic Generation** - No manual configuration
- **Accessibility** - High contrast for readability

---

## Phase Ordering

### Logical Timeline Order

Phases are organized to follow the natural PhD progression (not alphabetically):

#### Section 1: Setup & Proposal
1. **Project Metadata** - Initial project setup and documentation
2. **PhD Proposal** - Proposal development and defense
3. **Committee Management** - Committee coordination
4. **Microscope Setup** - Equipment assembly and configuration
5. **Laser System** - Laser system setup and alignment

#### Section 2: Research Aims
6. **Aim 1 - AAV-based Vascular Imaging** - First research aim
7. **Aim 2 - Dual-channel Imaging Platform** - Second research aim
8. **Aim 3 - Stroke Study & Analysis** - Third research aim
9. **Data Management & Analysis** - Data processing and validation

#### Section 3: Publications & Tools
10. **SLAVV-T Development** - Software tool development
11. **AR Platform Development** - AR visualization platform
12. **Research Paper** - Primary research publication
13. **Methodology Paper** - Methods paper
14. **Manuscript Submissions** - Publication milestones

#### Section 4: Dissertation & Defense
15. **Dissertation Writing** - Chapter writing and drafting
16. **Committee Review & Defense** - Defense preparation and execution
17. **Final Submission & Graduation** - Final revisions and completion

### Implementation

The phase ordering is defined in `internal/app/generator.go`:

```go
sections := []PhaseSection{
    {
        Name: "Setup \\& Proposal",
        Phases: []string{
            "Project Metadata",
            "PhD Proposal",
            "Committee Management",
            "Microscope Setup",
            "Laser System",
        },
    },
    {
        Name: "Research Aims",
        Phases: []string{
            "Aim 1 - AAV-based Vascular Imaging",
            "Aim 2 - Dual-channel Imaging Platform",
            "Aim 3 - Stroke Study & Analysis",
            "Data Management & Analysis",
        },
    },
    // ... more sections
}
```

### Fallback Handling

- Phases not in the predefined order are added alphabetically at the end
- Ensures all phases are included even if new ones are added to CSV files
- Maintains flexibility while providing logical default ordering
- No phases are ever excluded

---

## Visual Features

### Summary Section

Clean 2-column table at the top of the Task Index:

```latex
\noindent\begin{tabularx}{\linewidth}{@{}lX@{}}
\textbf{Data Sources:} & 4 CSV file(s) merged \\[2pt]
\textbf{Files:} & {\footnotesize file1.csv, file2.csv, ...} \\[2pt]
\textbf{Total Tasks:} & 108 tasks (32 milestones) | 10 completed \\
\end{tabularx}
```

**Features:**
- Left column: Bold labels
- Right column: Flexible width for content
- 2pt spacing between rows
- Footnotesize for file names

### Section Headers

```latex
\vspace{0.5cm}
{\LARGE\textbf{Setup \& Proposal}}
\vspace{0.2cm}
\noindent\rule{\linewidth}{0.8pt}
\vspace{0.3cm}
```

**Features:**
- Extra large font (`\LARGE`) for prominence
- Bold text for emphasis
- Horizontal rule (0.8pt thick) for separation
- Generous spacing:
  - 0.5cm above
  - 0.2cm between title and rule
  - 0.3cm below rule

### Phase Headers

```latex
\colorbox[RGB]{232,232,98}{\parbox{0.98\linewidth}{
  \vspace{2pt}
  \textbf{\large PhD Proposal}
  \hfill
  {\small 19 tasks, 6 milestones}
  \vspace{2pt}
}}
```

**Features:**
- Colored background (unique per phase)
- 98% line width (leaves small margins)
- Large bold phase name on left
- Small statistics on right (`\hfill` pushes to right)
- 2pt padding top and bottom
- Full-width colored bar

### Task Tables

**Structure:**
- Horizontal lines for clear boundaries
- Consistent spacing (0.5em and 0.8em)
- Compact footnotesize for dates
- Flexible task name column
- Visual indicators:
  - ★ for milestones
  - ✓ for completed tasks
  - Gray text for completed tasks

**Spacing:**
- 0.25cm between phase header and table
- 0.15cm between table and next phase
- 0.3cm between sections

---

## Statistics & Metrics

### Overall Statistics

- **4 sections** (major groupings)
- **17 phases** (work areas)
- **108 tasks** (individual activities)
- **32 milestones** (key achievements)
- **10 completed** (9.3% overall progress)
- **4 CSV files** merged
- **41 pages** in generated PDF

### Statistics by Section

#### Setup & Proposal
- **5 phases**
- **40 tasks** (37% of total)
- **10 milestones** (31% of total)
- **10 completed** (100% of this section)
- **Status:** Mostly complete

#### Research Aims
- **4 phases**
- **29 tasks** (27% of total)
- **4 milestones** (13% of total)
- **0 completed** (0%)
- **Status:** In progress

#### Publications & Tools
- **5 phases**
- **13 tasks** (12% of total)
- **5 milestones** (16% of total)
- **0 completed** (0%)
- **Status:** Planned

#### Dissertation & Defense
- **3 phases**
- **26 tasks** (24% of total)
- **13 milestones** (41% of total)
- **0 completed** (0%)
- **Status:** Future work

### Phase Statistics Examples

#### Project Metadata (100% complete)
- 10 tasks
- 2 milestones
- All tasks completed
- All shown in gray with ✓
- Spans: Jan 01 - Dec 31, 2025

#### PhD Proposal (in progress)
- 19 tasks
- 6 milestones
- Spans: Aug 29, 2025 - Jan 06, 2026
- Mix of regular tasks and milestones
- Critical path to candidacy

#### Dissertation Writing (future)
- 11 tasks
- 4 milestones
- Spans multiple months in 2027
- Critical path to graduation
- Includes chapter writing and revisions

---

## LaTeX Implementation

### Document Structure

```latex
\documentclass[9pt]{extarticle}
% Packages...
\begin{document}
  % Table of Contents (Task Index)
  \input{monthly.tex}
\end{document}
```

### Template System

Templates are located in `internal/templates/monthly/`:

- **document.tpl** - Document structure and preamble
- **toc.tpl** - Table of contents (Task Index)
- **calendar.tpl** - Calendar grid
- **page.tpl** - Page layout
- **header.tpl** - Page headers
- **body.tpl** - Calendar body
- **macros.tpl** - LaTeX macros

### Key LaTeX Packages

- **tabularx** - Flexible width tables
- **xcolor** - Color support
- **hyperref** - Clickable links
- **geometry** - Page layout
- **tikz** - Graphics and drawing
- **tcolorbox** - Colored boxes for tasks

---

## Code Organization

### Package Structure

#### internal/app/
- **cli.go** - Command-line interface setup
- **generator.go** - Main generation logic and orchestration
- **template_funcs.go** - Custom template functions

#### internal/calendar/
- **calendar.go** - Calendar data structures
- **cell_builder.go** - Day cell rendering
- **task_stacker.go** - Task layout and positioning

#### internal/core/
- **colors.go** - Color generation algorithms
- **config.go** - Configuration structures
- **config_manager.go** - Configuration loading and management
- **defaults.go** - Default configuration values
- **errors.go** - Custom error types
- **logger.go** - Logging system
- **reader.go** - CSV file reading
- **task.go** - Task data structures
- **validation.go** - Data validation

#### internal/templates/
- **monthly/** - Monthly calendar templates
- **rendering.go** - Template rendering engine

### Key Functions

#### Task Index Generation
```go
func createTableOfContentsModule(
    cfg core.Config,
    tasks []core.Task,
    templateName string,
    csvFiles []string,
) core.Module
```

#### CSV Merging
```go
func ReadTasksFromMultipleFiles(
    csvFiles []string,
) ([]core.Task, error)
```

#### Color Generation
```go
func GenerateCategoryColor(
    category string,
) string
```

---

## Design Principles

### 1. Separation of Concerns
- Input, output, and code clearly separated
- Each package has a single responsibility
- Minimal coupling between components

### 2. Data Flow
```
CSV Files → Reader → Tasks → Generator → LaTeX → PDF
```

### 3. Template-Driven
- LaTeX generation uses Go templates
- Separation of logic and presentation
- Easy to customize appearance

### 4. Error Handling
- Rich error context
- User-friendly error messages
- Validation before generation

### 5. Extensibility
- Easy to add new phases
- Easy to add new CSV files
- Easy to customize templates
- Easy to add new features

### 6. Performance
- In-memory processing
- No temporary files
- Efficient CSV parsing
- Minimal memory footprint

---

## Benefits Summary

### For Users
- ✅ Clear organization and navigation
- ✅ Visual hierarchy and color coding
- ✅ Complete information at a glance
- ✅ Professional appearance
- ✅ Easy to customize

### For Developers
- ✅ Clean code structure
- ✅ Logical organization
- ✅ Easy to understand
- ✅ Easy to extend
- ✅ Well-documented

### For Maintenance
- ✅ Minimal complexity
- ✅ Clear responsibilities
- ✅ Easy to find code
- ✅ Easy to debug
- ✅ Easy to test

---

**For more information, see [README.md](README.md) for user documentation.**
