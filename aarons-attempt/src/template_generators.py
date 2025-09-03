#!/usr/bin/env python3
"""
Template generators for different planner types inspired by latex-yearly-planner.
Provides modular template generation for various layouts and views.
"""

import logging
from typing import List, Dict, Any, Optional
from datetime import date, timedelta
from .models import Task, ProjectTimeline, MonthInfo
from .config_manager import ConfigManager, config_manager
from .latex_generator import LaTeXEscaper
from .utils import LaTeXUtilities


class BaseTemplateGenerator:
    """Base class for template generators."""
    
    def __init__(self, config_manager: ConfigManager = None):
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


class GanttTimelineGenerator(BaseTemplateGenerator):
    """Generator for Gantt timeline templates."""
    
    def _generate_document_content(self, timeline: ProjectTimeline, 
                                 config: Dict[str, Any]) -> str:
        """Generate Gantt timeline document."""
        template = config['template']
        device_profile = config['device_profile']
        colors = config['colors']
        
        # Document header
        content = self._generate_document_header(template, device_profile)
        
        # Title page
        content += self._generate_title_page(timeline, colors)
        
        # Timeline content
        content += self._generate_timeline_content(timeline, config)
        
        # Task table
        content += self._generate_task_table(timeline, colors)
        
        # Document footer
        content += self._generate_document_footer()
        
        return content
    
    def _generate_document_header(self, template: Any, device_profile: Any) -> str:
        """Generate LaTeX document header."""
        return LaTeXUtilities.generate_document_header(template, device_profile)
    

    
    def _generate_title_page(self, timeline: ProjectTimeline, colors: Dict[str, Any]) -> str:
        """Generate title page."""
        title = self.escaper.escape_latex(timeline.title)
        start_date = timeline.start_date.strftime("%B %d, %Y")
        end_date = timeline.end_date.strftime("%B %d, %Y")
        total_tasks = timeline.total_tasks
        
        return f"""
% Title page
\\begin{{titlepage}}
\\centering
\\vspace*{{2cm}}

{{\\Huge\\textbf{{{title}}}}}

\\vspace{{1cm}}

{{\\Large Project Timeline}}

\\vspace{{2cm}}

\\begin{{tabular}}{{ll}}
\\textbf{{Start Date:}} & {start_date} \\\\
\\textbf{{End Date:}} & {end_date} \\\\
\\textbf{{Total Tasks:}} & {total_tasks} \\\\
\\textbf{{Duration:}} & {timeline.total_duration_days} days \\\\
\\end{{tabular}}

\\vspace{{2cm}}

{{\\large Generated on \\today}}

\\end{{titlepage}}

\\newpage
"""
    
    def _generate_timeline_content(self, timeline: ProjectTimeline, config: Dict[str, Any]) -> str:
        """Generate timeline content."""
        # This would contain the main timeline visualization
        # For now, return a placeholder
        return """
% Timeline content
\\section*{Project Timeline}

\\begin{center}
\\begin{tikzpicture}[scale=0.8]
% Timeline visualization would go here
\\end{tikzpicture}
\\end{center}

\\newpage
"""
    
    def _generate_task_table(self, timeline: ProjectTimeline, colors: Dict[str, Any]) -> str:
        """Generate task table."""
        content = """
% Task table
\\section*{Task Details}

\\begin{longtable}{p{2cm}p{4cm}p{2cm}p{2cm}p{2cm}p{3cm}}
\\toprule
\\textbf{ID} & \\textbf{Task Name} & \\textbf{Start} & \\textbf{Due} & \\textbf{Milestone} & \\textbf{Category} \\\\
\\midrule
\\endhead

"""
        
        for task in timeline.tasks:
            task_id = self.escaper.escape_latex(task.id)
            task_name = self.escaper.escape_latex(task.name)
            start_date = task.start_date.strftime("%m/%d")
            due_date = task.due_date.strftime("%m/%d")
            category = self.escaper.escape_latex(task.category)
            milestone = "Yes" if task.is_milestone else "No"
            
            content += f"{task_id} & {task_name} & {start_date} & {due_date} & {milestone} & {category} \\\\\n"
        
        content += """
\\bottomrule
\\end{longtable}
"""
        
        return content
    
    def _generate_document_footer(self) -> str:
        """Generate document footer."""
        return LaTeXUtilities.generate_document_footer()


