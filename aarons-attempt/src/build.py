#!/usr/bin/env python3
"""
Enhanced build system inspired by latex-yearly-planner.
Provides automated building with multiple template and device support.
"""

import argparse
import logging
import os
import subprocess
import sys
from pathlib import Path
from datetime import datetime
from typing import List, Dict, Any, Optional
try:
    from tqdm import tqdm
except ImportError:
    # Fallback if tqdm is not available
    class tqdm:
        def __init__(self, iterable=None, **kwargs):
            self.iterable = iterable
        def __enter__(self):
            return self
        def __exit__(self, *args):
            pass
        def update(self, n=1):
            pass
        def __iter__(self):
            return iter(self.iterable) if self.iterable else iter([])

from .config_manager import config_manager
from .template_generators import TemplateGeneratorFactory
from .data_processor import DataProcessor
from .export_system import ExportSystem


class BuildSystem:
    """Enhanced build system for LaTeX Gantt chart generation."""
    
    def __init__(self):
        """Initialize enhanced build system."""
        self.setup_logging()
        self.logger = logging.getLogger(__name__)
        self.data_processor = DataProcessor()
        self.config_manager = config_manager
        self.export_system = ExportSystem()
        
        # Build configuration
        self.build_dir = Path("build")
        self.output_dir = Path("output")
        self.temp_dir = Path("temp")
        
        # Create directories
        self._create_directories()
    
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
    
    def _create_directories(self) -> None:
        """Create necessary directories."""
        for directory in [self.build_dir, self.output_dir, self.temp_dir]:
            directory.mkdir(exist_ok=True)
    
    def build_single(self, input_file: str, template_type: str = "gantt_timeline",
                    device_profile: str = None, color_scheme: str = None,
                    title: str = None, output_name: str = None) -> bool:
        """Build a single document with specified parameters."""
        try:
            self.logger.info(f"Building single document: {input_file}")
            
            # Generate timestamp for unique filenames
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            
            # Determine output filename
            if output_name:
                base_name = output_name
            else:
                base_name = f"Timeline_{timestamp}"
            
            # Process CSV data
            self.logger.info("Processing CSV data...")
            timeline = self.data_processor.process_csv_to_timeline(
                input_file, title or "Project Timeline"
            )
            
            # Generate LaTeX content
            self.logger.info("Generating LaTeX content...")
            generator = TemplateGeneratorFactory.create_generator(
                template_type, self.config_manager
            )
            latex_content = generator.generate_document(
                timeline, template_type, device_profile, color_scheme
            )
            
            # Write LaTeX file
            tex_file = self.output_dir / "tex" / f"{base_name}.tex"
            tex_file.parent.mkdir(exist_ok=True)
            
            with open(tex_file, 'w', encoding='utf-8') as f:
                f.write(latex_content)
            
            self.logger.info(f"LaTeX file written: {tex_file}")
            
            # Compile to PDF
            pdf_file = self.output_dir / "pdf" / f"{base_name}.pdf"
            pdf_file.parent.mkdir(exist_ok=True)
            
            success = self._compile_latex(tex_file, pdf_file)
            
            if success:
                self.logger.info(f"✅ Successfully built: {pdf_file}")
                return True
            else:
                self.logger.error(f"❌ Failed to compile: {tex_file}")
                return False
                
        except Exception as e:
            self.logger.error(f"❌ Error building document: {e}")
            return False
    
    def build_multiple(self, input_file: str, templates: List[str] = None,
                      devices: List[str] = None, color_schemes: List[str] = None,
                      title: str = None) -> Dict[str, bool]:
        """Build multiple documents with different configurations."""
        results = {}
        
        # Default configurations
        templates = templates or ["gantt_timeline"]
        devices = devices or [None]  # None means use default
        color_schemes = color_schemes or [None]  # None means use default
        
        self.logger.info(f"Building multiple documents with {len(templates)} templates")
        
        for template in templates:
            for device in devices:
                for color_scheme in color_schemes:
                    # Generate unique name
                    name_parts = [template]
                    if device:
                        name_parts.append(device)
                    if color_scheme:
                        name_parts.append(color_scheme)
                    
                    output_name = "_".join(name_parts)
                    
                    success = self.build_single(
                        input_file, template, device, color_scheme, title, output_name
                    )
                    
                    results[output_name] = success
        
        return results
    
    def build_all_templates(self, input_file: str, title: str = None) -> Dict[str, bool]:
        """Build all available templates."""
        templates = self.config_manager.list_templates()
        return self.build_multiple(input_file, templates, title=title)
    
    def build_all_devices(self, input_file: str, template_type: str = "gantt_timeline",
                         title: str = None) -> Dict[str, bool]:
        """Build for all available device profiles."""
        devices = self.config_manager.list_device_profiles()
        return self.build_multiple(input_file, [template_type], devices, title=title)
    
    def build_multiple_formats(self, input_file: str, template_type: str = "gantt_timeline",
                              device_profile: str = None, color_scheme: str = None,
                              title: str = None, formats: List[str] = None) -> Dict[str, bool]:
        """Build document in multiple export formats."""
        if formats is None:
            formats = ['pdf', 'svg', 'html', 'png']
        
        self.logger.info(f"Building multiple formats: {', '.join(formats)}")
        
        try:
            # Process CSV data
            with tqdm(total=100, desc="Processing data") as pbar:
                timeline = self.data_processor.process_csv_to_timeline(
                    input_file, title or "Project Timeline"
                )
                pbar.update(30)
            
            # Generate timestamp for unique filenames
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            base_name = f"Timeline_{timestamp}"
            
            results = {}
            
            # Export to each format
            for format_type in tqdm(formats, desc="Exporting formats"):
                output_path = self.output_dir / format_type / f"{base_name}.{format_type}"
                output_path.parent.mkdir(exist_ok=True)
                
                if format_type == 'pdf':
                    results[format_type] = self.export_system.export_to_pdf(
                        timeline, str(output_path), template_type, device_profile, color_scheme
                    )
                elif format_type == 'svg':
                    results[format_type] = self.export_system.export_to_svg(
                        timeline, str(output_path)
                    )
                elif format_type == 'html':
                    results[format_type] = self.export_system.export_to_html(
                        timeline, str(output_path)
                    )
                elif format_type == 'png':
                    results[format_type] = self.export_system.export_to_png(
                        timeline, str(output_path)
                    )
                else:
                    self.logger.warning(f"Unsupported format: {format_type}")
                    results[format_type] = False
            
            return results
            
        except Exception as e:
            self.logger.error(f"❌ Error building multiple formats: {e}")
            return {format_type: False for format_type in formats}
    
    def _compile_latex(self, tex_file: Path, pdf_file: Path) -> bool:
        """Compile LaTeX file to PDF."""
        try:
            self.logger.info(f"Compiling LaTeX: {tex_file} -> {pdf_file}")
            
            # Change to output directory for compilation
            original_cwd = os.getcwd()
            os.chdir(tex_file.parent)
            
            # Run pdflatex
            cmd = [
                "pdflatex",
                "-interaction=nonstopmode",
                str(tex_file.name)
            ]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            # Restore original directory
            os.chdir(original_cwd)
            
            if result.returncode == 0:
                self.logger.info("LaTeX compilation successful")
                # Move the generated PDF to the correct location
                generated_pdf = tex_file.with_suffix('.pdf')
                if generated_pdf.exists():
                    generated_pdf.rename(pdf_file)
                else:
                    self.logger.warning(f"Could not find generated PDF at {generated_pdf}")
                return True
            else:
                self.logger.error(f"LaTeX compilation failed. Stderr:\n{result.stderr}\nStdout:\n{result.stdout}")
                return False
                
        except FileNotFoundError:
            self.logger.error("pdflatex not found. Please install a LaTeX distribution.")
            return False
        except Exception as e:
            self.logger.error(f"Error compiling LaTeX: {e}")
            return False
    
    def clean(self) -> None:
        """Clean build artifacts."""
        self.logger.info("Cleaning build artifacts...")
        
        # Remove temporary files
        for pattern in ["*.aux", "*.log", "*.out", "*.toc", "*.fdb_latexmk", "*.fls", "*.synctex.gz"]:
            for file_path in self.temp_dir.glob(pattern):
                file_path.unlink()
        
        # Remove build directory
        if self.build_dir.exists():
            import shutil
            shutil.rmtree(self.build_dir)
            self.build_dir.mkdir(exist_ok=True)
        
        self.logger.info("Clean completed")
    
    def list_configurations(self) -> None:
        """List all available configurations."""
        print("Available Templates:")
        for template_id in self.config_manager.list_templates():
            template = self.config_manager.get_template(template_id)
            print(f"  {template_id}: {template.name}")
        
        print("\nAvailable Device Profiles:")
        for device_id in self.config_manager.list_device_profiles():
            device = self.config_manager.get_device_profile(device_id)
            print(f"  {device_id}: {device.name}")
        
        print("\nAvailable Color Schemes:")
        for scheme_id in self.config_manager.list_color_schemes():
            scheme = self.config_manager.get_color_scheme(scheme_id)
            print(f"  {scheme_id}: {scheme.name}")


