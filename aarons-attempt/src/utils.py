#!/usr/bin/env python3
"""
Utility functions for the LaTeX Gantt Chart Generator.
Contains helper functions and utilities.
"""

import argparse
import logging
from pathlib import Path
from typing import List, Tuple

from .config import config


def create_argument_parser() -> argparse.ArgumentParser:
    """Create argument parser."""
    parser = argparse.ArgumentParser(
        description="LaTeX Gantt Chart Generator",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s                                    # Use default files
  %(prog)s --input data.csv --output out.tex  # Specify custom files
  %(prog)s --title "My Project"               # Custom title
        """
    )
    
    parser.add_argument(
        "--input",
        default=config.default_input,
        help=f"Input CSV file (default: {config.default_input})"
    )
    parser.add_argument(
        "--output",
        default=config.default_output,
        help=f"Output LaTeX file (default: {config.default_output})"
    )
    parser.add_argument(
        "--title",
        default=config.latex.default_title,
        help=f"Document title (default: {config.latex.default_title})"
    )
    parser.add_argument("--prisma", action="store_true", help="Include PRISMA flow diagram")
    parser.add_argument("--verbose", "-v", action="store_true", help="Enable verbose logging")
    parser.add_argument("--quiet", "-q", action="store_true", help="Suppress all output except errors")
    
    return parser


def setup_logging(level: str = "INFO") -> None:
    """Setup logging configuration."""
    log_level = getattr(logging, level.upper(), logging.INFO)
    
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    
    console_handler = logging.StreamHandler()
    console_handler.setLevel(log_level)
    console_handler.setFormatter(formatter)
    
    root_logger = logging.getLogger()
    root_logger.setLevel(log_level)
    root_logger.addHandler(console_handler)


def validate_arguments(args: argparse.Namespace) -> Tuple[bool, List[str]]:
    """Validate command line arguments."""
    errors = []
    
    # Validate input file
    if not Path(args.input).exists():
        errors.append(f"Input file does not exist: {args.input}")
    
    # Validate output directory
    output_dir = Path(args.output).parent
    if not output_dir.exists():
        try:
            output_dir.mkdir(parents=True, exist_ok=True)
        except Exception as e:
            errors.append(f"Cannot create output directory {output_dir}: {e}")
    
    # Validate title
    if not args.title or not args.title.strip():
        errors.append("Title cannot be empty")
    
    return not errors, errors
