# Phase 2: Task Creation Template v5
**Agent**: Design Agent  
**Date**: 2025-01-27  
**Objective**: Template for consistent task creation in Research Timeline v5

## Template Overview
This template provides a standardized approach for creating tasks in Research Timeline v5, ensuring consistency, completeness, and quality across all task definitions.

## Task Creation Process

### Step 1: Task Identification
1. **Source Analysis**: Determine if task comes from data.cleaned.csv, v4, or is new
2. **Category Assignment**: Map to appropriate category (PROPOSAL, EQUIPMENT, RESEARCH, PUBLICATION, DISSERTATION, ADMIN)
3. **Phase Assignment**: Determine primary phase (1-4) based on timeline and dependencies
4. **Sub-Phase Assignment**: Assign to specific sub-phase within main phase

### Step 2: Task ID Assignment
1. **Phase Selection**: Choose appropriate phase (1-4)
2. **Number Assignment**: Assign next available number in sequence
3. **Milestone Check**: Determine if task is a milestone (use .M format)
4. **Uniqueness Verification**: Ensure ID is unique across entire timeline

### Step 3: Content Creation
1. **Task Name**: Create professional, descriptive name
2. **Objective**: Write clear, specific objective with technical requirements
3. **Dependencies**: Identify and format prerequisite tasks
4. **Dates**: Set appropriate start and end dates
5. **Priority**: Assign priority level based on importance and urgency

## Task Template

### Basic Task Template
```csv
Phase,Sub-Phase,Category,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Priority,Status,Source
[Phase],[Sub-Phase],[Category],[Task ID],[Dependencies],[Task Name],[Start Date],[End Date],[Objective],[Milestone],[Priority],[Status],[Source]
```

### Field Definitions

#### Phase
- **Type**: Integer (1-4)
- **Values**: 1 (Instrumentation & Proposal), 2 (Core Research), 3 (Publication), 4 (Dissertation)
- **Assignment**: Based on timeline position and research phase

#### Sub-Phase
- **Type**: String
- **Values**: Specific sub-phase within main phase
- **Examples**: "PhD Proposal", "Laser System", "Aim 1", "Methodology Paper", "Dissertation Writing"

#### Category
- **Type**: String
- **Values**: PROPOSAL, EQUIPMENT, RESEARCH, PUBLICATION, DISSERTATION, ADMIN
- **Assignment**: Based on original data.cleaned.csv category

#### Task ID
- **Type**: String
- **Format**: T[Phase].[Number] or T[Phase].M[Number] for milestones
- **Examples**: T1.1, T2.15, T1.M1, T4.M2

#### Dependencies
- **Type**: String
- **Format**: Comma-separated Task IDs
- **Examples**: "T1.1", "T1.2,T1.3", "T2.5,T2.8,T3.1"
- **Empty**: "" if no dependencies

#### Task
- **Type**: String
- **Format**: Professional, descriptive name
- **Length**: < 200 characters
- **Style**: Academic, research-appropriate tone

#### Start Date
- **Type**: Date
- **Format**: YYYY-MM-DD
- **Validation**: Must be valid date, aligned with timeline

#### End Date
- **Type**: Date
- **Format**: YYYY-MM-DD
- **Validation**: Must be >= Start Date

#### Objective
- **Type**: String
- **Format**: Clear, specific objective with technical requirements
- **Style**: Professional, measurable, achievable

#### Milestone
- **Type**: Boolean
- **Values**: true, false
- **Assignment**: true for key milestone tasks

#### Priority
- **Type**: String
- **Values**: Critical, High, Medium, Low
- **Assignment**: Based on importance and urgency

#### Status
- **Type**: String
- **Values**: Not Started, In Progress, Completed, Blocked, Cancelled
- **Default**: Not Started

#### Source
- **Type**: String
- **Values**: data.cleaned.csv, v4, new
- **Assignment**: Based on task origin

## Category-Specific Templates

### PROPOSAL Tasks
```csv
Phase,Sub-Phase,Category,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Priority,Status,Source
1,PhD Proposal,PROPOSAL,T1.X,[Dependencies],[Action] [Document/Process] [Requirements] [Timeline],[Start Date],[End Date],[Action] [Object] [Technical Requirements] [Success Criteria],false,[Priority],Not Started,data.cleaned.csv
```

**Example**:
```csv
1,PhD Proposal,PROPOSAL,T1.2,,Develop 1-page Specific Aims and detailed outline following BME format requirements for PhD proposal,2025-09-02,2025-09-08,Develop comprehensive proposal outline following BME format requirements for committee review,false,High,Not Started,data.cleaned.csv
```

### EQUIPMENT Tasks
```csv
Phase,Sub-Phase,Category,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Priority,Status,Source
1,Laser System,EQUIPMENT,T1.X,[Dependencies],[Action] [System] [Performance] [Purpose],[Start Date],[End Date],[Action] [Object] [Technical Requirements] [Success Criteria],false,[Priority],Not Started,data.cleaned.csv
```

**Example**:
```csv
1,Laser System,EQUIPMENT,T1.8,,Align seed laser to reach ≥30 mW output in fiber core (pre-pump),2025-09-02,2025-09-06,Align seed laser to achieve ≥30 mW output in fiber core for pre-pump optimization,false,High,Not Started,data.cleaned.csv
```

