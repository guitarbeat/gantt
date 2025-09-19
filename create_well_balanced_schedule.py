#!/usr/bin/env python3
"""
Create a well-balanced PhD schedule with work distributed evenly throughout the timeline.
"""

import csv
from datetime import datetime, timedelta
from collections import defaultdict

def create_well_balanced_schedule(input_file, output_file):
    """Create a well-balanced schedule by aggressively redistributing tasks."""
    
    # Read the original CSV
    tasks = []
    with open(input_file, 'r') as f:
        reader = csv.DictReader(f)
        for row in reader:
            if row['Start Date'] and row['Due Date']:
                tasks.append(row)
    
    # Convert dates
    for task in tasks:
        task['Start Date'] = datetime.strptime(task['Start Date'], '%Y-%m-%d')
        task['Due Date'] = datetime.strptime(task['Due Date'], '%Y-%m-%d')
        task['Duration'] = (task['Due Date'] - task['Start Date']).days + 1
    
    # Target: ~150-200 task-days per quarter (manageable workload)
    target_quarterly_workload = 180
    
    improved_tasks = []
    
    # * Phase 1: Redistribute dissertation work across 18 months (2026-Q3 to 2027-Q2)
    dissertation_tasks = [t for t in tasks if t['Category'] == 'DISSERTATION']
    other_tasks = [t for t in tasks if t['Category'] != 'DISSERTATION']
    
    # Spread dissertation work evenly from 2026-Q3 to 2027-Q2
    dissertation_start = datetime(2026, 7, 1)
    dissertation_end = datetime(2027, 6, 30)
    total_dissertation_days = 365
    
    # Redistribute dissertation tasks
    for i, task in enumerate(dissertation_tasks):
        new_task = task.copy()
        
        if task['Task ID'] == 'BI':  # Draft Introduction
            new_task['Start Date'] = datetime(2026, 7, 1)
            new_task['Due Date'] = datetime(2026, 10, 31)
        elif task['Task ID'] == 'BJ':  # Draft Aim 1 chapter
            new_task['Start Date'] = datetime(2026, 8, 1)
            new_task['Due Date'] = datetime(2026, 11, 30)
        elif task['Task ID'] == 'BK':  # Draft Aim 2 chapter
            new_task['Start Date'] = datetime(2026, 9, 1)
            new_task['Due Date'] = datetime(2026, 12, 31)
        elif task['Task ID'] == 'BL':  # Draft Aim 3 chapter
            new_task['Start Date'] = datetime(2026, 10, 1)
            new_task['Due Date'] = datetime(2027, 1, 31)
        elif task['Task ID'] == 'BM':  # Draft Conclusions
            new_task['Start Date'] = datetime(2026, 11, 1)
            new_task['Due Date'] = datetime(2027, 2, 28)
        elif task['Task ID'] == 'BN':  # Dissertation draft complete
            new_task['Start Date'] = datetime(2027, 2, 1)
            new_task['Due Date'] = datetime(2027, 3, 15)
        elif task['Task ID'] == 'BS':  # PhD Defense
            new_task['Start Date'] = datetime(2027, 4, 15)
            new_task['Due Date'] = datetime(2027, 4, 17)
        elif task['Task ID'] == 'BT':  # Revise dissertation
            new_task['Start Date'] = datetime(2027, 4, 18)
            new_task['Due Date'] = datetime(2027, 5, 31)
        elif task['Task ID'] == 'BU':  # Submit dissertation
            new_task['Start Date'] = datetime(2027, 6, 1)
            new_task['Due Date'] = datetime(2027, 6, 15)
        
        improved_tasks.append(new_task)
    
    # * Phase 2: Redistribute imaging work more evenly
    imaging_tasks = [t for t in other_tasks if t['Category'] == 'IMAGING']
    non_imaging_tasks = [t for t in other_tasks if t['Category'] != 'IMAGING']
    
    for task in imaging_tasks:
        new_task = task.copy()
        
        # * Spread pilot imaging work across 2026-Q1 to 2026-Q2
        if task['Task ID'] in ['W', 'Z', 'AA', 'AB', 'AC', 'AD', 'AE', 'AF', 'AG', 'AH']:
            if task['Task ID'] == 'W':  # Cranial window surgery #1
                new_task['Start Date'] = datetime(2026, 1, 15)
                new_task['Due Date'] = datetime(2026, 1, 20)
            elif task['Task ID'] == 'Z':  # Post-op recovery #1
                new_task['Start Date'] = datetime(2026, 1, 21)
                new_task['Due Date'] = datetime(2026, 1, 28)
            elif task['Task ID'] == 'AA':  # Cranial window surgery #2
                new_task['Start Date'] = datetime(2026, 2, 1)
                new_task['Due Date'] = datetime(2026, 2, 6)
            elif task['Task ID'] == 'AB':  # Post-op recovery #2
                new_task['Start Date'] = datetime(2026, 2, 7)
                new_task['Due Date'] = datetime(2026, 2, 14)
            elif task['Task ID'] == 'AC':  # Cranial window surgery #3
                new_task['Start Date'] = datetime(2026, 2, 15)
                new_task['Due Date'] = datetime(2026, 2, 20)
            elif task['Task ID'] == 'AD':  # Post-op recovery #3
                new_task['Start Date'] = datetime(2026, 2, 21)
                new_task['Due Date'] = datetime(2026, 2, 28)
            elif task['Task ID'] == 'AE':  # Pilot imaging session #1
                new_task['Start Date'] = datetime(2026, 3, 1)
                new_task['Due Date'] = datetime(2026, 3, 5)
            elif task['Task ID'] == 'AF':  # Pilot imaging session #2
                new_task['Start Date'] = datetime(2026, 3, 6)
                new_task['Due Date'] = datetime(2026, 3, 10)
            elif task['Task ID'] == 'AG':  # Pilot imaging session #3
                new_task['Start Date'] = datetime(2026, 3, 11)
                new_task['Due Date'] = datetime(2026, 3, 15)
            elif task['Task ID'] == 'AH':  # Pilot datasets complete
                new_task['Start Date'] = datetime(2026, 3, 16)
                new_task['Due Date'] = datetime(2026, 3, 18)
        
        # * Spread development work across 2026-Q2 to 2026-Q3
        elif task['Task ID'] in ['AJ', 'AK', 'AM', 'AI']:
            if task['Task ID'] == 'AI':  # Process pilot data
                new_task['Start Date'] = datetime(2026, 3, 19)
                new_task['Due Date'] = datetime(2026, 4, 5)
            elif task['Task ID'] == 'AJ':  # Develop U-Net pipeline
                new_task['Start Date'] = datetime(2026, 4, 1)
                new_task['Due Date'] = datetime(2026, 6, 30)
            elif task['Task ID'] == 'AK':  # Optimize imaging systems
                new_task['Start Date'] = datetime(2026, 4, 15)
                new_task['Due Date'] = datetime(2026, 7, 15)
            elif task['Task ID'] == 'AM':  # Order enhanced AAV
                new_task['Start Date'] = datetime(2026, 5, 1)
                new_task['Due Date'] = datetime(2026, 7, 31)
        
        # * Spread stroke imaging across 2026-Q3 to 2026-Q4
        elif task['Task ID'] in ['AR', 'AS', 'AT', 'AU', 'AV', 'AW', 'AX', 'AY', 'AZ', 'BA']:
            if task['Task ID'] == 'AR':  # Establish stroke protocol
                new_task['Start Date'] = datetime(2026, 7, 1)
                new_task['Due Date'] = datetime(2026, 7, 15)
            elif task['Task ID'] == 'AS':  # Induce stroke
                new_task['Start Date'] = datetime(2026, 7, 16)
                new_task['Due Date'] = datetime(2026, 7, 25)
            elif task['Task ID'] == 'AT':  # Acute-phase imaging
                new_task['Start Date'] = datetime(2026, 8, 1)
                new_task['Due Date'] = datetime(2026, 8, 10)
            elif task['Task ID'] == 'AU':  # Transition-phase imaging
                new_task['Start Date'] = datetime(2026, 8, 15)
                new_task['Due Date'] = datetime(2026, 8, 25)
            elif task['Task ID'] == 'AV':  # Stabilization-phase imaging
                new_task['Start Date'] = datetime(2026, 9, 1)
                new_task['Due Date'] = datetime(2026, 9, 10)
            elif task['Task ID'] == 'AW':  # Extended chronic imaging
                new_task['Start Date'] = datetime(2026, 9, 15)
                new_task['Due Date'] = datetime(2026, 9, 25)
            elif task['Task ID'] == 'AX':  # Refine ML pipeline
                new_task['Start Date'] = datetime(2026, 8, 1)
                new_task['Due Date'] = datetime(2026, 10, 31)
            elif task['Task ID'] == 'AY':  # Stroke data complete
                new_task['Start Date'] = datetime(2026, 9, 26)
                new_task['Due Date'] = datetime(2026, 9, 30)
            elif task['Task ID'] == 'AZ':  # Integrate flow data
                new_task['Start Date'] = datetime(2026, 10, 1)
                new_task['Due Date'] = datetime(2026, 11, 15)
            elif task['Task ID'] == 'BA':  # Analyze neurovascular coupling
                new_task['Start Date'] = datetime(2026, 11, 16)
                new_task['Due Date'] = datetime(2026, 12, 31)
        
        # * Keep other imaging tasks as is
        else:
            pass
        
        improved_tasks.append(new_task)
    
    # * Phase 3: Redistribute publication work across 2026-Q2 to 2026-Q4
    publication_tasks = [t for t in non_imaging_tasks if t['Category'] == 'PUBLICATION']
    other_tasks = [t for t in non_imaging_tasks if t['Category'] != 'PUBLICATION']
    
    for task in publication_tasks:
        new_task = task.copy()
        
        if task['Task ID'] == 'AP':  # Draft methodology paper
            new_task['Start Date'] = datetime(2026, 4, 1)
            new_task['Due Date'] = datetime(2026, 7, 31)
        elif task['Task ID'] == 'AQ':  # Submit methodology paper
            new_task['Start Date'] = datetime(2026, 8, 1)
            new_task['Due Date'] = datetime(2026, 8, 15)
        elif task['Task ID'] == 'BB':  # Prepare conference presentation
            new_task['Start Date'] = datetime(2026, 10, 1)
            new_task['Due Date'] = datetime(2026, 11, 15)
        elif task['Task ID'] == 'BC':  # Draft second manuscript
            new_task['Start Date'] = datetime(2026, 9, 1)
            new_task['Due Date'] = datetime(2026, 12, 15)
        elif task['Task ID'] == 'BD':  # Submit second manuscript
            new_task['Start Date'] = datetime(2026, 12, 16)
            new_task['Due Date'] = datetime(2026, 12, 31)
        
        improved_tasks.append(new_task)
    
    # * Phase 4: Keep other tasks mostly as is, but add some buffer
    for task in other_tasks:
        new_task = task.copy()
        
        # * Add buffer to critical milestones
        if 'MILESTONE' in task['Description']:
            if task['Task ID'] == 'G':  # Confirm exam date
                new_task['Due Date'] = task['Due Date'] + timedelta(days=7)
            elif task['Task ID'] == 'S':  # Send proposal to committee
                new_task['Due Date'] = task['Due Date'] + timedelta(days=10)
            elif task['Task ID'] == 'U':  # PhD Proposal Exam
                new_task['Start Date'] = task['Start Date'] - timedelta(days=7)
                new_task['Due Date'] = task['Due Date'] + timedelta(days=3)
        
        # * Extend short tasks to reduce daily intensity
        elif task['Duration'] <= 3 and task['Category'] in ['IMAGING', 'PUBLICATION']:
            new_task['Duration'] = max(5, task['Duration'] + 2)
            new_task['Due Date'] = new_task['Start Date'] + timedelta(days=new_task['Duration'] - 1)
        
        improved_tasks.append(new_task)
    
    # * Phase 5: Add strategic buffer and parallel tasks
    buffer_tasks = [
        {
            'Task ID': 'BUF1',
            'Task Name': 'Buffer period - Proposal preparation',
            'Parent Task ID': '',
            'Category': 'ADMIN',
            'Start Date': datetime(2025, 11, 1),
            'Due Date': datetime(2025, 11, 15),
            'Dependencies': 'F',
            'Description': 'Buffer time for proposal revisions and committee coordination.'
        },
        {
            'Task ID': 'BUF2',
            'Task Name': 'Buffer period - Early dissertation prep',
            'Parent Task ID': '',
            'Category': 'DISSERTATION',
            'Start Date': datetime(2026, 1, 1),
            'Due Date': datetime(2026, 2, 28),
            'Dependencies': 'V',
            'Description': 'Early literature review and dissertation planning to reduce final year pressure.'
        },
        {
            'Task ID': 'BUF3',
            'Task Name': 'Buffer period - Data analysis',
            'Parent Task ID': '',
            'Category': 'IMAGING',
            'Start Date': datetime(2026, 6, 1),
            'Due Date': datetime(2026, 6, 30),
            'Dependencies': 'AH',
            'Description': 'Buffer time for data processing and preliminary analysis.'
        },
        {
            'Task ID': 'BUF4',
            'Task Name': 'Buffer period - Final dissertation prep',
            'Parent Task ID': '',
            'Category': 'DISSERTATION',
            'Start Date': datetime(2027, 1, 1),
            'Due Date': datetime(2027, 1, 31),
            'Dependencies': 'BM',
            'Description': 'Final buffer period for dissertation revisions and preparation.'
        }
    ]
    
    # Add buffer tasks
    for buf_task in buffer_tasks:
        improved_tasks.append(buf_task)
    
    # * Phase 6: Add parallel development tasks to fill light periods
    parallel_tasks = [
        {
            'Task ID': 'PAR1',
            'Task Name': 'Develop data analysis scripts',
            'Parent Task ID': '',
            'Category': 'IMAGING',
            'Start Date': datetime(2025, 10, 1),
            'Due Date': datetime(2025, 11, 30),
            'Dependencies': 'L',
            'Description': 'Develop and test data analysis scripts in parallel with imaging setup.'
        },
        {
            'Task ID': 'PAR2',
            'Task Name': 'Literature review - stroke models',
            'Parent Task ID': '',
            'Category': 'IMAGING',
            'Start Date': datetime(2026, 1, 1),
            'Due Date': datetime(2026, 3, 31),
            'Dependencies': 'M',
            'Description': 'Comprehensive literature review of stroke imaging models and protocols.'
        },
        {
            'Task ID': 'PAR3',
            'Task Name': 'Develop presentation templates',
            'Parent Task ID': '',
            'Category': 'PUBLICATION',
            'Start Date': datetime(2026, 5, 1),
            'Due Date': datetime(2026, 6, 30),
            'Dependencies': 'AP',
            'Description': 'Create standardized presentation templates for conference talks and posters.'
        },
        {
            'Task ID': 'PAR4',
            'Task Name': 'Prepare conference abstracts',
            'Parent Task ID': '',
            'Category': 'PUBLICATION',
            'Start Date': datetime(2026, 8, 1),
            'Due Date': datetime(2026, 9, 30),
            'Dependencies': 'AH',
            'Description': 'Prepare conference abstracts and presentation materials in parallel with data collection.'
        }
    ]
    
    # Add parallel tasks
    for par_task in parallel_tasks:
        improved_tasks.append(par_task)
    
    # Write the improved schedule
    with open(output_file, 'w', newline='') as f:
        fieldnames = ['Task ID', 'Task Name', 'Parent Task ID', 'Category', 'Start Date', 'Due Date', 'Dependencies', 'Description']
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        
        for task in improved_tasks:
            # Convert dates back to strings and remove Duration field
            row = {k: v for k, v in task.items() if k != 'Duration'}
            row['Start Date'] = task['Start Date'].strftime('%Y-%m-%d')
            row['Due Date'] = task['Due Date'].strftime('%Y-%m-%d')
            writer.writerow(row)
    
    print(f"Well-balanced schedule written to {output_file}")
    print(f"Total tasks: {len(improved_tasks)}")
    print(f"Added {len(buffer_tasks)} buffer tasks and {len(parallel_tasks)} parallel tasks")

if __name__ == "__main__":
    create_well_balanced_schedule(
        '/Users/aaron/Downloads/gantt/input/data.cleaned.csv',
        '/Users/aaron/Downloads/gantt/input/data.well_balanced.csv'
    )
