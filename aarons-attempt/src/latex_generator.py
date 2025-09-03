#!/usr/bin/env python3
"""
LaTeX generation components for the Gantt chart generator.
Provides modular LaTeX generation for different document sections.
"""

from typing import List, Dict
from datetime import date, timedelta

from .models import Task, ProjectTimeline, MonthInfo
from .config import config


class LaTeXEscaper:
    """Handles LaTeX character escaping and text formatting."""
    
    @staticmethod
    def escape_latex(text: str) -> str:
        """Escape special LaTeX characters in text."""
        if not text:
            return ""
        
        # Handle backslashes first to avoid double-escaping
        text = text.replace('\\', r'\textbackslash{}')
        
        # Handle Unicode characters
        unicode_replacements = {
            '≥': r'$\geq$',
            '≤': r'$\leq$',
            '≠': r'$\neq$',
            '±': r'$\pm$',
            '×': r'$\times$',
            '÷': r'$\div$',
            '∞': r'$\infty$',
            'α': r'$\alpha$',
            'β': r'$\beta$',
            'γ': r'$\gamma$',
            'δ': r'$\delta$',
            'ε': r'$\varepsilon$',
            'θ': r'$\theta$',
            'λ': r'$\lambda$',
            'μ': r'$\mu$',
            'π': r'$\pi$',
            'σ': r'$\sigma$',
            'τ': r'$\tau$',
            'φ': r'$\phi$',
            'χ': r'$\chi$',
            'ψ': r'$\psi$',
            'ω': r'$\omega$',
        }
        
        for char, replacement in unicode_replacements.items():
            text = text.replace(char, replacement)
        
        # Handle other special characters
        replacements = {
            '&': r'\&',
            '%': r'\%',
            '$': r'\$',
            '#': r'\#',
            '^': r'\textasciicircum{}',
            '_': r'\_',
            '{': r'\{',
            '}': r'\}',
            '~': r'\textasciitilde{}',
        }
        
        for char, replacement in replacements.items():
            text = text.replace(char, replacement)
        
        # Fix any double-escaped characters that might cause issues
        text = text.replace(r'\textbackslash{}\#', r'\#')
        text = text.replace(r'\textbackslash{}\&', r'\&')
        
        return text
    
    @staticmethod
    def truncate_text(text: str, max_length: int) -> str:
        """Truncate text to maximum length with ellipsis."""
        if len(text) <= max_length:
            return text
        return text[:max_length-3] + "..."


class LaTeXDocumentGenerator:
    """Generates the main LaTeX document structure."""
    
    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()
    
    def generate_document_header(self) -> str:
        """Generate the LaTeX document header with packages and setup."""
        packages = '\n'.join(f"\\usepackage{{{pkg}}}" for pkg in config.latex.packages)
        
        return f"""\\documentclass[{config.calendar.page_orientation},{config.calendar.page_size}]{{{config.latex.document_class}}}
{packages}

% Page setup inspired by calendar.sty
\\pagestyle{{empty}}
\\setlength{{\\parskip}}{{0.5em}}

% Table formatting
\\setlength{{\\tabcolsep}}{{1pt}}
\\renewcommand{{\\arraystretch}}{{1.0}}

% Use Helvetica for sans-serif
\\renewcommand{{\\familydefault}}{{\\sfdefault}}

% Color definitions
{config.colors.to_latex_colors()}

\\begin{{document}}
"""
    
    def generate_document_footer(self) -> str:
        """Generate the LaTeX document footer."""
        return "\\end{document}\n"


class TitlePageGenerator:
    """Generates the title page for the LaTeX document."""
    
    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()
    
    def generate_title_page(self, timeline: ProjectTimeline) -> str:
        """Generate the title page for the timeline."""
        title = self.escaper.escape_latex(timeline.title)
        start_date_str = timeline.start_date.strftime('%B %d, %Y')
        end_date_str = timeline.end_date.strftime('%B %d, %Y')
        
        return f"""
% Title page inspired by calendar.sty
\\begin{{titlepage}}
\\centering
\\vspace*{{{config.calendar.title_spacing}}}

{{{config.calendar.title_font_size}\\textbf{{{title}}}}}

\\vspace{{{config.calendar.month_spacing}}}
{{{config.calendar.month_font_size} {config.latex.subtitle}}}

\\vspace{{{config.calendar.title_spacing}}}

\\begin{{minipage}}{{0.9\\textwidth}}
\\centering
\\textbf{{Timeline Period:}} {start_date_str} -- {end_date_str}\\\\
\\textbf{{Total Duration:}} {timeline.total_duration_days} days\\\\
\\textbf{{Total Tasks:}} {timeline.total_tasks} tasks\\\\
\\textbf{{Months Covered:}} {len(timeline.get_months_between())} months
\\end{{minipage}}

\\vfill

\\end{{titlepage}}
"""


