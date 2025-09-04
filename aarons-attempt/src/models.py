#!/usr/bin/env python3
"""
Data models for the LaTeX Gantt Chart Generator.
Contains Task and ProjectTimeline data classes.
"""

from dataclasses import dataclass
from datetime import date
from typing import List

from .config import config


@dataclass
class Task:
    """Represents a single task with all its properties."""
    
    id: str
    name: str
    start_date: date
    due_date: date
    category: str
    dependencies: str = ""
    notes: str = ""
    is_milestone: bool = False
    
    def __post_init__(self):
        """Validate and clean task data after initialization."""
        self._validate_dates()
        self._determine_milestone()
    
    def _validate_dates(self) -> None:
        """Validate that dates are reasonable."""
        if self.start_date > self.due_date:
            raise ValueError(f"Task '{self.name}': start date ({self.start_date}) cannot be after due date ({self.due_date})")
    
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
        return get_category_color(self.category)


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


def get_category_color(category: str) -> str:
    """Get the color name for a given task category."""
    category_upper = category.upper()
    
    return next(
        (color_name for color_name, keywords in config.tasks.category_keywords.items()
         if any(keyword in category_upper for keyword in keywords)),
        "other"
    )
