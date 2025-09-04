#!/usr/bin/env python3
"""
Data processing module for the LaTeX Gantt chart generator.
Handles CSV reading, data validation, and task processing.
"""

import csv
import logging
from datetime import datetime, date
from typing import List, Dict, Optional, Tuple
from pathlib import Path

from .models import Task, ProjectTimeline, TaskValidator
from .config import config


class CSVReader:
    """Enhanced CSV file reader with better error handling and validation."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.required_columns = [
            'Task ID', 'Task Name', 'Start Date', 'Due Date', 'Category'
        ]
        self.optional_columns = [
            'Parent Task ID', 'Dependencies', 'Description', 'Status', 'Priority', 'Notes'
        ]
    
    def read_csv_data(self, file_path: str) -> Tuple[List[Dict[str, str]], List[str]]:
        """Read CSV data with enhanced error handling and validation."""
        errors = []
        data = []
        
        try:
            file_path = Path(file_path)
            if not file_path.exists():
                errors.append(f"File not found: {file_path}")
                return data, errors
            
            with open(file_path, 'r', encoding='utf-8') as f:
                # Detect delimiter
                sample = f.read(1024)
                f.seek(0)
                delimiter = self._detect_delimiter(sample)
                
                reader = csv.DictReader(f, delimiter=delimiter)
                
                # Validate headers
                header_errors = self._validate_headers(reader.fieldnames)
                errors.extend(header_errors)
                
                if not header_errors:  # Only proceed if headers are valid
                    for row_num, row in enumerate(reader, start=2):  # Start at 2 (header is row 1)
                        row_errors = self._validate_row(row, row_num)
                        if row_errors:
                            errors.extend(row_errors)
                        else:
                            data.append(row)
                            
        except UnicodeDecodeError:
            errors.append("File encoding error. Please ensure the file is UTF-8 encoded.")
        except Exception as e:
            errors.append(f"Unexpected error reading file: {e}")
        
        self.logger.info(f"Successfully read {len(data)} rows from {file_path}")
        if errors:
            self.logger.warning(f"Found {len(errors)} validation issues")
        
        return data, errors
    
    def _detect_delimiter(self, sample: str) -> str:
        """Detect CSV delimiter."""
        delimiters = [',', ';', '\t', '|']
        delimiter_counts = {delim: sample.count(delim) for delim in delimiters}
        return max(delimiter_counts, key=delimiter_counts.get)
    
    def _validate_headers(self, headers: List[str]) -> List[str]:
        """Validate CSV headers."""
        errors = []
        if not headers:
            errors.append("No headers found in CSV file")
            return errors
        
        # Normalize headers (strip whitespace, convert to title case)
        normalized_headers = [h.strip().title() for h in headers]
        
        # Check for required columns (case-insensitive)
        missing_columns = []
        for required in self.required_columns:
            if not any(required.lower() == h.lower() for h in normalized_headers):
                missing_columns.append(required)
        errors.extend(f"Missing required column: {col}" for col in missing_columns)
        
        return errors
    
    def _validate_row(self, row: Dict[str, str], row_num: int) -> List[str]:
        """Validate a single CSV row."""
        errors = []
        
        # Check for empty required fields
        empty_fields = [field for field in ['Task Name', 'Start Date', 'Due Date'] 
                       if not row.get(field, '').strip()]
        errors.extend(f"Row {row_num}: Empty {field}" for field in empty_fields)
        
        # Validate dates
        start_date = self._parse_date(row.get('Start Date', ''))
        due_date = self._parse_date(row.get('Due Date', ''))
        
        if start_date and due_date:
            if start_date > due_date:
                errors.append(f"Row {row_num}: Start date ({start_date}) is after due date ({due_date})")
        
        return errors
    
    def _parse_date(self, date_str: str) -> Optional[date]:
        """Parse date string with multiple format support."""
        if not date_str:
            return None
        
        date_str = date_str.strip()
        date_formats = [
            '%Y-%m-%d',      # 2025-08-29
            '%m/%d/%Y',      # 08/29/2025
            '%d/%m/%Y',      # 29/08/2025
            '%B %d, %Y',     # August 29, 2025
            '%b %d, %Y',     # Aug 29, 2025
        ]
        
        for fmt in date_formats:
            try:
                return datetime.strptime(date_str, fmt).date()
            except ValueError:
                continue
        
        return None


class TaskProcessor:
    """Processes raw CSV data into structured Task objects."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.validator = TaskValidator()
    
    def parse_date(self, date_str: str) -> Optional[date]:
        """Parse date string to date object with error handling."""
        if not date_str or not date_str.strip():
            return None
        
        try:
            return datetime.strptime(date_str.strip(), config.tasks.date_format).date()
        except ValueError as e:
            self.logger.warning(f"Invalid date format '{date_str}': {e}")
            return None
    
    def process_csv_row(self, row: Dict[str, str]) -> Optional[Task]:
        """Process a single CSV row into a Task object."""
        # Validate the row first
        errors = self.validator.validate_csv_row(row)
        if errors:
            self.logger.warning(f"Skipping row due to validation errors: {errors}")
            return None
        
        # Extract and validate required fields
        task_id = row.get('Task ID', '').strip()
        task_name = row.get('Task Name', '').strip()
        start_date = self.parse_date(row.get('Start Date', ''))
        due_date = self.parse_date(row.get('Due Date', ''))
        
        if not all([task_id, task_name, start_date, due_date]):
            self.logger.warning(f"Skipping row with missing required fields: {row}")
            return None
        
        # Extract optional fields
        category = row.get('Category', '').strip()
        dependencies = row.get('Dependencies', '').strip()
        notes = row.get('Description', '').strip()
        
        try:
            task = Task(
                id=task_id,
                name=task_name,
                start_date=start_date,
                due_date=due_date,
                category=category,
                dependencies=dependencies,
                notes=notes
            )
            self.logger.debug(f"Successfully processed task: {task_id}")
            return task
        except ValueError as e:
            self.logger.warning(f"Skipping task {task_id} due to validation error: {e}")
            return None
    
    def process_tasks(self, csv_data: List[Dict[str, str]]) -> List[Task]:
        """Process CSV data into a list of Task objects."""
        tasks = []
        skipped_count = 0
        
        for row in csv_data:
            task = self.process_csv_row(row)
            if task:
                tasks.append(task)
            else:
                skipped_count += 1
        
        self.logger.info(f"Processed {len(tasks)} tasks, skipped {skipped_count} invalid rows")
        return tasks


