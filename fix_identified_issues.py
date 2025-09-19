#!/usr/bin/env python3
"""
Script to fix identified issues in the CSV data.
"""

import csv
from datetime import datetime, timedelta

def parse_date(date_str):
    """Parse date string in YYYY-MM-DD format."""
    return datetime.strptime(date_str, '%Y-%m-%d')

def format_date(date_obj):
    """Format datetime object to YYYY-MM-DD string."""
    return date_obj.strftime('%Y-%m-%d')

def fix_issues(csv_file):
    """
    Fix identified issues in the CSV data.
    """
    tasks = []
    fixes_applied = []
    
    # Read the CSV file
    with open(csv_file, 'r', newline='', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            if not row['Task ID']:  # Skip empty rows
                continue
            
            task_id = row['Task ID']
            task_name = row['Task Name']
            start_date = parse_date(row['Start Date'])
            due_date = parse_date(row['Due Date'])
            description = row['Description']
            
            # Check if this is a milestone task
            is_milestone = 'MILESTONE' in description
            
            # Check if this is a surgery task
            is_surgery = 'surgery' in task_name.lower()
            
            # Check if this is the parent dissertation task
            is_parent_dissertation = (task_id == 'BG' and 'PhD Dissertation & Defense' in task_name)
            
            original_duration = (due_date - start_date).days + 1
            
            # Apply fixes
            if (is_milestone or is_surgery) and original_duration > 1:
                # Fix milestone and surgery tasks to be single day
                new_due_date = start_date  # Same day start and finish
                row['Due Date'] = format_date(new_due_date)
                fixes_applied.append({
                    'task_id': task_id,
                    'task_name': task_name,
                    'type': 'Milestone/Surgery',
                    'old_duration': original_duration,
                    'new_duration': 1,
                    'old_due': due_date,
                    'new_due': new_due_date
                })
            
            elif is_parent_dissertation and original_duration > 90:
                # Keep parent dissertation task as is - it's meant to span the entire period
                # Just note it for reporting
                fixes_applied.append({
                    'task_id': task_id,
                    'task_name': task_name,
                    'type': 'Parent Task (kept as-is)',
                    'old_duration': original_duration,
                    'new_duration': original_duration,
                    'old_due': due_date,
                    'new_due': due_date
                })
            
            tasks.append(row)
    
    # Write the fixed CSV
    with open(csv_file, 'w', newline='', encoding='utf-8') as csvfile:
        fieldnames = ['Task ID', 'Task Name', 'Parent Task ID', 'Category', 'Start Date', 'Due Date', 'Dependencies', 'Description']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(tasks)
    
    return fixes_applied

if __name__ == "__main__":
    csv_file = "/Users/aaron/Downloads/gantt/input/data.cleaned.csv"
    
    print("=== FIXING IDENTIFIED ISSUES ===")
    fixes = fix_issues(csv_file)
    
    print(f"Applied {len(fixes)} fixes:")
    for fix in fixes:
        if fix['type'] == 'Parent Task (kept as-is)':
            print(f"  {fix['task_id']}: {fix['task_name']} - {fix['type']} ({fix['old_duration']} days)")
        else:
            print(f"  {fix['task_id']}: {fix['task_name']} - {fix['type']}")
            print(f"    Duration: {fix['old_duration']} → {fix['new_duration']} days")
            print(f"    Due date: {fix['old_due'].strftime('%Y-%m-%d')} → {fix['new_due'].strftime('%Y-%m-%d')}")
    
    print(f"\nCompleted! Applied {len(fixes)} fixes to the CSV data.")
