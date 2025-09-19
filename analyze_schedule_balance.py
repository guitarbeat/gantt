#!/usr/bin/env python3
"""
Analyze the current PhD schedule for workload distribution and identify imbalances.
"""

import pandas as pd
from datetime import datetime, timedelta
import matplotlib.pyplot as plt
import seaborn as sns
from collections import defaultdict

def analyze_schedule_balance(csv_file):
    """Analyze the current schedule for workload distribution."""
    
    # Read the CSV file
    df = pd.read_csv(csv_file)
    
    # Convert dates to datetime
    df['Start Date'] = pd.to_datetime(df['Start Date'])
    df['Due Date'] = pd.to_datetime(df['Due Date'])
    
    # Calculate task duration
    df['Duration'] = (df['Due Date'] - df['Start Date']).dt.days + 1
    
    # Create monthly workload analysis
    workload_by_month = defaultdict(int)
    workload_by_quarter = defaultdict(int)
    workload_by_category = defaultdict(int)
    
    # Analyze workload distribution
    for _, row in df.iterrows():
        if pd.notna(row['Start Date']) and pd.notna(row['Due Date']):
            # Monthly workload
            current_date = row['Start Date']
            while current_date <= row['Due Date']:
                month_key = current_date.strftime('%Y-%m')
                workload_by_month[month_key] += 1
                current_date += timedelta(days=1)
            
            # Quarterly workload
            start_quarter = f"{row['Start Date'].year}-Q{(row['Start Date'].month-1)//3 + 1}"
            end_quarter = f"{row['Due Date'].year}-Q{(row['Due Date'].month-1)//3 + 1}"
            
            if start_quarter == end_quarter:
                workload_by_quarter[start_quarter] += row['Duration']
            else:
                # Split workload across quarters
                quarter_start = row['Start Date']
                quarter_end = row['Due Date']
                
                # Calculate days in each quarter
                q1_end = datetime(row['Start Date'].year, 3, 31)
                q2_end = datetime(row['Start Date'].year, 6, 30)
                q3_end = datetime(row['Start Date'].year, 9, 30)
                q4_end = datetime(row['Start Date'].year, 12, 31)
                
                if row['Start Date'] <= q1_end and row['Due Date'] >= q1_end:
                    workload_by_quarter[start_quarter] += (q1_end - row['Start Date']).days + 1
                if row['Start Date'] <= q2_end and row['Due Date'] >= q2_end:
                    q2_start = max(row['Start Date'], datetime(row['Start Date'].year, 4, 1))
                    workload_by_quarter[f"{row['Start Date'].year}-Q2"] += (q2_end - q2_start).days + 1
                if row['Start Date'] <= q3_end and row['Due Date'] >= q3_end:
                    q3_start = max(row['Start Date'], datetime(row['Start Date'].year, 7, 1))
                    workload_by_quarter[f"{row['Start Date'].year}-Q3"] += (q3_end - q3_start).days + 1
                if row['Start Date'] <= q4_end and row['Due Date'] >= q4_end:
                    q4_start = max(row['Start Date'], datetime(row['Start Date'].year, 10, 1))
                    workload_by_quarter[f"{row['Start Date'].year}-Q4"] += (q4_end - q4_start).days + 1
            
            # Category workload
            if pd.notna(row['Category']):
                workload_by_category[row['Category']] += row['Duration']
    
    # Print analysis results
    print("=== SCHEDULE BALANCE ANALYSIS ===\n")
    
    print("1. QUARTERLY WORKLOAD DISTRIBUTION:")
    for quarter in sorted(workload_by_quarter.keys()):
        print(f"   {quarter}: {workload_by_quarter[quarter]} task-days")
    
    print(f"\n2. CATEGORY WORKLOAD DISTRIBUTION:")
    for category in sorted(workload_by_category.keys()):
        print(f"   {category}: {workload_by_category[category]} task-days")
    
    print(f"\n3. WORKLOAD PEAKS (Top 5 months):")
    sorted_months = sorted(workload_by_month.items(), key=lambda x: x[1], reverse=True)
    for month, workload in sorted_months[:5]:
        print(f"   {month}: {workload} task-days")
    
    # Identify problematic periods
    print(f"\n4. IDENTIFIED ISSUES:")
    
    # Check for workload concentration
    max_quarterly = max(workload_by_quarter.values()) if workload_by_quarter else 0
    min_quarterly = min(workload_by_quarter.values()) if workload_by_quarter else 0
    if max_quarterly > 0:
        imbalance_ratio = max_quarterly / min_quarterly if min_quarterly > 0 else float('inf')
        if imbalance_ratio > 2.0:
            print(f"   - High workload imbalance: {imbalance_ratio:.1f}x difference between heaviest and lightest quarters")
    
    # Check for consecutive heavy periods
    quarters = sorted(workload_by_quarter.keys())
    heavy_quarters = [q for q in quarters if workload_by_quarter[q] > max_quarterly * 0.7]
    if len(heavy_quarters) >= 3:
        print(f"   - Multiple consecutive heavy quarters: {', '.join(heavy_quarters)}")
    
    # Check for gaps in work
    light_quarters = [q for q in quarters if workload_by_quarter[q] < max_quarterly * 0.3]
    if light_quarters:
        print(f"   - Light workload quarters that could absorb more work: {', '.join(light_quarters)}")
    
    return {
        'quarterly_workload': workload_by_quarter,
        'monthly_workload': workload_by_month,
        'category_workload': workload_by_category,
        'dataframe': df
    }

if __name__ == "__main__":
    analysis = analyze_schedule_balance('/Users/aaron/Downloads/gantt/input/data.cleaned.csv')