class MonthlyCalendarGenerator(BaseTemplateGenerator):
    """Generator for monthly calendar templates."""
    
    def _generate_document_content(self, timeline: ProjectTimeline, 
                                 config: Dict[str, Any]) -> str:
        """Generate monthly calendar document."""
        template = config['template']
        device_profile = config['device_profile']
        colors = config['colors']
        
        # Document header
        content = self._generate_document_header(template, device_profile)
        
        # Generate monthly calendars
        months = timeline.get_months_between()
        for month_info in months:
            content += self._generate_monthly_calendar(month_info, timeline, colors)
            content += "\\newpage\n"
        
        # Document footer
        content += self._generate_document_footer()
        
        return content
    
    def _generate_document_header(self, template: Any, device_profile: Any) -> str:
        """Generate LaTeX document header for calendar."""
        page_size = device_profile.get_layout_value('page_size', 'a4paper')
        orientation = template.orientation
        margin = device_profile.get_layout_value('margin', '0.3in')
        
        return f"""\\documentclass[{orientation},{page_size}]{{article}}
\\usepackage{{[utf8]{{inputenc}}}}
\\usepackage{{[T1]{{fontenc}}}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage{{[{orientation},margin={margin}]{{geometry}}}}
\\usepackage{{tikz}}
\\usepackage{{xcolor}}
\\usepackage{{array}}
\\usepackage{{fancyhdr}}

% Page setup
\\pagestyle{{empty}}
\\setlength{{\\parskip}}{{0.5em}}

% Use Helvetica for sans-serif
\\renewcommand{{\\familydefault}}{{\\sfdefault}}

% Color definitions
{self._generate_color_definitions()}

\\begin{{document}}
"""
    
    def _generate_color_definitions(self) -> str:
        """Generate LaTeX color definitions."""
        colors = [
            ("task", (59, 130, 246)),
            ("milestone", (147, 51, 234)),
            ("completed", (34, 197, 94)),
            ("inprogress", (251, 146, 60)),
            ("blocked", (239, 68, 68)),
            ("grid", (200, 200, 200)),
        ]
        
        color_defs = []
        for name, rgb in colors:
            color_defs.append(f"\\definecolor{{{name}}}{{RGB}}{{{rgb[0]}, {rgb[1]}, {rgb[2]}}}")
        
        return '\n'.join(color_defs)
    
    def _generate_monthly_calendar(self, month_info: MonthInfo, timeline: ProjectTimeline, 
                                 colors: Dict[str, Any]) -> str:
        """Generate monthly calendar view."""
        month_name = month_info.month_name
        year = month_info.year
        tasks = timeline.get_tasks_for_month(month_info)
        
        content = f"""
% {month_name} {year} Calendar
\\section*{{{month_name} {year}}}

\\begin{{center}}
\\begin{{tikzpicture}}[scale=0.9]
% Calendar grid would go here
% Task overlays would be positioned on the grid
\\end{{tikzpicture}}
\\end{{center}}

\\vspace{{1cm}}

\\subsection*{{Tasks for {month_name} {year}}}
\\begin{{itemize}}
"""
        
        for task in tasks:
            task_name = self.escaper.escape_latex(task.name)
            start_date = task.start_date.strftime("%m/%d")
            due_date = task.due_date.strftime("%m/%d")
            content += f"\\item {task_name} ({start_date} - {due_date})\n"
        
        content += "\\end{itemize}\n"
        
        return content
    
    def _generate_document_footer(self) -> str:
        """Generate document footer."""
        return LaTeXUtilities.generate_document_footer()


