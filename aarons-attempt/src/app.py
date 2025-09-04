#!/usr/bin/env python3
"""
Main application for the LaTeX Gantt chart generator.
Provides the main entry point with comprehensive error handling and logging.
"""

import argparse
import logging
import os
import sys
from pathlib import Path
from typing import Optional

from .data_processor import DataProcessor
from .latex_generator import LaTeXGenerator
from .config import config
from .config_manager import ConfigManager, config_manager
from .template_generators import TemplateGeneratorFactory
from .session_manager import get_session_integration


class Application:
    """Main application class for the LaTeX Gantt chart generator."""
    
    def __init__(self):
        self.setup_logging()
        self.logger = logging.getLogger(__name__)
        self.data_processor = DataProcessor()
        self.latex_generator = LaTeXGenerator()
        self.config_manager = config_manager
        self.session_integration = get_session_integration()
    
    def setup_logging(self, level: str = "INFO") -> None:
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
    
    def validate_arguments(self, args: argparse.Namespace) -> tuple[bool, list[str]]:
        """Validate command line arguments."""
        errors = []
        
        # Validate input file
        is_valid, file_errors = self.data_processor.validate_input_file(args.input)
        if not is_valid:
            errors.extend(file_errors)
        
        # Validate output directory
        output_dir = Path(args.output).parent
        if not output_dir.exists():
            try:
                output_dir.mkdir(parents=True, exist_ok=True)
                self.logger.info(f"Created output directory: {output_dir}")
            except Exception as e:
                errors.append(f"Cannot create output directory {output_dir}: {e}")
        
        # Validate title
        if not args.title or not args.title.strip():
            errors.append("Title cannot be empty")
        
        return len(errors) == 0, errors
    
    def generate_latex_file(self, input_file: str, output_file: str, title: str, 
                          template_type: str = None, device_profile: str = None, 
                          color_scheme: str = None) -> bool:
        """Generate LaTeX file from CSV input with enhanced template support."""
        try:
            self.logger.info(f"Starting LaTeX generation: {input_file} -> {output_file}")
            
            # Process CSV data to timeline
            self.logger.info("Processing CSV data...")
            timeline = self.data_processor.process_csv_to_timeline(input_file, title)
            
            # Generate LaTeX content using template system
            self.logger.info("Generating LaTeX content with template system...")
            generator = TemplateGeneratorFactory.create_generator(
                template_type or 'gantt_timeline', 
                self.config_manager
            )
            latex_content = generator.generate_document(
                timeline, template_type, device_profile, color_scheme
            )
            
            # Write LaTeX file
            self.logger.info(f"Writing LaTeX file: {output_file}")
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(latex_content)
            
            self.logger.info(f"✅ Successfully generated LaTeX file: {output_file}")
            return True
            
        except Exception as e:
            self.logger.error(f"❌ Error generating LaTeX file: {e}")
            
            # * Handle error through APM debugging system
            try:
                error_context = {
                    "input_file": input_file,
                    "output_file": output_file,
                    "title": title,
                    "template_type": template_type,
                    "device_profile": device_profile,
                    "color_scheme": color_scheme
                }
                
                success, message, delegation_prompt = self.session_integration.handle_application_error(
                    e, error_context
                )
                
                if not success and delegation_prompt:
                    self.logger.info("Error requires Ad-Hoc debugging session")
                    self.logger.info("Delegation prompt generated for Ad-Hoc agent")
                    # * In a real implementation, this would trigger the Ad-Hoc session
                    # * For now, we just log the information
                
            except Exception as debug_error:
                self.logger.warning(f"Failed to handle error through debug system: {debug_error}")
            
            return False
    
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
            success = self.generate_latex_file(
                args.input, args.output, args.title,
                getattr(args, 'template', None),
                getattr(args, 'device', None),
                getattr(args, 'color_scheme', None)
            )
            if not success:
                return 1
            
            # Print compilation instructions
            print(f"\n📄 LaTeX file generated successfully!")
            print(f"📁 Output file: {args.output}")
            print(f"🔨 To compile: pdflatex {args.output}")
            print(f"📊 Timeline contains {len(self.data_processor.task_processor.process_tasks(self.data_processor.csv_reader.read_csv_data(args.input)))} tasks")
            
            return 0
            
        except KeyboardInterrupt:
            self.logger.info("Application interrupted by user")
            return 130
        except Exception as e:
            self.logger.error(f"Unexpected error: {e}")
            return 1


def create_argument_parser() -> argparse.ArgumentParser:
    """Create and configure the argument parser."""
    parser = argparse.ArgumentParser(
        description="Generate LaTeX calendar template file from CSV data",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s                                    # Use default files
  %(prog)s --input data.csv --output out.tex  # Specify custom files
  %(prog)s --title "My Project"               # Custom title
  %(prog)s --verbose                          # Enable debug logging
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
    
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Enable verbose logging"
    )
    
    parser.add_argument(
        "--quiet", "-q",
        action="store_true",
        help="Suppress all output except errors"
    )
    
    parser.add_argument(
        "--template", "-t",
        choices=["gantt_timeline", "monthly_calendar", "weekly_planner"],
        default="gantt_timeline",
        help="Template type to use (default: gantt_timeline)"
    )
    
    parser.add_argument(
        "--device", "-d",
        help="Device profile to use (e.g., supernote_a5x, remarkable_2, standard_print)"
    )
    
    parser.add_argument(
        "--color-scheme", "-c",
        help="Color scheme to use (e.g., academic, corporate, vibrant)"
    )
    
    parser.add_argument(
        "--list-templates",
        action="store_true",
        help="List available templates and exit"
    )
    
    parser.add_argument(
        "--list-devices",
        action="store_true",
        help="List available device profiles and exit"
    )
    
    parser.add_argument(
        "--list-color-schemes",
        action="store_true",
        help="List available color schemes and exit"
    )
    
    return parser


def main() -> int:
    """Main entry point for the application."""
    parser = create_argument_parser()
    args = parser.parse_args()
    
    # Create application instance
    app = Application()
    
    # Handle list commands
    if args.list_templates:
        print("Available templates:")
        for template_id in app.config_manager.list_templates():
            template = app.config_manager.get_template(template_id)
            print(f"  {template_id}: {template.name} - {template.description}")
        return 0
    
    if args.list_devices:
        print("Available device profiles:")
        for device_id in app.config_manager.list_device_profiles():
            device = app.config_manager.get_device_profile(device_id)
            print(f"  {device_id}: {device.name} - {device.description}")
        return 0
    
    if args.list_color_schemes:
        print("Available color schemes:")
        for scheme_id in app.config_manager.list_color_schemes():
            scheme = app.config_manager.get_color_scheme(scheme_id)
            print(f"  {scheme_id}: {scheme.name} - {scheme.description}")
        return 0
    
    # Adjust logging level based on arguments
    if args.verbose:
        app.setup_logging("DEBUG")
    elif args.quiet:
        app.setup_logging("ERROR")
    
    # Run the application
    return app.run(args)


if __name__ == "__main__":
    sys.exit(main())
