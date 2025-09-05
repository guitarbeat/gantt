% moved from templates/macro.tpl
\ExplSyntaxOn
\cs_new_eq:NN \Repeat \prg_replicate:nn
\ExplSyntaxOff

\NewDocumentCommand{\myMinLineHeight}{m}{\parbox{0pt}{\vskip#1}}
\NewDocumentCommand{\myDummyQ}{}{\textcolor{white}{Q}}

{{- $numbers := .Cfg.Layout.Numbers -}}
\NewDocumentCommand{\myNumArrayStretch}{}{ {{- $numbers.ArrayStretch -}} }

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
\setlength{\myLenTwoCol}{\dimexpr.5\linewidth-.5\myLenTwoColSep}
\setlength{\myLenFiveColSep}{ {{- $lengths.FiveColSep -}} }
\setlength{\myLenFiveCol}{\dimexpr.2\linewidth-\myLenFiveColSep}
\setlength{\myLenMonthlyCellHeight}{ {{- $lengths.MonthlyCellHeight -}} }
\setlength{\myLenTriColSep}{ {{- $lengths.TriColSep -}} }
\setlength{\myLenTriCol}{\dimexpr.333\linewidth-.667\myLenTriColSep}

\setlength{\myLenHeaderResizeBox}{ {{- $lengths.HeaderResizeBox -}} }
\setlength{\myLenHeaderSideMonthsWidth}{ {{- $lengths.HeaderSideMonthsWidth -}} }

\NewDocumentCommand{\myMonthlySpring}{}{ {{- $lengths.MonthlySpring -}} }
\NewDocumentCommand{\myColorGray}{}{ {{- .Cfg.Layout.Colors.Gray -}} }
\NewDocumentCommand{\myColorLightGray}{}{ {{- .Cfg.Layout.Colors.LightGray -}} }
\NewDocumentCommand{\myLinePlain}{}{\hrule width \linewidth height \myLenLineThicknessDefault}
\NewDocumentCommand{\myLineThick}{}{\hrule width \linewidth height \myLenLineThicknessThick}
\NewDocumentCommand{\myLineHeightButLine}{}{\myMinLineHeight{\myLenLineHeightButLine}}
\NewDocumentCommand{\myUnderline}{m}{#1\vskip1mm\myLineThick\par}
\NewDocumentCommand{\myLineColor}{m}{\textcolor{#1}{\myLinePlain}}
\NewDocumentCommand{\myLineGray}{}{\myLineColor{\myColorGray}}
\NewDocumentCommand{\myLineLightGray}{}{\myLineColor{\myColorLightGray}}
\NewDocumentCommand{\myLineGrayVskipBottom}{}{\myLineGray\vskip\myLenLineHeightButLine}
\NewDocumentCommand{\myLineGrayVskipTop}{}{\vskip\myLenLineHeightButLine\myLineGray}
\NewDocumentCommand{\myTodo}{}{\myLineHeightButLine$\square$\myLinePlain}
\NewDocumentCommand{\myTodoLineGray}{}{\myLineHeightButLine$\square$\myLineGray}
% Draw a dotted grid of size (#1 x #2) with 5mm spacing using LaTeX picture environment
% We wrap in a picture environment so that \put and \circle* are defined
\NewDocumentCommand{\myDotGrid}{mm}{%
  \leavevmode\begingroup
  \setlength{\unitlength}{1mm}% use millimeters for coordinates
  \begin{picture}(0,0)% zero-sized picture, we only place absolute dots
    \multido{\dC=0+5}{#1}{%
      \multido{\dR=0+5}{#2}{\put(\dR,\dC){\circle*{0.1}}}%
    }%
  \end{picture}%
  \endgroup
}
\NewDocumentCommand{\myMash}{O{}mm}{ {{- if $.Cfg.Dotted -}} \vskip\myLenLineHeightButLine#1\myDotGrid{#2}{#3} {{- else -}} \Repeat{#2}{\myLineGrayVskipTop} {{- end -}} }
\NewDocumentCommand{\remainingHeight}{}{
  \ifdim\pagegoal=\maxdimen
  \dimexpr\textheight-9.4pt\relax
  \else
  \dimexpr\pagegoal-\pagetotal-\lineskip-9.4pt\relax
  \fi
}

% * Calendar Task Rendering Macros (Google Calendar-style)
% ! Requires tcolorbox, tikz, xcolor

% Category color registry and palette
\ExplSyntaxOn
\prop_new:N \g_task_category_color_prop
\NewDocumentCommand{\DefineCategoryColor}{mm}{\prop_gput:Nnn \g_task_category_color_prop {#1} {#2}}
\NewDocumentCommand{\CategoryColor}{m}{\prop_get:NnN \g_task_category_color_prop {#1} \l_tmpa_tl \tl_if_blank:NF \l_tmpa_tl {\l_tmpa_tl}}
\ExplSyntaxOff

% Predefined accessible palette for categories
% Colors chosen for contrast on white with dark text
\providecolor{catProposal}{HTML}{2563EB}  % Indigo-600
\providecolor{catLaser}{HTML}{EA580C}     % Orange-600
\providecolor{catImaging}{HTML}{16A34A}   % Green-600
\providecolor{catAdmin}{HTML}{6B7280}     % Gray-500
\providecolor{catDissertation}{HTML}{7C3AED} % Violet-600
\providecolor{catResearch}{HTML}{0EA5E9}  % Sky-500
\providecolor{catPublication}{HTML}{DC2626} % Red-600

% Setup default category→color mapping
\NewDocumentCommand{\SetupDefaultCategoryPalette}{}{
  \DefineCategoryColor{PROPOSAL}{catProposal}
  \DefineCategoryColor{LASER}{catLaser}
  \DefineCategoryColor{IMAGING}{catImaging}
  \DefineCategoryColor{ADMIN}{catAdmin}
  \DefineCategoryColor{DISSERTATION}{catDissertation}
  \DefineCategoryColor{RESEARCH}{catResearch}
  \DefineCategoryColor{PUBLICATION}{catPublication}
}

% Variants
\NewDocumentCommand{\CategoryLight}{mO{20}}{\CategoryColor{#1}!#2}
\NewDocumentCommand{\CategoryDark}{mO{60}}{\CategoryColor{#1}!#2!black}

% Task text color chooser (placeholder: defaults to black for contrast)
% Usage: \TaskTextColor{<hex-or-xcolor>}{<text>}
\NewDocumentCommand{\TaskTextColor}{mm}{\textcolor{black}{#2}}

% Styling tokens
\newlength{\TaskBarCornerRadius}
\setlength{\TaskBarCornerRadius}{1pt}
\newlength{\TaskBorderWidth}
\setlength{\TaskBorderWidth}{0.6pt}
\newlength{\TaskPaddingH}
\setlength{\TaskPaddingH}{1.5mm}
\newlength{\TaskPaddingV}
\setlength{\TaskPaddingV}{0.5mm}

% Typography tokens
\NewDocumentCommand{\TaskTitleFont}{m}{{\centering\color{black}\textbf{\scriptsize #1}}}
\NewDocumentCommand{\TaskDescFont}{m}{{\color{black}\tiny #1}}
\NewDocumentCommand{\OverflowFont}{m}{{\centering\color{gray}\textbf{\tiny #1}}}

% Prominence levels: CRITICAL, HIGH, MEDIUM, LOW, MINIMAL
\ExplSyntaxOn
\prop_new:N \g_task_prominence_border_prop
\prop_gset_from_keyval:Nn \g_task_prominence_border_prop { CRITICAL = 2.0pt , HIGH = 1.4pt , MEDIUM = 1.0pt , LOW = 0.7pt , MINIMAL = 0.5pt }
\prop_new:N \g_task_prominence_fillL_prop
\prop_gset_from_keyval:Nn \g_task_prominence_fillL_prop { CRITICAL = 46 , HIGH = 40 , MEDIUM = 34 , LOW = 28 , MINIMAL = 22 }
\prop_new:N \g_task_prominence_fillR_prop
\prop_gset_from_keyval:Nn \g_task_prominence_fillR_prop { CRITICAL = 18 , HIGH = 12 , MEDIUM = 6 , LOW = 4 , MINIMAL = 2 }
\NewDocumentCommand{\PromBorder}{m}{\prop_item:Nn \g_task_prominence_border_prop {#1}}
\NewDocumentCommand{\PromFillLeft}{m}{\prop_item:Nn \g_task_prominence_fillL_prop {#1}}
\NewDocumentCommand{\PromFillRight}{m}{\prop_item:Nn \g_task_prominence_fillR_prop {#1}}
\ExplSyntaxOff

% Full overlay box used for a single-starting task with name/desc
% Args: 1=color, 2=name (already escaped), 3=desc (escaped, may be empty)
% Default prominence MEDIUM wrapper
\NewDocumentCommand{\TaskOverlayBox}{mmm}{\TaskOverlayBoxP{MEDIUM}{#1}{#2}{#3}}
% Prominence-aware overlay box
\NewDocumentCommand{\TaskOverlayBoxP}{mmmm}{%
  {\begingroup\setlength{\fboxsep}{0pt}%
    \begin{tcolorbox}[
      enhanced,
      boxrule=0pt,
      arc=\TaskBarCornerRadius,
      drop shadow,
      left=\TaskPaddingH, right=\TaskPaddingH, top=\TaskPaddingV, bottom=\TaskPaddingV,
      colback=#2!\PromFillLeft{#1},
      interior style={left color=#2!\PromFillLeft{#1}, right color=#2!\PromFillRight{#1}},
      borderline west={\PromBorder{#1}}{0pt}{#2!60!black},
      borderline east={\TaskBorderWidth}{0pt}{#2!45}
    ]
      {\hyphenpenalty=10000\exhyphenpenalty=10000\emergencystretch=2em\setstretch{0.75}%
        \TaskTitleFont{#3}% name line
        \\[-0.3ex]{\TaskDescFont{#4}}% desc line (prints nothing if empty)
      }
    \end{tcolorbox}\endgroup}
}

% Compact stacked bar used in multi-task overlay
% Args: 1=spacing, 2=height, 3=color, 4=inner text (tiny, centered)
% Default prominence MEDIUM wrapper
\NewDocumentCommand{\TaskCompactBox}{mmmm}{\TaskCompactBoxP{MEDIUM}{#1}{#2}{#3}{#4}}
% Prominence-aware compact bar
\NewDocumentCommand{\TaskCompactBoxP}{mmmmm}{%
  \vspace*{#2}{\begingroup\setlength{\fboxsep}{0pt}%
    \begin{tcolorbox}[
      enhanced,
      boxrule=0pt,
      arc=\TaskBarCornerRadius,
      left=1.0mm, right=1.0mm, top=0.2mm, bottom=0.2mm,
      height=#3,
      colback=#4!\PromFillLeft{#1},
      interior style={left color=#4!\PromFillLeft{#1}, right color=#4!\PromFillRight{#1}},
      borderline west={\PromBorder{#1}}{0pt}{#4!50!black}
    ]
      #5
    \end{tcolorbox}\endgroup}
}

% Continuation chevrons (left/right) for boundary indicators
% Args: 1=color
\NewDocumentCommand{\TaskContLeft}{m}{\textcolor{#1!60!black}{\small\ensuremath{\langle}}}
\NewDocumentCommand{\TaskContRight}{m}{\textcolor{#1!60!black}{\small\ensuremath{\rangle}}}

% Overflow indicator helper
\NewDocumentCommand{\TaskOverflow}{m}{\OverflowFont{#1}}

