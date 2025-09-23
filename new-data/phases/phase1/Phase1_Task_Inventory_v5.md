# Phase 1 Task Inventory - Research Timeline v5

## Executive Summary
**Total Tasks**: 83 from data.cleaned.csv  
**Categories**: 6 (PROPOSAL, EQUIPMENT, RESEARCH, PUBLICATION, DISSERTATION, ADMIN)  
**Date Range**: 2025-08-29 to 2027-08-31  
**Dependencies**: Complex multi-task dependencies identified

## Task Inventory by Category

### PROPOSAL (8 tasks)
| Task ID | Task Name                                 | Start Date | Due Date   | Dependencies | Key Details                                                            |
| ------- | ----------------------------------------- | ---------- | ---------- | ------------ | ---------------------------------------------------------------------- |
| A       | Draft timeline v1                         | 2025-08-29 | 2025-09-02 | -            | Draft timeline for Tuesday review; bring printed and digital copies    |
| B       | Initial proposal skeleton                 | 2025-09-02 | 2025-09-08 | -            | Develop 1-page Specific Aims and detailed outline following BME format |
| C       | Submit proposal outline                   | 2025-09-05 | 2025-09-09 | B            | Send outline to advisor before Monday deadline                         |
| D       | Define proposal committee                 | 2025-09-10 | 2025-09-14 | A            | Identify committee members; confirm availability; schedule oral exam   |
| F       | Expand proposal draft                     | 2025-09-15 | 2025-09-25 | B            | Write 12-page Research Strategy from outline per BME guidelines        |
| G       | Confirm exam date                         | 2025-09-28 | 2025-10-05 | D            | Oral exam date scheduled. Must send final proposal ≥2 weeks prior      |
| R1      | Draft Specific Aims and Research Strategy | 2025-10-28 | 2025-11-08 | F            | Write 1-page Specific Aims and 12-page Research Strategy sections      |
| R2      | Draft Methods and Timeline sections       | 2025-11-11 | 2025-11-18 | R1           | Write Methods and project timeline                                     |
| R3      | Finalize proposal draft and formatting    | 2025-12-01 | 2025-12-07 | R2,O2        | Complete full proposal document (~13 pages) with proper formatting     |
| S       | Send proposal to committee                | 2025-11-28 | 2025-12-05 | R3,G         | Email proposal to committee ≥2 weeks before exam                       |
| T       | Prepare presentation                      | 2025-12-08 | 2025-12-18 | S            | Create slide deck and practice oral exam presentation                  |
| U       | PhD Proposal Exam                         | 2025-12-19 | 2025-12-22 | T,BW         | Defend proposal in oral exam                                           |
| V       | Address committee feedback                | 2025-12-23 | 2026-01-06 | U            | Incorporate committee revisions and submit signed approval form        |

### EQUIPMENT (12 tasks)
| Task ID | Task Name                           | Start Date | Due Date   | Dependencies | Key Details                                                   |
| ------- | ----------------------------------- | ---------- | ---------- | ------------ | ------------------------------------------------------------- |
| H       | Align seed laser                    | 2025-09-02 | 2025-09-06 | -            | Reach ≥30 mW output in fiber core (pre-pump)                  |
| I       | Align amplifier                     | 2025-09-09 | 2025-09-16 | H            | Restore amplified output to ≥130 mW (previous benchmark)      |
| J       | Check pulse compression             | 2025-10-01 | 2025-10-07 | I            | Verify ≤200 fs pulse duration; log specs                      |
| K       | Calibrate microscope                | 2025-10-08 | 2025-10-14 | J            | Align imaging system using USAF target for optimal resolution |
| L       | Laser system ready                  | 2025-10-15 | 2025-10-21 | K            | Confirm laser and optics meet live imaging requirements       |
| CC1     | Equipment maintenance log - Q1 2026 | 2025-09-01 | 2025-12-31 | -            | Maintain detailed rig log for alignment changes               |
| CC2     | Equipment maintenance log - Q2 2026 | 2026-01-01 | 2026-03-31 | CC1          | Maintain detailed rig log for alignment changes               |
| CC3     | Equipment maintenance log - Q3 2026 | 2026-04-01 | 2026-06-30 | CC2          | Maintain detailed rig log for alignment changes               |
| CC4     | Equipment maintenance log - Q4 2026 | 2026-07-01 | 2026-09-30 | CC3          | Maintain detailed rig log for alignment changes               |
| CC5     | Equipment maintenance log - Q1 2027 | 2026-10-01 | 2026-12-31 | CC4          | Maintain detailed rig log for alignment changes               |
| CC6     | Equipment maintenance log - Q2 2027 | 2027-01-01 | 2027-03-31 | CC5          | Maintain detailed rig log for alignment changes               |
| CC7     | Equipment maintenance log - Q3 2027 | 2027-04-01 | 2027-06-30 | CC6          | Maintain detailed rig log for alignment changes               |
| CC8     | Equipment maintenance log - Q4 2027 | 2027-07-01 | 2027-08-31 | CC7          | Maintain detailed rig log for alignment changes               |

