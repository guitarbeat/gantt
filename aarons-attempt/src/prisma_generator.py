#!/usr/bin/env python3
"""
PRISMA Flow Diagram Generator for LaTeX.
Generates PRISMA 2020-compliant flow diagrams for systematic reviews and meta-analyses.
"""

import logging
from dataclasses import dataclass, field
from typing import List, Dict, Optional, Tuple
from datetime import date

from .config import config


@dataclass
class PRISMAData:
    """Data structure for PRISMA flow diagram."""
    
    # Identification phase
    records_identified_databases: int = 0
    records_identified_registers: int = 0
    records_identified_other: int = 0
    
    # Screening phase
    records_after_duplicates_removed: int = 0
    records_screened: int = 0
    records_excluded: int = 0
    
    # Eligibility phase
    full_text_articles_assessed: int = 0
    full_text_articles_excluded: int = 0
    exclusion_reasons: List[str] = field(default_factory=list)
    
    # Included phase
    studies_included_qualitative: int = 0
    studies_included_quantitative: int = 0
    studies_included_meta_analysis: int = 0
    
    # Additional data
    title: str = "PRISMA Flow Diagram"
    review_type: str = "Systematic Review"
    date_range: str = ""
    
    def total_identified(self) -> int:
        """Calculate total records identified."""
        return (self.records_identified_databases + 
                self.records_identified_registers + 
                self.records_identified_other)
    
    def total_included(self) -> int:
        """Calculate total studies included."""
        return (self.studies_included_qualitative + 
                self.studies_included_quantitative + 
                self.studies_included_meta_analysis)


class PRISMAGenerator:
    """Generates PRISMA flow diagrams in LaTeX format."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def generate_diagram(self, prisma_data: PRISMAData) -> str:
        """Generate complete PRISMA flow diagram."""
        content = self._generate_header()
        content += self._generate_title(prisma_data)
        content += self._generate_flow_diagram(prisma_data)
        content += self._generate_footer()
        return content
    
    def _generate_header(self) -> str:
        """Generate LaTeX header for PRISMA diagram."""
        return """\\documentclass[portrait,a4paper]{article}
\\usepackage[utf8]{inputenc}
\\usepackage[T1]{fontenc}
\\usepackage{lmodern}
\\usepackage{helvet}
\\usepackage[portrait,margin=0.5in]{geometry}
\\usepackage{tikz}
\\usepackage{xcolor}
\\usepackage{enumitem}
\\usepackage{hyperref}

% TikZ libraries for flow diagrams
\\usetikzlibrary{arrows.meta,shapes.geometric,positioning,calc,decorations.pathmorphing,patterns,shadows,fit,backgrounds,matrix,chains,scopes}

% Page setup
\\pagestyle{empty}
\\hypersetup{
    colorlinks=true,
    linkcolor=blue,
    urlcolor=blue,
    citecolor=blue,
    bookmarksopen=true,
    bookmarksnumbered=true
}

% Use Helvetica for sans-serif
\\renewcommand{\\familydefault}{\\sfdefault}

% PRISMA colors are defined in the main document header

\\begin{document}
"""
    
    def _generate_title(self, prisma_data: PRISMAData) -> str:
        """Generate title section."""
        title = self._escape_latex(prisma_data.title)
        review_type = self._escape_latex(prisma_data.review_type)
        date_range = self._escape_latex(prisma_data.date_range)
        
        return f"""
\\begin{{center}}
\\vspace*{{1cm}}
{{\\Huge\\bfseries {title}}}\\\\
\\vspace{{0.5cm}}
{{\\Large {review_type}}}\\\\
\\vspace{{0.3cm}}
{{\\large {date_range}}}\\\\
\\vspace{{1cm}}
\\end{{center}}
"""
    
    def _generate_flow_diagram(self, prisma_data: PRISMAData) -> str:
        """Generate the main PRISMA flow diagram."""
        content = """
\\begin{center}
\\begin{tikzpicture}[
    node distance=1.5cm,
    every node/.style={font=\\small},
    box/.style={rectangle, draw, fill=white, text width=3cm, text centered, minimum height=0.8cm},
    decision/.style={diamond, draw, fill=prismagray!20, text width=2.5cm, text centered, minimum height=0.8cm},
    exclusion/.style={rectangle, draw, fill=prismared!20, text width=2.5cm, text centered, minimum height=0.6cm},
    inclusion/.style={rectangle, draw, fill=prismagreen!20, text width=2.5cm, text centered, minimum height=0.8cm}
]

