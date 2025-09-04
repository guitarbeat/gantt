#!/usr/bin/env python3
"""
Core functionality for the LaTeX Gantt Chart Generator.
Consolidated data models, processing, and LaTeX generation.
"""

import csv
import logging
import argparse
from dataclasses import dataclass
from datetime import date, datetime
from typing import Optional, List, Dict, Tuple
from pathlib import Path

from .config import config
from .prisma_generator import PRISMAGenerator, PRISMAData, create_sample_prisma_data


@dataclass
class Task:
    """Represents a single task with all its properties."""
    
    id: str
    name: str
    start_date: date
    due_date: date
    category: str
    dependencies: str = ""
    notes: str = ""
    is_milestone: bool = False
    
    def __post_init__(self):
        """Validate and clean task data after initialization."""
        self._validate_dates()
        self._determine_milestone()
    
    
    def _validate_dates(self) -> None:
        """Validate that dates are reasonable."""
        if self.start_date > self.due_date:
            raise ValueError(f"Task '{self.name}': start date ({self.start_date}) cannot be after due date ({self.due_date})")
    
    def _determine_milestone(self) -> None:
        """Determine if this task is a milestone (same start and due date)."""
        self.is_milestone = (self.start_date == self.due_date)
    
    @property
    def duration_days(self) -> int:
        """Calculate task duration in days."""
        return (self.due_date - self.start_date).days + 1
    
    @property
    def category_color(self) -> str:
        """Get the color name for this task's category."""
        return get_category_color(self.category)


@dataclass
class ProjectTimeline:
    """Represents the complete project timeline with tasks and metadata."""
    
    tasks: List[Task]
    title: str
    start_date: date
    end_date: date
    
    def __post_init__(self):
        """Validate timeline data after initialization."""
        if not self.tasks:
            raise ValueError("Timeline must contain at least one task")
        
        # Recalculate date range from actual tasks
        self.start_date = min(task.start_date for task in self.tasks)
        self.end_date = max(task.due_date for task in self.tasks)
    
    @property
    def total_duration_days(self) -> int:
        """Calculate total project duration in days."""
        return (self.end_date - self.start_date).days
    
    @property
    def total_tasks(self) -> int:
        """Get total number of tasks."""
        return len(self.tasks)


class DataProcessor:
    """Main data processing coordinator."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def process_csv_to_timeline(self, csv_file_path: str, title: str = None) -> ProjectTimeline:
        """Complete pipeline from CSV file to ProjectTimeline."""
        self.logger.info(f"Processing CSV file: {csv_file_path}")
        
        # Read CSV data
        tasks = self._read_csv_data(csv_file_path)
        
        if not tasks:
            raise ValueError("No valid tasks found in CSV data")
        
        # Build timeline
        timeline_title = title or config.latex.default_title
        timeline = ProjectTimeline(
            tasks=tasks,
            title=timeline_title,
            start_date=date.today(),  # Will be recalculated
            end_date=date.today()     # Will be recalculated
        )
        
        self.logger.info(f"Processed {len(tasks)} tasks successfully")
        return timeline
    
    def _read_csv_data(self, file_path: str) -> List[Task]:
        """Read and process CSV data into Task objects."""
        tasks = []
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                reader = csv.DictReader(f)
                
                for row_num, row in enumerate(reader, start=2):
                    try:
                        task = self._parse_csv_row(row)
                        if task:
                            tasks.append(task)
                    except Exception as e:
                        self.logger.warning(f"Row {row_num}: {e}")
                        
        except Exception as e:
            self.logger.error(f"Error reading CSV file: {e}")
            raise
        
        return tasks
    
    def _parse_csv_row(self, row: Dict[str, str]) -> Optional[Task]:
        """Parse a single CSV row into a Task object."""
        # Extract required fields
        task_id = row.get('Task ID', '').strip()
        task_name = row.get('Task Name', '').strip()
        start_date_str = row.get('Start Date', '').strip()
        due_date_str = row.get('Due Date', '').strip()
        category = row.get('Category', '').strip()
        
        if not all([task_id, task_name, start_date_str, due_date_str]):
            return None
        
        # Parse dates
        if not (start_date := self._parse_date(start_date_str)) or not (due_date := self._parse_date(due_date_str)):
            return None
        
        # Extract optional fields
        dependencies = row.get('Dependencies', '').strip()
        notes = row.get('Description', '').strip()
        
        return Task(
            id=task_id,
            name=task_name,
            start_date=start_date,
            due_date=due_date,
            category=category,
            dependencies=dependencies,
            notes=notes
        )
    
    def _parse_date(self, date_str: str) -> Optional[date]:
        """Parse date string to date object."""
        if not date_str:
            return None
        
        # Try common date formats
        formats = ['%Y-%m-%d', '%m/%d/%Y', '%d/%m/%Y', '%B %d, %Y']
        
        for fmt in formats:
            try:
                return datetime.strptime(date_str, fmt).date()
            except ValueError:
                continue
        
        return None
    
    def validate_input_file(self, file_path: str) -> Tuple[bool, List[str]]:
        """Validate that the input file exists and is readable."""
        errors = []
        
        if not Path(file_path).exists():
            errors.append(f"File does not exist: {file_path}")
            return False, errors
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                # Just check if we can read the first few lines
                for i, line in enumerate(f):
                    if i >= 5:
                        break
        except Exception as e:
            errors.append(f"Cannot read file: {e}")
            return False, errors
        
        return True, errors


class LaTeXGenerator:
    """Generates LaTeX documents from project timelines."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def generate_document(self, timeline: ProjectTimeline, include_prisma: bool = False) -> str:
        """Generate complete LaTeX document."""
        content = self._generate_header()
        content += self._generate_title_page(timeline)
        content += self._generate_timeline_view(timeline)
        content += self._generate_task_list(timeline)
        
        if include_prisma:
            content += self._generate_prisma_section()
        
        content += self._generate_footer()
        return content
    
    def _generate_header(self) -> str:
        """Generate LaTeX document header."""
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
\\usepackage{{enumitem}}
\\usepackage{{longtable}}

