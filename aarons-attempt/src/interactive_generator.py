#!/usr/bin/env python3
"""
Interactive features generator for enhanced LaTeX documents.
Provides clickable elements, hyperlinks, and better navigation.
"""

import logging
from typing import List, Dict, Any, Optional
from datetime import date, timedelta
from .models import Task, ProjectTimeline, MonthInfo
from .config_manager import ConfigManager, config_manager
from .latex_generator import LaTeXEscaper
from .utils import LaTeXUtilities


class InteractiveElementGenerator:
    """Generates interactive elements for LaTeX documents."""
    
    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()
        self.logger = logging.getLogger(__name__)
    
    def generate_table_of_contents(self, timeline: ProjectTimeline) -> str:
        """Generate an interactive table of contents."""
        toc = """
% Interactive Table of Contents
\\section*{Table of Contents}
\\addcontentsline{toc}{section}{Table of Contents}

\\begin{enumerate}[leftmargin=1cm, itemsep=0.5em]
    \\item \\hyperref[sec:overview]{\\textbf{Project Overview}}
    \\item \\hyperref[sec:timeline]{\\textbf{Project Timeline}}
    \\item \\hyperref[sec:gantt]{\\textbf{Gantt Chart}}
    \\item \\hyperref[sec:calendar]{\\textbf{Monthly Calendar}}
    \\item \\hyperref[sec:tasks]{\\textbf{Task Details}}
    \\item \\hyperref[sec:legend]{\\textbf{Legend}}
\\end{enumerate}

\\vspace{1cm}
"""
        return toc
    
    def generate_navigation_bar(self, current_section: str = None) -> str:
        """Generate a navigation bar for the document."""
        nav = """
% Navigation Bar
\\begin{center}
\\begin{tikzpicture}
    \\node[rectangle, draw=primaryblue!30, fill=primaryblue!10, rounded corners=3pt, minimum width=0.8\\textwidth, minimum height=0.6cm] (nav) {};
    \\node[font=\\small\\bfseries, text=primaryblue] at (nav.center) {
        \\hyperref[sec:overview]{Overview} $\\bullet$ 
        \\hyperref[sec:timeline]{Timeline} $\\bullet$ 
        \\hyperref[sec:gantt]{Gantt} $\\bullet$ 
        \\hyperref[sec:calendar]{Calendar} $\\bullet$ 
        \\hyperref[sec:tasks]{Tasks} $\\bullet$ 
        \\hyperref[sec:legend]{Legend}
    };
\\end{tikzpicture}
\\end{center}

\\vspace{0.5cm}
"""
        return nav
    
    def generate_task_links(self, tasks: List[Task]) -> str:
        """Generate clickable task links."""
        links = """
% Task Quick Links
\\subsection*{Task Quick Links}
\\begin{itemize}[leftmargin=1cm, itemsep=0.3em]
"""
        
        for i, task in enumerate(tasks, 1):
            task_name = self.escaper.escape_latex(task.name)
            task_id = f"task_{i:03d}"
            links += f"    \\item \\hyperref[{task_id}]{{{task_name}}} - {task.start_date.strftime('%m/%d')} to {task.due_date.strftime('%m/%d')}\n"
        
        links += "\\end{itemize}\n"
        return links
    
    def generate_progress_summary(self, timeline: ProjectTimeline) -> str:
        """Generate an interactive progress summary."""
        total_tasks = len(timeline.tasks)
        completed_tasks = sum(1 for task in timeline.tasks if self._is_task_completed(task))
        in_progress_tasks = sum(1 for task in timeline.tasks if self._is_task_in_progress(task))
        blocked_tasks = sum(1 for task in timeline.tasks if self._is_task_blocked(task))
        
        progress_percentage = (completed_tasks / total_tasks * 100) if total_tasks > 0 else 0
        
        summary = f"""
% Interactive Progress Summary
\\section*{{Project Progress Summary}}
\\label{{sec:progress}}

\\begin{{tikzpicture}}[scale=0.8]
    % Progress bar background
    \\draw[progress bar background] (0,0) rectangle (10,1);
    
    % Progress bar fill
    \\draw[progress bar fill] (0,0) rectangle ({progress_percentage/10},1);
    
    % Progress percentage text
    \\node[font=\\large\\bfseries, text=white] at ({progress_percentage/20},0.5) {{{progress_percentage:.1f}\\%}};
    
    % Task status indicators
    \\node[legend item, fill=secondarygreen!30, draw=secondarygreen!60] at (0.5,2) {{}};
    \\node[legend text] at (0.8,2) {{Completed: {completed_tasks}}};
    
    \\node[legend item, fill=accentorange!30, draw=accentorange!60] at (0.5,1.7) {{}};
    \\node[legend text] at (0.8,1.7) {{In Progress: {in_progress_tasks}}};
    
    \\node[legend item, fill=warningred!30, draw=warningred!60] at (0.5,1.4) {{}};
    \\node[legend text] at (0.8,1.4) {{Blocked: {blocked_tasks}}};
    
    \\node[legend item, fill=neutralgray!30, draw=neutralgray!60] at (0.5,1.1) {{}};
    \\node[legend text] at (0.8,1.1) {{Total: {total_tasks}}};
\\end{{tikzpicture}}

\\vspace{{1cm}}
"""
        return summary
    
    def _is_task_completed(self, task: Task) -> bool:
        """Check if task is completed (placeholder implementation)."""
        # This would be enhanced with actual status checking
        return task.due_date < date.today()
    
    def _is_task_in_progress(self, task: Task) -> bool:
        """Check if task is in progress (placeholder implementation)."""
        today = date.today()
        return task.start_date <= today <= task.due_date
    
    def _is_task_blocked(self, task: Task) -> bool:
        """Check if task is blocked (placeholder implementation)."""
        # This would be enhanced with actual status checking
        return False


