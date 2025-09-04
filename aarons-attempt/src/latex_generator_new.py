#!/usr/bin/env python3
"""
Unified LaTeX generation system for the Gantt chart generator.
Combines all LaTeX generation components into a single, cohesive system.
"""

import logging
from typing import List, Dict, Any, Optional
from datetime import date, timedelta

from .models import Task, ProjectTimeline, MonthInfo
from .config import config, config_manager


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
            '≥': r'$\geq$', '≤': r'$\leq$', '≠': r'$\neq$', '±': r'$\pm$',
            '×': r'$\times$', '÷': r'$\div$', '∞': r'$\infty$',
            'α': r'$\alpha$', 'β': r'$\beta$', 'γ': r'$\gamma$', 'δ': r'$\delta$',
            'ε': r'$\varepsilon$', 'θ': r'$\theta$', 'λ': r'$\lambda$', 'μ': r'$\mu$',
            'π': r'$\pi$', 'σ': r'$\sigma$', 'τ': r'$\tau$', 'φ': r'$\phi$',
            'χ': r'$\chi$', 'ψ': r'$\psi$', 'ω': r'$\omega$',
        }

        for char, replacement in unicode_replacements.items():
            text = text.replace(char, replacement)

        # Handle other special characters
        replacements = {
            '&': r'\&', '%': r'\%', '$': r'\$', '#': r'\#',
            '^': r'\textasciicircum{}', '_': r'\_', '{': r'\{', '}': r'\}',
            '~': r'\textasciitilde{}',
        }

        for char, replacement in replacements.items():
            text = text.replace(char, replacement)

        return text


class BaseTemplateGenerator:
    """Base class for template generators."""
    
    def __init__(self, config_manager=None):
        """Initialize template generator."""
        self.logger = logging.getLogger(__name__)
        self.config_manager = config_manager or config_manager
        self.escaper = LaTeXEscaper()
    
    def generate_document(self, timeline: ProjectTimeline, 
                         template_id: str = None,
                         device_profile_id: str = None,
                         color_scheme_id: str = None) -> str:
        """Generate complete LaTeX document."""
        config = self.config_manager.get_combined_config(
            template_id, device_profile_id, color_scheme_id
        )
        
        return self._generate_document_content(timeline, config)
    
    def _generate_document_content(self, timeline: ProjectTimeline, 
                                 config: Dict[str, Any]) -> str:
        """Generate document content - to be implemented by subclasses."""
        raise NotImplementedError("Subclasses must implement _generate_document_content")


class LaTeXDocumentGenerator:
    """Generates LaTeX document structure and headers."""

    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()

    def generate_document_header(self) -> str:
        """Generate LaTeX document header with packages and styling."""
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
\\usepackage{{bookmark}}
\\usepackage{{enumitem}}
\\usepackage{{longtable}}
\\usepackage{{multirow}}
\\usepackage{{colortbl}}

% Enhanced TikZ libraries
\\usetikzlibrary{{{','.join(config.latex.get_tikz_libraries())}}}

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

{self._generate_tikz_styles()}

\\begin{{document}}
"""

    def _generate_tikz_styles(self) -> str:
        """Generate TikZ styles for enhanced graphics."""
        return """
% Enhanced TikZ styles
\\tikzset{
    task node/.style={
        rectangle, 
        rounded corners=4pt,
        draw=black!30,
        fill=white,
        drop shadow={shadow xshift=1pt, shadow yshift=-1pt, fill=black!20},
        minimum height=0.7cm,
        minimum width=2cm,
        font=\\small\\bfseries,
        align=center
    },
    milestone node/.style={
        diamond,
        draw=purple!60,
        fill=purple!20,
        drop shadow={shadow xshift=1pt, shadow yshift=-1pt, fill=purple!30},
        minimum size=1cm,
        font=\\small\\bfseries,
        align=center
    },
    timeline axis/.style={
        thick,
        line width=3pt,
        color=blue!70
    },
    task bar/.style={
        rounded corners=3pt,
        minimum height=0.6cm
    },
    progress bar/.style={
        rounded corners=2pt,
        minimum height=0.3cm
    }
}