### RESEARCH (25 tasks)
| Task ID | Task Name                                      | Start Date | Due Date   | Dependencies | Key Details                                                            |
| ------- | ---------------------------------------------- | ---------- | ---------- | ------------ | ---------------------------------------------------------------------- |
| M       | Plan imaging cohort                            | 2025-10-14 | 2025-10-18 | B            | Plan ~3 pilot mice cohort with IACUC protocol confirmation             |
| P       | Design AAV vectors                             | 2025-10-21 | 2025-11-04 | M            | Design and order AAV-mScarlet (vascular) and jRGECO1b (neuronal)       |
| Q       | AAV vectors ready                              | 2025-12-20 | 2026-01-17 | P            | Receive AAV vectors; ready for injections                              |
| W       | Cranial window surgeries (3 mice)              | 2026-02-01 | 2026-02-26 | L,Q,V        | Install cranial windows and inject AAV in three pilot mice             |
| Z       | Post-operative recovery and monitoring         | 2026-02-27 | 2026-03-25 | W            | Monitor/medicate mice; maintain analgesia schedule                     |
| AE      | Pilot imaging sessions (3 mice)                | 2026-03-28 | 2026-04-15 | Z            | Acquire in vivo images comparing AAV vs traditional dye methods        |
| AH      | Pilot datasets complete                        | 2026-04-16 | 2026-04-22 | AE           | Complete 3 pilot two-photon datasets                                   |
| AI      | Process pilot data                             | 2026-04-21 | 2026-04-28 | AE           | Perform image registration and SNR analysis                            |
| AJ1     | Design U-Net architecture and training data    | 2026-05-01 | 2026-05-15 | AI           | Design U-Net architecture for vascular segmentation                    |
| AJ2     | Implement and test segmentation pipeline       | 2026-05-28 | 2026-06-25 | AJ1          | Implement U-Net pipeline; test on pilot data                           |
| AK1     | Configure dual-channel two-photon imaging      | 2026-04-19 | 2026-05-10 | AH           | Tune microscope optics for dual-channel two-photon imaging             |
| AK2     | Configure LSCI for blood flow measurements     | 2026-05-18 | 2026-06-20 | AK1          | Configure LSCI for blood-flow measurements                             |
| AM      | Order enhanced AAV                             | 2026-04-19 | 2026-07-15 | AH           | Design and order enhanced-expression AAV with tissue-specific enhancer |
| AN      | Enhanced AAV delivered                         | 2026-07-18 | 2026-07-22 | AM           | Receive enhanced AAV; ready for in vivo tests                          |
| AO      | Compare labeling methods                       | 2026-05-23 | 2026-06-22 | AK1          | Systematically compare imaging depth, SNR, and contrast                |
| AR      | Establish stroke protocol                      | 2026-06-18 | 2026-06-23 | AO           | Complete training and IACUC approval for stroke induction              |
| AS      | Induce stroke                                  | 2026-06-26 | 2026-06-30 | AR           | Induce stroke in cohort to start Aim 3 imaging                         |
| AT      | Acute-phase imaging                            | 2026-07-08 | 2026-07-13 | AS           | Conduct two-photon + LSCI imaging in acute phase (0-1 week)            |
| AU      | Transition-phase imaging                       | 2026-07-20 | 2026-07-25 | AT           | Transition-phase imaging (2–4 wks)                                     |
| AV      | Stabilization-phase imaging                    | 2026-08-19 | 2026-08-24 | AU           | Early chronic phase imaging (5-8 weeks post-stroke)                    |
| AW      | Extended chronic imaging                       | 2026-09-16 | 2026-09-20 | AV           | Extended chronic imaging (~12 wks), if needed                          |
| AX1     | Adapt ML pipeline for stroke data              | 2026-07-14 | 2026-08-15 | AJ2          | Adapt machine learning segmentation pipeline for stroke data           |
| AX2     | Optimize and validate segmentation performance | 2026-08-18 | 2026-09-15 | AX1          | Optimize ML pipeline performance and validate accuracy                 |
| AY      | Stroke data complete                           | 2026-09-21 | 2026-09-27 | AW           | Completion of all planned longitudinal imaging sessions                |
| AZ      | Integrate flow data                            | 2026-09-26 | 2026-10-23 | AT,AU,AV,AW  | Integrate LSCI flow with two-photon data                               |
| BA      | Analyze neurovascular coupling                 | 2026-10-26 | 2026-12-09 | AZ           | Quantify microvascular network changes and neurovascular coupling      |