class EnhancedTemplateGenerator:
    """Enhanced template generator with interactive features."""
    
    def __init__(self, config_manager: ConfigManager = None):
        self.config_manager = config_manager or config_manager
        self.escaper = LaTeXEscaper()
        self.interactive_generator = InteractiveElementGenerator(self.escaper)
        self.logger = logging.getLogger(__name__)
    
    def generate_enhanced_document(self, timeline: ProjectTimeline, 
                                 template_id: str = None,
                                 device_profile_id: str = None,
                                 color_scheme_id: str = None) -> str:
        """Generate enhanced document with interactive features."""
        config = self.config_manager.get_combined_config(
            template_id, device_profile_id, color_scheme_id
        )
        
        # Document header with enhanced packages
        content = self._generate_enhanced_header(config)
        
        # Interactive table of contents
        content += self.interactive_generator.generate_table_of_contents(timeline)
        content += "\\newpage\n"
        
        # Navigation bar
        content += self.interactive_generator.generate_navigation_bar()
        
        # Project overview section
        content += self._generate_project_overview(timeline)
        content += "\\newpage\n"
        
        # Progress summary
        content += self.interactive_generator.generate_progress_summary(timeline)
        content += "\\newpage\n"
        
        # Enhanced timeline view
        content += self._generate_enhanced_timeline(timeline)
        content += "\\newpage\n"
        
        # Task quick links
        content += self.interactive_generator.generate_task_links(timeline.tasks)
        content += "\\newpage\n"
        
        # Enhanced task details
        content += self._generate_enhanced_task_details(timeline)
        content += "\\newpage\n"
        
        # Document footer
        content += self._generate_enhanced_footer()
        
        return content
    
    def _generate_enhanced_header(self, config: Dict[str, Any]) -> str:
        """Generate enhanced document header with interactive packages."""
        return """\\documentclass[portrait,a4paper]{article}
\\usepackage[utf8]{inputenc}
\\usepackage[T1]{fontenc}
\\usepackage{lmodern}
\\usepackage{helvet}
\\usepackage[portrait,margin=0.5in]{geometry}
\\usepackage{tikz}
\\usepackage{xcolor}
\\usepackage{array}
\\usepackage{fancyhdr}
\\usepackage{hyperref}
\\usepackage{bookmark}
\\usepackage{enumitem}
\\usepackage{longtable}
\\usepackage{multirow}
\\usepackage{colortbl}

% Enhanced TikZ libraries for interactive features
\\usetikzlibrary{arrows.meta,shapes.geometric,positioning,calc,decorations.pathmorphing,patterns,shadows,fit,backgrounds,matrix,chains,scopes,pgfgantt}

% Page setup
\\pagestyle{fancy}
\\fancyhf{}
\\fancyhead[L]{Project Timeline}
\\fancyhead[R]{\\today}
\\fancyfoot[C]{\\thepage}
\\renewcommand{\\headrulewidth}{0.4pt}

% Hyperlink setup
\\hypersetup{
    colorlinks=true,
    linkcolor=primaryblue,
    urlcolor=primaryblue,
    citecolor=primaryblue,
    bookmarksopen=true,
    bookmarksnumbered=true
}

% Use Helvetica for sans-serif
\\renewcommand{\\familydefault}{\\sfdefault}

% Enhanced TikZ styles (from latex_generator.py)
\\input{enhanced_styles.tex}

\\begin{document}
"""
    
    def _generate_project_overview(self, timeline: ProjectTimeline) -> str:
        """Generate project overview section."""
        title = self.escaper.escape_latex(timeline.title)
        start_date_str = timeline.start_date.strftime('%B %d, %Y')
        end_date_str = timeline.end_date.strftime('%B %d, %Y')
        
        return f"""
\\section{{Project Overview}}
\\label{{sec:overview}}

\\begin{{center}}
\\begin{{tikzpicture}}
    \\node[rectangle, draw=primaryblue!40, fill=primaryblue!5, rounded corners=5pt, minimum width=0.9\\textwidth, minimum height=3cm] (overview) {{}};
    \\node[font=\\Large\\bfseries, text=primaryblue] at (overview.north) [yshift=-0.5cm] {{{title}}};
    \\node[font=\\large, text=neutralgray] at (overview.center) {{
        \\textbf{{Timeline Period:}} {start_date_str} -- {end_date_str}\\\\
        \\textbf{{Total Duration:}} {timeline.total_duration_days} days\\\\
        \\textbf{{Total Tasks:}} {timeline.total_tasks} tasks\\\\
        \\textbf{{Months Covered:}} {len(timeline.get_months_between())} months
    }};
\\end{{tikzpicture}}
\\end{{center}}

\\vspace{{1cm}}
"""
    
    def _generate_enhanced_timeline(self, timeline: ProjectTimeline) -> str:
        """Generate enhanced timeline view."""
        from .latex_generator import GanttChartGenerator
        gantt_generator = GanttChartGenerator(self.escaper)
        
        return f"""
\\section{{Enhanced Project Timeline}}
\\label{{sec:timeline}}

{gantt_generator.generate_timeline_view(timeline)}

\\vspace{{1cm}}
"""
    
    def _generate_enhanced_task_details(self, timeline: ProjectTimeline) -> str:
        """Generate enhanced task details with interactive elements."""
        content = """
\\section{Enhanced Task Details}
\\label{sec:tasks}

\\begin{longtable}{|p{0.1\\textwidth}|p{0.4\\textwidth}|p{0.15\\textwidth}|p{0.15\\textwidth}|p{0.2\\textwidth}|}
\\hline
\\textbf{ID} & \\textbf{Task Name} & \\textbf{Start Date} & \\textbf{Due Date} & \\textbf{Status} \\\\
\\hline
\\endhead
"""
        
        for i, task in enumerate(timeline.tasks, 1):
            task_name = self.escaper.escape_latex(task.name)
            task_id = f"task_{i:03d}"
            start_date = task.start_date.strftime('%m/%d/%Y')
            due_date = task.due_date.strftime('%m/%d/%Y')
            status = self._get_task_status(task)
            status_color = self._get_status_color(status)
            
            content += f"""
\\hyperref[{task_id}]{{{i:03d}}} & {task_name} & {start_date} & {due_date} & \\textcolor{{{status_color}}}{{\\textbf{{{status}}}}} \\\\
\\hline
"""
        
        content += "\\end{longtable}\n"
        return content
    
    def _get_task_status(self, task: Task) -> str:
        """Get task status."""
        if self._is_task_completed(task):
            return "Completed"
        elif self._is_task_in_progress(task):
            return "In Progress"
        elif self._is_task_blocked(task):
            return "Blocked"
        else:
            return "Planned"
    
    def _get_status_color(self, status: str) -> str:
        """Get color for task status."""
        color_map = {
            "Completed": "secondarygreen",
            "In Progress": "accentorange", 
            "Blocked": "warningred",
            "Planned": "neutralgray"
        }
        return color_map.get(status, "neutralgray")
    
    def _is_task_completed(self, task: Task) -> bool:
        """Check if task is completed."""
        return task.due_date < date.today()
    
    def _is_task_in_progress(self, task: Task) -> bool:
        """Check if task is in progress."""
        today = date.today()
        return task.start_date <= today <= task.due_date
    
    def _is_task_blocked(self, task: Task) -> bool:
        """Check if task is blocked."""
        return False
    
    def _generate_enhanced_footer(self) -> str:
        """Generate enhanced document footer."""
        return """
\\section{Legend}
\\label{sec:legend}

\\subsection{Task Categories}
\\begin{itemize}[leftmargin=2cm]
    \\item[\\textcolor{primaryblue}{$\\bullet$}] Research Core - Proposals, dissertation work
    \\item[\\textcolor{secondarygreen}{$\\bullet$}] Experimental - Laser work, imaging, data collection
    \\item[\\textcolor{accentorange}{$\\bullet$}] Publications - Manuscripts and presentations
    \\item[\\textcolor{neutralgray}{$\\bullet$}] Administrative - Forms, applications, reviews
    \\item[\\textcolor{infopurple}{$\\bullet$}] Milestones - Key project deliverables
\\end{itemize}

\\subsection{Status Indicators}
\\begin{itemize}[leftmargin=2cm]
    \\item[\\textcolor{secondarygreen}{$\\bullet$}] Completed - Task finished
    \\item[\\textcolor{accentorange}{$\\bullet$}] In Progress - Currently being worked on
    \\item[\\textcolor{warningred}{$\\bullet$}] Blocked - Waiting for dependencies
    \\item[\\textcolor{neutralgray}{$\\bullet$}] Planned - Scheduled for future
\\end{itemize}

\\vfill

\\begin{center}
\\textcolor{neutralgray}{\\small Generated on \\today}
\\end{center}

\\end{document}
"""
