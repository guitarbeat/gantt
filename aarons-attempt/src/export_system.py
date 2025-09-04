#!/usr/bin/env python3
"""
Enhanced export system for multiple output formats.
Supports PDF, SVG, PNG, HTML, and other formats.
"""

import logging
import subprocess
import tempfile
import os
from pathlib import Path
from typing import Dict, Any, Optional, List
from datetime import datetime

from .models import ProjectTimeline
from .latex_generator import LaTeXGenerator
from .interactive_generator import EnhancedTemplateGenerator


class ExportSystem:
    """Enhanced export system for multiple output formats."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.latex_generator = LaTeXGenerator()
        self.enhanced_generator = EnhancedTemplateGenerator()
    
    def export_to_pdf(self, timeline: ProjectTimeline, output_path: str, 
                     template_type: str = "gantt_timeline",
                     device_profile: str = None,
                     color_scheme: str = None) -> bool:
        """Export timeline to PDF with enhanced quality."""
        try:
            self.logger.info(f"Exporting to PDF: {output_path}")
            
            # Generate LaTeX content
            latex_content = self.enhanced_generator.generate_enhanced_document(
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
    
    def export_multiple_formats(self, timeline: ProjectTimeline, 
                              base_path: str,
                              formats: List[str] = None) -> Dict[str, bool]:
        """Export timeline to multiple formats."""
        if formats is None:
            formats = ['pdf', 'svg', 'html', 'png']
        
        results = {}
        base_path = Path(base_path)
        
        for format_type in formats:
            output_path = base_path.with_suffix(f'.{format_type}')
            
            if format_type == 'pdf':
                results[format_type] = self.export_to_pdf(timeline, str(output_path))
            elif format_type == 'svg':
                results[format_type] = self.export_to_svg(timeline, str(output_path))
            elif format_type == 'html':
                results[format_type] = self.export_to_html(timeline, str(output_path))
            elif format_type == 'png':
                results[format_type] = self.export_to_png(timeline, str(output_path))
            else:
                self.logger.warning(f"Unsupported format: {format_type}")
                results[format_type] = False
        
        return results
    
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
