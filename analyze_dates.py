#!/usr/bin/env python3
"""
Analyze the CSV data to find the earliest and latest task dates
"""

import csv
from datetime import datetime
import sys

def analyze_dates(csv_file):
    """Extract earliest and latest dates from the CSV file"""
    earliest_date = None
    latest_date = None
    
    with open(csv_file, 'r') as f:
        reader = csv.DictReader(f)
        
        for row in reader:
            # Skip empty rows
            if not row.get('Task ID'):
                continue
                
            # Check start date
            start_date_str = row.get('Start Date', '').strip()
            if start_date_str:
                try:
                    start_date = datetime.strptime(start_date_str, '%Y-%m-%d')
                    if earliest_date is None or start_date < earliest_date:
                        earliest_date = start_date
                    if latest_date is None or start_date > latest_date:
                        latest_date = start_date
                except ValueError:
                    print(f"Warning: Invalid start date format: {start_date_str}")
            
            # Check due date
            due_date_str = row.get('Due Date', '').strip()
            if due_date_str:
                try:
                    due_date = datetime.strptime(due_date_str, '%Y-%m-%d')
                    if earliest_date is None or due_date < earliest_date:
                        earliest_date = due_date
                    if latest_date is None or due_date > latest_date:
                        latest_date = due_date
                except ValueError:
                    print(f"Warning: Invalid due date format: {due_date_str}")
    
    return earliest_date, latest_date

if __name__ == "__main__":
    csv_file = "aarons-attempt/input/data.cleaned.csv"
    
    earliest, latest = analyze_dates(csv_file)
    
    if earliest and latest:
        print(f"Earliest task date: {earliest.strftime('%Y-%m-%d')}")
        print(f"Latest task date: {latest.strftime('%Y-%m-%d')}")
        print(f"Duration: {(latest - earliest).days} days")
        print(f"Start year: {earliest.year}")
        print(f"End year: {latest.year}")
        print(f"Start month: {earliest.month}")
        print(f"End month: {latest.month}")
    else:
        print("No valid dates found in the CSV file")
        sys.exit(1)