% Color definitions
""" + config.colors.to_latex_colors()

    def generate_document_footer(self) -> str:
        """Generate LaTeX document footer."""
        return "\\end{document}"


class TitlePageGenerator:
    """Generates title pages and headers."""

    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()

    def generate_title_page(self, timeline: ProjectTimeline) -> str:
        """Generate title page for the document."""
        title = self.escaper.escape_latex(timeline.title)
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


class GanttChartGenerator:
    """Generates Gantt charts and timeline views."""

    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()

    def generate_timeline_view(self, timeline: ProjectTimeline) -> str:
        """Generate enhanced timeline view with modern TikZ features."""
        content = """
\\begin{center}
\\begin{tikzpicture}[scale=0.8]
    % Timeline axis
    \\draw[timeline axis] (0,0) -- (12,0);
    
    % Month markers
"""
        
        months = timeline.get_months_between()
        for i, month in enumerate(months):
            x_pos = (i / max(1, len(months) - 1)) * 12
            month_name = month.name.split()[0][:3]  # First 3 letters
            content += f"    \\node[font=\\small\\bfseries] at ({x_pos},-0.5) {{{month_name}}};\n"
        
        content += """
    % Task bars
"""
        
        y_pos = 1.5
        for i, task in enumerate(timeline.tasks):
            task_name = self.escaper.escape_latex(task.name)
            start_x = self._calculate_timeline_position(task.start_date, timeline.start_date)
            end_x = self._calculate_timeline_position(task.due_date, timeline.start_date)
            
            if task.is_milestone:
                content += f"    \\node[milestone node] at ({start_x},{y_pos}) {{{task_name}}};\n"
            else:
                color = self._get_category_color(task.category)
                content += f"    \\draw[task bar, fill={color}] ({start_x},{y_pos-0.3}) rectangle ({end_x},{y_pos+0.3});\n"
                content += f"    \\node[font=\\small\\bfseries, text=white] at ({(start_x+end_x)/2},{y_pos}) {{{task_name}}};\n"
            
            y_pos += 0.9
        
        content += """
\\end{tikzpicture}
\\end{center}
"""
        return content

    def _get_category_color(self, category: str) -> str:
        """Get color for task category."""
        from .config import get_category_color
        color_name = get_category_color(category)
        return color_name

    def _calculate_timeline_position(self, date: date, start_date: date) -> float:
        """Calculate position on timeline axis."""
        days_diff = (date - start_date).days
        return (days_diff / 365.0) * 12  # Scale to 12 units for year


class CalendarGenerator:
    """Generates calendar views and grids."""

    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()

    def generate_calendar_view(self, timeline: ProjectTimeline) -> str:
        """Generate calendar view for the timeline."""
        content = ""
        for month_info in timeline.get_months_between():
            content += self.generate_calendar_grid(month_info, timeline)
            content += "\\newpage\n"
        return content

    def generate_calendar_grid(self, month_info: MonthInfo, timeline: ProjectTimeline) -> str:
        """Generate calendar grid for a specific month."""
        tasks = timeline.get_tasks_for_month(month_info)
        
        content = f"""
\\section*{{{month_info.name}}}
\\begin{{center}}
\\begin{{tikzpicture}}[scale=0.9]
    % Calendar grid
    \\draw[step=1, gray, thin] (0,0) grid (7,6);
    
    % Day headers
    \\node[font=\\small\\bfseries] at (0.5,5.5) {{Sun}};
    \\node[font=\\small\\bfseries] at (1.5,5.5) {{Mon}};
    \\node[font=\\small\\bfseries] at (2.5,5.5) {{Tue}};
    \\node[font=\\small\\bfseries] at (3.5,5.5) {{Wed}};
    \\node[font=\\small\\bfseries] at (4.5,5.5) {{Thu}};
    \\node[font=\\small\\bfseries] at (5.5,5.5) {{Fri}};
    \\node[font=\\small\\bfseries] at (6.5,5.5) {{Sat}};
    
    % Day numbers and tasks
