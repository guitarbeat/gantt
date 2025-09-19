# 📚 Examples - PhD Dissertation Planner

Real-world examples and use cases for the PhD Dissertation Planner.

## 📋 Table of Contents

1. [Basic Examples](#basic-examples)
2. [Academic Examples](#academic-examples)
3. [Research Examples](#research-examples)
4. [Project Examples](#project-examples)
5. [Configuration Examples](#configuration-examples)

## 🚀 Basic Examples

### Simple Task List
**File**: `simple_tasks.csv`
```csv
Task Name,Start Date,Due Date,Category,Description,Priority
Literature Review,2024-01-15,2024-03-15,RESEARCH,Review relevant papers,3
Data Collection,2024-02-01,2024-04-01,RESEARCH,Collect experimental data,2
Analysis,2024-03-01,2024-05-01,ANALYSIS,Analyze collected data,1
Writing,2024-04-01,2024-06-01,WRITING,Write dissertation chapters,1
```

**Generate PDF**:
```bash
./scripts/simple.sh simple_tasks.csv simple_example
```

### Task with Dependencies
**File**: `dependent_tasks.csv`
```csv
Task Name,Start Date,Due Date,Category,Dependencies,Priority
Literature Review,2024-01-15,2024-03-15,RESEARCH,,3
Data Collection,2024-02-01,2024-04-01,RESEARCH,Literature Review,2
Analysis,2024-03-01,2024-05-01,ANALYSIS,Data Collection,1
Writing,2024-04-01,2024-06-01,WRITING,Analysis,1
```

## 🎓 Academic Examples

### PhD Dissertation Timeline
**File**: `phd_dissertation.csv`
```csv
Task Name,Start Date,Due Date,Category,Description,Priority,Status
Proposal Defense,2024-01-15,2024-01-30,MILESTONE,Defend research proposal,1,Planned
Literature Review,2024-02-01,2024-04-30,RESEARCH,Comprehensive literature review,2,Planned
Methodology Development,2024-03-01,2024-05-31,RESEARCH,Develop research methodology,2,Planned
Data Collection Phase 1,2024-04-01,2024-07-31,RESEARCH,Initial data collection,2,Planned
Data Collection Phase 2,2024-06-01,2024-09-30,RESEARCH,Extended data collection,2,Planned
Data Analysis,2024-08-01,2024-11-30,ANALYSIS,Statistical analysis of data,1,Planned
Chapter 1-3 Writing,2024-09-01,2024-12-31,WRITING,Introduction and methodology,1,Planned
Chapter 4-5 Writing,2024-10-01,2025-02-28,WRITING,Results and discussion,1,Planned
Chapter 6 Writing,2024-11-01,2025-03-31,WRITING,Conclusion and future work,1,Planned
Dissertation Defense,2025-04-15,2025-04-30,MILESTONE,Final dissertation defense,1,Planned
```

### Course Schedule
**File**: `course_schedule.csv`
```csv
Task Name,Start Date,Due Date,Category,Description,Priority,Status
Advanced Statistics,2024-01-15,2024-05-15,COURSE,Statistical methods course,2,Planned
Research Methods,2024-01-15,2024-05-15,COURSE,Research methodology course,2,Planned
Qualifying Exam,2024-06-01,2024-06-15,EXAM,Comprehensive qualifying exam,1,Planned
Thesis Proposal,2024-08-01,2024-08-31,WRITING,Write thesis proposal,1,Planned
```

## 🔬 Research Examples

### Experimental Research Timeline
**File**: `experimental_research.csv`
```csv
Task Name,Start Date,Due Date,Category,Description,Priority,Status
Protocol Development,2024-01-01,2024-02-28,RESEARCH,Develop experimental protocol,2,Planned
IRB Approval,2024-02-01,2024-03-31,ADMIN,Obtain IRB approval,1,Planned
Pilot Study,2024-03-01,2024-04-30,RESEARCH,Conduct pilot study,2,Planned
Main Study,2024-05-01,2024-10-31,RESEARCH,Conduct main study,1,Planned
Data Analysis,2024-09-01,2024-12-31,ANALYSIS,Analyze research data,1,Planned
Paper Writing,2024-11-01,2025-02-28,WRITING,Write research paper,1,Planned
Conference Presentation,2025-03-15,2025-03-20,PRESENTATION,Present at conference,2,Planned
```

### Grant Application Timeline
**File**: `grant_application.csv`
```csv
Task Name,Start Date,Due Date,Category,Description,Priority,Status
Grant Research,2024-01-01,2024-02-28,RESEARCH,Research funding opportunities,2,Planned
Proposal Writing,2024-02-01,2024-04-30,WRITING,Write grant proposal,1,Planned
Budget Development,2024-03-01,2024-04-15,ADMIN,Develop project budget,2,Planned
Review and Revision,2024-04-01,2024-05-15,WRITING,Review and revise proposal,1,Planned
Submission,2024-05-15,2024-05-15,MILESTONE,Submit grant application,1,Planned
Review Process,2024-06-01,2024-09-30,ADMIN,Grant review process,2,Planned
Award Notification,2024-10-01,2024-10-15,MILESTONE,Receive award notification,1,Planned
```

## 🏗️ Project Examples

### Software Development Project
**File**: `software_project.csv`
```csv
Task Name,Start Date,Due Date,Category,Description,Priority,Status
Requirements Analysis,2024-01-01,2024-01-31,PLANNING,Analyze project requirements,2,Planned
System Design,2024-02-01,2024-02-28,PLANNING,Design system architecture,2,Planned
Database Design,2024-02-15,2024-03-15,DEVELOPMENT,Design database schema,2,Planned
Backend Development,2024-03-01,2024-06-30,DEVELOPMENT,Develop backend services,1,Planned
Frontend Development,2024-04-01,2024-07-31,DEVELOPMENT,Develop user interface,1,Planned
Testing,2024-06-01,2024-08-31,TESTING,Comprehensive testing,1,Planned
Deployment,2024-08-15,2024-08-31,DEPLOYMENT,Deploy to production,1,Planned
```

### Conference Organization
**File**: `conference_organization.csv`
```csv
Task Name,Start Date,Due Date,Category,Description,Priority,Status
Venue Selection,2024-01-01,2024-02-28,PLANNING,Select conference venue,2,Planned
Call for Papers,2024-03-01,2024-03-31,COMMUNICATION,Release call for papers,1,Planned
Paper Review,2024-04-01,2024-06-30,REVIEW,Review submitted papers,1,Planned
Program Development,2024-07-01,2024-08-31,PLANNING,Develop conference program,1,Planned
Registration,2024-08-01,2024-09-30,ADMIN,Handle registrations,2,Planned
Event Execution,2024-10-15,2024-10-17,EVENT,Execute conference,1,Planned
```

## ⚙️ Configuration Examples

### Basic Configuration
**File**: `configs/basic.yaml`
```yaml
output:
  format: "pdf"
  quality: "high"
  
calendar:
  start_date: "2024-01-01"
  end_date: "2024-12-31"
  show_week_numbers: true
  
tasks:
  max_per_day: 5
  show_dependencies: true
  color_by_category: true
  
categories:
  RESEARCH:
    color: "blue"
    description: "Research activities"
  WRITING:
    color: "green"
    description: "Writing activities"
  MILESTONE:
    color: "red"
    description: "Important milestones"
```

### Advanced Configuration
**File**: `configs/advanced.yaml`
```yaml
output:
  format: "pdf"
  quality: "high"
  include_stats: true
  include_legend: true
  
calendar:
  start_date: "2024-01-01"
  end_date: "2024-12-31"
  show_week_numbers: true
  show_weekends: true
  
tasks:
  max_per_day: 8
  show_dependencies: true
  show_assignees: true
  color_by_category: true
  group_by_assignee: false
  
layout:
  task_height: 12.0
  task_spacing: 1.0
  compact_mode: false
  
categories:
  RESEARCH:
    color: "blue"
    description: "Research activities"
    priority: 2
  WRITING:
    color: "green"
    description: "Writing activities"
    priority: 1
  MILESTONE:
    color: "red"
    description: "Important milestones"
    priority: 1
  ADMIN:
    color: "orange"
    description: "Administrative tasks"
    priority: 3
```

## 🚀 Running Examples

### Generate All Examples
```bash
# Generate all example PDFs
for example in *.csv; do
    ./scripts/simple.sh "$example" "$(basename "$example" .csv)_example"
done
```

### Custom Examples
```bash
# Generate with custom configuration
./scripts/simple.sh phd_dissertation.csv phd_timeline --config configs/advanced.yaml

# Generate with specific date range
./scripts/simple.sh research_timeline.csv research_2024 --start-date 2024-01-01 --end-date 2024-12-31
```

## 📝 Tips for Creating Your Own Examples

1. **Start Simple**: Begin with basic task lists and gradually add complexity
2. **Use Categories**: Organize tasks by type (RESEARCH, WRITING, ADMIN, etc.)
3. **Set Priorities**: Use priority levels to indicate task importance
4. **Add Dependencies**: Show task relationships with dependency fields
5. **Include Milestones**: Mark important dates and deliverables
6. **Use Descriptions**: Add detailed descriptions for complex tasks

## 🔗 Related Resources

- [User Guide](../docs/user-guide/README.md) - Complete user documentation
- [Templates](../templates/README.md) - Template customization
- [Configuration Guide](../docs/user-guide/configuration.md) - Configuration options

---

*Need help with your specific use case? Check the [User Guide](../docs/user-guide/README.md) or create an issue with your example.*
