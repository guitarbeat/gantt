{%
{{ if not .Body.Large }} \renewcommand{\arraystretch}{\myNumArrayStretch} {{ end }}
\setlength{\tabcolsep}{\myLenTabColSep}
{{ $tbl := .Body.Month.DefineTable .Body.TableType .Body.Large }}
{{ if $tbl }}
{{$tbl}}
{{ end }}
  {{ $mname := .Body.Month.MaybeName .Body.Large }}
  {{ if $mname }}{{$mname}}{{ end }}
  {{ if $.Body.Large }} \hline {{ end }}
  {{ $wh := .Body.Month.WeekHeader .Body.Large }}
  {{ if $wh }}{{$wh}} \\ {{ if .Body.Large }} \noalign{\hrule height \myLenLineThicknessThick} {{ else }} \hline {{ end}}{{ end }}
  {{ range $i, $week := .Body.Month.Weeks }}
  {{ if $week.HasDays }}
  {{$week.WeekNumber $.Body.Large}} &
    {{ range $j, $day := $week.Days }}
      {{ $cell := $day.Day $.Body.Today $.Body.Large }}
      {{ if $cell }}
        {{$cell}}
      {{ end }}
      {{ if eq $j 6 }}
        \\ {{ if $.Body.Large }} \hline {{ end }}
      {{ else }} & {{ end }}
    {{ end }}
  {{ end }}
  {{ end }}
  {{ .Body.Month.EndTable .Body.TableType }}
}
% Simple macros template without problematic LaTeX commands
\ExplSyntaxOn
\cs_new_eq:NN \Repeat \prg_replicate:nn
\ExplSyntaxOff

{{- $numbers := .Cfg.Layout.Numbers -}}

\newlength{\myLenTabColSep}
\newlength{\myLenLineThicknessDefault}
\newlength{\myLenLineThicknessThick}
\newlength{\myLenLineHeightButLine}
\newlength{\myLenTwoColSep}
\newlength{\myLenTwoCol}
\newlength{\myLenTriColSep}
\newlength{\myLenTriCol}
\newlength{\myLenFiveColSep}
\newlength{\myLenFiveCol}
\newlength{\myLenMonthlyCellHeight}

\newlength{\myLenHeaderResizeBox}
\newlength{\myLenHeaderSideMonthsWidth}

\setlength{\myLenTabColSep}{3.5pt}
\setlength{\myLenLineThicknessDefault}{.4pt}
\setlength{\myLenLineThicknessThick}{.8pt}
\setlength{\myLenLineHeightButLine}{\dimexpr5mm-.4pt}
\setlength{\myLenTwoColSep}{5pt}
\setlength{\myLenTwoCol}{5pt}
\setlength{\myLenTriColSep}{5pt}
\setlength{\myLenTriCol}{5pt}
\setlength{\myLenFiveColSep}{5pt}
\setlength{\myLenFiveCol}{5pt}
\setlength{\myLenMonthlyCellHeight}{78pt}

\setlength{\myLenHeaderResizeBox}{6mm}
\setlength{\myLenHeaderSideMonthsWidth}{14.5cm}

% Simple task bar definitions
% * Define a fixed task font size macro
\newcommand{\TaskFontSize}{\footnotesize}
\newlength{\TaskBarHeight}
\setlength{\TaskBarHeight}{4mm}
\newlength{\TaskBorderWidth}
\setlength{\TaskBorderWidth}{0.6pt}
\newlength{\TaskPaddingH}
\setlength{\TaskPaddingH}{1.5mm}
\newlength{\TaskPaddingV}
\setlength{\TaskPaddingV}{0.5mm}
% * Global vertical nudge for task elements (push tasks slightly lower)
\newlength{\TaskVerticalOffset}
\setlength{\TaskVerticalOffset}{0.7mm}

% Array stretch macro
\newcommand{\myNumArrayStretch}{1.2}

% Line thickness macro
\newcommand{\myLineThick}{\rule{\linewidth}{\myLenLineThicknessThick}}

