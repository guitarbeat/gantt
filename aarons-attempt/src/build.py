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

from .config import config_manager
from .latex_generator import TemplateGeneratorFactory
from .data_processor import DataProcessor
from .models import ProjectTimeline


class BuildSystem:
    """Enhanced build system for LaTeX Gantt chart generation."""
    
    def __init__(self):
        """Initialize enhanced build system."""
        self.setup_logging()
        self.logger = logging.getLogger(__name__)
        self.data_processor = DataProcessor()
        self.config_manager = config_manager
        
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
                    results[format_type] = self.export_to_pdf(
                        timeline, str(output_path), template_type, device_profile, color_scheme
                    )
                elif format_type == 'svg':
                    results[format_type] = self.export_to_svg(
                        timeline, str(output_path)
                    )
                elif format_type == 'html':
                    results[format_type] = self.export_to_html(
                        timeline, str(output_path)
                    )
                elif format_type == 'png':
                    results[format_type] = self.export_to_png(
                        timeline, str(output_path)
                    )
                else:
                    self.logger.warning(f"Unsupported format: {format_type}")
                    results[format_type] = False
            
            return results
            
        except Exception as e:
            self.logger.error(f"❌ Error building multiple formats: {e}")
            return {format_type: False for format_type in formats}
    
    def export_to_pdf(self, timeline: ProjectTimeline, output_path: str, 
                     template_type: str = "gantt_timeline",
                     device_profile: str = None,
                     color_scheme: str = None) -> bool:
        """Export timeline to PDF with enhanced quality."""
        try:
            self.logger.info(f"Exporting to PDF: {output_path}")
            
            # Generate LaTeX content
            generator = TemplateGeneratorFactory.create_generator(
                template_type, self.config_manager
            )
            latex_content = generator.generate_document(
                timeline, template_type, device_profile, color_scheme
            )
            
            # Write LaTeX file
            tex_path = Path(output_path).with_suffix('.tex')
            with open(tex_path, 'w', encoding='utf-8') as f:
                f.write(latex_content)
            
            # Compile to PDF with enhanced settings
            success = self._compile_latex_to_pdf(tex_path, output_path)
            
            if success:
                self.logger.info(f"✅ Successfully exported PDF: {output_path}")
                return True
            else:
                self.logger.error(f"❌ Failed to compile PDF: {tex_path}")
                return False
                
        except Exception as e:
            self.logger.error(f"❌ Error exporting to PDF: {e}")
            return False
    
    def _compile_latex_to_pdf(self, tex_path: Path, pdf_path: str) -> bool:
        """Compile LaTeX to PDF with enhanced settings."""
        try:
            # Change to tex file directory
            original_cwd = Path.cwd()
            os.chdir(tex_path.parent)
            
            # Run pdflatex with enhanced settings
            cmd = [
                "pdflatex",
                "-interaction=nonstopmode",
                "-halt-on-error",
                "-file-line-error",
                "-output-directory", str(tex_path.parent),
                str(tex_path.name)
            ]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            # Restore original directory
            os.chdir(original_cwd)
            
            if result.returncode == 0:
                # Move the generated PDF to the correct location
                generated_pdf = tex_path.with_suffix('.pdf')
                if generated_pdf.exists():
                    generated_pdf.rename(pdf_path)
                    return True
                else:
                    self.logger.error(f"Generated PDF not found at {generated_pdf}")
                    return False
            else:
                self.logger.error(f"LaTeX compilation failed. Stderr:\n{result.stderr}")
                return False
                
        except FileNotFoundError:
            self.logger.error("pdflatex not found. Please install a LaTeX distribution.")
            return False
        except Exception as e:
            self.logger.error(f"Error compiling LaTeX: {e}")
            return False
    
    def export_to_svg(self, timeline: ProjectTimeline, output_path: str) -> bool:
        """Export timeline to SVG format."""
        try:
            self.logger.info(f"Exporting to SVG: {output_path}")
            
            # Generate LaTeX content optimized for SVG
            latex_content = self._generate_svg_optimized_latex(timeline)
            
            # Write LaTeX file
            tex_path = Path(output_path).with_suffix('.tex')
            with open(tex_path, 'w', encoding='utf-8') as f:
                f.write(latex_content)
            
            # Compile to SVG
            success = self._compile_latex_to_svg(tex_path, output_path)
            
            if success:
                self.logger.info(f"✅ Successfully exported SVG: {output_path}")
                return True
            else:
                self.logger.error(f"❌ Failed to compile SVG: {tex_path}")
                return False
                
        except Exception as e:
            self.logger.error(f"❌ Error exporting to SVG: {e}")
            return False
    
    def export_to_html(self, timeline: ProjectTimeline, output_path: str) -> bool:
        """Export timeline to HTML format with interactive features."""
        try:
            self.logger.info(f"Exporting to HTML: {output_path}")
            
            # Generate HTML content
            html_content = self._generate_html_content(timeline)
            
            # Write HTML file
            with open(output_path, 'w', encoding='utf-8') as f:
                f.write(html_content)
            
            self.logger.info(f"✅ Successfully exported HTML: {output_path}")
            return True
                
        except Exception as e:
            self.logger.error(f"❌ Error exporting to HTML: {e}")
            return False
    
    def export_to_png(self, timeline: ProjectTimeline, output_path: str, 
                     dpi: int = 300) -> bool:
        """Export timeline to PNG format."""
        try:
            self.logger.info(f"Exporting to PNG: {output_path}")
            
            # First export to PDF
            pdf_path = Path(output_path).with_suffix('.pdf')
            if not self.export_to_pdf(timeline, str(pdf_path)):
                return False
            
            # Convert PDF to PNG
            success = self._convert_pdf_to_png(pdf_path, output_path, dpi)
            
            if success:
                self.logger.info(f"✅ Successfully exported PNG: {output_path}")
                return True
            else:
                self.logger.error(f"❌ Failed to convert PDF to PNG")
                return False
                
        except Exception as e:
            self.logger.error(f"❌ Error exporting to PNG: {e}")
            return False
    
    def _compile_latex_to_svg(self, tex_path: Path, svg_path: str) -> bool:
        """Compile LaTeX to SVG format."""
        try:
            # First compile to PDF
            pdf_path = tex_path.with_suffix('.pdf')
            if not self._compile_latex_to_pdf(tex_path, str(pdf_path)):
                return False
            
            # Convert PDF to SVG using pdf2svg
            cmd = ["pdf2svg", str(pdf_path), svg_path]
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            if result.returncode == 0:
                return True
            else:
                self.logger.error(f"PDF to SVG conversion failed: {result.stderr}")
                return False
                
        except FileNotFoundError:
            self.logger.error("pdf2svg not found. Please install pdf2svg.")
            return False
        except Exception as e:
            self.logger.error(f"Error converting to SVG: {e}")
            return False
    
    def _convert_pdf_to_png(self, pdf_path: Path, png_path: str, dpi: int = 300) -> bool:
        """Convert PDF to PNG using ImageMagick."""
        try:
            cmd = [
                "convert",
                "-density", str(dpi),
                "-quality", "100",
                str(pdf_path),
                png_path
            ]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            if result.returncode == 0:
                return True
            else:
                self.logger.error(f"PDF to PNG conversion failed: {result.stderr}")
                return False
                
        except FileNotFoundError:
            self.logger.error("ImageMagick not found. Please install ImageMagick.")
            return False
        except Exception as e:
            self.logger.error(f"Error converting to PNG: {e}")
            return False
    
    def _generate_svg_optimized_latex(self, timeline: ProjectTimeline) -> str:
        """Generate LaTeX content optimized for SVG export."""
        return f"""\\documentclass[border=0pt]{{standalone}}
\\usepackage{{tikz}}
\\usepackage{{xcolor}}
\\usetikzlibrary{{arrows.meta,shapes.geometric,positioning,calc,decorations.pathmorphing,patterns,shadows,fit,backgrounds,matrix,chains,scopes}}

% SVG-optimized styles
\\tikzset{{
    task node/.style={{
        rectangle, 
        rounded corners=4pt,
        draw=black!30,
        fill=white,
        drop shadow={{shadow xshift=1pt, shadow yshift=-1pt, fill=black!20}},
        minimum height=0.7cm,
        minimum width=2cm,
        font=\\small\\bfseries,
        align=center
    }},
    milestone node/.style={{
        diamond,
        draw=purple!60,
        fill=purple!20,
        drop shadow={{shadow xshift=1pt, shadow yshift=-1pt, fill=purple!30}},
        minimum size=1cm,
        font=\\small\\bfseries,
        align=center
    }}
}}

\\begin{{document}}
\\begin{{tikzpicture}}[scale=1.2]
    % Timeline axis
    \\draw[thick, line width=3pt, color=blue!70] (0,0) -- (12,0);
    
    % Task bars
    y_pos = 1.5
    for i, task in enumerate(timeline.tasks):
        task_name = task.name.replace('_', '\\_')
        start_x = (task.start_date - timeline.start_date).days / 30.0 * 2
        end_x = (task.due_date - timeline.start_date).days / 30.0 * 2
        
        if task.is_milestone:
            # Milestone
            \\node[milestone node] at ({{start_x}},{{y_pos}}) {{{task_name}}};
        else:
            # Task bar
            \\draw[fill=blue!60, rounded corners=3pt] ({{start_x}},{{y_pos-0.3}}) rectangle ({{end_x}},{{y_pos+0.3}});
            \\node[font=\\small\\bfseries, text=white] at ({{(start_x+end_x)/2}},{{y_pos}}) {{{task_name}}};
        
        y_pos += 0.9
\\end{{tikzpicture}}
\\end{{document}}
"""
    
    def _generate_html_content(self, timeline: ProjectTimeline) -> str:
        """Generate HTML content with interactive features."""
        html_content = f"""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{timeline.title}</title>
    <style>
        body {{
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }}
        .container {{
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }}
        .header {{
            text-align: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 2px solid #3B82F6;
        }}
        .header h1 {{
            color: #1F2937;
            margin: 0;
            font-size: 2.5em;
        }}
        .header p {{
            color: #6B7280;
            margin: 10px 0 0 0;
        }}
        .timeline {{
            position: relative;
            margin: 40px 0;
        }}
        .timeline-axis {{
            height: 4px;
            background: linear-gradient(90deg, #3B82F6, #10B981);
            border-radius: 2px;
            position: relative;
        }}
        .task-bar {{
            position: absolute;
            height: 30px;
            background: #3B82F6;
            border-radius: 15px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            font-size: 12px;
            cursor: pointer;
            transition: all 0.3s ease;
        }}
        .task-bar:hover {{
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
        }}
        .milestone {{
            position: absolute;
            width: 20px;
            height: 20px;
            background: #9333EA;
            transform: rotate(45deg);
            border-radius: 3px;
            cursor: pointer;
        }}
        .task-details {{
            margin-top: 40px;
        }}
        .task-card {{
            background: #F9FAFB;
            border-left: 4px solid #3B82F6;
            padding: 15px;
            margin: 10px 0;
            border-radius: 0 8px 8px 0;
            transition: all 0.3s ease;
        }}
        .task-card:hover {{
            background: #F3F4F6;
            transform: translateX(5px);
        }}
        .task-name {{
            font-weight: bold;
            color: #1F2937;
            margin-bottom: 5px;
        }}
        .task-dates {{
            color: #6B7280;
            font-size: 0.9em;
        }}
        .progress-bar {{
            width: 100%;
            height: 8px;
            background: #E5E7EB;
            border-radius: 4px;
            overflow: hidden;
            margin: 20px 0;
        }}
        .progress-fill {{
            height: 100%;
            background: linear-gradient(90deg, #10B981, #3B82F6);
            width: 30%;
            transition: width 0.3s ease;
        }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{timeline.title}</h1>
            <p>Timeline: {timeline.start_date.strftime('%B %d, %Y')} - {timeline.end_date.strftime('%B %d, %Y')}</p>
            <p>Total Tasks: {len(timeline.tasks)} | Duration: {timeline.total_duration_days} days</p>
        </div>
        
        <div class="progress-bar">
            <div class="progress-fill"></div>
        </div>
        
        <div class="timeline">
            <div class="timeline-axis"></div>
            <!-- Task bars will be positioned here -->
        </div>
        
        <div class="task-details">
            <h2>Task Details</h2>
"""
        
        for i, task in enumerate(timeline.tasks):
            task_name = task.name.replace('"', '&quot;')
            html_content += f"""
            <div class="task-card">
                <div class="task-name">{task_name}</div>
                <div class="task-dates">{task.start_date.strftime('%B %d, %Y')} - {task.due_date.strftime('%B %d, %Y')} ({task.duration_days} days)</div>
                {f'<div style="color: #6B7280; font-size: 0.9em; margin-top: 5px;">{task.notes}</div>' if task.notes else ''}
            </div>
"""
        
        html_content += """
        </div>
    </div>
    
    <script>
        // Add interactive features
        document.addEventListener('DOMContentLoaded', function() {
            // Add click handlers for task cards
            const taskCards = document.querySelectorAll('.task-card');
            taskCards.forEach(card => {
                card.addEventListener('click', function() {
                    this.style.background = '#DBEAFE';
                    setTimeout(() => {
                        this.style.background = '#F9FAFB';
                    }, 200);
                });
            });
        });
    </script>
</body>
</html>
"""
        
        return html_content
    
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
