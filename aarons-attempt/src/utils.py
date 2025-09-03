#!/usr/bin/env python3
"""
Shared utility functions and classes to reduce code duplication.
Provides common functionality used across multiple modules.
"""

import logging
import sys
from typing import Any, Dict, List, Tuple
from dataclasses import dataclass


class LoggingSetup:
    """Centralized logging setup to eliminate duplication."""
    
    @staticmethod
    def setup_logging(level: str = "INFO") -> None:
        """Setup logging configuration."""
        log_level = getattr(logging, level.upper(), logging.INFO)
        
        # Create formatter
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
        
        # Setup console handler
        console_handler = logging.StreamHandler(sys.stdout)
        console_handler.setLevel(log_level)
        console_handler.setFormatter(formatter)
        
        # Setup root logger
        root_logger = logging.getLogger()
        root_logger.setLevel(log_level)
        root_logger.addHandler(console_handler)
        
        # Reduce noise from some libraries
        logging.getLogger('matplotlib').setLevel(logging.WARNING)
        logging.getLogger('PIL').setLevel(logging.WARNING)


class LaTeXUtilities:
    """Shared LaTeX generation utilities."""
    
    # Standard color definitions used across templates
    STANDARD_COLORS = [
        ("accountability", (139, 92, 246)),
        ("administrative", (107, 114, 128)),
        ("blocked", (239, 68, 68)),
        ("completed", (34, 197, 94)),
        ("inprogress", (251, 146, 60)),
        ("milestone", (147, 51, 234)),
        ("other", (156, 163, 175)),
        ("planned", (59, 130, 246)),
        ("researchcore", (59, 130, 246)),
        ("researchexp", (16, 185, 129)),
        ("researchout", (245, 158, 11)),
        ("service", (236, 72, 153)),
    ]
    
    # Simplified color definitions for calendar/planner templates
    SIMPLIFIED_COLORS = [
        ("task", (59, 130, 246)),
        ("milestone", (147, 51, 234)),
        ("completed", (34, 197, 94)),
        ("inprogress", (251, 146, 60)),
        ("blocked", (239, 68, 68)),
        ("grid", (200, 200, 200)),
    ]
    
    @staticmethod
    def generate_color_definitions(color_set: str = "standard") -> str:
        """Generate LaTeX color definitions."""
        colors = LaTeXUtilities.STANDARD_COLORS if color_set == "standard" else LaTeXUtilities.SIMPLIFIED_COLORS
        
        color_defs = []
        for name, rgb in colors:
            color_defs.append(f"\\definecolor{{{name}}}{{RGB}}{{{rgb[0]}, {rgb[1]}, {rgb[2]}}}")
        
        return '\n'.join(color_defs)
    
    @staticmethod
    def generate_document_header(template: Any, device_profile: Any, 
                               packages: List[str] = None) -> str:
        """Generate LaTeX document header with customizable packages."""
        page_size = device_profile.get_layout_value('page_size', 'a4paper')
        orientation = template.orientation
        margin = device_profile.get_layout_value('margin', '0.5in')
        
        # Default packages
        if packages is None:
            packages = [
                "tikz", "pgfplots", "xcolor", "enumitem", "booktabs", 
                "array", "longtable", "fancyhdr", "graphicx", "amsmath", 
                "amsfonts", "amssymb", "ragged2e"
            ]
        
        # Generate package includes
        package_includes = []
        for package in packages:
            package_includes.append(f"\\usepackage{{{package}}}")
        
        return f"""\\documentclass[{orientation},{page_size}]{{article}}
\\usepackage{{[utf8]{{inputenc}}}}
\\usepackage{{[T1]{{fontenc}}}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage{{[{orientation},margin={margin}]{{geometry}}}}
{chr(10).join(package_includes)}

% Page setup
\\pagestyle{{empty}}
\\setlength{{\\parskip}}{{0.5em}}

% Table formatting
\\setlength{{\\tabcolsep}}{{1pt}}
\\renewcommand{{\\arraystretch}}{{1.0}}

% Use Helvetica for sans-serif
\\renewcommand{{\\familydefault}}{{\\sfdefault}}

% Color definitions
{LaTeXUtilities.generate_color_definitions()}

\\begin{{document}}
"""
    
    @staticmethod
    def generate_simple_document_header(template: Any, device_profile: Any, 
                                      packages: List[str] = None) -> str:
        """Generate simplified LaTeX document header for calendar/planner templates."""
        page_size = device_profile.get_layout_value('page_size', 'a4paper')
        orientation = template.orientation
        margin = device_profile.get_layout_value('margin', '0.5in')
        
        # Default packages for simple templates
        if packages is None:
            packages = ["tikz", "xcolor", "array"]
        
        # Generate package includes
        package_includes = []
        for package in packages:
            package_includes.append(f"\\usepackage{{{package}}}")
        
        return f"""\\documentclass[{orientation},{page_size}]{{article}}
\\usepackage{{[utf8]{{inputenc}}}}
\\usepackage{{[T1]{{fontenc}}}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage{{[{orientation},margin={margin}]{{geometry}}}}
{chr(10).join(package_includes)}

% Page setup
\\pagestyle{{empty}}
\\setlength{{\\parskip}}{{0.5em}}

% Use Helvetica for sans-serif
\\renewcommand{{\\familydefault}}{{\\sfdefault}}

% Color definitions
{LaTeXUtilities.generate_color_definitions("simplified")}

\\begin{{document}}
"""
    
    @staticmethod
    def generate_document_footer() -> str:
        """Generate LaTeX document footer."""
        return "\\end{document}"


@dataclass
class CommonImports:
    """Common import patterns to reduce duplication."""
    
    # Standard imports used across multiple modules
    STANDARD_IMPORTS = [
        "import argparse",
        "import logging", 
        "import os",
        "import sys",
        "from pathlib import Path",
        "from typing import Optional, List, Dict, Any",
        "from datetime import datetime, date, timedelta"
    ]
    
    # Application-specific imports
    APP_IMPORTS = [
        "from .data_processor import DataProcessor",
        "from .latex_generator import LaTeXGenerator", 
        "from .config import config",
        "from .config_manager import ConfigManager, config_manager",
        "from .template_generators import TemplateGeneratorFactory"
    ]
    
    # Build system imports
    BUILD_IMPORTS = [
        "from src.config_manager import config_manager",
        "from src.template_generators import TemplateGeneratorFactory",
        "from src.data_processor import DataProcessor"
    ]
    
    # Template generator imports
    TEMPLATE_IMPORTS = [
        "from .models import Task, ProjectTimeline, MonthInfo",
        "from .config_manager import ConfigManager, config_manager",
        "from .latex_generator import LaTeXEscaper"
    ]


class DirectoryManager:
    """Centralized directory management to reduce duplication."""
    
    @staticmethod
    def create_directories(directories: List[str]) -> None:
        """Create multiple directories if they don't exist."""
        for directory in directories:
            Path(directory).mkdir(exist_ok=True)
    
    @staticmethod
    def get_standard_directories() -> Dict[str, str]:
        """Get standard directory paths used across the application."""
        return {
            "build": "build",
            "output": "output", 
            "temp": "temp",
            "config": "config",
            "input": "input"
        }
