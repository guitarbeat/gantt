#!/usr/bin/env python3
"""
Main application module for the LaTeX Gantt Chart Generator.
Contains the Application class and main entry point.
"""

import argparse
import logging
import sys
from pathlib import Path
from typing import List, Tuple

from .config import config
from .generator import LaTeXGenerator
from .models import ProjectTimeline
from .processor import DataProcessor
from .utils import create_argument_parser, setup_logging, validate_arguments


class Application:
    """Main application class."""
    
    def __init__(self):
        self.setup_logging()
        self.logger = logging.getLogger(__name__)
        self.data_processor = DataProcessor()
        self.latex_generator = LaTeXGenerator()
    
    def setup_logging(self, level: str = "INFO") -> None:
        """Setup logging configuration."""
        setup_logging(level)
    
    def run(self, args: argparse.Namespace) -> int:
        """Run the main application."""
        try:
            self.logger.info("Starting LaTeX Gantt Chart Generator")
            
            # Validate arguments
            is_valid, errors = self.validate_arguments(args)
            if not is_valid:
                self.logger.error("Argument validation failed:")
                for error in errors:
                    self.logger.error(f"  - {error}")
                return 1
            
            # Generate LaTeX file
            success = self.generate_latex_file(args.input, args.output, args.title, args.prisma)
            if not success:
                return 1
            
            print(f"\n📄 LaTeX file generated successfully!")
            print(f"📁 Output file: {args.output}")
            print(f"🔨 To compile: pdflatex {args.output}")
            
            return 0
            
        except KeyboardInterrupt:
            self.logger.info("Application interrupted by user")
            return 130
        except Exception as e:
            self.logger.error(f"Unexpected error: {e}")
            return 1
    
    def validate_arguments(self, args: argparse.Namespace) -> Tuple[bool, List[str]]:
        """Validate command line arguments."""
        errors = []
        
        # Validate input file
        is_valid, file_errors = self.data_processor.validate_input_file(args.input)
        if not is_valid:
            errors.extend(file_errors)
        
        # Validate output directory and title
        errors.extend(self._validate_output_directory(args.output))
        errors.extend(self._validate_title(args.title))
        
        return not errors, errors
    
    def _validate_output_directory(self, output_path: str) -> List[str]:
        """Validate output directory."""
        errors = []
        output_dir = Path(output_path).parent
        if not output_dir.exists():
            try:
                output_dir.mkdir(parents=True, exist_ok=True)
                self.logger.info(f"Created output directory: {output_dir}")
            except Exception as e:
                errors.append(f"Cannot create output directory {output_dir}: {e}")
        return errors
    
    def _validate_title(self, title: str) -> List[str]:
        """Validate title."""
        errors = []
        if not title or not title.strip():
            errors.append("Title cannot be empty")
        return errors
    
    def generate_latex_file(self, input_file: str, output_file: str, title: str, include_prisma: bool = False) -> bool:
        """Generate LaTeX file from CSV input."""
        try:
            self.logger.info(f"Starting LaTeX generation: {input_file} -> {output_file}")
            
            # Process CSV data to timeline
            self.logger.info("Processing CSV data...")
            timeline = self.data_processor.process_csv_to_timeline(input_file, title)
            
            # Generate LaTeX content
            self.logger.info("Generating LaTeX content...")
            latex_content = self.latex_generator.generate_document(timeline, include_prisma)
            
            # Write LaTeX file
            self.logger.info(f"Writing LaTeX file: {output_file}")
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(latex_content)
            
            self.logger.info(f"✅ Successfully generated LaTeX file: {output_file}")
            return True
            
        except Exception as e:
            self.logger.error(f"❌ Error generating LaTeX file: {e}")
            return False


def main() -> int:
    """Main entry point."""
    parser = create_argument_parser()
    args = parser.parse_args()
    
    # Create application instance
    app = Application()
    
    # Adjust logging level based on arguments
    if args.verbose:
        app.setup_logging("DEBUG")
    elif args.quiet:
        app.setup_logging("ERROR")
    
    # Run the application
    return app.run(args)


if __name__ == "__main__":
    sys.exit(main())
