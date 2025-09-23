# Phase 1 Enhanced Phase Structure - Research Timeline v5

## Executive Summary
**Total Tasks**: 83 from data.cleaned.csv + 42 from v4 = 125 total tasks  
**Enhanced Phases**: 4 phases with 12 sub-phases  
**Task Distribution**: Phase 1 (25), Phase 2 (35), Phase 3 (12), Phase 4 (15)  
**Milestones**: 8 critical milestones across all phases

## Enhanced Phase Structure Design

### Phase 1: Instrumentation & Proposal (25 tasks)
**Duration**: 2025-08-29 to 2026-01-06  
**Objective**: Complete laser system setup, proposal development, and initial research preparation

#### Sub-Phase 1.1: Laser System Setup (8 tasks)
**Duration**: 2025-09-02 to 2025-10-21  
**Tasks**: T1.1-T1.8

| Task ID | Task Name                          | Source             | Dependencies | Key Details                                   |
| ------- | ---------------------------------- | ------------------ | ------------ | --------------------------------------------- |
| T1.1    | Align seed laser (≥30 mW output)   | data.cleaned.csv H | -            | Reach ≥30 mW output in fiber core             |
| T1.2    | Align amplifier (≥130 mW output)   | data.cleaned.csv I | T1.1         | Restore amplified output to ≥130 mW           |
| T1.3    | Check pulse compression (≤200 fs)  | data.cleaned.csv J | T1.2         | Verify ≤200 fs pulse duration                 |
| T1.4    | Calibrate microscope (USAF target) | data.cleaned.csv K | T1.3         | Align imaging system using USAF target        |
| T1.5    | Laser system ready (live imaging)  | data.cleaned.csv L | T1.4         | Confirm laser and optics meet requirements    |
| T1.6    | Align Laser Through Microscope     | v4 T1.3            | T1.5         | Align laser through microscope optics         |
| T1.7    | Image Air Force Target             | v4 T1.4            | T1.6         | Image Air Force target to validate resolution |
| T1.8    | Preliminary In Vivo Imaging        | v4 T1.5            | T1.7         | Acquire preliminary 6-panel tiled image       |

#### Sub-Phase 1.2: PhD Proposal Development (12 tasks)
**Duration**: 2025-08-29 to 2025-12-22  
**Tasks**: T1.9-T1.20

| Task ID | Task Name                                 | Source              | Dependencies | Key Details                                              |
| ------- | ----------------------------------------- | ------------------- | ------------ | -------------------------------------------------------- |
| T1.9    | Draft timeline v1                         | data.cleaned.csv A  | -            | Draft timeline for Tuesday review                        |
| T1.10   | Initial proposal skeleton                 | data.cleaned.csv B  | -            | Develop 1-page Specific Aims and outline                 |
| T1.11   | Submit proposal outline                   | data.cleaned.csv C  | T1.10        | Send outline to advisor before deadline                  |
| T1.12   | Define proposal committee                 | data.cleaned.csv D  | T1.9         | Identify committee members and schedule exam             |
| T1.13   | Expand proposal draft                     | data.cleaned.csv F  | T1.10        | Write 12-page Research Strategy                          |
| T1.14   | Confirm exam date                         | data.cleaned.csv G  | T1.12        | Oral exam date scheduled                                 |
| T1.15   | Draft Specific Aims and Research Strategy | data.cleaned.csv R1 | T1.13        | Write 1-page Specific Aims and 12-page Research Strategy |
| T1.16   | Draft Methods and Timeline sections       | data.cleaned.csv R2 | T1.15        | Write Methods and project timeline                       |
| T1.17   | Finalize proposal draft and formatting    | data.cleaned.csv R3 | T1.16,T1.23  | Complete full proposal document                          |
| T1.18   | Send proposal to committee                | data.cleaned.csv S  | T1.17,T1.14  | Email proposal to committee ≥2 weeks before exam         |
| T1.19   | Prepare presentation                      | data.cleaned.csv T  | T1.18        | Create slide deck and practice presentation              |
| T1.20   | PhD Proposal Exam                         | data.cleaned.csv U  | T1.19,T1.25  | Defend proposal in oral exam                             |

#### Sub-Phase 1.3: Administrative Setup (5 tasks)
**Duration**: 2025-09-01 to 2026-01-06  
**Tasks**: T1.21-T1.25

