#!/usr/bin/env python3
"""
Final dependency fix after milestone/surgery task corrections.
"""

import csv
from datetime import datetime, timedelta

def parse_date(date_str):
    """Parse date string in YYYY-MM-DD format."""
    return datetime.strptime(date_str, '%Y-%m-%d')

def format_date(date_obj):
    """Format datetime object to YYYY-MM-DD string."""
    return date_obj.strftime('%Y-%m-%d')

def parse_dependencies(deps_str):
    """Parse dependencies string and return list of task IDs."""
    if not deps_str or deps_str.strip() == '':
        return []
    deps_str = deps_str.strip('"')
    return [dep.strip() for dep in deps_str.split(',') if dep.strip()]

def fix_dependencies_final(input_file, output_file):
    """
    Fix dependencies after milestone/surgery corrections.
    """
    tasks = {}
    task_list = []
    
    # Read all tasks
    with open(input_file, 'r', newline='', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            if not row['Task ID']:
                continue
            
            task_id = row['Task ID']
            tasks[task_id] = {
                'data': row,
                'start_date': parse_date(row['Start Date']),
                'due_date': parse_date(row['Due Date']),
                'dependencies': parse_dependencies(row['Dependencies'])
            }
            task_list.append(task_id)
    
    print("=== FIXING DEPENDENCIES AFTER MILESTONE CORRECTIONS ===")
    
    # Find and fix dependency issues
    issues_found = []
    for task_id, task_info in tasks.items():
        for dep_id in task_info['dependencies']:
            if dep_id in tasks:
                dep_finish_date = tasks[dep_id]['due_date']
                task_start_date = task_info['start_date']
                
                if task_start_date <= dep_finish_date:
                    issues_found.append({
                        'task_id': task_id,
                        'task_name': task_info['data']['Task Name'],
                        'dep_id': dep_id,
                        'dep_name': tasks[dep_id]['data']['Task Name'],
                        'dep_finish': dep_finish_date,
                        'task_start': task_start_date
                    })
    
    if not issues_found:
        print("✅ No dependency issues found!")
        return []
    
    print(f"Found {len(issues_found)} dependency issues to fix:")
    for issue in issues_found:
        print(f"  - {issue['task_id']} ({issue['task_name']}) starts {issue['task_start'].strftime('%Y-%m-%d')} but depends on {issue['dep_id']} ({issue['dep_name']}) which finishes {issue['dep_finish'].strftime('%Y-%m-%d')}")
    
    # Fix issues
    fixed_tasks = []
    for issue in issues_found:
        task_id = issue['task_id']
        dep_finish_date = issue['dep_finish']
        
        # Move task start date to day after dependency finishes
        new_start_date = dep_finish_date + timedelta(days=1)
        
        # Calculate new due date maintaining the same duration
        original_duration = (tasks[task_id]['due_date'] - tasks[task_id]['start_date']).days + 1
        new_due_date = new_start_date + timedelta(days=original_duration - 1)
        
        # Update task data
        tasks[task_id]['start_date'] = new_start_date
        tasks[task_id]['due_date'] = new_due_date
        tasks[task_id]['data']['Start Date'] = format_date(new_start_date)
        tasks[task_id]['data']['Due Date'] = format_date(new_due_date)
        
        fixed_tasks.append({
            'task_id': task_id,
            'task_name': tasks[task_id]['data']['Task Name'],
            'old_start': issue['task_start'],
            'new_start': new_start_date,
            'old_due': tasks[task_id]['due_date'],
            'new_due': new_due_date
        })
        
        print(f"Fixed {task_id}: {tasks[task_id]['data']['Task Name']}")
        print(f"  Start: {issue['task_start'].strftime('%Y-%m-%d')} → {new_start_date.strftime('%Y-%m-%d')}")
        print(f"  Due: {tasks[task_id]['due_date'].strftime('%Y-%m-%d')} → {new_due_date.strftime('%Y-%m-%d')}")
    
    # Write the fixed CSV
    with open(output_file, 'w', newline='', encoding='utf-8') as csvfile:
        fieldnames = ['Task ID', 'Task Name', 'Parent Task ID', 'Category', 'Start Date', 'Due Date', 'Dependencies', 'Description']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        
        for task_id in task_list:
            writer.writerow(tasks[task_id]['data'])
    
    print(f"\nFixed {len(fixed_tasks)} tasks with dependency issues")
    return fixed_tasks

if __name__ == "__main__":
    input_file = "/Users/aaron/Downloads/gantt/input/data.cleaned.csv"
    output_file = "/Users/aaron/Downloads/gantt/input/data.cleaned.csv"
    
    print("Fixing dependencies after milestone corrections...")
    fixed_tasks = fix_dependencies_final(input_file, output_file)
    
    if fixed_tasks:
        print(f"\nCompleted! Fixed {len(fixed_tasks)} tasks.")
    else:
        print("\nNo dependency issues found - all tasks are properly scheduled!")
