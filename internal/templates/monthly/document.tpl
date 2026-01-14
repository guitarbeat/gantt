% moved from templates/document.tpl
{{/* same content */}}
\documentclass[{{.Cfg.Layout.LaTeX.Document.FontSize}}]{extarticle}

% Core packages (load early)
\usepackage{expl3}
\usepackage{xparse}
\usepackage{calc}
\usepackage{geometry}

% Font configuration - use modern sans-serif font with Unicode support
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{lmodern}
\renewcommand{\familydefault}{\sfdefault}

% Unicode character support
\usepackage{textcomp}
\usepackage{gensymb}

% Color and graphics
\usepackage[table]{xcolor}
\usepackage{graphicx}
\usepackage{tikz}
\usepackage{adjustbox}

% Table and array packages
\usepackage{array}
\usepackage{tabularx}
\usepackage{multirow}
\usepackage{makecell}
\usepackage{ragged2e}

% Layout and spacing
\usepackage{setspace}
\usepackage{leading}
\usepackage{dashrule}
\usepackage{varwidth}
\usepackage{wrapfig}
\usepackage{marginnote}
\usepackage{fancyhdr}

% Math and symbols
\usepackage{mathtools}
\usepackage{amssymb}

% Special features
\usepackage{multido}
\usepackage{pgffor}
\usepackage[most]{tcolorbox}
\usepackage{enumitem}
\usepackage{blindtext}
% Hyperlink support
\usepackage{hyperref}
\usepackage{bookmark}

{{if $.Cfg.Debug.ShowFrame}}\usepackage{showframe}{{end}}

\hypersetup{
    pdftitle={PhD Dissertation Planner {{.Cfg.Year}}},
    pdfauthor={PlannerGen},
    pdfsubject={PhD Dissertation Timeline},
    pdfkeywords={PhD, Dissertation, Planner, Timeline, {{.Cfg.Year}}},
    pdfcreator={PlannerGen},
{{- if not .Cfg.Debug.ShowLinks}}
    hidelinks,
    colorlinks=false,
    linkbordercolor={1 1 1},
    citebordercolor={1 1 1},
    filebordercolor={1 1 1},
    urlbordercolor={1 1 1},
    pdfborderstyle={/S/U/W 0},
    pdfborder={0 0 0}
{{- end}}
}

\geometry{verbose=false,paperwidth={{.Cfg.Layout.Paper.Width}}, paperheight={{.Cfg.Layout.Paper.Height}}}
\geometry{
  top={{.Cfg.Layout.Paper.Margin.Top}},
  bottom={{.Cfg.Layout.Paper.Margin.Bottom}},
  left={{.Cfg.Layout.Paper.Margin.Left}},
  right={{.Cfg.Layout.Paper.Margin.Right}},
  marginparwidth={{.Cfg.Layout.Paper.MarginParWidth}},
  marginparsep={{.Cfg.Layout.Paper.MarginParSep}}
}

\pagestyle{empty}
{{if $.Cfg.Layout.Paper.ReverseMargins}}\reversemarginpar{{end}}
\newcolumntype{Y}{>{\centering\arraybackslash}X}
\parindent={{.Cfg.Layout.LaTeX.Document.ParIndent}}
\fboxsep0pt

% Suppress verbose output
\hoffset=0pt
\voffset=0pt

\begin{document}

{{template "macros.tpl" .}}

  {{range .Pages -}}
    \include{ {{- .Name -}} .tex}
  {{end}}
\end{document}