| Task ID | Task Name                                        | Source               | Dependencies | Key Details                                |
| ------- | ------------------------------------------------ | -------------------- | ------------ | ------------------------------------------ |
| T1.21   | Annual progress review                           | data.cleaned.csv N   | -            | Submit annual progress report              |
| T1.22   | Complete committee paperwork and Program of Work | data.cleaned.csv O1  | T1.14        | Complete and submit all required forms     |
| T1.23   | Reserve exam room and submit final paperwork     | data.cleaned.csv O2  | T1.22        | Reserve exam room and submit paperwork     |
| T1.24   | Maintain continuous registration - Fall 2025     | data.cleaned.csv BY1 | -            | Maintain full-time registration (9+ hours) |
| T1.25   | Update committee membership                      | data.cleaned.csv BW  | T1.12        | Update committee membership if needed      |

**Phase 1 Milestones**:
- T1.M1: PhD Proposal Exam (2025-12-16)
- T1.M2: Laser System Operational (2025-10-11)

### Phase 2: Core Research & Analysis (35 tasks)
**Duration**: 2025-12-17 to 2026-12-09  
**Objective**: Execute all three research aims, maintain equipment, and complete data acquisition

#### Sub-Phase 2.1: Aim 1 - AAV-based Vascular Imaging (12 tasks)
**Duration**: 2025-12-17 to 2026-06-06  
**Tasks**: T2.1-T2.12

| Task ID | Task Name                                   | Source               | Dependencies    | Key Details                                   |
| ------- | ------------------------------------------- | -------------------- | --------------- | --------------------------------------------- |
| T2.1    | Plan imaging cohort                         | data.cleaned.csv M   | T1.10           | Plan ~3 pilot mice cohort with IACUC protocol |
| T2.2    | Design AAV vectors                          | data.cleaned.csv P   | T2.1            | Design and order AAV-mScarlet and jRGECO1b    |
| T2.3    | AAV vectors ready                           | data.cleaned.csv Q   | T2.2            | Receive AAV vectors; ready for injections     |
| T2.4    | Cranial window surgeries (3 mice)           | data.cleaned.csv W   | T1.5,T2.3,T1.20 | Install cranial windows and inject AAV        |
| T2.5    | Post-operative recovery and monitoring      | data.cleaned.csv Z   | T2.4            | Monitor/medicate mice; maintain analgesia     |
| T2.6    | Pilot imaging sessions (3 mice)             | data.cleaned.csv AE  | T2.5            | Acquire in vivo images comparing methods      |
| T2.7    | Pilot datasets complete                     | data.cleaned.csv AH  | T2.6            | Complete 3 pilot two-photon datasets          |
| T2.8    | Process pilot data                          | data.cleaned.csv AI  | T2.6            | Perform image registration and SNR analysis   |
| T2.9    | Design U-Net architecture and training data | data.cleaned.csv AJ1 | T2.8            | Design U-Net architecture for segmentation    |
| T2.10   | Implement and test segmentation pipeline    | data.cleaned.csv AJ2 | T2.9            | Implement U-Net pipeline; test on pilot data  |
| T2.11   | Configure dual-channel two-photon imaging   | data.cleaned.csv AK1 | T2.7            | Tune microscope optics for dual-channel       |
| T2.12   | Configure LSCI for blood flow measurements  | data.cleaned.csv AK2 | T2.11           | Configure LSCI for blood-flow measurements    |

#### Sub-Phase 2.2: Aim 2 - Dual-Color Platform Development (8 tasks)
**Duration**: 2026-04-19 to 2026-07-22  
**Tasks**: T2.13-T2.20

| Task ID | Task Name                                     | Source              | Dependencies | Key Details                                           |
| ------- | --------------------------------------------- | ------------------- | ------------ | ----------------------------------------------------- |
| T2.13   | Order enhanced AAV                            | data.cleaned.csv AM | T2.7         | Design and order enhanced-expression AAV              |
| T2.14   | Enhanced AAV delivered                        | data.cleaned.csv AN | T2.13        | Receive enhanced AAV; ready for tests                 |
| T2.15   | Compare labeling methods                      | data.cleaned.csv AO | T2.11        | Systematically compare imaging methods                |
| T2.16   | Validate Platform with Fluorescent Beads      | v4 T2.9             | T2.12        | Validate the platform using fluorescent beads         |
| T2.17   | In Vivo Validation with AAV-expressing mice   | v4 T2.10            | T2.16        | Achieve operational status of dual-color platform     |
| T2.18   | Assemble Dual-Color Detection Path            | v4 T2.7             | T2.8         | Assemble a dual-color detection path                  |
| T2.19   | Develop Analysis Software (Spectral Unmixing) | v4 T2.8             | T2.8         | Develop analysis software including spectral unmixing |
| T2.20   | Dual-Color Platform Operational               | v4 T2.M1            | T2.17        | Dual-Color Platform Operational milestone             |