% Essential TikZ libraries
\\usetikzlibrary{{arrows.meta,shapes.geometric,positioning,calc}}

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

{self._generate_color_definitions()}

% PRISMA color scheme
\\definecolor{{prismablue}}{{RGB}}{{59, 130, 246}}
\\definecolor{{prismagray}}{{RGB}}{{107, 114, 128}}
\\definecolor{{prismagreen}}{{RGB}}{{34, 197, 94}}
\\definecolor{{prismared}}{{RGB}}{{239, 68, 68}}

\\begin{{document}}
"""
    
    def _generate_color_definitions(self) -> str:
        """Generate LaTeX color definitions."""
        colors = []
        for attr_name in dir(config.colors):
            if not attr_name.startswith('_') and not callable(getattr(config.colors, attr_name)):
                rgb = getattr(config.colors, attr_name)
                colors.append(f"\\definecolor{{{attr_name}}}{{RGB}}{{{rgb[0]}, {rgb[1]}, {rgb[2]}}}")
        return '\n'.join(colors)
    
    def _generate_title_page(self, timeline: ProjectTimeline) -> str:
        """Generate title page."""
        title = self._escape_latex(timeline.title)
        start_date = timeline.start_date.strftime('%B %d, %Y')
        end_date = timeline.end_date.strftime('%B %d, %Y')
        
        return f"""
\\begin{{center}}
\\vspace*{{2cm}}
{{\\Huge\\bfseries {title}}}\\\\
\\vspace{{1cm}}
{{\\Large PhD Research Timeline}}\\\\
\\vspace{{0.5cm}}
{{\\large {start_date} -- {end_date}}}\\\\
\\vspace{{1cm}}
{{\\large Total Tasks: {timeline.total_tasks} | Duration: {timeline.total_duration_days} days}}
\\end{{center}}

\\vspace{{2cm}}
"""
    
    def _generate_timeline_view(self, timeline: ProjectTimeline) -> str:
        """Generate timeline view."""
        content = """
\\section*{Project Timeline}
\\begin{center}
\\begin{tikzpicture}[scale=0.8]
    % Timeline axis
    \\draw[thick, line width=3pt, color=blue!70] (0,0) -- (12,0);
    
    % Task bars