% Identification phase
\\node[box] (ident1) at (0,8) {Records identified\\\\(n = """ + str(prisma_data.total_identified()) + """)};
\\node[box, below=of ident1] (ident2) {Records from databases\\\\(n = """ + str(prisma_data.records_identified_databases) + """)};
\\node[box, left=of ident2] (ident3) {Records from registers\\\\(n = """ + str(prisma_data.records_identified_registers) + """)};
\\node[box, right=of ident2] (ident4) {Records from other sources\\\\(n = """ + str(prisma_data.records_identified_other) + """)};

% Duplicates removal
\\node[box, below=of ident2] (dup) {Records after duplicates removed\\\\(n = """ + str(prisma_data.records_after_duplicates_removed) + """)};

% Screening phase
\\node[box, below=of dup] (screen) {Records screened\\\\(n = """ + str(prisma_data.records_screened) + """)};
\\node[exclusion, right=of screen] (excl1) {Records excluded\\\\(n = """ + str(prisma_data.records_excluded) + """)};

% Eligibility phase
\\node[box, below=of screen] (elig) {Full-text articles assessed for eligibility\\\\(n = """ + str(prisma_data.full_text_articles_assessed) + """)};
\\node[exclusion, right=of elig] (excl2) {Full-text articles excluded\\\\(n = """ + str(prisma_data.full_text_articles_excluded) + """)};

% Inclusion phase
\\node[box, below=of elig] (incl) {Studies included in review\\\\(n = """ + str(prisma_data.total_included()) + """)};

% Qualitative synthesis
\\node[inclusion, below left=of incl] (qual) {Studies included in qualitative synthesis\\\\(n = """ + str(prisma_data.studies_included_qualitative) + """)};

% Quantitative synthesis
\\node[inclusion, below right=of incl] (quant) {Studies included in quantitative synthesis\\\\(n = """ + str(prisma_data.studies_included_quantitative) + """)};

% Meta-analysis
\\node[inclusion, below=of quant] (meta) {Studies included in meta-analysis\\\\(n = """ + str(prisma_data.studies_included_meta_analysis) + """)};

% Arrows
\\draw[->] (ident1) -- (ident2);
\\draw[->] (ident3) -- (ident2);
\\draw[->] (ident4) -- (ident2);
\\draw[->] (ident2) -- (dup);
\\draw[->] (dup) -- (screen);
\\draw[->] (screen) -- (excl1);
\\draw[->] (screen) -- (elig);
\\draw[->] (elig) -- (excl2);
\\draw[->] (elig) -- (incl);
\\draw[->] (incl) -- (qual);
\\draw[->] (incl) -- (quant);
\\draw[->] (quant) -- (meta);

% Exclusion reasons box
"""
        
        if prisma_data.exclusion_reasons:
            content += self._generate_exclusion_reasons(prisma_data.exclusion_reasons)
        
        content += """
\\end{tikzpicture}
\\end{center}
"""
        return content
    
    def _generate_exclusion_reasons(self, reasons: List[str]) -> str:
        """Generate exclusion reasons box."""
        content = """
\\node[decision, below=of excl2, text width=4cm] (reasons) {Exclusion reasons:};
"""
        
        for i, reason in enumerate(reasons[:5]):  # Limit to 5 reasons
            escaped_reason = self._escape_latex(reason)
            content += f"\\node[below=of reasons, text width=4cm, font=\\tiny] (reason{i}) {{{escaped_reason}}};\\n"
        
        return content
    
    def _generate_footer(self) -> str:
        """Generate document footer."""
        return """
\\vspace{1cm}
\\begin{center}
\\textit{PRISMA 2020 Flow Diagram for new systematic reviews which included searches of databases, registers and other sources}
\\end{center}

\\end{document}
"""
    
    def _escape_latex(self, text: str) -> str:
        """Escape special LaTeX characters."""
        if not text:
            return ""
        
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
            '\\': r'\textbackslash{}'
        }
        
        for char, replacement in replacements.items():
            text = text.replace(char, replacement)
        
        return text


def create_sample_prisma_data() -> PRISMAData:
    """Create sample PRISMA data for testing."""
    return PRISMAData(
        records_identified_databases=1250,
        records_identified_registers=45,
        records_identified_other=12,
        records_after_duplicates_removed=1180,
        records_screened=1180,
        records_excluded=1100,
        full_text_articles_assessed=80,
        full_text_articles_excluded=65,
        exclusion_reasons=[
            "Not relevant to research question",
            "Insufficient data reported",
            "Wrong study design",
            "Duplicate publication",
            "Language not English"
        ],
        studies_included_qualitative=15,
        studies_included_quantitative=12,
        studies_included_meta_analysis=8,
        title="Systematic Review of Project Management Tools",
        review_type="Systematic Review and Meta-Analysis",
        date_range="January 2020 - December 2024"
    )
