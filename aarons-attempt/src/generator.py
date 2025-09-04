#!/usr/bin/env python3
"""
LaTeX generation module for the LaTeX Gantt Chart Generator.
Handles LaTeX document generation and formatting.
"""

import logging
from datetime import date

from .config import config
from .models import ProjectTimeline
from .prisma_generator import PRISMAGenerator, create_sample_prisma_data


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
\\usepackage{{colortbl}}
\\usepackage{{booktabs}}

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
        """Generate beautiful timeline list view."""
        content = """
\\section*{Project Timeline}
\\vspace{0.5cm}

\\begin{enumerate}[leftmargin=0pt, itemindent=0pt, labelsep=0pt, labelwidth=0pt]
"""
        
        for i, task in enumerate(timeline.tasks, 1):
            task_name = self._escape_latex(task.name)
            start_date = task.start_date.strftime('%b %d, %Y')
            due_date = task.due_date.strftime('%b %d, %Y')
            category = self._escape_latex(task.category)
            duration = task.duration_days
            
            # Create timeline item
            if task.is_milestone:
                milestone_indicator = "\\textcolor{red}{\\textbf{★}}"
                content += f"\\item[{milestone_indicator}] \\textbf{{{task_name}}}\n"
                content += f"    \\begin{{itemize}}\n"
                content += f"        \\item \\textit{{Milestone: {start_date}}}\n"
                content += f"        \\item \\textcolor{{gray}}{{Category: {category}}}\n"
                content += f"    \\end{{itemize}}\n"
            else:
                content += f"\\item[{i:02d}] \\textbf{{{task_name}}}\n"
                content += f"    \\begin{{itemize}}\n"
                content += f"        \\item \\textit{{{start_date} -- {due_date} ({duration} days)}}\n"
                content += f"        \\item \\textcolor{{gray}}{{Category: {category}}}\n"
                content += f"    \\end{{itemize}}\n"
            
            content += "\\vspace{0.3cm}\n"
        
        content += """
\\end{enumerate}

\\vspace{1cm}
\\begin{center}
\\textit{Timeline spans from """ + timeline.start_date.strftime('%B %d, %Y') + """ to """ + timeline.end_date.strftime('%B %d, %Y') + """}
\\end{center}
"""
        return content
    
    def _generate_task_list(self, timeline: ProjectTimeline) -> str:
        """Generate beautiful task list."""
        content = """
\\newpage
\\section*{Detailed Task List}
\\vspace{0.5cm}

% Enhanced table styling
\\renewcommand{\\arraystretch}{1.3}
\\begin{longtable}{|c|p{0.5\\textwidth}|c|c|c|c|}
\\hline
\\rowcolor{blue!10}
\\textbf{ID} & \\textbf{Task Name} & \\textbf{Start Date} & \\textbf{Due Date} & \\textbf{Category} & \\textbf{Days} \\\\
\\hline
\\endhead

\\hline
\\multicolumn{6}{|c|}{\\textit{Continued on next page...}} \\\\
\\hline
\\endfoot

\\hline
\\endlastfoot
"""
        
        for i, task in enumerate(timeline.tasks, 1):
            task_name = self._escape_latex(task.name)
            start_date = task.start_date.strftime('%b %d, %Y')
            due_date = task.due_date.strftime('%b %d, %Y')
            category = self._escape_latex(task.category)
            duration = task.duration_days
            
            # Add milestone indicator
            milestone_indicator = "\\textcolor{red}{\\textbf{★}}" if task.is_milestone else ""
            task_display = f"{milestone_indicator} {task_name}" if milestone_indicator else task_name
            
            # Color coding for categories
            category_color = self._get_category_color_latex(task.category)
            
            # Ensure proper table row formatting
            content += f"{i:03d} & {task_display} & {start_date} & {due_date} & \\textcolor{{{category_color}}}{{{category}}} & {duration} \\\\\n"
            content += "\\hline\n"
        
        content += "\\end{longtable}\n"
        return content
    
    def _get_category_color_latex(self, category: str) -> str:
        """Get LaTeX color name for category."""
        category_upper = category.upper()
        
        # Map categories to colors
        color_map = {
            'RESEARCH': 'blue',
            'WRITING': 'green', 
            'ANALYSIS': 'orange',
            'MEETING': 'purple',
            'REVIEW': 'red',
            'PLANNING': 'teal',
            'DATA': 'brown',
            'EXPERIMENT': 'pink',
            'OTHER': 'gray'
        }
        
        for keyword, color in color_map.items():
            if keyword in category_upper:
                return color
        
        return 'gray'
    
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
