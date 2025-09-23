#!/usr/bin/env python3
"""
Timeline Validation Script
Validates task IDs and dependencies in the research timeline CSV
"""

import csv
import sys
from collections import defaultdict


def validate_timeline(csv_file):
    """Validate task IDs and dependencies in the timeline CSV"""

    print("🔍 Starting Timeline Validation...")
    print("=" * 50)

    # Read the CSV file
    tasks = {}
    task_ids = set()
    errors = []
    warnings = []

    try:
        with open(csv_file, 'r', encoding='utf-8') as file:
            reader = csv.DictReader(file)

            # Start at 2 for header
            for row_num, row in enumerate(reader, start=2):
                task_id = row['Task ID'].strip()
                dependencies = row['Dependencies'].strip()

                # Store task info
                tasks[task_id] = {
                    'row': row_num,
                    'Task': row['Task'],
                    'Dependencies': dependencies,
                    'Phase': row['Phase'],
                    'Sub-Phase': row['Sub-Phase'],
                    'Start Date': row['Start Date'],
                    'End Date': row['End Date'],
                    'Objective': row['Objective'],
                    'Milestone': row['Milestone'],
                    'Status': row['Status']
                }
                task_ids.add(task_id)

    except FileNotFoundError:
        print(f"❌ Error: File '{csv_file}' not found")
        return False
    except Exception as e:
        print(f"❌ Error reading file: {e}")
        return False

    print(f"📊 Found {len(tasks)} tasks")
    print()

    # Validation 1: Check for duplicate task IDs
    print("1️⃣ Checking for duplicate Task IDs...")
    seen_ids = set()
    duplicates = []

    for task_id in task_ids:
        if task_id in seen_ids:
            duplicates.append(task_id)
        seen_ids.add(task_id)

    if duplicates:
        errors.append(f"Duplicate Task IDs found: {duplicates}")
        print(f"❌ Found {len(duplicates)} duplicate Task IDs: {duplicates}")
    else:
        print("✅ No duplicate Task IDs found")

    print()

    # Validation 2: Check for missing dependencies
    print("2️⃣ Checking for missing dependencies...")
    missing_deps = []

    for task_id, task_info in tasks.items():
        if not task_info['Dependencies']:
            continue

        # Parse dependencies (handle comma-separated and quoted lists)
        deps = []
        if ',' in task_info['Dependencies']:
            # Handle quoted comma-separated values
            deps = [dep.strip().strip('"')
                    for dep in task_info['Dependencies'].split(',')]
        else:
            deps = [task_info['Dependencies'].strip()]

        for dep in deps:
            if dep and dep not in task_ids:
                missing_deps.append({
                    'task': task_id,
                    'missing_dep': dep,
                    'row': task_info['row']
                })

    if missing_deps:
        errors.append(f"Missing dependencies found: {len(missing_deps)}")
        print(f"❌ Found {len(missing_deps)} missing dependencies:")
        for issue in missing_deps:
            print(
                f"   Row {issue['row']}: Task '{issue['task']}' references missing task '{issue['missing_dep']}'")
    else:
        print("✅ All dependencies reference existing tasks")

    print()

    # Validation 3: Check for circular dependencies
    print("3️⃣ Checking for circular dependencies...")
    circular_deps = []

    def has_circular_dependency(task_id, visited=None, path=None):
        if visited is None:
            visited = set()
        if path is None:
            path = []

        if task_id in path:
            return path[path.index(task_id):] + [task_id]

        if task_id in visited:
            return None

        visited.add(task_id)
        path.append(task_id)

        if task_id not in tasks:
            return None

        task_info = tasks[task_id]
        if not task_info['Dependencies']:
            path.pop()
            return None

        # Parse dependencies
        deps = []
        if ',' in task_info['Dependencies']:
            deps = [dep.strip().strip('"')
                    for dep in task_info['Dependencies'].split(',')]
        else:
            deps = [task_info['Dependencies'].strip()]

        for dep in deps:
            if dep:
                result = has_circular_dependency(dep, visited, path)
                if result:
                    return result

        path.pop()
        return None

    for task_id in task_ids:
        cycle = has_circular_dependency(task_id)
        if cycle:
            circular_deps.append(cycle)

    if circular_deps:
        errors.append(f"Circular dependencies found: {len(circular_deps)}")
        print(f"❌ Found {len(circular_deps)} circular dependencies:")
        for i, cycle in enumerate(circular_deps, 1):
            print(f"   Cycle {i}: {' → '.join(cycle)}")
    else:
        print("✅ No circular dependencies found")

    print()

    # Validation 4: Check for orphaned tasks (no dependencies and no dependents)
    print("4️⃣ Checking for orphaned tasks...")
    dependents = defaultdict(list)

    for task_id, task_info in tasks.items():
        if not task_info['Dependencies']:
            continue

        deps = []
        if ',' in task_info['Dependencies']:
            deps = [dep.strip().strip('"')
                    for dep in task_info['Dependencies'].split(',')]
        else:
            deps = [task_info['Dependencies'].strip()]

        for dep in deps:
            if dep:
                dependents[dep].append(task_id)

    orphaned = []
    for task_id in task_ids:
        if (not tasks[task_id]['Dependencies'] and
            task_id not in dependents and
            not task_id.endswith('.M1') and  # Exclude milestones
            not task_id.endswith('.M2') and
                not task_id.endswith('.M3')):
            orphaned.append(task_id)

    if orphaned:
        warnings.append(f"Orphaned tasks found: {len(orphaned)}")
        print(
            f"⚠️  Found {len(orphaned)} orphaned tasks (no dependencies, no dependents):")
        for task_id in orphaned:
            print(f"   {task_id}: {tasks[task_id]['Task']}")
    else:
        print("✅ No orphaned tasks found")

    print()

    # Validation 5: Check task ID format consistency
    print("5️⃣ Checking task ID format consistency...")
    format_issues = []

    for task_id in task_ids:
        if not task_id:
            format_issues.append(f"Empty task ID found")
            continue

        # Check for basic format patterns
        if not (task_id.startswith('T') and
                ('.' in task_id or task_id.endswith('M1') or task_id.endswith('M2') or task_id.endswith('M3'))):
            format_issues.append(f"Unusual format: {task_id}")

    if format_issues:
        warnings.append(f"Task ID format issues: {len(format_issues)}")
        print(f"⚠️  Found {len(format_issues)} task ID format issues:")
        for issue in format_issues:
            print(f"   {issue}")
    else:
        print("✅ Task ID formats look consistent")

    print()

    # Validation 6: Check timeline logic based on dependencies
    print("6️⃣ Checking timeline logic based on dependencies...")
    timeline_issues = []

    def parse_date(date_str):
        """Parse date string and return comparable value"""
        try:
            from datetime import datetime
            return datetime.strptime(date_str, '%Y-%m-%d')
        except:
            return None

    for task_id, task_info in tasks.items():
        if not task_info['Dependencies']:
            continue

        # Parse dependencies
        deps = []
        if ',' in task_info['Dependencies']:
            deps = [dep.strip().strip('"')
                    for dep in task_info['Dependencies'].split(',')]
        else:
            deps = [task_info['Dependencies'].strip()]

        # Check if required columns exist
        if 'Start Date' not in task_info or 'End Date' not in task_info:
            timeline_issues.append({
                'task': task_id,
                'issue': 'Missing date columns',
                'details': f"Available columns: {list(task_info.keys())}"
            })
            continue

        task_start = parse_date(task_info['Start Date'])
        task_end = parse_date(task_info['End Date'])

        if not task_start or not task_end:
            timeline_issues.append({
                'task': task_id,
                'issue': 'Invalid date format',
                'details': f"Start: '{task_info['Start Date']}', End: '{task_info['End Date']}'"
            })
            continue

        # Check each dependency
        for dep in deps:
            if not dep or dep not in tasks:
                continue

            dep_info = tasks[dep]
            dep_start = parse_date(dep_info['Start Date'])
            dep_end = parse_date(dep_info['End Date'])

            if not dep_start or not dep_end:
                timeline_issues.append({
                    'task': task_id,
                    'issue': 'Dependency has invalid date format',
                    'details': f"Dependency {dep}: Start: '{dep_info['Start Date']}', End: '{dep_info['End Date']}'"
                })
                continue

            # Check if task starts before dependency ends
            if task_start < dep_end:
                timeline_issues.append({
                    'task': task_id,
                    'issue': 'Task starts before dependency ends',
                    'details': f"Task starts {task_info['Start Date']} but dependency {dep} ends {dep_info['End Date']}"
                })

            # Check if task ends before dependency ends
            if task_end < dep_end:
                timeline_issues.append({
                    'task': task_id,
                    'issue': 'Task ends before dependency ends',
                    'details': f"Task ends {task_info['End Date']} but dependency {dep} ends {dep_info['End Date']}"
                })

        # Check if task duration is reasonable (not negative)
        if task_start > task_end:
            timeline_issues.append({
                'task': task_id,
                'issue': 'Task starts after it ends',
                'details': f"Start: {task_info['Start Date']}, End: {task_info['End Date']}"
            })

    if timeline_issues:
        errors.append(f"Timeline logic issues found: {len(timeline_issues)}")
        print(f"❌ Found {len(timeline_issues)} timeline logic issues:")
        for issue in timeline_issues:
            print(f"   {issue['task']}: {issue['issue']}")
            print(f"      {issue['details']}")
    else:
        print("✅ Timeline logic looks consistent")

    print()

    # Validation 7: Check for overlapping tasks in same phase/sub-phase
    print("7️⃣ Checking for overlapping tasks in same phase/sub-phase...")
    overlap_issues = []

    # Group tasks by phase and sub-phase
    phase_groups = {}
    for task_id, task_info in tasks.items():
        phase_key = f"{task_info['Phase']}|{task_info['Sub-Phase']}"
        if phase_key not in phase_groups:
            phase_groups[phase_key] = []
        phase_groups[phase_key].append((task_id, task_info))

    for phase_key, phase_tasks in phase_groups.items():
        if len(phase_tasks) < 2:
            continue

        # Check for overlaps within the same phase/sub-phase
        for i, (task1_id, task1_info) in enumerate(phase_tasks):
            for j, (task2_id, task2_info) in enumerate(phase_tasks[i+1:], i+1):
                start1 = parse_date(task1_info['Start Date'])
                end1 = parse_date(task1_info['End Date'])
                start2 = parse_date(task2_info['Start Date'])
                end2 = parse_date(task2_info['End Date'])

                if not all([start1, end1, start2, end2]):
                    continue

                # Check for overlap
                if not (end1 <= start2 or end2 <= start1):
                    overlap_issues.append({
                        'phase': phase_key,
                        'task1': task1_id,
                        'task2': task2_id,
                        'overlap': f"{task1_info['Start Date']} to {task1_info['End Date']} overlaps with {task2_info['Start Date']} to {task2_info['End Date']}"
                    })

    if overlap_issues:
        warnings.append(f"Overlapping tasks found: {len(overlap_issues)}")
        print(f"⚠️  Found {len(overlap_issues)} overlapping tasks:")
        for issue in overlap_issues:
            print(
                f"   Phase {issue['phase']}: {issue['task1']} and {issue['task2']}")
            print(f"      {issue['overlap']}")
    else:
        print("✅ No overlapping tasks found")

    print()

    # Validation 8: Check for gaps in sequential tasks
    print("8️⃣ Checking for gaps in sequential tasks...")
    gap_issues = []

    # Find sequential task chains
    for task_id, task_info in tasks.items():
        if not task_info['Dependencies']:
            continue

        deps = []
        if ',' in task_info['Dependencies']:
            deps = [dep.strip().strip('"')
                    for dep in task_info['Dependencies'].split(',')]
        else:
            deps = [task_info['Dependencies'].strip()]

        for dep in deps:
            if not dep or dep not in tasks:
                continue

            dep_info = tasks[dep]
            task_start = parse_date(task_info['Start Date'])
            dep_end = parse_date(dep_info['End Date'])

            if not task_start or not dep_end:
                continue

            # Check for gaps (more than 7 days between tasks)
            gap_days = (task_start - dep_end).days
            if gap_days > 7:
                gap_issues.append({
                    'task': task_id,
                    'dependency': dep,
                    'gap': f"{gap_days} days between {dep_info['End Date']} and {task_info['Start Date']}"
                })

    if gap_issues:
        warnings.append(f"Large gaps found: {len(gap_issues)}")
        print(
            f"⚠️  Found {len(gap_issues)} large gaps between sequential tasks:")
        for issue in gap_issues:
            print(
                f"   {issue['task']} after {issue['dependency']}: {issue['gap']}")
    else:
        print("✅ No large gaps found between sequential tasks")

    print()

    # Summary
    print("📋 VALIDATION SUMMARY")
    print("=" * 50)

    if errors:
        print(f"❌ ERRORS: {len(errors)}")
        for error in errors:
            print(f"   • {error}")
    else:
        print("✅ NO ERRORS FOUND")

    if warnings:
        print(f"⚠️  WARNINGS: {len(warnings)}")
        for warning in warnings:
            print(f"   • {warning}")
    else:
        print("✅ NO WARNINGS")

    print()

    if errors:
        print("🔧 RECOMMENDATIONS:")
        print("   • Fix missing dependencies")
        print("   • Resolve circular dependencies")
        print("   • Check for typos in task IDs")
        return False
    else:
        print("🎉 Timeline validation passed!")
        return True


if __name__ == "__main__":
    csv_file = "reference/source-data/Research Timeline v5 - Comprehensive.csv"

    if len(sys.argv) > 1:
        csv_file = sys.argv[1]

    success = validate_timeline(csv_file)
    sys.exit(0 if success else 1)
