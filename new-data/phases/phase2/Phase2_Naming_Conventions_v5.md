# Phase 2: Naming Conventions v5
**Agent**: Design Agent  
**Date**: 2025-01-27  
**Objective**: Define professional naming conventions for Research Timeline v5

## Naming Philosophy
Combine the technical precision of data.cleaned.csv with the professional polish of v4 to create clear, descriptive, and actionable task names.

## Task Naming Standards

### General Principles
- **Professional tone**: Academic and research-appropriate language
- **Technical accuracy**: Include specific metrics and requirements
- **Clarity**: Immediately understandable purpose and scope
- **Consistency**: Uniform structure across all tasks
- **Brevity**: Concise but complete descriptions

### Naming Structure
**Format**: `[Action] [Object] [Specifications] [Context]`

### Action Verbs
- **Setup/Configure**: Initial system preparation
- **Build/Develop**: Creation and construction
- **Align/Calibrate**: Precision adjustment
- **Design/Plan**: Planning and design work
- **Implement/Execute**: Active implementation
- **Analyze/Process**: Data analysis and processing
- **Draft/Write**: Document creation
- **Submit/Complete**: Final delivery
- **Maintain/Monitor**: Ongoing maintenance

### Object Descriptions
- **Specific equipment**: "Laser System", "Microscope", "AAV vectors"
- **Documents**: "Proposal", "Manuscript", "Dissertation"
- **Processes**: "Imaging Protocol", "Data Analysis", "Surgical Training"
- **Administrative**: "Committee Paperwork", "Progress Review"

### Technical Specifications
- **Performance metrics**: "≥30 mW output", "≤200 fs pulse duration"
- **Quantities**: "3 pilot mice", "12-page Research Strategy"
- **Timeframes**: "≥2 weeks prior", "nightly backups"
- **Standards**: "BME format requirements", "IACUC protocol"

### Context Information
- **Phase indicators**: "for Aim 1", "post-stroke", "acute phase"
- **Dependencies**: "following BME guidelines", "with advisor approval"
- **Deliverables**: "ready for injections", "for committee review"

## Category-Specific Naming

### PROPOSAL Tasks
**Pattern**: `[Action] [Document/Process] [Requirements] [Timeline]`
- ✅ "Draft timeline v1 for Tuesday review; bring printed and digital copies"
- ✅ "Develop 1-page Specific Aims and detailed outline following BME format requirements"
- ✅ "Write 12-page Research Strategy from outline per BME guidelines"

### EQUIPMENT Tasks
**Pattern**: `[Action] [System] [Performance] [Purpose]`
- ✅ "Align seed laser to reach ≥30 mW output in fiber core (pre-pump)"
- ✅ "Restore amplified output to ≥130 mW (previous benchmark performance level)"
- ✅ "Verify ≤200 fs pulse duration; log specifications"

### RESEARCH Tasks
**Pattern**: `[Action] [Process] [Specifications] [Phase/Context]`
- ✅ "Plan ~3 pilot mice cohort with IACUC protocol confirmation and surgery slot booking"
- ✅ "Design and order AAV-mScarlet (vascular) and jRGECO1b (neuronal) vectors"
- ✅ "Acquire in vivo images for all three mice comparing AAV fluorescence vs traditional dye injection methods"

### PUBLICATION Tasks
**Pattern**: `[Action] [Document] [Content] [Target]`
- ✅ "Write manuscript on AAV-based vascular imaging methodology and pilot results from Aim 1 studies"
- ✅ "Prepare conference talk/poster with results"
- ✅ "Write second research paper covering dual-color imaging platform and initial stroke study findings"

### DISSERTATION Tasks
**Pattern**: `[Action] [Chapter/Section] [Content] [Requirements]`
- ✅ "Write dissertation Introduction chapter including literature review and study rationale for PhD thesis"
- ✅ "Write chapters detailing all three research aims including methods, experiments, data, and findings for dissertation"
- ✅ "Complete PhD dissertation draft compiled and ready for committee review"

