#!/usr/bin/env python3
"""
Data processing module for the LaTeX Gantt Chart Generator.
Handles CSV reading, parsing, and validation.
"""

import csv
import logging
from datetime import date, datetime
from pathlib import Path
from typing import Dict, List, Optional, Tuple

from .config import config
from .models import Task, ProjectTimeline


class DataProcessor:
    """Main data processing coordinator."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def process_csv_to_timeline(self, csv_file_path: str, title: str = None) -> ProjectTimeline:
        """Complete pipeline from CSV file to ProjectTimeline."""
        self.logger.info(f"Processing CSV file: {csv_file_path}")
        
        # Read CSV data
        tasks = self._read_csv_data(csv_file_path)
        
        if not tasks:
            raise ValueError("No valid tasks found in CSV data")
        
        # Build timeline
        timeline_title = title or config.latex.default_title
        timeline = ProjectTimeline(
            tasks=tasks,
            title=timeline_title,
            start_date=date.today(),  # Will be recalculated
            end_date=date.today()     # Will be recalculated
        )
        
        self.logger.info(f"Processed {len(tasks)} tasks successfully")
        return timeline
    
    def _read_csv_data(self, file_path: str) -> List[Task]:
        """Read and process CSV data into Task objects."""
        tasks = []
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                reader = csv.DictReader(f)
                
                for row_num, row in enumerate(reader, start=2):
                    try:
                        task = self._parse_csv_row(row)
                        if task:
                            tasks.append(task)
                    except Exception as e:
                        self.logger.warning(f"Row {row_num}: {e}")
                        
        except Exception as e:
            self.logger.error(f"Error reading CSV file: {e}")
            raise
        
        return tasks
    
    def _parse_csv_row(self, row: Dict[str, str]) -> Optional[Task]:
        """Parse a single CSV row into a Task object."""
        # Extract required fields
        task_id = row.get('Task ID', '').strip()
        task_name = row.get('Task Name', '').strip()
        start_date_str = row.get('Start Date', '').strip()
        due_date_str = row.get('Due Date', '').strip()
        category = row.get('Category', '').strip()
        
        if not all([task_id, task_name, start_date_str, due_date_str]):
            return None
        
        # Parse dates
        start_date = self._parse_date(start_date_str)
        due_date = self._parse_date(due_date_str)
        if not start_date or not due_date:
            return None
        
        # Extract optional fields
        dependencies = row.get('Dependencies', '').strip()
        notes = row.get('Description', '').strip()
        
        return Task(
            id=task_id,
            name=task_name,
            start_date=start_date,
            due_date=due_date,
            category=category,
            dependencies=dependencies,
            notes=notes
        )
    
    def _parse_date(self, date_str: str) -> Optional[date]:
        """Parse date string to date object."""
        if not date_str:
            return None
        
        # Try common date formats
        formats = ['%Y-%m-%d', '%m/%d/%Y', '%d/%m/%Y', '%B %d, %Y']
        
        for fmt in formats:
            try:
                return datetime.strptime(date_str, fmt).date()
            except ValueError:
                continue
        
        return None
    
    def validate_input_file(self, file_path: str) -> Tuple[bool, List[str]]:
        """Validate that the input file exists and is readable."""
        errors = []
        
        if not Path(file_path).exists():
            errors.append(f"File does not exist: {file_path}")
            return False, errors
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                # Just check if we can read the first few lines
                for i, line in enumerate(f):
                    if i >= 5:
                        break
        except Exception as e:
            errors.append(f"Cannot read file: {e}")
            return False, errors
        
        return True, errors