"""
        
        y_pos = 1.5
        for i, task in enumerate(timeline.tasks):
            task_name = self._escape_latex(task.name)
            start_x = self._calculate_timeline_position(task.start_date, timeline.start_date)
            end_x = self._calculate_timeline_position(task.due_date, timeline.start_date)
            
            if task.is_milestone:
                content += f"    \\node[diamond, draw=purple!60, fill=purple!20, minimum size=1cm, font=\\small\\bfseries] at ({start_x},{y_pos}) {{{task_name}}};\n"
            else:
                color = task.category_color
                content += f"    \\draw[fill={color}, rounded corners=3pt] ({start_x},{y_pos-0.3}) rectangle ({end_x},{y_pos+0.3});\n"
                content += f"    \\node[font=\\small\\bfseries, text=white] at ({(start_x+end_x)/2},{y_pos}) {{{task_name}}};\n"
            
            y_pos += 0.9
        
        content += """
\\end{tikzpicture}
\\end{center}
"""
        return content
    
    def _generate_task_list(self, timeline: ProjectTimeline) -> str:
        """Generate task list."""
        content = """
\\newpage
\\section*{Task List}
\\begin{longtable}{|p{0.1\\textwidth}|p{0.4\\textwidth}|p{0.15\\textwidth}|p{0.15\\textwidth}|p{0.2\\textwidth}|}
\\hline
\\textbf{ID} & \\textbf{Task Name} & \\textbf{Start Date} & \\textbf{Due Date} & \\textbf{Category} \\\\
\\hline
\\endhead
"""
        
        for i, task in enumerate(timeline.tasks, 1):
            task_name = self._escape_latex(task.name)
            start_date = task.start_date.strftime('%m/%d/%Y')
            due_date = task.due_date.strftime('%m/%d/%Y')
            category = self._escape_latex(task.category)
            
            # Ensure proper table row formatting
            content += f"{i:03d} & {task_name} & {start_date} & {due_date} & {category} \\\\\n"
            content += "\\hline\n"
        
        content += "\\end{longtable}\n"
        return content
    
    def _generate_prisma_section(self) -> str:
        """Generate PRISMA flow diagram section."""
        prisma_generator = PRISMAGenerator()
        sample_data = create_sample_prisma_data()
        
        # Generate the PRISMA diagram content
        prisma_content = prisma_generator._generate_flow_diagram(sample_data)
        
        return f"""
\\newpage
\\section*{{PRISMA Flow Diagram}}

{prisma_content}

\\vspace{{1cm}}
\\begin{{center}}
\\textit{{PRISMA 2020 Flow Diagram for systematic reviews}}
\\end{{center}}
"""
    
    def _generate_footer(self) -> str:
        """Generate document footer."""
        return "\\end{document}"
    
    def _escape_latex(self, text: str) -> str:
        """Escape special LaTeX characters."""
        if not text:
            return ""
        
        # Handle LaTeX special characters that need escaping
        replacements = {
            '&': r'\&',           # Table column separator
            '%': r'\%',           # Comment character
            '$': r'\$',           # Math mode
            '#': r'\#',           # Macro parameter
            '^': r'\textasciicircum{}',  # Superscript
            '_': r'\_',           # Subscript
            '{': r'\{',           # Group start
            '}': r'\}',           # Group end
            '~': r'\textasciitilde{}',   # Non-breaking space
            '\\': r'\textbackslash{}'    # Backslash
        }
        
        # Apply replacements
        for char, replacement in replacements.items():
            text = text.replace(char, replacement)
        
        return text
    
    def _calculate_timeline_position(self, date: date, start_date: date) -> float:
        """Calculate position on timeline axis."""
        days_diff = (date - start_date).days
        return (days_diff / 365.0) * 12  # Scale to 12 units for year


class Application:
    """Main application class."""
    
    def __init__(self):
        self.setup_logging()
        self.logger = logging.getLogger(__name__)
        self.data_processor = DataProcessor()
        self.latex_generator = LaTeXGenerator()
    
    def setup_logging(self, level: str = "INFO") -> None:
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
        
        return not errors, errors
    
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


def get_category_color(category: str) -> str:
    """Get the color name for a given task category."""
    category_upper = category.upper()
    
    return next(
        (color_name for color_name, keywords in config.tasks.category_keywords.items()
         if any(keyword in category_upper for keyword in keywords)),
        "other"
    )


if __name__ == "__main__":
    import sys
    sys.exit(main())