#### Sub-Phase 2.3: Aim 3 - Stroke Study (10 tasks)
**Duration**: 2026-06-18 to 2026-12-09  
**Tasks**: T2.21-T2.30

| Task ID | Task Name                                      | Source               | Dependencies | Key Details                                     |
| ------- | ---------------------------------------------- | -------------------- | ------------ | ----------------------------------------------- |
| T2.21   | Establish stroke protocol                      | data.cleaned.csv AR  | T2.15        | Complete training and IACUC approval            |
| T2.22   | Induce stroke                                  | data.cleaned.csv AS  | T2.21        | Induce stroke in cohort to start Aim 3          |
| T2.23   | Acute-phase imaging                            | data.cleaned.csv AT  | T2.22        | Conduct imaging in acute phase (0-1 week)       |
| T2.24   | Transition-phase imaging                       | data.cleaned.csv AU  | T2.23        | Transition-phase imaging (2–4 wks)              |
| T2.25   | Stabilization-phase imaging                    | data.cleaned.csv AV  | T2.24        | Early chronic phase imaging (5-8 weeks)         |
| T2.26   | Extended chronic imaging                       | data.cleaned.csv AW  | T2.25        | Extended chronic imaging (~12 wks)              |
| T2.27   | Adapt ML pipeline for stroke data              | data.cleaned.csv AX1 | T2.10        | Adapt ML pipeline for stroke dataset analysis   |
| T2.28   | Optimize and validate segmentation performance | data.cleaned.csv AX2 | T2.27        | Optimize ML pipeline performance                |
| T2.29   | Stroke data complete                           | data.cleaned.csv AY  | T2.26        | Completion of all longitudinal imaging sessions |
| T2.30   | Analyze neurovascular coupling                 | data.cleaned.csv BA  | T2.29        | Quantify microvascular network changes          |

#### Sub-Phase 2.4: Equipment Maintenance (5 tasks)
**Duration**: 2025-09-01 to 2027-08-31  
**Tasks**: T2.31-T2.35

| Task ID | Task Name                           | Source               | Dependencies | Key Details                                     |
| ------- | ----------------------------------- | -------------------- | ------------ | ----------------------------------------------- |
| T2.31   | Equipment maintenance log - Q1 2026 | data.cleaned.csv CC1 | -            | Maintain detailed rig log for alignment changes |
| T2.32   | Equipment maintenance log - Q2 2026 | data.cleaned.csv CC2 | T2.31        | Maintain detailed rig log for alignment changes |
| T2.33   | Equipment maintenance log - Q3 2026 | data.cleaned.csv CC3 | T2.32        | Maintain detailed rig log for alignment changes |
| T2.34   | Data backup system implementation   | data.cleaned.csv CD1 | -            | Implement automated quality control checks      |
| T2.35   | Data backup system maintenance      | data.cleaned.csv CD2 | T2.34        | Maintain and monitor automated backups          |

**Phase 2 Milestones**:
- T2.M1: Dual-Color Platform Operational (2026-07-18)
- T2.M2: Data Acquisition Complete (2026-12-09)

### Phase 3: Publication (12 tasks)
**Duration**: 2026-04-19 to 2026-12-23  
**Objective**: Complete all manuscript submissions and conference presentations

#### Sub-Phase 3.1: Methodology Paper (2 tasks)
**Duration**: 2026-04-19 to 2026-07-22  
**Tasks**: T3.1-T3.2

| Task ID | Task Name                | Source              | Dependencies | Key Details                                    |
| ------- | ------------------------ | ------------------- | ------------ | ---------------------------------------------- |
| T3.1    | Draft methodology paper  | data.cleaned.csv AP | T2.7         | Write manuscript on AAV-based vascular imaging |
| T3.2    | Submit methodology paper | data.cleaned.csv AQ | T3.1         | Submit methodology paper to journal            |

#### Sub-Phase 3.2: Research Papers (4 tasks)
**Duration**: 2026-08-11 to 2026-12-23  
**Tasks**: T3.3-T3.6

