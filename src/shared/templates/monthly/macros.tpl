% Simple macros template without problematic LaTeX commands
\ExplSyntaxOn
\cs_new_eq:NN \Repeat \prg_replicate:nn
\ExplSyntaxOff

{{- $numbers := .Cfg.Layout.Numbers -}}

% Task colors are now generated algorithmically - no need for predefined colors

\newlength{\myLenTabColSep}
\newlength{\myLenLineThicknessDefault}
\newlength{\myLenLineThicknessThick}
\newlength{\myLenLineHeightButLine}
\newlength{\myLenColSep}
\newlength{\myLenCol}
\newlength{\myLenMonthlyCellHeight}

\newlength{\myLenHeaderResizeBox}
\newlength{\myLenHeaderSideMonthsWidth}

\setlength{\myLenTabColSep}{ {{.Cfg.Layout.LaTeX.TabColSep}} }
\setlength{\myLenLineThicknessDefault}{ {{.Cfg.Layout.LaTeX.LineThicknessDefault}} }
\setlength{\myLenLineThicknessThick}{ {{.Cfg.Layout.LaTeX.LineThicknessThick}} }
\setlength{\myLenLineHeightButLine}{\dimexpr5mm-.4pt}
\setlength{\myLenColSep}{ {{.Cfg.Layout.LaTeX.ColSep}} }
\setlength{\myLenCol}{ {{.Cfg.Layout.LaTeX.Spacing.Col}} }
\setlength{\myLenMonthlyCellHeight}{ {{.Cfg.Layout.LaTeX.MonthlyCellHeight}} }

\setlength{\myLenHeaderResizeBox}{ {{.Cfg.Layout.LaTeX.HeaderResizeBox}} }
\setlength{\myLenHeaderSideMonthsWidth}{ {{.Cfg.Layout.LaTeX.HeaderSideMonthsWidth}} }

% Simple task bar definitions
% * Define fixed font size macros for task title and body
\newcommand{\TaskTitleSize}{ {{.Cfg.Layout.LaTeX.TaskFontSize}} }
\newcommand{\TaskFontSize}{\footnotesize}
\newlength{\TaskBarHeight}
\setlength{\TaskBarHeight}{ {{.Cfg.Layout.LaTeX.TaskBarHeight}} }
\newlength{\TaskBorderWidth}
\setlength{\TaskBorderWidth}{ {{.Cfg.Layout.LaTeX.TaskBorderWidth}} }
\newlength{\TaskPaddingH}
\setlength{\TaskPaddingH}{ {{.Cfg.Layout.LaTeX.TaskPaddingH}} }
\newlength{\TaskPaddingV}
\setlength{\TaskPaddingV}{ {{.Cfg.Layout.LaTeX.TaskPaddingV}} }
% * Global vertical nudge for task elements (push tasks slightly lower)
\newlength{\TaskVerticalOffset}
\setlength{\TaskVerticalOffset}{ {{.Cfg.Layout.LaTeX.TaskVerticalOffset}} }

% Array stretch macro
\newcommand{\myNumArrayStretch}{ {{.Cfg.Layout.LaTeX.ArrayStretch}} }

% Line thickness macro
\newcommand{\myLineThick}{\rule{\linewidth}{\myLenLineThicknessThick}}