class WeeklyPlannerGenerator(BaseTemplateGenerator):
    """Generator for weekly planner templates."""
    
    def _generate_document_content(self, timeline: ProjectTimeline, 
                                 config: Dict[str, Any]) -> str:
        """Generate weekly planner document."""
        template = config['template']
        device_profile = config['device_profile']
        colors = config['colors']
        
        # Document header
        content = self._generate_document_header(template, device_profile)
        
        # Generate weekly planners
        weeks = self._get_weeks_between(timeline.start_date, timeline.end_date)
        for week_start in weeks:
            content += self._generate_weekly_planner(week_start, timeline, colors)
            content += "\\newpage\n"
        
        # Document footer
        content += self._generate_document_footer()
        
        return content
    
    def _get_weeks_between(self, start_date: date, end_date: date) -> List[date]:
        """Get list of week start dates between start and end dates."""
        weeks = []
        current_date = start_date
        
        # Find the Monday of the first week
        while current_date.weekday() != 0:  # Monday is 0
            current_date -= timedelta(days=1)
        
        while current_date <= end_date:
            weeks.append(current_date)
            current_date += timedelta(days=7)
        
        return weeks
    
    def _generate_weekly_planner(self, week_start: date, timeline: ProjectTimeline, 
                               colors: Dict[str, Any]) -> str:
        """Generate weekly planner view."""
        week_end = week_start + timedelta(days=6)
        week_tasks = []
        
        for task in timeline.tasks:
            if task.overlaps_with_range(week_start, week_end):
                week_tasks.append(task)
        
        content = f"""
% Week of {week_start.strftime('%B %d, %Y')}
\\section*{{Week of {week_start.strftime('%B %d, %Y')}}}

\\begin{{center}}
\\begin{{tikzpicture}}[scale=0.8]
% Weekly grid would go here
% Time slots and task scheduling would be positioned
\\end{{tikzpicture}}
\\end{{center}}

\\vspace{{1cm}}

\\subsection*{{Tasks for This Week}}
\\begin{{itemize}}
"""
        
        for task in week_tasks:
            task_name = self.escaper.escape_latex(task.name)
            start_date = task.start_date.strftime("%m/%d")
            due_date = task.due_date.strftime("%m/%d")
            content += f"\\item {task_name} ({start_date} - {due_date})\n"
        
        content += "\\end{itemize}\n"
        
        return content
    
    def _generate_document_header(self, template: Any, device_profile: Any) -> str:
        """Generate LaTeX document header for weekly planner."""
        page_size = device_profile.get_layout_value('page_size', 'a4paper')
        orientation = template.orientation
        margin = device_profile.get_layout_value('margin', '0.4in')
        
        return f"""\\documentclass[{orientation},{page_size}]{{article}}
\\usepackage{{[utf8]{{inputenc}}}}
\\usepackage{{[T1]{{fontenc}}}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage{{[{orientation},margin={margin}]{{geometry}}}}
\\usepackage{{tikz}}
\\usepackage{{xcolor}}
\\usepackage{{array}}

% Page setup
\\pagestyle{{empty}}
\\setlength{{\\parskip}}{{0.5em}}

% Use Helvetica for sans-serif
\\renewcommand{{\\familydefault}}{{\\sfdefault}}

% Color definitions
{self._generate_color_definitions()}

\\begin{{document}}
"""
    
    def _generate_color_definitions(self) -> str:
        """Generate LaTeX color definitions."""
        colors = [
            ("task", (59, 130, 246)),
            ("milestone", (147, 51, 234)),
            ("completed", (34, 197, 94)),
            ("inprogress", (251, 146, 60)),
            ("blocked", (239, 68, 68)),
            ("grid", (200, 200, 200)),
        ]
        
        color_defs = []
        for name, rgb in colors:
            color_defs.append(f"\\definecolor{{{name}}}{{RGB}}{{{rgb[0]}, {rgb[1]}, {rgb[2]}}}")
        
        return '\n'.join(color_defs)
    
    def _generate_document_footer(self) -> str:
        """Generate document footer."""
        return LaTeXUtilities.generate_document_footer()


class TemplateGeneratorFactory:
    """Factory for creating template generators."""
    
    @staticmethod
    def create_generator(template_type: str, config_manager: ConfigManager = None) -> BaseTemplateGenerator:
        """Create appropriate template generator based on type."""
        generators = {
            'gantt_timeline': GanttTimelineGenerator,
            'monthly_calendar': MonthlyCalendarGenerator,
            'weekly_planner': WeeklyPlannerGenerator,
        }
        
        generator_class = generators.get(template_type, GanttTimelineGenerator)
        return generator_class(config_manager)
