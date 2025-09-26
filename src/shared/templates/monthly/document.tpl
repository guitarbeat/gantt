% moved from templates/document.tpl
{{/* same content */}}
\documentclass[{{.Cfg.Layout.LaTeX.Document.FontSize}}]{extarticle}

% Core packages (load early)
\usepackage{expl3}
\usepackage{xparse}
\usepackage{calc}
\usepackage{geometry}

% Font configuration - use modern sans-serif font
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{lmodern}
\renewcommand{\familydefault}{\sfdefault}

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
\usepackage{blindtext}
% Do not load hyperref; define safe no-op link commands
\providecommand{\href}[2]{#2}
\providecommand{\hyperlink}[2]{#2}
\providecommand{\hypertarget}[2]{#2}
\providecommand{\hypersetup}[1]{}

{{if $.Cfg.Debug.ShowFrame}}\usepackage{showframe}{{end}}

{{- if not .Cfg.Debug.ShowLinks}}
\hypersetup{hidelinks,colorlinks=false,urlcolor=black,linkcolor=black,citecolor=black,pdfborder={0 0 0},pdfborderstyle={}}
{{- end}}

\geometry{paperwidth={{.Cfg.Layout.Paper.Width}}, paperheight={{.Cfg.Layout.Paper.Height}}}
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

\begin{document}

{{template "macros.tpl" .}}

  {{range .Pages -}}
    \include{ {{- .Name -}} }
  {{end}}
\end{document}
