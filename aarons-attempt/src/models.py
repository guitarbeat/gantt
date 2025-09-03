#!/usr/bin/env python3
"""
Data models and validation for the LaTeX Gantt chart generator.
Provides structured data classes and validation logic.
"""

from dataclasses import dataclass
from datetime import date, datetime, timedelta
from typing import Optional, List, Dict, Any
import re

from .config import config


@dataclass
class Task:
    """Represents a single task with all its properties."""
    
    id: str
    name: str
    start_date: date
    due_date: date
    category: str
    dependencies: str
    notes: str
    is_milestone: bool = False
    
    def __post_init__(self):
        """Validate and clean task data after initialization."""
        self._clean_name()
        self._validate_dates()
        self._determine_milestone()
    
    def _clean_name(self) -> None:
        """Clean task name by removing common prefixes."""
        for prefix in config.tasks.name_prefixes_to_remove:
            if self.name.startswith(prefix):
                self.name = self.name[len(prefix):]
                break
    
    def _validate_dates(self) -> None:
        """Validate that dates are reasonable."""
        if self.start_date > self.due_date:
            raise ValueError(f"Task '{self.name}': start date ({self.start_date}) cannot be after due date ({self.due_date})")
        
        # Check for reasonable date ranges (not too far in past/future)
        today = date.today()
        if self.start_date < date(2020, 1, 1):
            raise ValueError(f"Task '{self.name}': start date ({self.start_date}) is too far in the past")
        
        if self.due_date > date(2030, 12, 31):
            raise ValueError(f"Task '{self.name}': due date ({self.due_date}) is too far in the future")
    
    def _determine_milestone(self) -> None:
        """Determine if this task is a milestone (same start and due date)."""
        self.is_milestone = (self.start_date == self.due_date)
    
    @property
    def duration_days(self) -> int:
        """Calculate task duration in days."""
        return (self.due_date - self.start_date).days + 1
    
    @property
    def category_color(self) -> str:
        """Get the color name for this task's category."""
        from .config import get_category_color
        return get_category_color(self.category)
    
    @property
    def marker(self) -> str:
        """Get the LaTeX marker for this task."""
        from .config import get_task_marker
        return get_task_marker(self.is_milestone)
    
    def overlaps_with_date(self, check_date: date) -> bool:
        """Check if this task overlaps with a given date."""
        return self.start_date <= check_date <= self.due_date
    
    def overlaps_with_range(self, start: date, end: date) -> bool:
        """Check if this task overlaps with a given date range."""
        return not (self.due_date < start or self.start_date > end)


@dataclass
class MonthInfo:
    """Represents information about a month in the calendar."""
    
    name: str
    start_date: date
    end_date: date
    first_weekday: int  # 0 = Sunday, 1 = Monday, etc.
    num_days: int
    
    @classmethod
    def from_date(cls, month_start: date) -> 'MonthInfo':
        """Create MonthInfo from the first day of a month."""
        # Calculate last day of month
        if month_start.month == 12:
            next_month = date(month_start.year + 1, 1, 1)
        else:
            next_month = date(month_start.year, month_start.month + 1, 1)
        month_end = next_month - timedelta(days=1)
        
        # Calculate first weekday (0 = Sunday)
        first_weekday = month_start.weekday() + 1
        if first_weekday == 7:
            first_weekday = 0
        
        # Calculate number of days
        num_days = (month_end - month_start).days + 1
        
        # Create month name
        month_name = month_start.strftime('%B %Y')
        
        return cls(
            name=month_name,
            start_date=month_start,
            end_date=month_end,
            first_weekday=first_weekday,
            num_days=num_days
        )