% Category palette setup macro
\newcommand{\SetupDefaultCategoryPalette}[1]{#1}

% Simple task rendering
\newcommand{\SimpleTaskBar}[4]{%
  \vspace*{\TaskVerticalOffset}%
  \fbox{\parbox{\dimexpr#3-2\TaskPaddingH\relax}{%
    \vspace{\TaskPaddingV}%
    {\TaskTitleSize\raggedright\textbf{#1}\par}%
    \vspace{\TaskPaddingV}%
  }}%
}

% Task overlay box macros - pill shaped with rounded corners
\newcommand{\TaskOverlayBox}[3]{%
  \definecolor{taskbgcolor}{RGB}{#1}%
  \definecolor{taskfgcolor}{RGB}{#1}%
  \vspace*{\TaskVerticalOffset}%
  \begin{tcolorbox}[enhanced, boxrule={{.Cfg.Layout.LaTeX.TColorBox.BoxRule}}, arc={{.Cfg.Layout.LaTeX.TColorBox.Arc}},
    left={{.Cfg.Layout.LaTeX.TColorBox.Left}}, right={{.Cfg.Layout.LaTeX.TColorBox.Right}}, top={{.Cfg.Layout.LaTeX.TColorBox.Top}}, bottom={{.Cfg.Layout.LaTeX.TColorBox.Bottom}},
    colback=taskbgcolor!{{.Cfg.Layout.LaTeX.TaskBackgroundOpacity}}, colframe=taskfgcolor!{{.Cfg.Layout.LaTeX.TaskBorderOpacity}},
    width=\linewidth, halign=left, before skip=0pt, after skip=0pt]
    {\sloppy\hyphenpenalty={{.Cfg.Layout.LaTeX.Typography.HyphenPenalty}}\tolerance={{.Cfg.Layout.LaTeX.Typography.Tolerance}}\emergencystretch={{.Cfg.Layout.LaTeX.Typography.EmergencyStretch}}%
     \TaskTitleSize\textbf{#2}\par
     \vspace{ {{.Cfg.Layout.LaTeX.Spacing.TaskContentVspace}} }%
     {\TaskFontSize\raggedright #3\par}}%
  \end{tcolorbox}%
}

% Task overlay box without vertical offset - for stacked tasks that should touch
\newcommand{\TaskOverlayBoxNoOffset}[3]{%
  \definecolor{taskbgcolor}{RGB}{#1}%
  \definecolor{taskfgcolor}{RGB}{#1}%
  \begin{tcolorbox}[enhanced, boxrule={{.Cfg.Layout.LaTeX.TColorBox.BoxRule}}, arc={{.Cfg.Layout.LaTeX.TColorBox.Arc}},
    left={{.Cfg.Layout.LaTeX.TColorBox.Left}}, right={{.Cfg.Layout.LaTeX.TColorBox.Right}}, top=0pt, bottom=0pt,
    colback=taskbgcolor!{{.Cfg.Layout.LaTeX.TaskBackgroundOpacity}}, colframe=taskfgcolor!{{.Cfg.Layout.LaTeX.TaskBorderOpacity}},
    width=\linewidth, halign=left, before skip=0pt, after skip=0pt, overlay={
      \begin{tikzpicture}[remember picture, overlay]
        \fill[taskbgcolor!{{.Cfg.Layout.LaTeX.TaskBackgroundOpacity}}] (0,0) rectangle (\linewidth, \baselineskip);
      \end{tikzpicture}
    }]
    {\sloppy\hyphenpenalty={{.Cfg.Layout.LaTeX.Typography.HyphenPenalty}}\tolerance={{.Cfg.Layout.LaTeX.Typography.Tolerance}}\emergencystretch={{.Cfg.Layout.LaTeX.Typography.EmergencyStretch}}%
     \TaskTitleSize\textbf{#2}\par
     \vspace{ {{.Cfg.Layout.LaTeX.Spacing.TaskContentVspace}} }%
     {\TaskFontSize\raggedright #3\par}}%
  \end{tcolorbox}%
}

% Multi-day task bar drawing macro to centralize styling
% Args: 1=x(pt), 2=y(pt), 3=width(pt), 4=height(pt), 5=color, 6=label
\newcommand{\DrawTaskBar}[6]{%
  \definecolor{taskbarcolor}{RGB}{#5}%
  \begin{tikzpicture}[overlay]
    \node[anchor=north west, inner sep=0pt] at (#1,#2) {
      \begin{tcolorbox}[enhanced, boxrule=0pt, arc={ {{.Cfg.Layout.LaTeX.Spacing.TaskOverlayArc}} },
        left={{.Cfg.Layout.LaTeX.TColorBox.Left}}, right={{.Cfg.Layout.LaTeX.TColorBox.Right}}, top={{.Cfg.Layout.LaTeX.TColorBox.Top}}, bottom={{.Cfg.Layout.LaTeX.TColorBox.Bottom}},
        width=#3pt, height=#4pt,
        colback=taskbarcolor]
        {\sloppy\hyphenpenalty={{.Cfg.Layout.LaTeX.Typography.HyphenPenalty}}\tolerance={{.Cfg.Layout.LaTeX.Typography.Tolerance}}\emergencystretch={{.Cfg.Layout.LaTeX.Typography.EmergencyStretch}}%
         \footnotesize \raggedright #6}
      \end{tcolorbox}
    };
  \end{tikzpicture}%
}

\newcommand{\TaskOverlayBoxP}[3]{%
  \definecolor{taskoverlaypbgcolor}{RGB}{#2}%
  \definecolor{taskoverlaypfgcolor}{RGB}{#2}%
  \vspace*{\TaskVerticalOffset}%
  \begin{tcolorbox}[enhanced, boxrule={{.Cfg.Layout.LaTeX.TColorBox.BoxRule}}, arc={{.Cfg.Layout.LaTeX.TColorBox.Arc}},
    left={{.Cfg.Layout.LaTeX.TColorBox.Left}}, right={{.Cfg.Layout.LaTeX.TColorBox.Right}}, top={{.Cfg.Layout.LaTeX.TColorBox.Top}}, bottom={{.Cfg.Layout.LaTeX.TColorBox.Bottom}},
    colback=taskoverlaypbgcolor!{{.Cfg.Layout.LaTeX.TaskBackgroundOpacity}}, colframe=taskoverlaypfgcolor!{{.Cfg.Layout.LaTeX.TaskBorderOpacity}},
    width=\linewidth, halign=left]
    {\sloppy\hyphenpenalty={{.Cfg.Layout.LaTeX.Typography.HyphenPenalty}}\tolerance={{.Cfg.Layout.LaTeX.Typography.Tolerance}}\emergencystretch={{.Cfg.Layout.LaTeX.Typography.EmergencyStretch}}%
     \TaskTitleSize\textbf{#1}\par
     \vspace{ {{.Cfg.Layout.LaTeX.Spacing.TaskContentVspace}} }%
     {\TaskFontSize\raggedright #3\par}}%
  \end{tcolorbox}%
}

% Task compact box macro with pill shape and better spacing
\newcommand{\TaskCompactBox}[4]{%
  \definecolor{taskcompactbgcolor}{RGB}{#3}%
  \definecolor{taskcompactfgcolor}{RGB}{#3}%
  \vspace*{#1}%
  \vspace*{\TaskVerticalOffset}%
  \begin{tcolorbox}[enhanced, boxrule={{.Cfg.Layout.LaTeX.TColorBox.TaskBoxRule}}, arc={{.Cfg.Layout.LaTeX.TColorBox.TaskArc}},
    left={{.Cfg.Layout.LaTeX.TColorBox.TaskLeft}}, right={{.Cfg.Layout.LaTeX.TColorBox.TaskRight}}, top={{.Cfg.Layout.LaTeX.TColorBox.TaskTop}}, bottom={{.Cfg.Layout.LaTeX.TColorBox.TaskBottom}},
    colback=taskcompactbgcolor!{{.Cfg.Layout.LaTeX.TaskBackgroundOpacity}}, colframe=taskcompactfgcolor!{{.Cfg.Layout.LaTeX.TaskBorderOpacity}},
    width=\linewidth, halign=left, height=#2]
    {\sloppy\hyphenpenalty={{.Cfg.Layout.LaTeX.Typography.HyphenPenalty}}\tolerance={{.Cfg.Layout.LaTeX.Typography.Tolerance}}\emergencystretch={{.Cfg.Layout.LaTeX.Typography.EmergencyStretch}}%
     \TaskTitleSize\raggedright\textbf{#4}\par}
  \end{tcolorbox}%
}

% Underline macro
\newcommand{\myUnderline}[1]{%
  \underline{\textbf{#1}}%
}

% Colored circle macro for legend - handles hex colors
\newcommand{\ColorCircle}[1]{%
  \definecolor{circlecolor}{RGB}{#1}%
  \textcolor{circlecolor}{\Large$\bullet$}%
}



% Color legend macro for task categories - uses algorithmic colors
\newcommand{\ColorLegend}{%
  {\centering
    \textcolor[RGB]{ {{- .Cfg.Layout.AlgorithmicColors.Proposal -}} }{\Large$\bullet$}~\small{Proposal}%
    \quad
    \textcolor[RGB]{ {{- .Cfg.Layout.AlgorithmicColors.Laser -}} }{\Large$\bullet$}~\small{Laser}%
    \quad
    \textcolor[RGB]{ {{- .Cfg.Layout.AlgorithmicColors.Imaging -}} }{\Large$\bullet$}~\small{Imaging}%
    \quad
    \textcolor[RGB]{ {{- .Cfg.Layout.AlgorithmicColors.Admin -}} }{\Large$\bullet$}~\small{Admin}%
    \quad
    \textcolor[RGB]{ {{- .Cfg.Layout.AlgorithmicColors.Dissertation -}} }{\Large$\bullet$}~\small{Dissertation}%
    \quad
    \textcolor[RGB]{ {{- .Cfg.Layout.AlgorithmicColors.Research -}} }{\Large$\bullet$}~\small{Research}%
    \quad
    \textcolor[RGB]{ {{- .Cfg.Layout.AlgorithmicColors.Publication -}} }{\Large$\bullet$}~\small{Publication}%
  \par}
}