% Category palette setup macro
\newcommand{\SetupDefaultCategoryPalette}[1]{#1}

% Simple task rendering
\newcommand{\SimpleTaskBar}[4]{%
  \vspace*{\TaskVerticalOffset}%
  \fbox{\parbox{\dimexpr#3-2\TaskPaddingH\relax}{%
    \vspace{\TaskPaddingV}%
    \centering\small\textbf{#1}%
    \vspace{\TaskPaddingV}%
  }}%
}

% Task overlay box macros - pill shaped with rounded corners
\newcommand{\TaskOverlayBox}[3]{%
  \vspace*{\TaskVerticalOffset}%
  \begin{tcolorbox}[enhanced, boxrule=0.9pt, arc=9pt, drop shadow={0.6pt}{-0.6pt}{0pt}{black!20},
    left=2.8mm, right=2.8mm, top=1.8mm, bottom=1.8mm,
    colback=#1!20, colframe=#1!80,
    width=\linewidth, halign=center]
    \TaskFontSize\textbf{#2}\\#3%
  \end{tcolorbox}%
}

% Multi-day task bar drawing macro to centralize styling
% Args: 1=x(pt), 2=y(pt), 3=width(pt), 4=height(pt), 5=color, 6=label
\newcommand{\DrawTaskBar}[6]{%
  \begin{tikzpicture}[overlay]
    \node[anchor=north west, inner sep=0pt] at (#1,#2) {
      \begin{tcolorbox}[enhanced, boxrule=0pt, arc=2pt, drop shadow,
        left=1.5mm, right=1.5mm, top=0.5mm, bottom=0.5mm,
        width=#3pt, height=#4pt,
        colback=#5,
        borderline west={1.4pt}{0pt}{#5!60!black},
        borderline east={1.0pt}{0pt}{#5!45}]
        {\footnotesize #6}
      \end{tcolorbox}
    };
  \end{tikzpicture}%
}

\newcommand{\TaskOverlayBoxP}[3]{%
  \vspace*{\TaskVerticalOffset}%
  \begin{tcolorbox}[enhanced, boxrule=0.9pt, arc=9pt, drop shadow={0.6pt}{-0.6pt}{0pt}{black!20},
    left=3mm, right=3mm, top=1.8mm, bottom=1.8mm,
    colback=#2!20, colframe=#2!80,
    width=\linewidth, halign=center]
    \TaskFontSize\textbf{#1}\\#3%
  \end{tcolorbox}%
}

% Task compact box macro with pill shape and better spacing
\newcommand{\TaskCompactBox}[4]{%
  \vspace*{#1}%
  \vspace*{\TaskVerticalOffset}%
  \begin{tcolorbox}[enhanced, boxrule=0.7pt, arc=8pt, drop shadow={0.4pt}{-0.4pt}{0pt}{black!15},
    left=2mm, right=2mm, top=1.4mm, bottom=1.4mm,
    colback=#3!20, colframe=#3!70,
    width=\linewidth, halign=center, height=#2]
    \vfil
    \TaskFontSize\textbf{#4}%
    \vfil
  \end{tcolorbox}%
}

% Underline macro
\newcommand{\myUnderline}[1]{%
  \underline{\textbf{#1}}%
}

% Colored circle macro for legend - bigger circles
\newcommand{\ColorCircle}[1]{%
  \textcolor{#1}{\Large$\bullet$}%
}



% Color legend macro for task categories - uses circles instead of boxes
\newcommand{\ColorLegend}{%
  {\centering
    \ColorCircle{blue}~\small Proposal%
    \hspace{1.5em}%
    \ColorCircle{orange}~\small Laser%
    \hspace{1.5em}%
    \ColorCircle{green}~\small Imaging%
    \hspace{1.5em}%
    \ColorCircle{purple}~\small Admin%
    \hspace{1.5em}%
    \ColorCircle{red}~\small Dissertation%
    \hspace{1.5em}%
    \ColorCircle{teal}~\small Research%
    \hspace{1.5em}%
    \ColorCircle{gray}~\small Publication%
  \par}
}
% Simple macros template without problematic LaTeX commands
\ExplSyntaxOn
\cs_new_eq:NN \Repeat \prg_replicate:nn
\ExplSyntaxOff

{{- $numbers := .Cfg.Layout.Numbers -}}

\newlength{\myLenTabColSep}
\newlength{\myLenLineThicknessDefault}
\newlength{\myLenLineThicknessThick}
\newlength{\myLenLineHeightButLine}
\newlength{\myLenTwoColSep}
\newlength{\myLenTwoCol}
\newlength{\myLenTriColSep}
\newlength{\myLenTriCol}
\newlength{\myLenFiveColSep}
\newlength{\myLenFiveCol}
\newlength{\myLenMonthlyCellHeight}

\newlength{\myLenHeaderResizeBox}
\newlength{\myLenHeaderSideMonthsWidth}

{{- $lengths := .Cfg.Layout.Lengths -}}
\setlength{\myLenTabColSep}{ {{- $lengths.TabColSep -}} }
\setlength{\myLenLineThicknessDefault}{ {{- $lengths.LineThicknessDefault -}} }
\setlength{\myLenLineThicknessThick}{ {{- $lengths.LineThicknessThick -}} }
\setlength{\myLenLineHeightButLine}{ {{- $lengths.LineHeightButLine -}} }
\setlength{\myLenTwoColSep}{ {{- $lengths.TwoColSep -}} }
\setlength{\myLenTwoCol}{ {{- $lengths.TwoCol -}} }
\setlength{\myLenTriColSep}{ {{- $lengths.TriColSep -}} }
\setlength{\myLenTriCol}{ {{- $lengths.TriCol -}} }
\setlength{\myLenFiveColSep}{ {{- $lengths.FiveColSep -}} }
\setlength{\myLenFiveCol}{ {{- $lengths.FiveCol -}} }
\setlength{\myLenMonthlyCellHeight}{ {{- $lengths.MonthlyCellHeight -}} }

\setlength{\myLenHeaderResizeBox}{ {{- $lengths.HeaderResizeBox -}} }
\setlength{\myLenHeaderSideMonthsWidth}{ {{- $lengths.HeaderSideMonthsWidth -}} }

% Simple task bar definitions
\newlength{\TaskBarHeight}
\setlength{\TaskBarHeight}{4mm}
\newlength{\TaskBorderWidth}
\setlength{\TaskBorderWidth}{0.6pt}
\newlength{\TaskPaddingH}
\setlength{\TaskPaddingH}{1.5mm}
\newlength{\TaskPaddingV}
\setlength{\TaskPaddingV}{0.5mm}

% Simple task rendering
\newcommand{\SimpleTaskBar}[4]{%
  \fbox{\parbox{\dimexpr#3-2\TaskPaddingH\relax}{%
    \vspace{\TaskPaddingV}%
    \centering\small\textbf{#1}%
    \vspace{\TaskPaddingV}%
  }}%
}
% moved from templates/document.tpl
{{/* same content */}}
\documentclass[9pt]{extarticle}

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
\usepackage{hyperref}

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
\parindent=0pt
\fboxsep0pt

\begin{document}

{{template "macros.tpl" .}}

  {{range .Pages -}}
    \include{ {{- .Name -}} }
  {{end}}
\end{document}
% Setup category palette for this month
\SetupDefaultCategoryPalette

{{- template "calendar_table.tpl" dict "Cfg" .Cfg "Body" .Body -}}

% Legend at bottom of page
\vfill
{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
{\centering
{{- range $color, $category := $taskColors -}}
\ColorCircle{ {{- $color -}} } \small{ {{- $category -}} }%
\hspace{1.5em}%
{{- end -}}
\par}
{{- else -}}
\ColorLegend
{{- end -}}
{{ template "page_header.tpl" dict "Cfg" .Cfg "Body" .Body }}
{{ template "monthly_body.tpl" dict "Cfg" .Cfg "Body" .Body }}

\pagebreak
{\noindent\Large\renewcommand{\arraystretch}{\myNumArrayStretch}
{{- .Body.Breadcrumb -}}
\hfill%
{{ .Body.Extra.Table false -}}
}
\myLineThick