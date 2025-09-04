#!/usr/bin/env python3
"""
Shared utilities for the LaTeX Gantt chart generator.
Essential shared functions and common patterns.
"""

import logging
import sys
from typing import List, Dict, Any
from pathlib import Path


class LoggingSetup:
    """Setup logging configuration for the application."""
    
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
    
    @staticmethod
    def generate_color_definitions(color_set: str = "standard") -> str:
        """Generate LaTeX color definitions."""
        colors = LaTeXUtilities.STANDARD_COLORS
        color_defs = []
        for name, rgb in colors:
            color_defs.append(f"\\definecolor{{{name}}}{{RGB}}{{{rgb[0]}, {rgb[1]}, {rgb[2]}}}")
        return '\n'.join(color_defs)
    
    @staticmethod
    def generate_document_header(template: Any, device_profile: Any, 
                               packages: List[str] = None) -> str:
        """Generate LaTeX document header with packages and styling."""
        if packages is None:
            packages = [
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
        
        return f"""\\documentclass[portrait,a4paper]{{article}}
\\usepackage[utf8]{{inputenc}}
\\usepackage[T1]{{fontenc}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage[portrait,margin=0.5in]{{geometry}}
\\usepackage{{tikz}}
\\usepackage{{xcolor}}
\\usepackage{{array}}
\\usepackage{{fancyhdr}}
\\usepackage{{hyperref}}
\\usepackage{{bookmark}}
\\usepackage{{enumitem}}
\\usepackage{{longtable}}
\\usepackage{{multirow}}
\\usepackage{{colortbl}}

% Enhanced TikZ libraries
\\usetikzlibrary{{arrows.meta,shapes.geometric,positioning,calc,decorations.pathmorphing,patterns,shadows,fit,backgrounds,matrix,chains,scopes,pgfgantt}}

% Page setup
\\pagestyle{{fancy}}
\\fancyhf{{}}
\\fancyhead[L]{{Project Timeline}}
\\fancyhead[R]{{\\today}}
\\fancyfoot[C]{{\\thepage}}
\\renewcommand{{\\headrulewidth}}{{0.4pt}}

% Hyperlink setup
\\hypersetup{{
    colorlinks=true,
    linkcolor=blue,
    urlcolor=blue,
    citecolor=blue,
    bookmarksopen=true,
    bookmarksnumbered=true
}}

% Use Helvetica for sans-serif
\\renewcommand{{\\familydefault}}{{\\sfdefault}}

{LaTeXUtilities.generate_color_definitions()}

\\begin{{document}}
"""
    
    @staticmethod
    def generate_document_footer() -> str:
        """Generate LaTeX document footer."""
        return "\\end{document}"


class DirectoryManager:
    """Directory management utilities."""
    
    @staticmethod
    def create_directories(directories: List[str]) -> None:
        """Create directories if they don't exist."""
        for directory in directories:
            Path(directory).mkdir(parents=True, exist_ok=True)
    
    @staticmethod
    def get_standard_directories() -> Dict[str, str]:
        """Get standard directory structure."""
        return {
            "build": "build",
            "output": "output",
            "output_tex": "output/tex",
            "output_pdf": "output/pdf",
            "temp": "temp"
        }