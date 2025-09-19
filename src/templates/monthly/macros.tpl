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
\setlength{\myLenMonthlyCellHeight}{70pt}

\setlength{\myLenHeaderResizeBox}{6mm}
\setlength{\myLenHeaderSideMonthsWidth}{14.5cm}

% Simple task bar definitions
\newlength{\TaskBarHeight}
\setlength{\TaskBarHeight}{4mm}
\newlength{\TaskBorderWidth}
\setlength{\TaskBorderWidth}{0.6pt}
\newlength{\TaskPaddingH}
\setlength{\TaskPaddingH}{1.5mm}
\newlength{\TaskPaddingV}
\setlength{\TaskPaddingV}{0.5mm}

% Array stretch macro
\newcommand{\myNumArrayStretch}{1.2}

% Line thickness macro
\newcommand{\myLineThick}{\rule{\linewidth}{\myLenLineThicknessThick}}

% Category palette setup macro
\newcommand{\SetupDefaultCategoryPalette}[1]{#1}

% Simple task rendering
\newcommand{\SimpleTaskBar}[4]{%
  \fbox{\parbox{\dimexpr#3-2\TaskPaddingH\relax}{%
    \vspace{\TaskPaddingV}%
    \centering\small\textbf{#1}%
    \vspace{\TaskPaddingV}%
  }}%
}

% Task overlay box macros
\newcommand{\TaskOverlayBox}[3]{%
  \fcolorbox{#1}{#1!20}{\parbox{\linewidth}{%
    \centering\small\textbf{#2}\\#3%
  }}%
}

\newcommand{\TaskOverlayBoxP}[3]{%
  \fcolorbox{#2}{#2!20}{\parbox{\linewidth}{%
    \centering\small\textbf{#1}\\#3%
  }}%
}

% Task compact box macro
\newcommand{\TaskCompactBox}[4]{%
  \vspace*{#1}%
  \fcolorbox{#3}{#3!20}{\parbox{\linewidth}{%
    \centering\small\textbf{#4}%
  }}%
}

% Underline macro
\newcommand{\myUnderline}[1]{%
  \underline{\textbf{#1}}%
}
