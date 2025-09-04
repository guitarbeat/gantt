#!/usr/bin/env python3
"""Simple LaTeX generator to test the basic functionality."""

import csv
from datetime import datetime, date
from pathlib import Path

def read_csv_data(file_path):
    """Read CSV data and return as list of dictionaries."""
    with open(file_path, 'r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        return list(reader)

def escape_latex(text):
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

def generate_latex_document(csv_data, title="Project Timeline"):
    """Generate a simple LaTeX document."""
    
    # Document header
    latex = f"""\\documentclass[landscape,a4paper]{{article}}
\\usepackage{{[utf8]{{inputenc}}}}
\\usepackage{{[T1]{{fontenc}}}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage{{[landscape,margin=0.5in]{{geometry}}}}
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
\\definecolor{{researchcore}}{{RGB}}{{59, 130, 246}}
\\definecolor{{researchexp}}{{RGB}}{{16, 185, 129}}
\\definecolor{{researchout}}{{RGB}}{{245, 158, 11}}
\\definecolor{{administrative}}{{RGB}}{{107, 114, 128}}
\\definecolor{{milestone}}{{RGB}}{{147, 51, 234}}

\\begin{{document}}

% Title page
\\begin{{titlepage}}
\\centering
\\vspace*{{2cm}}

{{\\LARGE\\textbf{{{escape_latex(title)}}}}}

\\vspace{{1cm}}
{{\\large PhD Research Calendar}}

\\vspace{{2cm}}

\\begin{{minipage}}{{0.9\\textwidth}}
\\centering
\\textbf{{Total Tasks:}} {len(csv_data)} tasks\\\\
\\textbf{{Generated:}} {datetime.now().strftime('%B %d, %Y')}
\\end{{minipage}}

\\vfill

\\end{{titlepage}}

\\newpage

% Task list
\\section{{Complete Task List}}
\\vspace{{0.5cm}}

\\begin{{enumerate}}[leftmargin=1.5cm, itemsep=1em]
"""
    
    # Add tasks
    for i, row in enumerate(csv_data, 1):
        task_name = escape_latex(row.get('Task Name', ''))
        category = row.get('Category', '')
        start_date = row.get('Start Date', '')
        due_date = row.get('Due Date', '')
        description = escape_latex(row.get('Description', ''))
        
        # Determine color based on category
        color_map = {
            'PROPOSAL': 'researchcore',
            'LASER': 'researchexp',
            'EXPERIMENTAL': 'researchexp',
            'PUBLICATION': 'researchout',
            'ADMINISTRATIVE': 'administrative'
        }
        color = color_map.get(category, 'black')
        
        # Check if it's a milestone
        is_milestone = description.startswith('MILESTONE:') if description else False
        
        latex += f"""
    \\item \\textcolor{{{color}}}{{\\textbf{{\\large {task_name}}}}}
          \\hfill \\textcolor{{black!60}}{{\\small [{category}]}}
          
          \\vspace{{0.2em}}
          \\textcolor{{black!70}}{{\\textbf{{Duration:}} {start_date} -- {due_date}}}
"""
        
        if is_milestone:
            latex += f"          \\textcolor{{orange}}{{\\textbf{{ [MILESTONE]}}}}\\n"
        
        if description:
            latex += f"""
          \\vspace{{0.4em}}
          \\begin{{minipage}}[t]{{0.9\\textwidth}}
          \\textcolor{{black!85}}{{{description}}}
          \\end{{minipage}}
"""
        
        latex += "\n"
    
    # Document footer
    latex += """\\end{enumerate}

\\end{document}
"""
    
    return latex

def main():
    """Main function."""
    # Read CSV data
    csv_file = "../input/data.cleaned.csv"
    csv_data = read_csv_data(csv_file)
    
    # Generate LaTeX
    latex_content = generate_latex_document(csv_data, "Proposal Timeline")
    
    # Write to file
    output_file = "output/tex/proposal-timeline.tex"
    Path(output_file).parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(latex_content)
    
    print(f"✅ LaTeX file generated: {output_file}")
    print(f"📊 Timeline contains {len(csv_data)} tasks")
    print(f"🔨 To compile: pdflatex {output_file}")

if __name__ == "__main__":
    main()
