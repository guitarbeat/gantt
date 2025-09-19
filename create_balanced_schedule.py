#!/usr/bin/env python3
"""
Create a more balanced and manageable PhD schedule with work distributed evenly.
"""

import csv
from datetime import datetime, timedelta
from collections import defaultdict

def create_balanced_schedule(input_file, output_file):
    """Create a more balanced schedule by redistributing tasks."""
    
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
    
    # Strategy for rebalancing:
    # 1. Start dissertation work earlier (2026-Q3 instead of 2026-Q4)
    # 2. Distribute heavy imaging work more evenly
    # 3. Add buffer periods for critical milestones
    # 4. Create more parallel work opportunities
    # 5. Extend some tasks to reduce daily intensity
    
    improved_tasks = []
    
    for task in tasks:
        new_task = task.copy()
        
        # * Rebalance dissertation work - start earlier and spread out
        if task['Category'] == 'DISSERTATION':
            if task['Task ID'] == 'BI':  # Draft Introduction
                new_task['Start Date'] = datetime(2026, 9, 1)  # Start 3 months earlier
                new_task['Due Date'] = datetime(2027, 1, 15)   # Extend duration
            elif task['Task ID'] == 'BJ':  # Draft Aim 1 chapter
                new_task['Start Date'] = datetime(2026, 10, 1)  # Start earlier
                new_task['Due Date'] = datetime(2027, 2, 28)    # Extend duration
            elif task['Task ID'] == 'BK':  # Draft Aim 2 chapter
                new_task['Start Date'] = datetime(2026, 11, 1)  # Start earlier
                new_task['Due Date'] = datetime(2027, 3, 31)    # Extend duration
            elif task['Task ID'] == 'BL':  # Draft Aim 3 chapter
                new_task['Start Date'] = datetime(2026, 12, 1)  # Start earlier
                new_task['Due Date'] = datetime(2027, 4, 30)    # Extend duration
            elif task['Task ID'] == 'BM':  # Draft Conclusions
                new_task['Start Date'] = datetime(2027, 1, 1)   # Start earlier
                new_task['Due Date'] = datetime(2027, 5, 31)    # Keep same end
            elif task['Task ID'] == 'BN':  # Dissertation draft complete
                new_task['Start Date'] = datetime(2027, 5, 1)   # Start earlier
                new_task['Due Date'] = datetime(2027, 5, 31)    # Shorter duration
            elif task['Task ID'] == 'BS':  # PhD Defense
                new_task['Start Date'] = datetime(2027, 6, 15)  # Move earlier
                new_task['Due Date'] = datetime(2027, 6, 17)    # Keep same
            elif task['Task ID'] == 'BT':  # Revise dissertation
                new_task['Start Date'] = datetime(2027, 6, 18)  # Start immediately after defense
                new_task['Due Date'] = datetime(2027, 7, 15)    # Extend revision time
            elif task['Task ID'] == 'BU':  # Submit dissertation
                new_task['Start Date'] = datetime(2027, 7, 16)  # Start after revisions
                new_task['Due Date'] = datetime(2027, 7, 31)    # Earlier submission
        
        # * Rebalance imaging work - spread out heavy periods
        elif task['Category'] == 'IMAGING':
            if task['Task ID'] in ['AJ', 'AK', 'AM']:  # Parallel development tasks
                # Start these earlier to spread workload
                if task['Task ID'] == 'AJ':  # Develop U-Net pipeline
                    new_task['Start Date'] = datetime(2026, 3, 15)  # Start earlier
                    new_task['Due Date'] = datetime(2026, 5, 15)    # Extend duration
                elif task['Task ID'] == 'AK':  # Optimize imaging systems
                    new_task['Start Date'] = datetime(2026, 3, 20)  # Start earlier
                    new_task['Due Date'] = datetime(2026, 6, 15)    # Extend duration
                elif task['Task ID'] == 'AM':  # Order enhanced AAV
                    new_task['Start Date'] = datetime(2026, 3, 25)  # Start earlier
                    new_task['Due Date'] = datetime(2026, 6, 15)    # Extend duration
            
            # * Spread out stroke imaging sessions
            elif task['Task ID'] in ['AT', 'AU', 'AV', 'AW']:  # Stroke imaging phases
                if task['Task ID'] == 'AT':  # Acute-phase imaging
                    new_task['Start Date'] = datetime(2026, 6, 20)  # Slightly later
                    new_task['Due Date'] = datetime(2026, 6, 25)    # Extend duration
                elif task['Task ID'] == 'AU':  # Transition-phase imaging
                    new_task['Start Date'] = datetime(2026, 7, 10)  # More spacing
                    new_task['Due Date'] = datetime(2026, 7, 15)    # Extend duration
                elif task['Task ID'] == 'AV':  # Stabilization-phase imaging
                    new_task['Start Date'] = datetime(2026, 8, 5)   # More spacing
                    new_task['Due Date'] = datetime(2026, 8, 10)    # Extend duration
                elif task['Task ID'] == 'AW':  # Extended chronic imaging
                    new_task['Start Date'] = datetime(2026, 9, 1)   # More spacing
                    new_task['Due Date'] = datetime(2026, 9, 5)     # Extend duration
        
        # * Add buffer periods for critical milestones
        elif 'MILESTONE' in task['Description']:
            # Add 1-2 week buffer before major milestones
            if task['Task ID'] == 'G':  # Confirm exam date
                new_task['Due Date'] = task['Due Date'] + timedelta(days=7)
            elif task['Task ID'] == 'S':  # Send proposal to committee
                new_task['Due Date'] = task['Due Date'] + timedelta(days=10)
            elif task['Task ID'] == 'U':  # PhD Proposal Exam
                new_task['Start Date'] = task['Start Date'] - timedelta(days=7)
                new_task['Due Date'] = task['Due Date'] + timedelta(days=3)
        
        # * Extend some short tasks to reduce daily intensity
        elif task['Duration'] <= 3 and task['Category'] in ['IMAGING', 'PUBLICATION']:
            new_task['Duration'] = max(5, task['Duration'] + 2)  # Add 2 days minimum
            new_task['Due Date'] = new_task['Start Date'] + timedelta(days=new_task['Duration'] - 1)
        
        # * Create more parallel work opportunities
        elif task['Category'] == 'PUBLICATION' and task['Task ID'] in ['AP', 'BC']:
            # Start writing papers earlier to parallel with data collection
            if task['Task ID'] == 'AP':  # Draft methodology paper
                new_task['Start Date'] = datetime(2026, 3, 1)  # Start much earlier
                new_task['Due Date'] = datetime(2026, 6, 30)   # Extend duration
            elif task['Task ID'] == 'BC':  # Draft second manuscript
                new_task['Start Date'] = datetime(2026, 9, 1)  # Start earlier
                new_task['Due Date'] = datetime(2026, 12, 31)  # Extend duration
        
        # * Add buffer periods in light quarters
        elif task['Category'] == 'ADMIN' and task['Task ID'] in ['N', 'BE']:
            # Extend annual progress reviews to fill light periods
            new_task['Duration'] = 14  # 2 weeks instead of 1 week
            new_task['Due Date'] = new_task['Start Date'] + timedelta(days=13)
        
        improved_tasks.append(new_task)
    
    # * Add new buffer and preparation tasks
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
            'Task Name': 'Buffer period - Imaging preparation',
            'Parent Task ID': '',
            'Category': 'IMAGING',
            'Start Date': datetime(2026, 1, 15),
            'Due Date': datetime(2026, 1, 31),
            'Dependencies': 'Q',
            'Description': 'Buffer time for imaging protocol refinement and equipment testing.'
        },
        {
            'Task ID': 'BUF3',
            'Task Name': 'Buffer period - Data analysis',
            'Parent Task ID': '',
            'Category': 'IMAGING',
            'Start Date': datetime(2026, 8, 15),
            'Due Date': datetime(2026, 8, 31),
            'Dependencies': 'AY',
            'Description': 'Buffer time for data processing and preliminary analysis.'
        },
        {
            'Task ID': 'BUF4',
            'Task Name': 'Buffer period - Dissertation preparation',
            'Parent Task ID': '',
            'Category': 'DISSERTATION',
            'Start Date': datetime(2026, 8, 1),
            'Due Date': datetime(2026, 8, 31),
            'Dependencies': 'V',
            'Description': 'Buffer time for dissertation planning and literature review.'
        }
    ]
    
    # Add buffer tasks
    for buf_task in buffer_tasks:
        improved_tasks.append(buf_task)
    
    # * Add parallel development tasks to fill light periods
    parallel_tasks = [
        {
            'Task ID': 'PAR1',
            'Task Name': 'Develop data analysis scripts',
            'Parent Task ID': '',
            'Category': 'IMAGING',
            'Start Date': datetime(2025, 10, 1),
            'Due Date': datetime(2025, 10, 31),
            'Dependencies': 'L',
            'Description': 'Develop and test data analysis scripts in parallel with imaging setup.'
        },
        {
            'Task ID': 'PAR2',
            'Task Name': 'Literature review - stroke models',
            'Parent Task ID': '',
            'Category': 'IMAGING',
            'Start Date': datetime(2026, 1, 1),
            'Due Date': datetime(2026, 2, 28),
            'Dependencies': 'M',
            'Description': 'Comprehensive literature review of stroke imaging models and protocols.'
        },
        {
            'Task ID': 'PAR3',
            'Task Name': 'Develop presentation templates',
            'Parent Task ID': '',
            'Category': 'PUBLICATION',
            'Start Date': datetime(2026, 6, 1),
            'Due Date': datetime(2026, 6, 30),
            'Dependencies': 'AP',
            'Description': 'Create standardized presentation templates for conference talks and posters.'
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
    
    print(f"Improved schedule written to {output_file}")
    print(f"Total tasks: {len(improved_tasks)}")
    print(f"Added {len(buffer_tasks)} buffer tasks and {len(parallel_tasks)} parallel tasks")

if __name__ == "__main__":
    create_balanced_schedule(
        '/Users/aaron/Downloads/gantt/input/data.cleaned.csv',
        '/Users/aaron/Downloads/gantt/input/data.balanced.csv'
    )
