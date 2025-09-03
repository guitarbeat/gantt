#!/usr/bin/env python3
"""
Configuration management for the LaTeX Gantt chart generator.
Centralizes all styling, colors, and settings for easy customization.
"""

from dataclasses import dataclass
from typing import Dict, List, Tuple


@dataclass
class ColorScheme:
    """Color definitions for different task categories and statuses."""
    
    # Task category colors
    milestone: Tuple[int, int, int] = (147, 51, 234)
    researchcore: Tuple[int, int, int] = (59, 130, 246)
    researchexp: Tuple[int, int, int] = (16, 185, 129)
    researchout: Tuple[int, int, int] = (245, 158, 11)
    administrative: Tuple[int, int, int] = (107, 114, 128)
    accountability: Tuple[int, int, int] = (139, 92, 246)
    service: Tuple[int, int, int] = (236, 72, 153)
    other: Tuple[int, int, int] = (156, 163, 175)
    
    # Status colors
    completed: Tuple[int, int, int] = (34, 197, 94)
    inprogress: Tuple[int, int, int] = (251, 146, 60)
    blocked: Tuple[int, int, int] = (239, 68, 68)
    planned: Tuple[int, int, int] = (59, 130, 246)
    
    def to_latex_colors(self) -> str:
        """Generate LaTeX color definitions."""
        colors = []
        for attr_name in dir(self):
            if not attr_name.startswith('_') and not callable(getattr(self, attr_name)):
                rgb = getattr(self, attr_name)
                colors.append(f"\\definecolor{{{attr_name}}}{{RGB}}{{{rgb[0]}, {rgb[1]}, {rgb[2]}}}")
        return '\n'.join(colors)


@dataclass
class CalendarConfig:
    """Configuration for calendar layout and styling."""
    
    # Page layout
    page_orientation: str = "landscape"
    page_size: str = "a4paper"
    margin: str = "0.5in"
    
    # Calendar grid
    calendar_scale: float = 1.0
    calendar_width: int = 7
    calendar_height: int = 6
    max_tasks_per_day: int = 2
    max_task_name_length: int = 18
    
    # Typography
    title_font_size: str = "\\LARGE"
    month_font_size: str = "\\large"
    day_font_size: str = "\\large"
    task_font_size: str = "\\tiny"
    
    # Spacing
    title_spacing: str = "1cm"
    month_spacing: str = "0.5cm"
    task_spacing: str = "0.5cm"


@dataclass
class TaskConfig:
    """Configuration for task processing and categorization."""
    
    # Task name cleaning patterns
    name_prefixes_to_remove: List[str] = None
    
    # Category mapping
    category_keywords: Dict[str, List[str]] = None
    
    # Date format
    date_format: str = "%Y-%m-%d"
    display_date_format: str = "%m/%d"
    
    def __post_init__(self):
        if self.name_prefixes_to_remove is None:
            self.name_prefixes_to_remove = [
                "Milestone: ",
                "Draft ",
                "Complete "
            ]
        
        if self.category_keywords is None:
            self.category_keywords = {
                "researchcore": ["PROPOSAL"],
                "researchexp": ["LASER", "EXPERIMENTAL"],
                "researchout": ["PUBLICATION", "PRESENTATION"],
                "administrative": ["ADMINISTRATIVE"],
                "accountability": ["ACCOUNTABILITY"],
                "service": ["SERVICE"]
            }


@dataclass
class LaTeXConfig:
    """Configuration for LaTeX document generation."""
    
    # Document class and packages
    document_class: str = "article"
    packages: List[str] = None
    
    # Document metadata
    default_title: str = "Project Timeline"
    subtitle: str = "PhD Research Calendar"
    
    def __post_init__(self):
        if self.packages is None:
            self.packages = [
                "[utf8]{inputenc}",
                "[T1]{fontenc}",
                "lmodern",
                "helvet",
                "[landscape,margin=0.5in]{geometry}",
                "tikz",
                "pgfplots",
                "xcolor",
                "enumitem",
                "booktabs",
                "array",
                "longtable",
                "fancyhdr",
                "graphicx",
                "amsmath",
                "amsfonts",
                "amssymb",
                "ragged2e"
            ]


@dataclass
class AppConfig:
    """Main application configuration combining all sub-configurations."""
    
    colors: ColorScheme = None
    calendar: CalendarConfig = None
    tasks: TaskConfig = None
    latex: LaTeXConfig = None
    
    # File paths
    default_input: str = "input/data.cleaned.csv"
    default_output: str = "output/tex/Calendar_template.tex"
    
    def __post_init__(self):
        if self.colors is None:
            self.colors = ColorScheme()
        if self.calendar is None:
            self.calendar = CalendarConfig()
        if self.tasks is None:
            self.tasks = TaskConfig()
        if self.latex is None:
            self.latex = LaTeXConfig()


# Global configuration instance
config = AppConfig()


def get_category_color(category: str) -> str:
    """Get the color name for a given task category."""
    category_upper = category.upper()
    
    for color_name, keywords in config.tasks.category_keywords.items():
        if any(keyword in category_upper for keyword in keywords):
            return color_name
    
    return "other"


def get_task_marker(is_milestone: bool) -> str:
    """Get the LaTeX marker for a task based on its type."""
    return "$\\diamond$" if is_milestone else "$\\bullet$"