"""
        
        current_date = month_info.start_date
        week = 5
        day_of_week = month_info.first_weekday
        
        for day in range(1, month_info.num_days + 1):
            if day_of_week == 7:
                day_of_week = 0
                week -= 1
            
            day_tasks = timeline.get_tasks_for_date(current_date)
            task_text = self._generate_day_task_text(day_tasks)
            
            content += f"    \\node[font=\\small] at ({day_of_week + 0.5},{week + 0.3}) {{{day}}};\n"
            if task_text:
                content += f"    \\node[font=\\tiny] at ({day_of_week + 0.5},{week - 0.1}) {{{task_text}}};\n"
            
            current_date += timedelta(days=1)
            day_of_week += 1
        
        content += """
\\end{tikzpicture}
\\end{center}
"""
        return content

    def _generate_day_task_text(self, tasks: List[Task]) -> str:
        """Generate task text for a single calendar day cell."""
        if not tasks:
            return ""
        
        if len(tasks) == 1:
            return self.escaper.escape_latex(tasks[0].name[:15])
        else:
            return f"{len(tasks)} tasks"


class TaskListGenerator:
    """Generates task lists and detailed views."""

    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()

    def generate_comprehensive_task_list(self, timeline: ProjectTimeline) -> str:
        """Generate comprehensive task list with details."""
        content = """
\\section*{Task List}
\\begin{longtable}{|p{0.1\\textwidth}|p{0.4\\textwidth}|p{0.15\\textwidth}|p{0.15\\textwidth}|p{0.2\\textwidth}|}
\\hline
\\textbf{ID} & \\textbf{Task Name} & \\textbf{Start Date} & \\textbf{Due Date} & \\textbf{Category} \\\\
\\hline
\\endhead
"""
        
        for i, task in enumerate(timeline.tasks, 1):
            task_name = self.escaper.escape_latex(task.name)
            start_date = task.start_date.strftime('%m/%d/%Y')
            due_date = task.due_date.strftime('%m/%d/%Y')
            category = self.escaper.escape_latex(task.category)
            
            content += f"""
{i:03d} & {task_name} & {start_date} & {due_date} & {category} \\\\
\\hline
"""
        
        content += "\\end{longtable}\n"
        return content


class LegendGenerator:
    """Generates legends and documentation."""

    def __init__(self, escaper: LaTeXEscaper = None):
        self.escaper = escaper or LaTeXEscaper()

    def generate_legend(self) -> str:
        """Generate legend for the document."""
        return """
\\section*{Legend}

\\subsection*{Task Categories}
\\begin{itemize}[leftmargin=2cm]
    \\item[\\textcolor{researchcore}{$\\bullet$}] Research Core - Proposals, dissertation work
    \\item[\\textcolor{researchexp}{$\\bullet$}] Experimental - Laser work, imaging, data collection
    \\item[\\textcolor{researchout}{$\\bullet$}] Publications - Manuscripts and presentations
    \\item[\\textcolor{administrative}{$\\bullet$}] Administrative - Forms, applications, reviews
    \\item[\\textcolor{milestone}{$\\bullet$}] Milestones - Key project deliverables
\\end{itemize}

\\subsection*{Status Indicators}
\\begin{itemize}[leftmargin=2cm]
    \\item[\\textcolor{completed}{$\\bullet$}] Completed - Task finished
    \\item[\\textcolor{inprogress}{$\\bullet$}] In Progress - Currently being worked on
    \\item[\\textcolor{blocked}{$\\bullet$}] Blocked - Waiting for dependencies
    \\item[\\textcolor{planned}{$\\bullet$}] Planned - Scheduled for future
