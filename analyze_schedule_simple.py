#!/usr/bin/env python3
"""
Simple analysis of the current PhD schedule for workload distribution.
"""

import csv
from datetime import datetime, timedelta
from collections import defaultdict


def analyze_schedule_balance(csv_file):
    """Analyze the current schedule for workload distribution."""

    # Read the CSV file
    tasks = []
    with open(csv_file, 'r') as f:
        reader = csv.DictReader(f)
        for row in reader:
            if row['Start Date'] and row['Due Date']:
                tasks.append(row)

    # Convert dates and calculate durations
    for task in tasks:
        task['Start Date'] = datetime.strptime(task['Start Date'], '%Y-%m-%d')
        task['Due Date'] = datetime.strptime(task['Due Date'], '%Y-%m-%d')
        task['Duration'] = (task['Due Date'] - task['Start Date']).days + 1

    # Analyze workload distribution
    workload_by_quarter = defaultdict(int)
    workload_by_category = defaultdict(int)
    workload_by_month = defaultdict(int)

    for task in tasks:
        # Quarterly workload
        start_quarter = f"{task['Start Date'].year}-Q{(task['Start Date'].month-1)//3 + 1}"
        end_quarter = f"{task['Due Date'].year}-Q{(task['Due Date'].month-1)//3 + 1}"

        if start_quarter == end_quarter:
            workload_by_quarter[start_quarter] += task['Duration']
        else:
            # Split workload across quarters
            current_date = task['Start Date']
            while current_date <= task['Due Date']:
                quarter = f"{current_date.year}-Q{(current_date.month-1)//3 + 1}"
                workload_by_quarter[quarter] += 1
                current_date += timedelta(days=1)

        # Category workload
        if task['Category']:
            workload_by_category[task['Category']] += task['Duration']

        # Monthly workload
        current_date = task['Start Date']
        while current_date <= task['Due Date']:
            month_key = current_date.strftime('%Y-%m')
            workload_by_month[month_key] += 1
            current_date += timedelta(days=1)

    # Print analysis results
    print("=== SCHEDULE BALANCE ANALYSIS ===\n")

    print("1. QUARTERLY WORKLOAD DISTRIBUTION:")
    for quarter in sorted(workload_by_quarter.keys()):
        print(f"   {quarter}: {workload_by_quarter[quarter]} task-days")

    print(f"\n2. CATEGORY WORKLOAD DISTRIBUTION:")
    for category in sorted(workload_by_category.keys()):
        print(f"   {category}: {workload_by_category[category]} task-days")

    print(f"\n3. WORKLOAD PEAKS (Top 10 months):")
    sorted_months = sorted(workload_by_month.items(),
                           key=lambda x: x[1], reverse=True)
    for month, workload in sorted_months[:10]:
        print(f"   {month}: {workload} task-days")

    # Identify problematic periods
    print(f"\n4. IDENTIFIED ISSUES:")

    # Check for workload concentration
    if workload_by_quarter:
        max_quarterly = max(workload_by_quarter.values())
        min_quarterly = min(workload_by_quarter.values())
        imbalance_ratio = max_quarterly / \
            min_quarterly if min_quarterly > 0 else float('inf')
        if imbalance_ratio > 2.0:
            print(
                f"   - High workload imbalance: {imbalance_ratio:.1f}x difference between heaviest and lightest quarters")

    # Check for consecutive heavy periods
    quarters = sorted(workload_by_quarter.keys())
    heavy_quarters = [
        q for q in quarters if workload_by_quarter[q] > max_quarterly * 0.7]
    if len(heavy_quarters) >= 3:
        print(
            f"   - Multiple consecutive heavy quarters: {', '.join(heavy_quarters)}")

    # Check for gaps in work
    light_quarters = [
        q for q in quarters if workload_by_quarter[q] < max_quarterly * 0.3]
    if light_quarters:
        print(
            f"   - Light workload quarters that could absorb more work: {', '.join(light_quarters)}")

    # Identify specific problematic periods
    print(f"\n5. SPECIFIC PROBLEMATIC PERIODS:")

    # Find the heaviest quarters
    heaviest_quarters = sorted(
        workload_by_quarter.items(), key=lambda x: x[1], reverse=True)[:3]
    print(
        f"   - Heaviest quarters: {[f'{q} ({days} days)' for q, days in heaviest_quarters]}")

    # Find the lightest quarters
    lightest_quarters = sorted(
        workload_by_quarter.items(), key=lambda x: x[1])[:3]
    print(
        f"   - Lightest quarters: {[f'{q} ({days} days)' for q, days in lightest_quarters]}")

    return {
        'quarterly_workload': workload_by_quarter,
        'monthly_workload': workload_by_month,
        'category_workload': workload_by_category,
        'tasks': tasks
    }


if __name__ == "__main__":
    analysis = analyze_schedule_balance(
        '/Users/aaron/Downloads/gantt/input/data.final_balanced.csv')