class TimelineBuilder:
    """Builds ProjectTimeline objects from processed tasks."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.validator = TaskValidator()
    
    def build_timeline(self, tasks: List[Task], title: str = None) -> ProjectTimeline:
        """Build a ProjectTimeline from a list of tasks."""
        if not tasks:
            raise ValueError("Cannot build timeline from empty task list")
        
        timeline_title = title or config.latex.default_title
        
        try:
            timeline = ProjectTimeline(
                tasks=tasks,
                title=timeline_title,
                start_date=date.today(),  # Will be recalculated in __post_init__
                end_date=date.today()     # Will be recalculated in __post_init__
            )
            
            # Validate the complete timeline
            errors = self.validator.validate_timeline(timeline)
            if errors:
                self.logger.warning(f"Timeline validation warnings: {errors}")
            
            self.logger.info(f"Built timeline '{timeline_title}' with {len(tasks)} tasks")
            return timeline
            
        except Exception as e:
            self.logger.error(f"Error building timeline: {e}")
            raise


class DataProcessor:
    """Main data processing coordinator."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.csv_reader = CSVReader()
        self.task_processor = TaskProcessor()
        self.timeline_builder = TimelineBuilder()
    
    def process_csv_to_timeline(self, csv_file_path: str, title: str = None) -> ProjectTimeline:
        """Complete pipeline from CSV file to ProjectTimeline with enhanced error handling."""
        self.logger.info(f"Starting enhanced data processing pipeline for {csv_file_path}")
        
        # Read CSV data with validation
        csv_data, errors = self.csv_reader.read_csv_data(csv_file_path)
        
        if errors:
            self.logger.warning(f"CSV validation found {len(errors)} issues:")
            for error in errors:
                self.logger.warning(f"  - {error}")
        
        if not csv_data:
            raise ValueError(f"No valid data found in CSV file. Errors: {errors}")
        
        # Process tasks
        tasks = self.task_processor.process_tasks(csv_data)
        
        if not tasks:
            raise ValueError("No valid tasks found in CSV data after processing")
        
        # Build timeline
        timeline = self.timeline_builder.build_timeline(tasks, title)
        
        self.logger.info(f"Enhanced data processing pipeline completed successfully")
        self.logger.info(f"Processed {len(tasks)} tasks from {len(csv_data)} CSV rows")
        return timeline
    
    def validate_input_file(self, file_path: str) -> Tuple[bool, List[str]]:
        """Validate that the input file exists and is readable."""
        errors = []
        
        if not Path(file_path).exists():
            errors.append(f"File does not exist: {file_path}")
            return False, errors
        
        if not Path(file_path).is_file():
            errors.append(f"Path is not a file: {file_path}")
            return False, errors
        
        try:
            # Try to read the file
            with open(file_path, 'r', encoding='utf-8') as f:
                # Just check if we can read the first few lines
                for i, line in enumerate(f):
                    if i >= 5:  # Check first 5 lines
                        break
        except Exception as e:
            errors.append(f"Cannot read file: {e}")
            return False, errors
        
        return True, errors
