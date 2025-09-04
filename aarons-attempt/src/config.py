#!/usr/bin/env python3
"""
Consolidated configuration for the LaTeX Gantt Chart Generator.
"""

from dataclasses import dataclass, field
from typing import Dict, List, Tuple


@dataclass
class ColorScheme:
    """Color definitions for different task categories."""
    
    # Task category colors
    milestone: Tuple[int, int, int] = (147, 51, 234)
    researchcore: Tuple[int, int, int] = (59, 130, 246)
    researchexp: Tuple[int, int, int] = (16, 185, 129)
    researchout: Tuple[int, int, int] = (245, 158, 11)
    administrative: Tuple[int, int, int] = (107, 114, 128)
    accountability: Tuple[int, int, int] = (139, 92, 246)
    service: Tuple[int, int, int] = (236, 72, 153)
    other: Tuple[int, int, int] = (156, 163, 175)


@dataclass
class TaskConfig:
    """Configuration for task processing and categorization."""
    
    # Category mapping
    category_keywords: Dict[str, List[str]] = field(default_factory=lambda: {
        "researchcore": ["PROPOSAL"],
        "researchexp": ["LASER", "EXPERIMENTAL"],
        "researchout": ["PUBLICATION", "PRESENTATION"],
        "administrative": ["ADMINISTRATIVE"],
        "accountability": ["ACCOUNTABILITY"],
        "service": ["SERVICE"]
    })


@dataclass
class LaTeXConfig:
    """Configuration for LaTeX document generation."""
    
    # Document metadata
    default_title: str = "Project Timeline"
    subtitle: str = "PhD Research Timeline"


@dataclass
class AppConfig:
    """Main application configuration."""
    
    colors: ColorScheme = field(default_factory=ColorScheme)
    tasks: TaskConfig = field(default_factory=TaskConfig)
    latex: LaTeXConfig = field(default_factory=LaTeXConfig)
    
    # File paths
    default_input: str = "../input/data.cleaned.csv"
    default_output: str = "output/tex/Timeline_template.tex"


# Global configuration instance
config = AppConfig()