### ADMIN Tasks
**Pattern**: `[Action] [Process] [Requirements] [Timeline]`
- ✅ "Submit annual progress report (due early September)"
- ✅ "Complete and submit all required committee forms and develop detailed Program of Work document"
- ✅ "Maintain full-time registration (9+ hours) for Fall 2025 semester"

## Task ID Naming System

### Format
**T[Phase].[Number]**

### Phase 1: Instrumentation & Proposal (T1.1-T1.25)
- **T1.1-T1.10**: Proposal Development
- **T1.11-T1.20**: Equipment Setup
- **T1.21-T1.25**: Committee & Administrative
- **T1.M1**: PhD Proposal Exam Milestone

### Phase 2: Core Research & Analysis (T2.1-T2.35)
- **T2.1-T2.10**: Aim 1 - AAV-based Vascular Imaging
- **T2.11-T2.20**: Aim 2 - Dual-channel Imaging Platform
- **T2.21-T2.30**: Aim 3 - Stroke Study & Analysis
- **T2.31-T2.35**: Data Management & Analysis
- **T2.M1**: Pilot Data Complete Milestone
- **T2.M2**: Stroke Study Complete Milestone

### Phase 3: Publication (T3.1-T3.12)
- **T3.1-T3.6**: Methodology Paper
- **T3.7-T3.12**: Research Paper & Presentations

### Phase 4: Dissertation & Graduation (T4.1-T4.15)
- **T4.1-T4.5**: Dissertation Writing
- **T4.6-T4.10**: Committee Review & Defense
- **T4.11-T4.15**: Final Submission & Graduation
- **T4.M1**: Dissertation Complete Milestone
- **T4.M2**: PhD Defense Complete Milestone

## Objective Writing Standards

### Structure
**Format**: `[Action] [Object] [Technical Requirements] [Success Criteria]`

### Components
1. **Action**: Clear, measurable action verb
2. **Object**: Specific deliverable or outcome
3. **Technical Requirements**: Specific metrics, standards, or specifications
4. **Success Criteria**: How completion is determined

### Examples
- **Good**: "Align seed laser to achieve ≥30 mW output in fiber core for pre-pump optimization"
- **Good**: "Write 12-page Research Strategy section following BME format requirements for committee review"
- **Good**: "Acquire in vivo two-photon images from 3 pilot mice with AAV and dye labeling for methodology comparison"

### Quality Guidelines
- **Specific**: Include exact metrics and requirements
- **Measurable**: Clear success criteria
- **Achievable**: Realistic scope and timeline
- **Relevant**: Aligned with research objectives
- **Time-bound**: Clear completion timeline

## Dependency Format Standards

### Reference Format
**Use new Task IDs, not original data.cleaned.csv letters**

### Examples
- **Original**: "B" → **New**: "T1.2"
- **Original**: "R2,O2" → **New**: "T1.8,T1.16"
- **Original**: "AT,AU,AV,AW" → **New**: "T2.25,T2.26,T2.27,T2.28"

### Format Rules
- **Comma-separated**: No spaces around commas
- **No gaps**: Sequential numbering within phases
- **Consistent**: Always use T[Phase].[Number] format
- **Validated**: All referenced Task IDs must exist

## Quality Assurance

### Naming Checklist
- [ ] Professional tone maintained
- [ ] Technical specifications included
- [ ] Clear and actionable
- [ ] Consistent with category standards
- [ ] Appropriate length (not too verbose)
- [ ] No ambiguous language

### Task ID Checklist
- [ ] Unique identifier
- [ ] Follows T[Phase].[Number] format
- [ ] Sequential numbering within phase
- [ ] Milestones use T[Phase].M format
- [ ] No gaps in numbering

### Objective Checklist
- [ ] Clear action verb
- [ ] Specific deliverable
- [ ] Technical requirements included
- [ ] Success criteria defined
- [ ] Appropriate scope

### Dependency Checklist
- [ ] Uses new Task ID format
- [ ] References valid Task IDs
- [ ] Comma-separated format
- [ ] No spaces around commas
- [ ] Logical dependency relationships

These naming conventions ensure consistency, clarity, and professionalism across all tasks in Research Timeline v5 while preserving the technical accuracy and detail from the original data sources.
