# Legend Fix Summary

## Problem
The legend at the bottom of each monthly calendar page was only showing **one sub-phase** instead of showing **all sub-phases** that appear in that month.

For example, September 2025 has tasks from:
- **Phase 1**: PhD Proposal, Laser System, Committee Management, Microscope Setup (4 sub-phases)
- **Phase 2**: Data Management & Analysis (1 sub-phase)  
- **Phase 4**: Final Submission & Graduation (1 sub-phase)

But the PDF was only showing one of these.

## Root Cause
The issue wasn't with the `GetTaskColorsByPhase()` function - it was working correctly! 

The problem was that **the CSV file path environment variable wasn't being set**, so:
1. No tasks were being loaded from the CSV file
2. `ApplySpanningTasksToMonth()` was never being called
3. Months had empty `day.Tasks` arrays
4. `GetTaskColorsByPhase()` found 0 tasks to process
5. The fallback legend generation was used, which only showed generic/default categories

## Solution
Created `build.sh` script that:
1. Sets the required environment variable: `PLANNER_CSV_FILE="input_data/research_timeline_v5_comprehensive.csv"`
2. Compiles the Go code
3. Generates LaTeX files with task data
4. Compiles the PDF

## How to Build
```bash
./build.sh
```

This will:
- ✅ Load all 107 tasks from CSV
- ✅ Generate monthly calendars with proper task placement
- ✅ Create legends showing ALL sub-phases grouped by phase
- ✅ Produce a 31-page PDF (231KB)

## Verification
The September 2025 page now shows:

```
Phase 1: Proposal & Setup
○ PhD Proposal  ○ Laser System  ○ Committee Management  ○ Microscope Setup

Phase 2: Research & Data Collection
○ Data Management & Analysis

Phase 4: Dissertation
○ Final Submission & Graduation
```

Each month's legend dynamically shows only the sub-phases that have tasks in that specific month.

## Technical Details
- `GetTaskColorsByPhase()` in `src/calendar/calendar.go` collects unique Phase/SubPhase combinations per month
- Returns a `[]PhaseGroup` structure with nested `[]SubPhaseLegendItem` arrays
- Template in `src/shared/templates/monthly/body.tpl` renders this as grouped legend items
- Each sub-phase gets its own color circle and label

## Files Modified
- `build.sh` (new) - Build script with environment variable setup
- `src/calendar/calendar.go` - Cleaned up debug statements
- `TASK_SPANNING_FIX_LOG.md` - Updated with fix details

## Environment Variables
- `PLANNER_CSV_FILE` - **Required** - Path to the CSV file with task data
- If not set, planner generates empty calendars with default categories

## Future Improvements
Consider:
1. Adding `PLANNER_CSV_FILE` to a `.env` file for easier configuration
2. Making CSV path configurable via command-line flag
3. Adding validation to warn if CSV file is not found or tasks fail to load
