# Phase 2: Schema Specification v5
**Agent**: Design Agent  
**Date**: 2025-01-27  
**Objective**: Complete schema specification for Research Timeline v5

## Schema Overview
This document defines the complete schema for Research Timeline v5, accommodating all 83 tasks from data.cleaned.csv while maintaining professional polish from v4.

## Complete Column Set

### Core Identification
- **Phase**: Primary phase identifier (1-4)
- **Sub-Phase**: Detailed sub-phase within main phase
- **Category**: Original category from data.cleaned.csv
- **Task ID**: Unique identifier following T1.1-T4.15 format
- **Dependencies**: Reference to prerequisite Task IDs

### Task Content
- **Task**: Professional task name combining v4 polish with technical details
- **Start Date**: Task start date (ISO format: YYYY-MM-DD)
- **End Date**: Task completion date (ISO format: YYYY-MM-DD)
- **Objective**: Clear, specific objective with technical requirements
- **Milestone**: Boolean flag for milestone tasks

### Management Fields
- **Priority**: Task priority level (Critical, High, Medium, Low)
- **Status**: Current task status (Not Started, In Progress, Completed, Blocked, Cancelled)
- **Source**: Origin tracking (data.cleaned.csv, v4, new)

## Category Mapping

### PROPOSAL → PhD Proposal
- **Phase 1**: PhD Proposal
- **Sub-Phases**: 
  - Proposal Development
  - Committee Formation
  - Exam Preparation
  - Defense & Approval

### EQUIPMENT → Equipment Management
- **Phase 1**: Laser System, Microscope Setup
- **Phase 2**: Equipment Maintenance
- **Sub-Phases**:
  - Laser System Build & Alignment
  - Microscope Configuration
  - Equipment Maintenance
  - System Optimization

### RESEARCH → Core Research
- **Phase 2**: Core Research & Analysis
- **Sub-Phases**:
  - Aim 1: AAV-based Vascular Imaging
  - Aim 2: Dual-channel Imaging Platform
  - Aim 3: Stroke Study & Analysis
  - Data Management & Analysis

### PUBLICATION → Publication
- **Phase 3**: Publication
- **Sub-Phases**:
  - Methodology Paper
  - Research Paper
  - Conference Presentations
  - Manuscript Submissions

### DISSERTATION → Dissertation
- **Phase 4**: Dissertation & Graduation
- **Sub-Phases**:
  - Dissertation Writing
  - Committee Review
  - Defense Preparation
  - Final Submission

### ADMIN → Administrative
- **Distributed across all phases**
- **Sub-Phases**:
  - Committee Management
  - Registration & Compliance
  - Progress Reviews
  - Graduation Requirements

## Priority System

### Critical Priority
- PhD Proposal Exam (U)
- PhD Defense (BS)
- Key milestone deliverables
- Regulatory compliance deadlines

### High Priority
- Research tasks with specific technical requirements
- Manuscript submissions
- Equipment setup and calibration
- Data collection milestones

### Medium Priority
- Equipment maintenance
- Data analysis tasks
- Administrative paperwork
- Training and preparation

### Low Priority
- Routine maintenance
- Optional activities
- Background tasks
- Future planning

## Status Values

### Not Started
- Tasks not yet begun
- Default status for new tasks
- Dependencies not yet met

### In Progress
- Tasks currently being worked on
- Active development phase
- Regular status updates required

### Completed
- Tasks successfully finished
- All deliverables met
- Quality standards achieved

### Blocked
- Tasks waiting on external dependencies
- Resource constraints
- Requires resolution before proceeding

### Cancelled
- Tasks no longer needed
- Superseded by other approaches
- Scope changes

## Source Tracking

### data.cleaned.csv
- Original 83 tasks from source file
- Technical specifications preserved
- Dependency relationships maintained

### v4
- Professional naming conventions
- Polish and clarity improvements
- MD alignment maintained

### new
- Tasks created during integration
- Bridge tasks for dependencies
- Quality improvements

## Schema Validation Rules

### Required Fields
- Phase, Sub-Phase, Category, Task ID, Task, Start Date, End Date, Priority, Status, Source

### Optional Fields
- Dependencies, Objective, Milestone

### Data Types
- **Phase**: Integer (1-4)
- **Sub-Phase**: String
- **Category**: String (PROPOSAL, EQUIPMENT, RESEARCH, PUBLICATION, DISSERTATION, ADMIN)
- **Task ID**: String (T1.1-T4.15 format)
- **Dependencies**: String (comma-separated Task IDs)
- **Task**: String
- **Start Date**: Date (YYYY-MM-DD)
- **End Date**: Date (YYYY-MM-DD)
- **Objective**: String
- **Milestone**: Boolean
- **Priority**: String (Critical, High, Medium, Low)
- **Status**: String (Not Started, In Progress, Completed, Blocked, Cancelled)
- **Source**: String (data.cleaned.csv, v4, new)

### Constraints
- Task IDs must be unique
- Dependencies must reference valid Task IDs
- End Date must be >= Start Date
- Status must be valid enum value
- Priority must be valid enum value
- Source must be valid enum value

## Implementation Notes

### CSV Format
- UTF-8 encoding
- Comma-separated values
- Header row included
- Quote strings containing commas
- ISO date format (YYYY-MM-DD)

### Task ID Format
- Phase.Number (e.g., T1.1, T2.15, T4.8)
- Milestones: Phase.M (e.g., T1.M1, T2.M1)
- Sequential numbering within each phase
- No gaps in numbering sequence

### Dependency Format
- Comma-separated list of Task IDs
- No spaces around commas
- Example: "T1.1,T1.3,T2.5"
- Empty string if no dependencies

This schema provides the foundation for creating a comprehensive, professional, and maintainable Research Timeline v5 that accommodates all requirements while maintaining data integrity and usability.
