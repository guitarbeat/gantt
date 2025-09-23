#!/usr/bin/env python3
"""
Task Mapping Script for Research Timeline v5
Maps all 83 tasks from data.cleaned.csv to v5 format with proper Task IDs
"""

import csv
from datetime import datetime

def create_task_mapping():
    """Create comprehensive mapping of all 83 tasks from data.cleaned.csv to v5 format"""
    
    # Read original data.cleaned.csv
    with open('input/data.cleaned.csv', 'r') as f:
        reader = csv.DictReader(f)
        original_tasks = list(reader)
    
    print(f"Original tasks: {len(original_tasks)}")
    
    # Define the mapping from original Task ID to new Task ID
    task_mapping = {
        # Phase 1: Instrumentation & Proposal (T1.1-T1.25)
        'A': 'T1.1',    # Draft timeline v1
        'B': 'T1.2',    # Initial proposal skeleton  
        'C': 'T1.3',    # Submit proposal outline
        'D': 'T1.4',    # Define proposal committee
        'F': 'T1.5',    # Expand proposal draft
        'G': 'T1.6',    # Confirm exam date
        'H': 'T1.13',   # Align seed laser
        'I': 'T1.14',   # Align amplifier
        'J': 'T1.15',   # Check pulse compression
        'K': 'T1.16',   # Calibrate microscope
        'L': 'T1.17',   # Laser system ready
        'M': 'T2.1',    # Plan imaging cohort
        'N': 'T1.21',   # Annual progress review
        'O1': 'T1.22',  # Complete committee paperwork and Program of Work
        'O2': 'T1.23',  # Reserve exam room and submit final paperwork
        'P': 'T2.2',    # Design AAV vectors
        'Q': 'T2.3',    # AAV vectors ready
        'R1': 'T1.7',   # Draft Specific Aims and Research Strategy
        'R2': 'T1.8',   # Draft Methods and Timeline sections
        'R3': 'T1.9',   # Finalize proposal draft and formatting
        'S': 'T1.10',   # Send proposal to committee
        'T': 'T1.11',   # Prepare presentation
        'U': 'T1.M1',   # PhD Proposal Exam (milestone)
        'V': 'T1.12',   # Address committee feedback
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
        'BB': 'T3.5',   # Prepare conference presentation
        'BC': 'T3.3',   # Draft second manuscript
        'BD': 'T3.4',   # Submit second manuscript
        'BE': 'T4.8',   # Annual progress review
        'BG': 'T4.1',   # PhD Dissertation & Defense (mapped to T4.1 for Introduction)
        'BI': 'T4.1',   # Draft Introduction and Literature Review
        'BJ': 'T4.2',   # Draft Methods and Results chapters (Aims 1-3)
        'BK': 'T4.3',   # Draft Discussion and Conclusions
        'BN': 'T4.4',   # Dissertation draft complete
        'BS': 'T4.5',   # PhD Defense
        'BT': 'T4.6',   # Revise dissertation
        'BU': 'T4.7',   # Submit dissertation
        'BV': 'T4.9',   # Complete TA requirement
        'BW': 'T1.24',  # Update committee membership
        'BX': 'T1.25',  # Advance to candidacy
        'BY1': 'T1.25', # Maintain continuous registration - Fall 2025
        'BY2': 'T4.10', # Maintain continuous registration - Spring 2026
        'BY3': 'T4.11', # Maintain continuous registration - Fall 2026
        'BY4': 'T4.12', # Maintain continuous registration - Spring 2027
        'BY5': 'T4.13', # Maintain continuous registration - Summer 2027
        'BZ': 'T4.14',  # Apply for graduation
        'CA': 'T4.15',  # Request final oral exam
        'CB': 'T4.16',  # SPIE chapter activities
        'CC1': 'T2.27', # Equipment maintenance log - Q1 2026
        'CC2': 'T2.28', # Equipment maintenance log - Q2 2026
        'CC3': 'T2.29', # Equipment maintenance log - Q3 2026
        'CC4': 'T2.30', # Equipment maintenance log - Q4 2026
        'CC5': 'T2.31', # Equipment maintenance log - Q1 2027 (mapped to T2.31)
        'CC6': 'T2.32', # Equipment maintenance log - Q2 2027 (mapped to T2.32)
        'CC7': 'T2.33', # Equipment maintenance log - Q3 2027 (mapped to T2.33)
        'CC8': 'T2.34', # Equipment maintenance log - Q4 2027 (mapped to T2.34)
        'CD1': 'T2.31', # Data backup system implementation
        'CD2': 'T2.32', # Data backup system maintenance
        'CE': 'T2.33',  # Surgical training
        'CF': 'T2.34',  # Lab culture responsibilities
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
    
    return task_mapping

if __name__ == "__main__":
    mapping = create_task_mapping()
    print(f"\nTotal mappings: {len(mapping)}")