def create_argument_parser() -> argparse.ArgumentParser:
    """Create and configure the argument parser."""
    parser = argparse.ArgumentParser(
        description="Enhanced build system for LaTeX Gantt chart generation",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s single data.csv                           # Build single document
  %(prog)s single data.csv -t monthly_calendar       # Build with specific template
  %(prog)s single data.csv -d supernote_a5x          # Build for specific device
  %(prog)s multiple data.csv                         # Build all templates
  %(prog)s all-templates data.csv                    # Build all templates
  %(prog)s all-devices data.csv                      # Build for all devices
  %(prog)s clean                                      # Clean build artifacts
  %(prog)s list                                       # List configurations
        """
    )
    
    subparsers = parser.add_subparsers(dest='command', help='Available commands')
    
    # Single build command
    single_parser = subparsers.add_parser('single', help='Build single document')
    single_parser.add_argument('input', help='Input CSV file')
    single_parser.add_argument('-t', '--template', default='gantt_timeline',
                              help='Template type to use')
    single_parser.add_argument('-d', '--device', help='Device profile to use')
    single_parser.add_argument('-c', '--color-scheme', help='Color scheme to use')
    single_parser.add_argument('--title', help='Document title')
    single_parser.add_argument('-o', '--output', help='Output filename (without extension)')
    
    # Multiple build command
    multiple_parser = subparsers.add_parser('multiple', help='Build multiple documents')
    multiple_parser.add_argument('input', help='Input CSV file')
    multiple_parser.add_argument('-t', '--templates', nargs='+', help='Templates to build')
    multiple_parser.add_argument('-d', '--devices', nargs='+', help='Device profiles to build')
    multiple_parser.add_argument('-c', '--color-schemes', nargs='+', help='Color schemes to build')
    multiple_parser.add_argument('--title', help='Document title')
    
    # All templates command
    all_templates_parser = subparsers.add_parser('all-templates', help='Build all templates')
    all_templates_parser.add_argument('input', help='Input CSV file')
    all_templates_parser.add_argument('--title', help='Document title')
    
    # All devices command
    all_devices_parser = subparsers.add_parser('all-devices', help='Build for all devices')
    all_devices_parser.add_argument('input', help='Input CSV file')
    all_devices_parser.add_argument('-t', '--template', default='gantt_timeline',
                                   help='Template type to use')
    all_devices_parser.add_argument('--title', help='Document title')
    
    # Multiple formats command
    multiple_formats_parser = subparsers.add_parser('multiple-formats', help='Build in multiple export formats')
    multiple_formats_parser.add_argument('input', help='Input CSV file')
    multiple_formats_parser.add_argument('-t', '--template', default='gantt_timeline',
                                        help='Template type to use')
    multiple_formats_parser.add_argument('-d', '--device', help='Device profile to use')
    multiple_formats_parser.add_argument('-c', '--color-scheme', help='Color scheme to use')
    multiple_formats_parser.add_argument('--title', help='Document title')
    multiple_formats_parser.add_argument('--formats', nargs='+', 
                                        choices=['pdf', 'svg', 'html', 'png'],
                                        default=['pdf', 'svg', 'html', 'png'],
                                        help='Export formats to generate')
    
    # Clean command
    subparsers.add_parser('clean', help='Clean build artifacts')
    
    # List command
    subparsers.add_parser('list', help='List available configurations')
    
    # Global options
    parser.add_argument('--verbose', '-v', action='store_true', help='Enable verbose logging')
    parser.add_argument('--quiet', '-q', action='store_true', help='Suppress all output except errors')
    
    return parser


def main() -> int:
    """Main entry point for the build system."""
    parser = create_argument_parser()
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        return 1
    
    # Create build system
    build_system = BuildSystem()
    
    # Adjust logging level
    if args.verbose:
        build_system.setup_logging("DEBUG")
    elif args.quiet:
        build_system.setup_logging("ERROR")
    
    try:
        if args.command == 'single':
            success = build_system.build_single(
                args.input, args.template, args.device, 
                args.color_scheme, args.title, args.output
            )
            return 0 if success else 1
            
        elif args.command == 'multiple':
            results = build_system.build_multiple(
                args.input, args.templates, args.devices, 
                args.color_schemes, args.title
            )
            success_count = sum(1 for success in results.values() if success)
            total_count = len(results)
            print(f"Built {success_count}/{total_count} documents successfully")
            return 0 if success_count == total_count else 1
            
        elif args.command == 'all-templates':
            results = build_system.build_all_templates(args.input, args.title)
            success_count = sum(1 for success in results.values() if success)
            total_count = len(results)
            print(f"Built {success_count}/{total_count} templates successfully")
            return 0 if success_count == total_count else 1
            
        elif args.command == 'all-devices':
            results = build_system.build_all_devices(args.input, args.template, args.title)
            success_count = sum(1 for success in results.values() if success)
            total_count = len(results)
            print(f"Built {success_count}/{total_count} device profiles successfully")
            return 0 if success_count == total_count else 1
            
        elif args.command == 'multiple-formats':
            results = build_system.build_multiple_formats(
                args.input, args.template, args.device, 
                args.color_scheme, args.title, args.formats
            )
            success_count = sum(1 for success in results.values() if success)
            total_count = len(results)
            print(f"Built {success_count}/{total_count} formats successfully")
            for format_type, success in results.items():
                status = "✅" if success else "❌"
                print(f"  {status} {format_type.upper()}")
            return 0 if success_count == total_count else 1
            
        elif args.command == 'clean':
            build_system.clean()
            return 0
            
        elif args.command == 'list':
            build_system.list_configurations()
            return 0
            
        else:
            parser.print_help()
            return 1
            
    except KeyboardInterrupt:
        print("\nBuild interrupted by user")
        return 130
    except Exception as e:
        print(f"Unexpected error: {e}")
        return 1


if __name__ == "__main__":
    sys.exit(main())