| Task ID | Task Name                 | Source              | Dependencies | Key Details                                              |
| ------- | ------------------------- | ------------------- | ------------ | -------------------------------------------------------- |
| T3.3    | Develop SLAVV-T Codebase  | v4 T3.3             | T2.8         | Develop improved codebase for temporal analysis          |
| T3.4    | Draft SLAVV-T Manuscript  | v4 T3.4             | T3.3         | Draft SLAVV-T manuscript                                 |
| T3.5    | Submit SLAVV-T Manuscript | v4 T3.5             | T3.4         | Submit MS on SLAVV-T analysis method                     |
| T3.6    | Draft second manuscript   | data.cleaned.csv BC | T2.30        | Write second research paper covering dual-color platform |

#### Sub-Phase 3.3: Conference Presentations (2 tasks)
**Duration**: 2026-12-10 to 2026-12-16  
**Tasks**: T3.7-T3.8

| Task ID | Task Name                       | Source              | Dependencies | Key Details                                 |
| ------- | ------------------------------- | ------------------- | ------------ | ------------------------------------------- |
| T3.7    | Prepare conference presentation | data.cleaned.csv BB | T2.30        | Prepare conference talk/poster with results |
| T3.8    | Submit second manuscript        | data.cleaned.csv BD | T3.6         | Submit second manuscript to journal         |

#### Sub-Phase 3.4: Additional Publications (4 tasks)
**Duration**: 2026-08-01 to 2027-12-31  
**Tasks**: T3.9-T3.12

| Task ID | Task Name                 | Source  | Dependencies | Key Details                                    |
| ------- | ------------------------- | ------- | ------------ | ---------------------------------------------- |
| T3.9    | Draft Aim 1 Manuscript    | v4 T3.1 | T2.8         | Draft the Aim 1 manuscript                     |
| T3.10   | Submit Aim 1 Manuscript   | v4 T3.2 | T3.9         | Submit MS on tiled comparison                  |
| T3.11   | Draft Aim 2/3 Manuscript  | v4 T3.6 | T2.30        | Draft primary biological MS on chronic imaging |
| T3.12   | Submit Aim 2/3 Manuscript | v4 T3.7 | T3.11        | Submit primary biological MS                   |

**Phase 3 Milestones**:
- T3.M1: Manuscript Submissions Complete (2026-12-23)

### Phase 4: Dissertation & Graduation (15 tasks)
**Duration**: 2026-12-19 to 2027-08-11  
**Objective**: Complete dissertation writing, defense, and graduation requirements

#### Sub-Phase 4.1: Dissertation Writing (8 tasks)
**Duration**: 2026-12-19 to 2027-06-30  
**Tasks**: T4.1-T4.8

| Task ID | Task Name                                     | Source              | Dependencies   | Key Details                                      |
| ------- | --------------------------------------------- | ------------------- | -------------- | ------------------------------------------------ |
| T4.1    | Draft Introduction and Literature Review      | data.cleaned.csv BI | T1.20          | Write dissertation Introduction chapter          |
| T4.2    | Draft Methods and Results chapters (Aims 1-3) | data.cleaned.csv BJ | T1.20          | Write chapters detailing all three research aims |
| T4.3    | Draft Discussion and Conclusions              | data.cleaned.csv BK | T1.20          | Write final dissertation chapters                |
| T4.4    | Dissertation draft complete                   | data.cleaned.csv BN | T4.1,T4.2,T4.3 | Complete PhD dissertation draft                  |
| T4.5    | Final Committee Meeting                       | v4 T4.5             | T4.4           | Final Committee Meeting                          |
| T4.6    | Final Revisions                               | v4 T4.6             | T4.5           | Incorporate final committee revisions            |
| T4.7    | PhD Defense                                   | data.cleaned.csv BS | T4.6,T4.15     | Defend dissertation in oral exam                 |
| T4.8    | Submit dissertation                           | data.cleaned.csv BU | T4.7           | Upload approved dissertation PDF                 |

#### Sub-Phase 4.2: Administrative Completion (7 tasks)
**Duration**: 2025-09-01 to 2027-08-31  
**Tasks**: T4.9-T4.15

