#!/usr/bin/env python3
"""Simple LaTeX generator to test the basic functionality."""

import csv
import re
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

def build_task_hierarchy(csv_data):
    """Build a hierarchical structure of tasks based on parent relationships."""
    tasks = {}
    children = {}
    
    # First pass: create task lookup and identify children
    for row in csv_data:
        task_id = row.get('Task ID', '')
        parent_id = row.get('Parent Task ID', '').strip()
        
        tasks[task_id] = row
        
        if parent_id:
            if parent_id not in children:
                children[parent_id] = []
            children[parent_id].append(task_id)
    
    return tasks, children

def generate_latex_document(csv_data, title="Project Timeline"):
    """Generate a simple LaTeX document with hierarchical subtasks."""
    
    # Document header - PORTRAIT orientation with NO MARGINS for maximum space
    latex = f"""\\documentclass[a4paper]{{article}}
\\usepackage[utf8]{{inputenc}}
\\usepackage[T1]{{fontenc}}
\\usepackage{{lmodern}}
\\usepackage{{helvet}}
\\usepackage[portrait,margin=0.2in,top=0.3in,bottom=0.3in]{{geometry}}
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

\\begin{{enumerate}}[leftmargin=0.3cm, itemsep=0.3em, parsep=0.1em, topsep=0.1em]
"""
    
    # Build task hierarchy
    tasks, children = build_task_hierarchy(csv_data)
    
    # Add tasks with hierarchical subtasks - only show parent tasks, subtasks will be indented
    task_counter = 1
    for row in csv_data:
        task_id = row.get('Task ID', '')
        parent_id = row.get('Parent Task ID', '').strip()
        
        # Skip subtasks here - they'll be handled under their parents
        if parent_id:
            continue
            
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
            'ADMINISTRATIVE': 'administrative',
            'DISSERTATION': 'milestone'
        }
        color = color_map.get(category, 'black')
        
        # Check if it's a milestone
        is_milestone = description.startswith('MILESTONE:') if description else False
        
        # Main task formatting with better visual hierarchy
        latex += f"""
    \\item \\textcolor{{{color}}}{{\\textbf{{\\Large {task_name}}}}}
          \\hfill \\cellcolor{{{color}!15}}{{\\textbf{{\\small {category}}}}}
          \\\\ \\textcolor{{black!60}}{{\\textbf{{Duration:}} {start_date} -- {due_date}}}"""
        
        if is_milestone:
            latex += f" \\\\hfill \\textcolor{{milestone}}{{\\textbf{{$\\star$ MILESTONE}}}}"
        
        if description:
            # Remove MILESTONE: prefix if present for cleaner display
            clean_description = description.replace('MILESTONE:', '').strip()
            if clean_description:
                latex += f" \\\\ \\textcolor{{black!80}}{{\\small {clean_description}}}"
        
        # Add subtasks if they exist with better visual hierarchy
        if task_id in children:
            latex += " \\\\ \\vspace{{0.2em}} \\\\ \\textcolor{{black!50}}{{\\small \\textbf{{Subtasks:}}}}"
            latex += f" \\\\ \\begin{{itemize}}[leftmargin=0.8cm, itemsep=0.1em, parsep=0.05em, label=\\textcolor{{{color}!60}}{{\\tiny$\\bullet$}}]"
            for child_id in children[task_id]:
                child_task = tasks[child_id]
                child_name = escape_latex(child_task.get('Task Name', ''))
                child_start = child_task.get('Start Date', '')
                child_due = child_task.get('Due Date', '')
                child_desc = escape_latex(child_task.get('Description', ''))
                child_milestone = child_desc.startswith('MILESTONE:') if child_desc else False
                
                # Better subtask formatting with clear visual hierarchy
                latex += f" \\\\item \\textcolor{{{color}!70}}{{\\textbf{{\\small {child_name}}}}}"
                latex += f" \\\\hfill \\textcolor{{black!50}}{{\\tiny {child_start} -- {child_due}}}"
                
                if child_milestone:
                    latex += f" \\\\hfill \\textcolor{{milestone}}{{\\textbf{{\\tiny$\\star$}}}}"
                
                if child_desc:
                    clean_child_desc = child_desc.replace('MILESTONE:', '').strip()
                    if clean_child_desc:
                        latex += f" \\\\ \\\\textcolor{{black!70}}{{\\tiny {clean_child_desc}}}"
            
            latex += " \\\\end{{itemize}}"
        
        latex += "\n"
        task_counter += 1
    
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