@dataclass
class ProjectTimeline:
    """Represents the complete project timeline with tasks and metadata."""
    
    tasks: List[Task]
    title: str
    start_date: date
    end_date: date
    
    def __post_init__(self):
        """Validate timeline data after initialization."""
        if not self.tasks:
            raise ValueError("Timeline must contain at least one task")
        
        # Recalculate date range from actual tasks
        self.start_date = min(task.start_date for task in self.tasks)
        self.end_date = max(task.due_date for task in self.tasks)
    
    @property
    def total_duration_days(self) -> int:
        """Calculate total project duration in days."""
        return (self.end_date - self.start_date).days
    
    @property
    def total_tasks(self) -> int:
        """Get total number of tasks."""
        return len(self.tasks)
    
    @property
    def milestones(self) -> List[Task]:
        """Get all milestone tasks."""
        return [task for task in self.tasks if task.is_milestone]
    
    @property
    def regular_tasks(self) -> List[Task]:
        """Get all non-milestone tasks."""
        return [task for task in self.tasks if not task.is_milestone]
    
    def get_tasks_for_month(self, month_info: MonthInfo) -> List[Task]:
        """Get all tasks that overlap with a given month."""
        return [
            task for task in self.tasks 
            if task.overlaps_with_range(month_info.start_date, month_info.end_date)
        ]
    
    def get_tasks_for_date(self, check_date: date) -> List[Task]:
        """Get all tasks that overlap with a given date."""
        return [task for task in self.tasks if task.overlaps_with_date(check_date)]
    
    def get_months_between(self) -> List[MonthInfo]:
        """Get all months between start and end dates."""
        months = []
        current = date(self.start_date.year, self.start_date.month, 1)
        end_month = date(self.end_date.year, self.end_date.month, 1)
        
        while current <= end_month:
            month_info = MonthInfo.from_date(current)
            
            # Only include months that overlap with our task range
            if not (month_info.end_date < self.start_date or month_info.start_date > self.end_date):
                months.append(month_info)
            
            # Move to next month
            if current.month == 12:
                current = date(current.year + 1, 1, 1)
            else:
                current = date(current.year, current.month + 1, 1)
        
        return months


class TaskValidator:
    """Validates task data and provides helpful error messages."""
    
    @staticmethod
    def validate_csv_row(row: Dict[str, str]) -> List[str]:
        """Validate a single CSV row and return list of errors."""
        errors = []
        
        # Check required fields
        required_fields = ['Task ID', 'Task Name', 'Start Date', 'Due Date']
        for field in required_fields:
            if not row.get(field, '').strip():
                errors.append(f"Missing required field: {field}")
        
        # Validate dates
        if row.get('Start Date') and row.get('Due Date'):
            try:
                start_date = datetime.strptime(row['Start Date'], config.tasks.date_format).date()
                due_date = datetime.strptime(row['Due Date'], config.tasks.date_format).date()
                
                if start_date > due_date:
                    errors.append(f"Start date ({start_date}) cannot be after due date ({due_date})")
                    
            except ValueError as e:
                errors.append(f"Invalid date format: {e}")
        
        # Validate task ID format
        task_id = row.get('Task ID', '').strip()
        if task_id and not re.match(r'^[A-Za-z0-9_-]+$', task_id):
            errors.append(f"Task ID '{task_id}' contains invalid characters (only letters, numbers, hyphens, underscores allowed)")
        
        return errors
    
    @staticmethod
    def validate_timeline(timeline: ProjectTimeline) -> List[str]:
        """Validate a complete timeline and return list of errors."""
        errors = []
        
        # Check for duplicate task IDs
        task_ids = [task.id for task in timeline.tasks]
        if len(task_ids) != len(set(task_ids)):
            duplicates = [tid for tid in set(task_ids) if task_ids.count(tid) > 1]
            errors.append(f"Duplicate task IDs found: {duplicates}")
        
        # Check for tasks with dependencies that don't exist
        all_task_ids = set(task.id for task in timeline.tasks)
        for task in timeline.tasks:
            if task.dependencies:
                # Simple dependency parsing - could be enhanced
                dep_ids = [dep.strip() for dep in task.dependencies.split(',')]
                for dep_id in dep_ids:
                    if dep_id and dep_id not in all_task_ids:
                        errors.append(f"Task '{task.name}' depends on non-existent task ID: {dep_id}")
        
        return errors