class CalendarGenerator:
    """Generates calendar views for the LaTeX document."""
    
    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()
    
    def generate_calendar_grid(self, month_info: MonthInfo, tasks: List[Task]) -> str:
        """Generate the TikZ calendar grid for a month."""
        grid = f"""
% Calendar grid with better proportions
\\begin{{tikzpicture}}[scale={config.calendar.calendar_scale}]
    % Main calendar border
    \\draw[thick] (0,0) rectangle ({config.calendar.calendar_width},{config.calendar.calendar_height});
    
    % Day headers with better styling
    \\node[font=\\bfseries{config.calendar.day_font_size}] at (0.5,5.5) {{Sun}};
    \\node[font=\\bfseries{config.calendar.day_font_size}] at (1.5,5.5) {{Mon}};
    \\node[font=\\bfseries{config.calendar.day_font_size}] at (2.5,5.5) {{Tue}};
    \\node[font=\\bfseries{config.calendar.day_font_size}] at (3.5,5.5) {{Wed}};
    \\node[font=\\bfseries{config.calendar.day_font_size}] at (4.5,5.5) {{Thu}};
    \\node[font=\\bfseries{config.calendar.day_font_size}] at (5.5,5.5) {{Fri}};
    \\node[font=\\bfseries{config.calendar.day_font_size}] at (6.5,5.5) {{Sat}};
    
    % Vertical lines
    \\foreach \\x in {{1,2,3,4,5,6}} {{
        \\draw[thick] (\\x,0) -- (\\x,5);
    }}
    
    % Horizontal lines
    \\foreach \\y in {{1,2,3,4,5}} {{
        \\draw[thick] (0,\\y) -- ({config.calendar.calendar_width},\\y);
    }}
"""
        
        # Add day numbers and tasks
        current_day = 1
        for week in range(6):  # Maximum 6 weeks
            for day in range(7):  # 7 days per week
                if week == 0 and day < month_info.first_weekday:
                    continue  # Skip days before month starts
                if current_day > month_info.num_days:
                    break
                
                x_pos = day + 0.5
                y_pos = 4.5 - week
                
                # Add day number
                grid += f"    \\node[font=\\bfseries{config.calendar.day_font_size}, anchor=north west] at ({day+0.05},{y_pos+0.4}) {{{current_day}}};\n"
                
                # Find tasks for this day
                day_date = month_info.start_date + timedelta(days=current_day - 1)
                day_tasks = [t for t in tasks if t.overlaps_with_date(day_date)]
                
                # Add task content in the day cell
                if day_tasks:
                    task_text = self._generate_day_task_text(day_tasks)
                    grid += f"    \\node[font={config.calendar.task_font_size}, anchor=north west, text width=0.9cm] at ({day+0.05},{y_pos-0.1}) {{{task_text}}};\n"
                
                current_day += 1
            if current_day > month_info.num_days:
                break
        
        grid += "\\end{tikzpicture}\n"
        return grid
    
    def _generate_day_task_text(self, tasks: List[Task]) -> str:
        """Generate task text for a single day cell."""
        task_text = ""
        limited_tasks = tasks[:config.calendar.max_tasks_per_day]
        
        for i, task in enumerate(limited_tasks):
            task_name = self.escaper.escape_latex(task.name)
            task_name = self.escaper.truncate_text(task_name, config.calendar.max_task_name_length)
            
            if task.is_milestone:
                task_text += f"\\textcolor{{{task.category_color}}}{{\\textbf{{$\\diamond$ {task_name}}}}}"
            else:
                task_text += f"\\textcolor{{{task.category_color}}}{{\\textbf{{$\\bullet$ {task_name}}}}}"
            
            if i < len(limited_tasks) - 1:
                task_text += "\\\\"
        
        return task_text
    
    def generate_month_page(self, month_info: MonthInfo, tasks: List[Task]) -> str:
        """Generate a complete calendar page for a month."""
        if not tasks:
            return ""
        
        page = f"""
\\newpage
\\pagestyle{{empty}}

\\begin{{center}}
{{{config.calendar.month_font_size}\\textbf{{{month_info.name}}}}}
\\end{{center}}

\\vspace{{{config.calendar.month_spacing}}}
"""
        
        # Add calendar grid
        page += self.generate_calendar_grid(month_info, tasks)
        
        # Add detailed task list for this month
        page += f"\\vspace{{{config.calendar.task_spacing}}}\n"
        page += f"\\subsection{{Task Details for {month_info.name}}}\n"
        page += "\\begin{itemize}[leftmargin=1cm]\n"
        
        for task in tasks:
            task_name = self.escaper.escape_latex(task.name)
            description = self.escaper.escape_latex(task.notes)
            start_date_str = task.start_date.strftime(config.tasks.display_date_format)
            due_date_str = task.due_date.strftime(config.tasks.display_date_format)
            
            page += f"    \\item[\\textcolor{{{task.category_color}}}{{{task.marker}}}] \\textbf{{{task_name}}} ({start_date_str} - {due_date_str})\\\\ {description}\n"
        
        page += "\\end{itemize}\n"
        return page
    
    def generate_calendar_view(self, timeline: ProjectTimeline) -> str:
        """Generate the complete calendar view for all months."""
        if not timeline.tasks:
            return ""
        
        months = timeline.get_months_between()
        calendar_pages = ""
        
        for month_info in months:
            month_tasks = timeline.get_tasks_for_month(month_info)
            if month_tasks:
                calendar_pages += self.generate_month_page(month_info, month_tasks)
        
        return calendar_pages


