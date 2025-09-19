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
  \begin{tcolorbox}[enhanced, boxrule=0.8pt, arc=6pt, drop shadow={0.2pt}{-0.2pt}{0pt}{black!20},
    left=2.2mm, right=2.2mm, top=1.4mm, bottom=1.4mm,
    colback=#1!20, colframe=#1!80,
    width=\linewidth, halign=center]
    \TaskFontSize\textbf{#2}\\#3%
  \end{tcolorbox}%
}

\newcommand{\TaskOverlayBoxP}[3]{%
  \vspace*{\TaskVerticalOffset}%
  \begin{tcolorbox}[enhanced, boxrule=0.8pt, arc=6pt, drop shadow={0.2pt}{-0.2pt}{0pt}{black!20},
    left=2.4mm, right=2.4mm, top=1.4mm, bottom=1.4mm,
    colback=#2!20, colframe=#2!80,
    width=\linewidth, halign=center]
    \TaskFontSize\textbf{#1}\\#3%
  \end{tcolorbox}%
}

% Task compact box macro with pill shape and better spacing
\newcommand{\TaskCompactBox}[4]{%
  \vspace*{#1}%
  \vspace*{\TaskVerticalOffset}%
  \begin{tcolorbox}[enhanced, boxrule=0.6pt, arc=5pt, drop shadow={0.1pt}{-0.1pt}{0pt}{black!15},
    left=1.6mm, right=1.6mm, top=1mm, bottom=1mm,
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
