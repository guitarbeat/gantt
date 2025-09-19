#!/usr/bin/env python3
"""
Detailed line-by-line analysis of the CSV data.
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

def analyze_csv_detailed(csv_file):
    """
    Detailed analysis of CSV data line by line.
    """
    tasks = {}
    
    print("=== DETAILED LINE-BY-LINE ANALYSIS ===\n")
    
    with open(csv_file, 'r', newline='', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)
        for row_num, row in enumerate(reader, start=2):
            if not row['Task ID']:
                continue
            
            task_id = row['Task ID']
            task_name = row['Task Name']
            start_date = parse_date(row['Start Date'])
            due_date = parse_date(row['Due Date'])
            dependencies = parse_dependencies(row['Dependencies'])
            category = row['Category']
            parent = row['Parent Task ID']
            
            duration = (due_date - start_date).days + 1
            
            print(f"Line {row_num}: {task_id} - {task_name}")
            print(f"  Category: {category}")
            print(f"  Duration: {duration} days ({start_date.strftime('%Y-%m-%d')} → {due_date.strftime('%Y-%m-%d')})")
            print(f"  Dependencies: {dependencies if dependencies else 'None'}")
            print(f"  Parent: {parent if parent else 'None'}")
            
            # Check for specific issues
            issues = []
            
            # 1. Surgery tasks should be single day
            if 'surgery' in task_name.lower() and duration > 1:
                issues.append(f"⚠️  Surgery task has {duration} days (should be 1 day)")
            
            # 2. Milestone tasks should be single day
            if 'MILESTONE' in row['Description'] and duration > 1:
                issues.append(f"⚠️  Milestone task has {duration} days (should be 1 day)")
            
            # 3. Very short durations
            if duration < 3:
                issues.append(f"⚠️  Very short duration: {duration} days")
            
            # 4. Very long durations
            if duration > 90:
                issues.append(f"⚠️  Very long duration: {duration} days")
            
            # 5. Check dependencies
            for dep_id in dependencies:
                if dep_id not in tasks:
                    issues.append(f"❌ Depends on non-existent task: {dep_id}")
            
            # 6. Check for logical issues
            if category == 'IMAGING' and 'recovery' in task_name.lower() and not dependencies:
                issues.append("❌ Recovery task should have dependencies")
            
            if issues:
                print("  Issues found:")
                for issue in issues:
                    print(f"    {issue}")
            else:
                print("  ✅ No issues")
            
            print()
            
            # Store task for further analysis
            tasks[task_id] = {
                'row_num': row_num,
                'name': task_name,
                'start_date': start_date,
                'due_date': due_date,
                'dependencies': dependencies,
                'category': category,
                'parent': parent,
                'duration': duration
            }
    
    return tasks

def check_overlapping_tasks(tasks):
    """
    Check for overlapping tasks within the same category.
    """
    print("=== OVERLAPPING TASKS ANALYSIS ===\n")
    
    categories = {}
    for task_id, task_info in tasks.items():
        cat = task_info['category']
        if cat not in categories:
            categories[cat] = []
        categories[cat].append((task_id, task_info))
    
    for cat, cat_tasks in categories.items():
        print(f"{cat} Category ({len(cat_tasks)} tasks):")
        
        # Check for overlaps
        overlaps = []
        for i, (task_id1, task1) in enumerate(cat_tasks):
            for j, (task_id2, task2) in enumerate(cat_tasks[i+1:], i+1):
                # Check if tasks overlap
                if not (task1['due_date'] < task2['start_date'] or 
                       task1['start_date'] > task2['due_date']):
                    overlap_days = min(task1['due_date'], task2['due_date']) - max(task1['start_date'], task2['start_date']) + 1
                    overlaps.append((task_id1, task1['name'], task_id2, task2['name'], overlap_days))
        
        if overlaps:
            print(f"  Found {len(overlaps)} overlapping pairs:")
            for task_id1, name1, task_id2, name2, overlap_days in overlaps:
                print(f"    {task_id1} ({name1}) overlaps with {task_id2} ({name2}) by {overlap_days} days")
        else:
            print("  ✅ No overlapping tasks")
        print()

def check_dependency_chains(tasks):
    """
    Check dependency chains for logical consistency.
    """
    print("=== DEPENDENCY CHAIN ANALYSIS ===\n")
    
    # Find chains
    chains = []
    processed = set()
    
    for task_id, task_info in tasks.items():
        if task_id in processed:
            continue
        
        # Start a new chain
        chain = []
        current_id = task_id
        
        while current_id and current_id not in processed:
            if current_id not in tasks:
                break
            
            current_task = tasks[current_id]
            chain.append((current_id, current_task))
            processed.add(current_id)
            
            # Find next task in chain (task that depends on current)
            next_id = None
            for other_id, other_task in tasks.items():
                if current_id in other_task['dependencies'] and other_id not in processed:
                    next_id = other_id
                    break
            
            current_id = next_id
        
        if len(chain) > 1:
            chains.append(chain)
    
    for i, chain in enumerate(chains, 1):
        print(f"Chain {i} ({len(chain)} tasks):")
        for j, (task_id, task_info) in enumerate(chain):
            print(f"  {j+1}. {task_id}: {task_info['name']} ({task_info['start_date'].strftime('%Y-%m-%d')} → {task_info['due_date'].strftime('%Y-%m-%d')})")
        print()

if __name__ == "__main__":
    csv_file = "/Users/aaron/Downloads/gantt/input/data.cleaned.csv"
    
    tasks = analyze_csv_detailed(csv_file)
    check_overlapping_tasks(tasks)
    check_dependency_chains(tasks)
    
    print(f"Analysis complete. Processed {len(tasks)} tasks.")