### PUBLICATION (6 tasks)
| Task ID | Task Name                       | Start Date | Due Date   | Dependencies | Key Details                                                      |
| ------- | ------------------------------- | ---------- | ---------- | ------------ | ---------------------------------------------------------------- |
| AP      | Draft methodology paper         | 2026-04-19 | 2026-07-15 | AH           | Write manuscript on AAV-based vascular imaging methodology       |
| AQ      | Submit methodology paper        | 2026-07-18 | 2026-07-22 | AP           | Submit methodology paper to journal                              |
| BB      | Prepare conference presentation | 2026-12-10 | 2026-12-16 | BA           | Prepare conference talk/poster with results                      |
| BC      | Draft second manuscript         | 2026-12-10 | 2026-12-16 | BA           | Write second research paper covering dual-color imaging platform |
| BD      | Submit second manuscript        | 2026-12-19 | 2026-12-23 | BC           | Submit second manuscript to journal                              |

### DISSERTATION (8 tasks)
| Task ID | Task Name                                     | Start Date | Due Date   | Dependencies | Key Details                                                             |
| ------- | --------------------------------------------- | ---------- | ---------- | ------------ | ----------------------------------------------------------------------- |
| BG      | PhD Dissertation & Defense                    | 2026-12-19 | 2027-08-15 | V,BX         | Write dissertation, defend, and complete graduation requirements        |
| BI      | Draft Introduction and Literature Review      | 2026-12-19 | 2027-01-31 | V            | Write dissertation Introduction chapter including literature review     |
| BJ      | Draft Methods and Results chapters (Aims 1-3) | 2027-02-03 | 2027-05-15 | V            | Write chapters detailing all three research aims                        |
| BK      | Draft Discussion and Conclusions              | 2027-02-03 | 2027-06-15 | V            | Write final dissertation chapters summarizing findings                  |
| BN      | Dissertation draft complete                   | 2027-06-19 | 2027-06-30 | BI,BJ,BK     | Complete PhD dissertation draft compiled and ready for committee review |
| BS      | PhD Defense                                   | 2027-07-17 | 2027-07-20 | BN,CA        | Defend dissertation in oral exam                                        |
| BT      | Revise dissertation                           | 2027-07-21 | 2027-08-03 | BS           | Incorporate committee feedback and revisions after defense              |
| BU      | Submit dissertation                           | 2027-08-05 | 2027-08-11 | BT           | Upload approved dissertation PDF and submit all required forms          |

