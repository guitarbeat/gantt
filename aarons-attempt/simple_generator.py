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
    """Escape special LaTeX characters and handle Unicode."""
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
        '\\': r'\textbackslash{}',
        '≥': r'$\geq$',
        '≤': r'$\leq$',
        '★': r'$\star$'
    }
    
    for char, replacement in replacements.items():
        text = text.replace(char, replacement)
    
    return text

def generate_latex_document(csv_data, title="Project Timeline"):
    """Generate a simple LaTeX document."""
    
    # Document header - PORTRAIT orientation with NO MARGINS for maximum space
    latex = f"""\\documentclass[portrait,a4paper]{{article}}
\\usepackage{{[utf8]{{inputenc}}}}
\\usepackage{{[T1]{{fontenc}}}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage{{[portrait,margin=0.2in,top=0.3in,bottom=0.3in]{{geometry}}}}
\\usepackage{{tikz}}
\\usepackage{{xcolor}}
\\usepackage{{array}}
\\usepackage{{fancyhdr}}
\\usepackage{{hyperref}}
\\usepackage{{enumitem}}
\\usepackage{{parskip}}
\\usepackage{{amsmath}}

% Page setup - MAXIMUM space usage
\\pagestyle{{empty}}
\\setlength{{\\parskip}}{{0.1em}}
\\setlength{{\\parsep}}{{0.05em}}
\\setlength{{\\itemsep}}{{0.05em}}
\\setlength{{\\topsep}}{{0.1em}}

% Use Helvetica for sans-serif
\\renewcommand{{\\familydefault}}{{\\sfdefault}}

% Color definitions inspired by @latex/ helpers
\\definecolor{{researchcore}}{{RGB}}{{59, 130, 246}}
\\definecolor{{researchexp}}{{RGB}}{{16, 185, 129}}
\\definecolor{{researchout}}{{RGB}}{{245, 158, 11}}
\\definecolor{{administrative}}{{RGB}}{{107, 114, 128}}
\\definecolor{{milestone}}{{RGB}}{{147, 51, 234}}

% LaTeX helper functions inspired by @latex/texfuncs.go
\\newcommand{{\\cellcolor}}[2]{{\\colorbox{{#1}}{{#2}}}}
\\newcommand{{\\textcolor}}[2]{{\\textcolor{{#1}}{{#2}}}}
\\newcommand{{\\hyperlink}}[2]{{\\hyperlink{{#1}}{{#2}}}}
\\newcommand{{\\hypertarget}}[2]{{\\hypertarget{{#1}}{{#2}}}}

\\begin{{document}}

% Compact title page
\\begin{{titlepage}}
\\centering
\\vspace*{{0.5cm}}

{{\\LARGE\\textbf{{{escape_latex(title)}}}}}

\\vspace{{0.3cm}}
{{\\large PhD Research Calendar}}

\\vspace{{0.5cm}}

\\begin{{minipage}}{{0.95\\textwidth}}
\\centering
\\textbf{{Total Tasks:}} {len(csv_data)} tasks \\hfill \\textbf{{Generated:}} {datetime.now().strftime('%B %d, %Y')}
\\end{{minipage}}

\\vfill

\\end{{titlepage}}

\\newpage

% Task list with MAXIMUM space usage - no margins
\\section{{Complete Task List}}
\\vspace{{0.1cm}}

\\begin{{enumerate}}[leftmargin=0.3cm, itemsep=0.1em, parsep=0.05em, topsep=0.1em]
"""
    
    # Add tasks with improved formatting for portrait mode
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
        
        # Ultra-compact formatting with NO margins - inspired by @latex/texfuncs.go
        latex += f"""
    \\item \\textcolor{{{color}}}{{\\textbf{{\\large {task_name}}}}}
          \\hfill \\cellcolor{{{color}!20}}{{\\textbf{{\\small {category}}}}}
          \\textcolor{{black!70}}{{\\textbf{{Duration:}} {start_date} -- {due_date}}}"""
        
        if is_milestone:
            latex += f" \\hfill \\textcolor{{milestone}}{{\\textbf{{$\\star$ MILESTONE}}}}"
        
        if description:
            # Remove MILESTONE: prefix if present for cleaner display
            clean_description = description.replace('MILESTONE:', '').strip()
            if clean_description:
                latex += f" \\\\ \\textcolor{{black!85}}{{\\small {clean_description}}}"
        
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
