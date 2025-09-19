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

% Task compact box macro with better spacing and height
\newcommand{\TaskCompactBox}[4]{%
  \vspace*{#1}%
  \fcolorbox{#3}{#3!20}{\parbox{\dimexpr\linewidth-2\fboxsep\relax}{%
    \vbox to #2{%
      \vfil
      \centering\small\textbf{#4}%
      \vfil
    }%
  }}%
  \vspace*{0.1ex}%
}

% Underline macro
\newcommand{\myUnderline}[1]{%
  \underline{\textbf{#1}}%
}

% Color legend macro for task categories
\newcommand{\ColorLegend}{%
  \vspace*{-0.5ex}%
  \begin{center}%
    \small\textbf{Task Categories:}%
    \hspace{1em}%
    \fcolorbox{#4A90E2}{#4A90E2}{\makebox[0.8em][c]{\textcolor{white}{\tiny P}}}%
    \hspace{0.3em}\tiny Proposal%
    \hspace{0.8em}%
    \fcolorbox{#F5A623}{#F5A623}{\makebox[0.8em][c]{\textcolor{white}{\tiny L}}}%
    \hspace{0.3em}\tiny Laser%
    \hspace{0.8em}%
    \fcolorbox{#7ED321}{#7ED321}{\makebox[0.8em][c]{\textcolor{white}{\tiny I}}}%
    \hspace{0.3em}\tiny Imaging%
    \hspace{0.8em}%
    \fcolorbox{#BD10E0}{#BD10E0}{\makebox[0.8em][c]{\textcolor{white}{\tiny A}}}%
    \hspace{0.3em}\tiny Admin%
    \hspace{0.8em}%
    \fcolorbox{#D0021B}{#D0021B}{\makebox[0.8em][c]{\textcolor{white}{\tiny D}}}%
    \hspace{0.3em}\tiny Dissertation%
    \hspace{0.8em}%
    \fcolorbox{#50E3C2}{#50E3C2}{\makebox[0.8em][c]{\textcolor{white}{\tiny R}}}%
    \hspace{0.3em}\tiny Research%
    \hspace{0.8em}%
    \fcolorbox{#B8E986}{#B8E986}{\makebox[0.8em][c]{\textcolor{white}{\tiny P}}}%
    \hspace{0.3em}\tiny Publication%
  \end{center}%
  \vspace*{0.3ex}%
}