\\end{itemize}
"""


# Template Generators
class GanttTimelineGenerator(BaseTemplateGenerator):
    """Enhanced Gantt timeline generator with modern TikZ features."""
    
    def __init__(self, config_manager=None):
        super().__init__(config_manager)
        self.document_generator = LaTeXDocumentGenerator(self.escaper)
        self.title_generator = TitlePageGenerator(self.escaper)
        self.gantt_generator = GanttChartGenerator(self.escaper)
        self.calendar_generator = CalendarGenerator(self.escaper)
        self.task_list_generator = TaskListGenerator(self.escaper)
        self.legend_generator = LegendGenerator(self.escaper)
    
    def _generate_document_content(self, timeline: ProjectTimeline, 
                                 config: Dict[str, Any]) -> str:
        """Generate enhanced Gantt timeline document."""
        content = self.document_generator.generate_document_header()
        content += self.title_generator.generate_title_page(timeline)
        content += "\\newpage\n"
        content += self.task_list_generator.generate_comprehensive_task_list(timeline)
        content += "\\newpage\n"
        content += "\\section*{Enhanced Project Timeline}\n"
        content += self.gantt_generator.generate_timeline_view(timeline)
        content += "\\newpage\n"
        content += "\\section*{Monthly Calendar Views}\n"
        content += self.calendar_generator.generate_calendar_view(timeline)
        content += self.legend_generator.generate_legend()
        content += self.document_generator.generate_document_footer()
        return content


class MonthlyCalendarGenerator(BaseTemplateGenerator):
    """Monthly calendar generator."""
    
    def _generate_document_content(self, timeline: ProjectTimeline, 
                                 config: Dict[str, Any]) -> str:
        """Generate monthly calendar document."""
        content = self._generate_document_header(config['template'], config['device_profile'])
        content += self._generate_color_definitions()
        
        for month_info in timeline.get_months_between():
            content += self._generate_monthly_calendar(month_info, timeline, config['colors'])
            content += "\\newpage\n"
        
        content += self._generate_document_footer()
        return content
    
    def _generate_document_header(self, template, device_profile) -> str:
        """Generate document header for monthly calendar."""
        return f"""\\documentclass[landscape,a4paper]{{article}}