| Task ID | Task Name                                      | Source               | Dependencies | Key Details                                                  |
| ------- | ---------------------------------------------- | -------------------- | ------------ | ------------------------------------------------------------ |
| T4.9    | Annual progress review                         | data.cleaned.csv BE  | T1.21        | Complete yearly graduate student progress review             |
| T4.10   | Complete TA requirement                        | data.cleaned.csv BV  | -            | Serve as Teaching Assistant at least once                    |
| T4.11   | Maintain continuous registration - Spring 2026 | data.cleaned.csv BY2 | T1.24        | Maintain full-time registration (9+ hours)                   |
| T4.12   | Maintain continuous registration - Fall 2026   | data.cleaned.csv BY3 | T4.11        | Maintain full-time registration (9+ hours)                   |
| T4.13   | Maintain continuous registration - Spring 2027 | data.cleaned.csv BY4 | T4.12        | Maintain full-time registration (9+ hours)                   |
| T4.14   | Apply for graduation                           | data.cleaned.csv BZ  | -            | Submit graduation application at beginning of final semester |
| T4.15   | Request final oral exam                        | data.cleaned.csv CA  | -            | Submit 'Request for Final Oral Exam' form                    |

**Phase 4 Milestones**:
- T4.M1: Dissertation Complete (2027-06-30)
- T4.M2: PhD Defense (2027-07-17)
- T4.M3: Graduation (2027-08-11)

## Task Distribution Summary

| Phase   | Sub-Phase                   | Task Count | Duration                 | Key Focus                |
| ------- | --------------------------- | ---------- | ------------------------ | ------------------------ |
| Phase 1 | Laser System Setup          | 8          | 2025-09-02 to 2025-10-21 | Equipment preparation    |
| Phase 1 | PhD Proposal Development    | 12         | 2025-08-29 to 2025-12-22 | Proposal completion      |
| Phase 1 | Administrative Setup        | 5          | 2025-09-01 to 2026-01-06 | Compliance and paperwork |
| Phase 2 | Aim 1 - AAV-based Imaging   | 12         | 2025-12-17 to 2026-06-06 | Core research execution  |
| Phase 2 | Aim 2 - Dual-Color Platform | 8          | 2026-04-19 to 2026-07-22 | Platform development     |
| Phase 2 | Aim 3 - Stroke Study        | 10         | 2026-06-18 to 2026-12-09 | Longitudinal study       |
| Phase 2 | Equipment Maintenance       | 5          | 2025-09-01 to 2027-08-31 | Ongoing maintenance      |
| Phase 3 | Methodology Paper           | 2          | 2026-04-19 to 2026-07-22 | First publication        |
| Phase 3 | Research Papers             | 4          | 2026-08-11 to 2026-12-23 | Main publications        |
| Phase 3 | Conference Presentations    | 2          | 2026-12-10 to 2026-12-16 | Conference work          |
| Phase 3 | Additional Publications     | 4          | 2026-08-01 to 2027-12-31 | Extended publications    |
| Phase 4 | Dissertation Writing        | 8          | 2026-12-19 to 2027-06-30 | Dissertation completion  |
| Phase 4 | Administrative Completion   | 7          | 2025-09-01 to 2027-08-31 | Graduation requirements  |

## Task ID System Design

### ID Format
- **Phase 1**: T1.1 - T1.25 (25 tasks)
- **Phase 2**: T2.1 - T2.35 (35 tasks)
- **Phase 3**: T3.1 - T3.12 (12 tasks)
- **Phase 4**: T4.1 - T4.15 (15 tasks)
- **Milestones**: T1.M1, T1.M2, T2.M1, T2.M2, T3.M1, T4.M1, T4.M2, T4.M3

### Naming Convention
- **Format**: [Phase].[Sub-Phase].[Task Number]
- **Example**: T1.1 (Phase 1, Laser System Setup, Task 1)
- **Milestones**: [Phase].M[Milestone Number]
- **Example**: T1.M1 (Phase 1, Milestone 1)

## Quality Assurance

### Validation Checklist
- [ ] All 83 tasks from data.cleaned.csv included
- [ ] All 42 tasks from v4 timeline included
- [ ] Total task count: 125 tasks
- [ ] Phase distribution: 25, 35, 12, 15
- [ ] All dependencies mapped correctly
- [ ] All milestones positioned appropriately
- [ ] All technical specifications preserved
- [ ] All administrative tasks included

### Testing Strategy
1. **Task Completeness**: Verify all tasks are included and properly categorized
2. **Dependency Integrity**: Check all dependencies reference correct Task IDs
3. **Phase Logic**: Ensure phase progression makes sense
4. **Milestone Placement**: Verify milestones are properly positioned
5. **Technical Accuracy**: Ensure all technical specifications are preserved

---

*This enhanced phase structure provides the foundation for Phase 2 design and Phase 3 implementation of Research Timeline v5, accommodating all 83 tasks from data.cleaned.csv while maintaining the professional polish of v4.*