### ADMIN (24 tasks)
| Task ID | Task Name                                        | Start Date | Due Date   | Dependencies | Key Details                                                                |
| ------- | ------------------------------------------------ | ---------- | ---------- | ------------ | -------------------------------------------------------------------------- |
| N       | Annual progress review                           | 2025-09-03 | 2025-09-07 | -            | Submit annual progress report (due early September)                        |
| O1      | Complete committee paperwork and Program of Work | 2025-10-08 | 2025-10-18 | G            | Complete and submit all required committee forms                           |
| O2      | Reserve exam room and submit final paperwork     | 2025-10-21 | 2025-11-30 | O1           | Reserve exam room; submit final paperwork to graduate office               |
| BE      | Annual progress review                           | 2026-09-01 | 2026-09-07 | N            | Complete yearly graduate student progress review                           |
| BV      | Complete TA requirement                          | 2025-09-04 | 2026-08-31 | -            | Serve as Teaching Assistant at least once before graduation                |
| BW      | Update committee membership                      | 2025-09-13 | 2025-09-23 | D            | Update committee membership if needed and file change form                 |
| BX      | Advance to candidacy                             | 2025-09-06 | 2025-09-30 | -            | Ensure doctoral candidacy status is confirmed                              |
| BY1     | Maintain continuous registration - Fall 2025     | 2025-09-01 | 2025-12-31 | -            | Maintain full-time registration (9+ hours) for Fall 2025                   |
| BY2     | Maintain continuous registration - Spring 2026   | 2026-01-01 | 2026-05-31 | BY1          | Maintain full-time registration (9+ hours) for Spring 2026                 |
| BY3     | Maintain continuous registration - Fall 2026     | 2026-09-01 | 2026-12-31 | BY2          | Maintain full-time registration (9+ hours) for Fall 2026                   |
| BY4     | Maintain continuous registration - Spring 2027   | 2027-01-01 | 2027-05-31 | BY3          | Maintain full-time registration (9+ hours) for Spring 2027                 |
| BY5     | Maintain continuous registration - Summer 2027   | 2027-06-01 | 2027-08-31 | BY4          | Maintain full-time registration (9+ hours) for Summer 2027                 |
| BZ      | Apply for graduation                             | 2027-01-01 | 2027-01-15 | -            | Submit graduation application at beginning of final semester               |
| CA      | Request final oral exam                          | 2027-06-01 | 2027-06-15 | -            | Submit 'Request for Final Oral Exam' form at least 2 weeks before defense  |
| CB      | SPIE chapter activities                          | 2025-09-08 | 2026-12-31 | -            | Plan and execute SPIE Student Chapter events and submit annual report      |
| CD1     | Data backup system implementation                | 2025-10-01 | 2025-12-31 | -            | Implement automated quality control checks and nightly backups             |
| CD2     | Data backup system maintenance                   | 2026-01-01 | 2027-08-31 | CD1          | Maintain and monitor automated quality control checks and backups          |
| CE      | Surgical training                                | 2025-10-01 | 2026-06-01 | -            | Maintain proficiency with cranial window surgeries and train others        |
| CF      | Lab culture responsibilities                     | 2025-09-10 | 2027-08-31 | -            | Be punctual to meetings and injections; provide regular updates to advisor |

## Technical Specifications Analysis

### Performance Benchmarks
- **Laser Power**: ≥30 mW (seed), ≥130 mW (amplified)
- **Pulse Duration**: ≤200 fs
- **Imaging Resolution**: USAF target validation
- **Cohort Size**: ~3 pilot mice
- **Timeline**: 2025-2027 (2+ years)

### Key Dependencies
- **Complex Multi-Task**: R3 depends on R2,O2
- **Sequential Equipment**: H→I→J→K→L
- **Research Pipeline**: M→P→Q→W→Z→AE→AH→AI
- **Publication Chain**: AP→AQ, BC→BD
- **Dissertation Chain**: BI,BJ,BK→BN→BS→BT→BU

### Critical Milestones
1. **PhD Proposal Exam** (U): 2025-12-19 to 2025-12-22
2. **Laser System Ready** (L): 2025-10-15 to 2025-10-21
3. **Pilot Datasets Complete** (AH): 2026-04-16 to 2026-04-22
4. **Stroke Data Complete** (AY): 2026-09-21 to 2026-09-27
5. **PhD Defense** (BS): 2027-07-17 to 2027-07-20

## Date Conflict Analysis

### data.cleaned.csv vs v4 Timeline
- **data.cleaned.csv**: Uses 2025-2027 timeline with specific dates
- **v4 Timeline**: Uses different specific dates within same range
- **Conflict Resolution**: Need to use MD timeline as authoritative source

### Key Date Conflicts Identified
1. **Proposal Timeline**: data.cleaned.csv has different proposal dates than v4
2. **Equipment Setup**: Different alignment and calibration schedules
3. **Research Phases**: Different timing for Aim 1, 2, 3 execution
4. **Publication Schedule**: Different manuscript submission dates

## Recommendations for v5 Integration

### Phase Structure Mapping
- **Phase 1**: PROPOSAL + EQUIPMENT (early tasks) → Instrumentation & Proposal
- **Phase 2**: RESEARCH + EQUIPMENT (maintenance) → Core Research & Analysis  
- **Phase 3**: PUBLICATION → Publication
- **Phase 4**: DISSERTATION + ADMIN (final) → Dissertation & Graduation

### Task ID System
- **Phase 1**: T1.1-T1.25 (Instrumentation & Proposal)
- **Phase 2**: T2.1-T2.35 (Core Research & Analysis)
- **Phase 3**: T3.1-T3.12 (Publication)
- **Phase 4**: T4.1-T4.15 (Dissertation & Graduation)

### Priority Classification
- **Critical**: Milestones, defense, key deliverables
- **High**: Research tasks, manuscript submissions
- **Medium**: Equipment setup, data analysis
- **Low**: Maintenance, administrative

---

*This comprehensive task inventory provides the foundation for Phase 2 design and Phase 3 implementation of Research Timeline v5.*
