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
    """Handles CSV file reading and basic validation."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def read_csv_data(self, file_path: str) -> List[Dict[str, str]]:
        """Read CSV data and return as list of dictionaries."""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                reader = csv.DictReader(f)
                data = list(reader)
                self.logger.info(f"Successfully read {len(data)} rows from {file_path}")
                return data
        except FileNotFoundError:
            self.logger.error(f"CSV file not found: {file_path}")
            raise
        except Exception as e:
            self.logger.error(f"Error reading CSV file {file_path}: {e}")
            raise


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
        """Complete pipeline from CSV file to ProjectTimeline."""
        self.logger.info(f"Starting data processing pipeline for {csv_file_path}")
        
        # Read CSV data
        csv_data = self.csv_reader.read_csv_data(csv_file_path)
        
        # Process tasks
        tasks = self.task_processor.process_tasks(csv_data)
        
        if not tasks:
            raise ValueError("No valid tasks found in CSV data")
        
        # Build timeline
        timeline = self.timeline_builder.build_timeline(tasks, title)
        
        self.logger.info(f"Data processing pipeline completed successfully")
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
