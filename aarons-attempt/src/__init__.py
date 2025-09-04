#!/usr/bin/env python3
"""
LaTeX Gantt Chart Generator Package

A clean, consolidated tool for generating publication-quality LaTeX timelines
from CSV data. Perfect for PhD research, formal reports, and advisor meetings.
"""

__version__ = "2.0.0"
__author__ = "LaTeX Gantt Chart Generator"
__description__ = "Generate publication-quality LaTeX timelines from CSV data"

# Core imports - only import what's actually needed
from .core import Task, ProjectTimeline, DataProcessor, LaTeXGenerator, Application, main
from .config import config

__all__ = [
    # Core classes
    'Task', 'ProjectTimeline', 'DataProcessor', 'LaTeXGenerator', 'Application',
    
    # Configuration
    'config',
    
    # Main entry point
    'main',
]