class LegendGenerator:
    """Generates the legend section for the LaTeX document."""
    
    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()
    
    def generate_legend(self) -> str:
        """Generate the complete legend section."""
        return """
% Legend
\\section{Legend}

\\subsection{Task Categories}
\\begin{itemize}[leftmargin=2cm]
    \\item[\\textcolor{researchcore}{$\\bullet$}] PROPOSAL - Research proposals and dissertation work
    \\item[\\textcolor{researchexp}{$\\bullet$}] LASER - Laser alignment and experimental work
    \\item[\\textcolor{researchexp}{$\\bullet$}] EXPERIMENTAL - Imaging, surgery, and data collection
    \\item[\\textcolor{researchout}{$\\bullet$}] PUBLICATION - Manuscripts and presentations
    \\item[\\textcolor{administrative}{$\\bullet$}] ADMINISTRATIVE - Forms, applications, and reviews
    \\item[\\textcolor{accountability}{$\\bullet$}] ACCOUNTABILITY - Meetings and responsibilities
    \\item[\\textcolor{service}{$\\bullet$}] SERVICE - SPIE chapter and other service activities
\\end{itemize}

\\subsection{Task Information}
\\begin{itemize}[leftmargin=2cm]
    \\item[\\textbf{Dependencies}] Tasks that must be completed before this task
    \\item[\\textbf{Description}] Additional details and context for each task
    \\item[\\textbf{Milestones}] Tasks with same start and due date (marked with diamond shapes)
\\end{itemize}
"""


class LaTeXGenerator:
    """Main LaTeX generator that coordinates all components."""
    
    def __init__(self):
        self.escaper = LaTeXEscaper()
        self.document_generator = LaTeXDocumentGenerator(self.escaper)
        self.title_generator = TitlePageGenerator(self.escaper)
        self.calendar_generator = CalendarGenerator(self.escaper)
        self.legend_generator = LegendGenerator(self.escaper)
    
    def generate_complete_document(self, timeline: ProjectTimeline) -> str:
        """Generate the complete LaTeX document."""
        latex_content = self.document_generator.generate_document_header()
        latex_content += self.title_generator.generate_title_page(timeline)
        latex_content += self.calendar_generator.generate_calendar_view(timeline)
        latex_content += self.legend_generator.generate_legend()
        latex_content += self.document_generator.generate_document_footer()
        
        return latex_content
