#!/usr/bin/env python3
"""
LaTeX Gantt Chart Generator Package

A modular, well-structured tool for generating publication-quality LaTeX timelines
from CSV data. Perfect for PhD research, formal reports, and advisor meetings.
"""

__version__ = "2.0.0"
__author__ = "LaTeX Gantt Chart Generator"
__description__ = "Generate publication-quality LaTeX timelines from CSV data"

# Import main components for easy access
from .config import config, AppConfig, ColorScheme, CalendarConfig, TaskConfig, LaTeXConfig
from .models import Task, ProjectTimeline, MonthInfo, TaskValidator
from .data_processor import DataProcessor, CSVReader, TaskProcessor, TimelineBuilder
from .latex_generator import LaTeXGenerator, LaTeXEscaper, LaTeXDocumentGenerator
from .app import Application, main

__all__ = [
    # Configuration
    'config', 'AppConfig', 'ColorScheme', 'CalendarConfig', 'TaskConfig', 'LaTeXConfig',
    
    # Data Models
    'Task', 'ProjectTimeline', 'MonthInfo', 'TaskValidator',
    
    # Data Processing
    'DataProcessor', 'CSVReader', 'TaskProcessor', 'TimelineBuilder',
    
    # LaTeX Generation
    'LaTeXGenerator', 'LaTeXEscaper', 'LaTeXDocumentGenerator',
    
    # Application
    'Application', 'main',
]
