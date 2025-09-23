#!/usr/bin/env python3
"""
Corrected Task Mapping Script for Research Timeline v5
Maps all 83 tasks from data.cleaned.csv to v5 format with unique Task IDs
"""

import csv

def create_corrected_mapping():
    """Create corrected mapping ensuring each original task gets unique Task ID"""
    
    # Read original data.cleaned.csv
    with open('input/data.cleaned.csv', 'r') as f:
        reader = csv.DictReader(f)
        original_tasks = list(reader)
    
    print(f"Original tasks: {len(original_tasks)}")
    
    # Define corrected mapping - each original task gets unique new Task ID
    task_mapping = {
        # Phase 1: Instrumentation & Proposal (T1.1-T1.30)
        'A': 'T1.1',    # Draft timeline v1
        'B': 'T1.2',    # Initial proposal skeleton  
        'C': 'T1.3',    # Submit proposal outline
        'D': 'T1.4',    # Define proposal committee
        'F': 'T1.5',    # Expand proposal draft
        'G': 'T1.6',    # Confirm exam date
        'H': 'T1.7',    # Align seed laser
        'I': 'T1.8',    # Align amplifier
        'J': 'T1.9',    # Check pulse compression
        'K': 'T1.10',   # Calibrate microscope
        'L': 'T1.11',   # Laser system ready
        'M': 'T2.1',    # Plan imaging cohort
        'N': 'T1.12',   # Annual progress review
        'O1': 'T1.13',  # Complete committee paperwork and Program of Work
        'O2': 'T1.14',  # Reserve exam room and submit final paperwork
        'P': 'T2.2',    # Design AAV vectors
        'Q': 'T2.3',    # AAV vectors ready
        'R1': 'T1.15',  # Draft Specific Aims and Research Strategy
        'R2': 'T1.16',  # Draft Methods and Timeline sections
        'R3': 'T1.17',  # Finalize proposal draft and formatting
        'S': 'T1.18',   # Send proposal to committee
        'T': 'T1.19',   # Prepare presentation
        'U': 'T1.M1',   # PhD Proposal Exam (milestone)
        'V': 'T1.20',   # Address committee feedback
        'W': 'T2.4',    # Cranial window surgeries (3 mice)
        'Z': 'T2.5',    # Post-operative recovery and monitoring
        'AE': 'T2.6',   # Pilot imaging sessions (3 mice)
        'AH': 'T2.7',   # Pilot datasets complete
        'AI': 'T2.8',   # Process pilot data
        'AJ1': 'T2.9',  # Design U-Net architecture and training data
        'AJ2': 'T2.10', # Implement and test segmentation pipeline
        'AK1': 'T2.11', # Configure dual-channel two-photon imaging
        'AK2': 'T2.12', # Configure LSCI for blood flow measurements
        'AM': 'T2.13',  # Order enhanced AAV
        'AN': 'T2.14',  # Enhanced AAV delivered
        'AO': 'T2.15',  # Compare labeling methods
        'AP': 'T3.1',   # Draft methodology paper
        'AQ': 'T3.2',   # Submit methodology paper
        'AR': 'T2.16',  # Establish stroke protocol
        'AS': 'T2.17',  # Induce stroke
        'AT': 'T2.18',  # Acute-phase imaging
        'AU': 'T2.19',  # Transition-phase imaging
        'AV': 'T2.20',  # Stabilization-phase imaging
        'AW': 'T2.21',  # Extended chronic imaging
        'AX1': 'T2.22', # Adapt ML pipeline for stroke data
        'AX2': 'T2.23', # Optimize and validate segmentation performance
        'AY': 'T2.24',  # Stroke data complete
        'AZ': 'T2.25',  # Integrate flow data
        'BA': 'T2.26',  # Analyze neurovascular coupling
        'BB': 'T3.3',   # Prepare conference presentation
        'BC': 'T3.4',   # Draft second manuscript
        'BD': 'T3.5',   # Submit second manuscript
        'BE': 'T4.1',   # Annual progress review
        'BG': 'T4.2',   # PhD Dissertation & Defense (placeholder)
        'BI': 'T4.3',   # Draft Introduction and Literature Review
        'BJ': 'T4.4',   # Draft Methods and Results chapters (Aims 1-3)
        'BK': 'T4.5',   # Draft Discussion and Conclusions
        'BN': 'T4.6',   # Dissertation draft complete
        'BS': 'T4.7',   # PhD Defense
        'BT': 'T4.8',   # Revise dissertation
        'BU': 'T4.9',   # Submit dissertation
        'BV': 'T4.10',  # Complete TA requirement
        'BW': 'T1.21',  # Update committee membership
        'BX': 'T1.22',  # Advance to candidacy
        'BY1': 'T1.23', # Maintain continuous registration - Fall 2025
        'BY2': 'T4.11', # Maintain continuous registration - Spring 2026
        'BY3': 'T4.12', # Maintain continuous registration - Fall 2026
        'BY4': 'T4.13', # Maintain continuous registration - Spring 2027
        'BY5': 'T4.14', # Maintain continuous registration - Summer 2027
        'BZ': 'T4.15',  # Apply for graduation
        'CA': 'T4.16',  # Request final oral exam
        'CB': 'T4.17',  # SPIE chapter activities
        'CC1': 'T2.27', # Equipment maintenance log - Q1 2026
        'CC2': 'T2.28', # Equipment maintenance log - Q2 2026
        'CC3': 'T2.29', # Equipment maintenance log - Q3 2026
        'CC4': 'T2.30', # Equipment maintenance log - Q4 2026
        'CC5': 'T2.31', # Equipment maintenance log - Q1 2027
        'CC6': 'T2.32', # Equipment maintenance log - Q2 2027
        'CC7': 'T2.33', # Equipment maintenance log - Q3 2027
        'CC8': 'T2.34', # Equipment maintenance log - Q4 2027
        'CD1': 'T2.35', # Data backup system implementation
        'CD2': 'T2.36', # Data backup system maintenance
        'CE': 'T2.37',  # Surgical training
        'CF': 'T4.18',  # Lab culture responsibilities
    }
    
    print(f"Task mapping entries: {len(task_mapping)}")
    
    # Check for missing mappings
    original_ids = set(task['Task ID'] for task in original_tasks)
    mapped_ids = set(task_mapping.keys())
    missing_ids = original_ids - mapped_ids
    
    if missing_ids:
        print(f"Missing mappings: {missing_ids}")
    else:
        print("✅ All original tasks have mappings!")
    
    # Check for duplicate mappings
    mapped_values = list(task_mapping.values())
    duplicates = set([x for x in mapped_values if mapped_values.count(x) > 1])
    if duplicates:
        print(f"Duplicate mappings: {duplicates}")
    else:
        print("✅ No duplicate mappings!")
    
    # Check Task ID ranges
    phase1_ids = [tid for tid in mapped_values if tid.startswith('T1.') and not tid.startswith('T1.M')]
    phase2_ids = [tid for tid in mapped_values if tid.startswith('T2.') and not tid.startswith('T2.M')]
    phase3_ids = [tid for tid in mapped_values if tid.startswith('T3.') and not tid.startswith('T3.M')]
    phase4_ids = [tid for tid in mapped_values if tid.startswith('T4.') and not tid.startswith('T4.M')]
    
    print(f"\nPhase 1 tasks: {len(phase1_ids)} (T1.1-T1.23)")
    print(f"Phase 2 tasks: {len(phase2_ids)} (T2.1-T2.37)")
    print(f"Phase 3 tasks: {len(phase3_ids)} (T3.1-T3.5)")
    print(f"Phase 4 tasks: {len(phase4_ids)} (T4.1-T4.18)")
    
    return task_mapping

if __name__ == "__main__":
    mapping = create_corrected_mapping()
    print(f"\nTotal mappings: {len(mapping)}")