\\usepackage[utf8]{{inputenc}}
\\usepackage[T1]{{fontenc}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage[landscape,margin=0.3in]{{geometry}}
\\usepackage{{tikz}}
\\usepackage{{xcolor}}
\\usepackage{{array}}

% Enhanced TikZ libraries
\\usetikzlibrary{{{','.join(config.latex.get_tikz_libraries())}}}

\\begin{{document}}
"""
    
    def _generate_color_definitions(self) -> str:
        """Generate color definitions."""
        return config.colors.to_latex_colors()
    
    def _generate_monthly_calendar(self, month_info: MonthInfo, timeline: ProjectTimeline, 
                                 colors: Dict[str, Any]) -> str:
        """Generate monthly calendar."""
        tasks = timeline.get_tasks_for_month(month_info)
        
        return f"""
\\section*{{{month_info.name}}}
\\begin{{center}}
\\begin{{tikzpicture}}[scale=1.2]
    % Calendar grid
    \\draw[step=1, gray, thin] (0,0) grid (7,6);
    
    % Day headers
    \\node[font=\\large\\bfseries] at (0.5,5.5) {{Sun}};
    \\node[font=\\large\\bfseries] at (1.5,5.5) {{Mon}};
    \\node[font=\\large\\bfseries] at (2.5,5.5) {{Tue}};
    \\node[font=\\large\\bfseries] at (3.5,5.5) {{Wed}};
    \\node[font=\\large\\bfseries] at (4.5,5.5) {{Thu}};
    \\node[font=\\large\\bfseries] at (5.5,5.5) {{Fri}};
    \\node[font=\\large\\bfseries] at (6.5,5.5) {{Sat}};
    
    % Day numbers and tasks
    % [Calendar generation logic would go here]
\\end{{tikzpicture}}
\\end{{center}}
"""
    
    def _generate_document_footer(self) -> str:
        """Generate document footer."""
        return "\\end{document}"


class WeeklyPlannerGenerator(BaseTemplateGenerator):
    """Weekly planner generator."""
    
    def _generate_document_content(self, timeline: ProjectTimeline, 
                                 config: Dict[str, Any]) -> str:
        """Generate weekly planner document."""
        content = self._generate_document_header(config['template'], config['device_profile'])
        content += self._generate_color_definitions()
        
        weeks = self._get_weeks_between(timeline.start_date, timeline.end_date)
        for week_start in weeks:
            content += self._generate_weekly_planner(week_start, timeline, config['colors'])
            content += "\\newpage\n"
        
        content += self._generate_document_footer()
        return content
    
    def _get_weeks_between(self, start_date: date, end_date: date) -> List[date]:
        """Get list of week start dates between start and end dates."""
        weeks = []
        current = start_date
        while current <= end_date:
            weeks.append(current)
            current += timedelta(days=7)
        return weeks
    
    def _generate_weekly_planner(self, week_start: date, timeline: ProjectTimeline, 
                               colors: Dict[str, Any]) -> str:
        """Generate weekly planner."""
        return f"""
\\section*{{Week of {week_start.strftime('%B %d, %Y')}}}
\\begin{{center}}
\\begin{{tikzpicture}}[scale=1.0]
    % Weekly grid
    \\draw[step=1, gray, thin] (0,0) grid (7,12);
    
    % Day headers
    days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
    for i, day in enumerate(days):
        \\node[font=\\large\\bfseries] at ({{i + 0.5}},11.5) {{{day}}};
    
    % Time slots
    for hour in range(8, 20):
        \\node[font=\\small] at (-0.5,{{hour - 8 + 0.5}}) {{{hour}:00}};
\\end{{tikzpicture}}
\\end{{center}}
"""
    
    def _generate_document_header(self, template, device_profile) -> str:
        """Generate document header for weekly planner."""
        return f"""\\documentclass[landscape,a4paper]{{article}}
\\usepackage[utf8]{{inputenc}}
\\usepackage[T1]{{fontenc}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage[landscape,margin=0.4in]{{geometry}}
\\usepackage{{tikz}}
\\usepackage{{xcolor}}

% Enhanced TikZ libraries
\\usetikzlibrary{{{','.join(config.latex.get_tikz_libraries())}}}

\\begin{{document}}
"""
    
    def _generate_color_definitions(self) -> str:
        """Generate color definitions."""
        return config.colors.to_latex_colors()
    
    def _generate_document_footer(self) -> str:
        """Generate document footer."""
        return "\\end{document}"


class TemplateGeneratorFactory:
    """Factory for creating template generators."""
    
    _generators = {
        'gantt_timeline': GanttTimelineGenerator,
        'monthly_calendar': MonthlyCalendarGenerator,
        'weekly_planner': WeeklyPlannerGenerator,
    }
    
    @staticmethod
    def create_generator(template_type: str, config_manager=None) -> BaseTemplateGenerator:
        """Create a template generator for the specified type."""
        if template_type not in TemplateGeneratorFactory._generators:
            raise ValueError(f"Unknown template type: {template_type}")
        
        generator_class = TemplateGeneratorFactory._generators[template_type]
        return generator_class(config_manager)


class LaTeXGenerator:
    """Main LaTeX generator that coordinates all components."""

    def __init__(self):
        self.escaper = LaTeXEscaper()
        self.document_generator = LaTeXDocumentGenerator(self.escaper)
        self.title_generator = TitlePageGenerator(self.escaper)
        self.calendar_generator = CalendarGenerator(self.escaper)
        self.legend_generator = LegendGenerator(self.escaper)
        self.task_list_generator = TaskListGenerator(self.escaper)

    def generate_complete_document(self, timeline: ProjectTimeline) -> str:
        """Generate the complete LaTeX document."""
        latex_content = self.document_generator.generate_document_header()
        latex_content += self.title_generator.generate_title_page(timeline)
        latex_content += self.task_list_generator.generate_comprehensive_task_list(timeline)
        latex_content += self.calendar_generator.generate_calendar_view(timeline)
        latex_content += self.legend_generator.generate_legend()
        latex_content += self.document_generator.generate_document_footer()

        return latex_content
