#!/usr/bin/env python3
"""
Script to analyze CSV data line by line for potential issues.
"""

import csv
from datetime import datetime, timedelta

def parse_date(date_str):
    """Parse date string in YYYY-MM-DD format."""
    return datetime.strptime(date_str, '%Y-%m-%d')

def parse_dependencies(deps_str):
    """Parse dependencies string and return list of task IDs."""
    if not deps_str or deps_str.strip() == '':
        return []
    deps_str = deps_str.strip('"')
    return [dep.strip() for dep in deps_str.split(',') if dep.strip()]

def analyze_csv_issues(csv_file):
    """
    Analyze CSV data line by line for potential issues.
    """
    tasks = {}
    issues = []
    
    # Read all tasks
    with open(csv_file, 'r', newline='', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)
        for row_num, row in enumerate(reader, start=2):  # Start at 2 for header
            if not row['Task ID']:  # Skip empty rows
                continue
            
            task_id = row['Task ID']
            task_name = row['Task Name']
            start_date = parse_date(row['Start Date'])
            due_date = parse_date(row['Due Date'])
            dependencies = parse_dependencies(row['Dependencies'])
            
            tasks[task_id] = {
                'row_num': row_num,
                'name': task_name,
                'start_date': start_date,
                'due_date': due_date,
                'dependencies': dependencies,
                'category': row['Category'],
                'parent': row['Parent Task ID']
            }
    
    print("=== LINE-BY-LINE ANALYSIS ===")
    
    # Check for various issues
    for task_id, task_info in tasks.items():
        row_num = task_info['row_num']
        issues_found = []
        
        # 1. Check if start date is after due date
        if task_info['start_date'] > task_info['due_date']:
            issues_found.append(f"Start date ({task_info['start_date'].strftime('%Y-%m-%d')}) is after due date ({task_info['due_date'].strftime('%Y-%m-%d')})")
        
        # 2. Check for very short durations (less than 3 days)
        duration = (task_info['due_date'] - task_info['start_date']).days + 1
        if duration < 3:
            issues_found.append(f"Duration is only {duration} days (less than 3-day minimum)")
        
        # 3. Check for very long durations (more than 1 year)
        if duration > 365:
            issues_found.append(f"Duration is {duration} days (more than 1 year)")
        
        # 4. Check dependencies
        for dep_id in task_info['dependencies']:
            if dep_id not in tasks:
                issues_found.append(f"Depends on non-existent task '{dep_id}'")
            else:
                dep_task = tasks[dep_id]
                # Check if dependency finishes before this task starts
                if dep_task['due_date'] >= task_info['start_date']:
                    gap = (task_info['start_date'] - dep_task['due_date']).days - 1
                    if gap < 0:
                        issues_found.append(f"Overlaps with dependency '{dep_id}' by {abs(gap)} days")
                    elif gap == 0:
                        issues_found.append(f"Starts same day dependency '{dep_id}' finishes (no buffer)")
        
        # 5. Check for logical inconsistencies
        if task_info['category'] == 'MILESTONE' and duration > 1:
            issues_found.append(f"Milestone task has duration of {duration} days (should be 1 day)")
        
        # 6. Check for parent-child relationship issues
        if task_info['parent'] and task_info['parent'] in tasks:
            parent_task = tasks[task_info['parent']]
            if task_info['start_date'] < parent_task['start_date']:
                issues_found.append(f"Starts before parent task '{task_info['parent']}'")
            if task_info['due_date'] > parent_task['due_date']:
                issues_found.append(f"Ends after parent task '{task_info['parent']}'")
        
        # 7. Check for unrealistic scheduling
        if task_info['category'] == 'IMAGING' and 'surgery' in task_info['name'].lower():
            # Surgery tasks should be single day
            if duration > 1:
                issues_found.append(f"Surgery task has {duration} days duration (should be 1 day)")
        
        # 8. Check for overlapping tasks in same category
        overlapping_tasks = []
        for other_id, other_task in tasks.items():
            if (other_id != task_id and 
                other_task['category'] == task_info['category'] and
                not (task_info['due_date'] < other_task['start_date'] or 
                     task_info['start_date'] > other_task['due_date'])):
                overlapping_tasks.append(other_id)
        
        if overlapping_tasks:
            issues_found.append(f"Overlaps with other {task_info['category']} tasks: {', '.join(overlapping_tasks)}")
        
        # 9. Check for tasks that should have dependencies but don't
        if (task_info['category'] == 'IMAGING' and 
            'recovery' in task_info['name'].lower() and 
            not task_info['dependencies']):
            issues_found.append("Recovery task should depend on surgery task")
        
        # 10. Check for tasks that depend on themselves
        if task_id in task_info['dependencies']:
            issues_found.append("Task depends on itself")
        
        # Report issues for this task
        if issues_found:
            issues.append({
                'task_id': task_id,
                'row_num': row_num,
                'name': task_name,
                'issues': issues_found
            })
    
    # Print all issues
    if issues:
        print(f"Found {len(issues)} tasks with potential issues:\n")
        for issue in issues:
            print(f"Line {issue['row_num']}: {issue['task_id']} - {issue['name']}")
            for problem in issue['issues']:
                print(f"  ❌ {problem}")
            print()
    else:
        print("✅ No issues found in the data!")
    
    return issues

def check_timeline_consistency(csv_file):
    """
    Check for timeline consistency issues.
    """
    tasks = {}
    
    with open(csv_file, 'r', newline='', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            if not row['Task ID']:
                continue
            tasks[row['Task ID']] = {
                'name': row['Task Name'],
                'start_date': parse_date(row['Start Date']),
                'due_date': parse_date(row['Due Date']),
                'category': row['Category']
            }
    
    print("\n=== TIMELINE CONSISTENCY CHECK ===")
    
    # Check for tasks that seem to be in wrong chronological order
    categories = {}
    for task_id, task_info in tasks.items():
        cat = task_info['category']
        if cat not in categories:
            categories[cat] = []
        categories[cat].append((task_id, task_info))
    
    for cat, cat_tasks in categories.items():
        # Sort by start date
        cat_tasks.sort(key=lambda x: x[1]['start_date'])
        
        print(f"\n{cat} tasks (chronological order):")
        for i, (task_id, task_info) in enumerate(cat_tasks):
            print(f"  {i+1}. {task_id}: {task_info['name']} ({task_info['start_date'].strftime('%Y-%m-%d')} → {task_info['due_date'].strftime('%Y-%m-%d')})")

if __name__ == "__main__":
    csv_file = "/Users/aaron/Downloads/gantt/input/data.cleaned.csv"
    
    issues = analyze_csv_issues(csv_file)
    check_timeline_consistency(csv_file)
    
    if issues:
        print(f"\n📊 SUMMARY: Found {len(issues)} tasks with potential issues")
        print("Review the issues above and consider fixing them for better project planning.")
    else:
        print("\n🎉 All tasks look good! No issues found.")