### RESEARCH Tasks
```csv
Phase,Sub-Phase,Category,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Priority,Status,Source
2,Core Research,RESEARCH,T2.X,[Dependencies],[Action] [Process] [Specifications] [Phase/Context],[Start Date],[End Date],[Action] [Object] [Technical Requirements] [Success Criteria],false,[Priority],Not Started,data.cleaned.csv
```

**Example**:
```csv
2,Aim 1,RESEARCH,T2.1,T1.8,Plan ~3 pilot mice cohort with IACUC protocol confirmation and surgery slot booking,2025-10-14,2025-10-18,Plan pilot mouse cohort with IACUC protocol confirmation and surgery scheduling for Aim 1 studies,false,High,Not Started,data.cleaned.csv
```

### PUBLICATION Tasks
```csv
Phase,Sub-Phase,Category,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Priority,Status,Source
3,Publication,PUBLICATION,T3.X,[Dependencies],[Action] [Document] [Content] [Target],[Start Date],[End Date],[Action] [Object] [Technical Requirements] [Success Criteria],false,[Priority],Not Started,data.cleaned.csv
```

**Example**:
```csv
3,Methodology Paper,PUBLICATION,T3.1,T2.15,Write manuscript on AAV-based vascular imaging methodology and pilot results from Aim 1 studies,2026-04-19,2026-07-15,Write comprehensive methodology manuscript covering AAV-based vascular imaging approach and pilot study results,false,High,Not Started,data.cleaned.csv
```

### DISSERTATION Tasks
```csv
Phase,Sub-Phase,Category,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Priority,Status,Source
4,Dissertation Writing,DISSERTATION,T4.X,[Dependencies],[Action] [Chapter/Section] [Content] [Requirements],[Start Date],[End Date],[Action] [Object] [Technical Requirements] [Success Criteria],false,[Priority],Not Started,data.cleaned.csv
```

**Example**:
```csv
4,Dissertation Writing,DISSERTATION,T4.1,,Write dissertation Introduction chapter including literature review and study rationale for PhD thesis,2026-12-19,2027-01-31,Write comprehensive Introduction chapter with literature review and study rationale for PhD dissertation,false,Critical,Not Started,data.cleaned.csv
```

### ADMIN Tasks
```csv
Phase,Sub-Phase,Category,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Priority,Status,Source
[Phase],Administrative,ADMIN,T[Phase].X,[Dependencies],[Action] [Process] [Requirements] [Timeline],[Start Date],[End Date],[Action] [Object] [Technical Requirements] [Success Criteria],false,[Priority],Not Started,data.cleaned.csv
```

**Example**:
```csv
1,Committee Management,ADMIN,T1.16,,Complete and submit all required committee forms and develop detailed Program of Work document,2025-10-08,2025-10-18,Complete all required committee paperwork and develop comprehensive Program of Work document,false,High,Not Started,data.cleaned.csv
```

## Milestone Template

### Milestone Task Template
```csv
Phase,Sub-Phase,Category,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Priority,Status,Source
[Phase],[Sub-Phase],[Category],T[Phase].M[Number],[Dependencies],[Milestone Name],[Start Date],[End Date],[Milestone Description],true,Critical,Not Started,new
```

**Example**:
```csv
1,PhD Proposal,PROPOSAL,T1.M1,T1.8,T1.15,PhD Proposal Exam,2025-12-19,2025-12-22,Successfully defend PhD proposal in oral examination,true,Critical,Not Started,new
```

## Quality Checklist

### Pre-Creation Checklist
- [ ] Source identified and verified
- [ ] Category assignment confirmed
- [ ] Phase assignment validated
- [ ] Sub-phase assignment appropriate
- [ ] Dependencies identified and validated

### Content Creation Checklist
- [ ] Task name follows naming conventions
- [ ] Objective is clear and specific
- [ ] Technical requirements included
- [ ] Success criteria defined
- [ ] Priority level appropriate
- [ ] Dates aligned with timeline

### Post-Creation Checklist
- [ ] Task ID is unique
- [ ] All required fields populated
- [ ] Dependencies reference valid Task IDs
- [ ] Dates are valid and logical
- [ ] Content follows quality standards
- [ ] Source tracking accurate

## Validation Rules

### Data Validation
- **Task ID**: Must be unique, follow T[Phase].[Number] format
- **Dependencies**: Must reference existing Task IDs
- **Dates**: End date >= Start date, valid ISO format
- **Priority**: Must be valid enum value
- **Status**: Must be valid enum value
- **Source**: Must be valid enum value

### Content Validation
- **Task name**: Professional tone, technical accuracy, appropriate length
- **Objective**: Clear, specific, measurable, achievable
- **Category**: Appropriate for task content
- **Phase**: Logical placement in timeline
- **Dependencies**: Accurate and complete

### Integration Validation
- **Timeline alignment**: Dates align with overall project timeline
- **Dependency integrity**: All dependencies are logical and necessary
- **Phase consistency**: Task fits logically within assigned phase
- **Category consistency**: Task content matches assigned category

This template ensures consistent, high-quality task creation across all phases of Research Timeline v5 while maintaining the technical accuracy and professional polish required for the project.
