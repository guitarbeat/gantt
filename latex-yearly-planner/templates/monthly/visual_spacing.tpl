% Visual Spacing Optimization Template
% This template provides enhanced spacing and alignment for professional appearance

% Enhanced spacing macros with professional quality
\ExplSyntaxOn

% Spacing calculation functions
\cs_new:Nn \calculate_optimal_spacing:nnn {
  % Args: base_spacing, density_multiplier, content_multiplier
  \fp_set:Nn \l_tmpa_fp { #1 }
  \fp_mul:Nn \l_tmpa_fp { #2 }
  \fp_mul:Nn \l_tmpa_fp { #3 }
  \fp_use:N \l_tmpa_fp
}

% Responsive spacing based on content density
\NewDocumentCommand{\ResponsiveSpacing}{mmm}{
  % Args: base_value, density_level, content_length
  \fp_set:Nn \l_tmpa_fp { #1 }
  \str_case:nnF { #2 } {
    {low} { \fp_mul:Nn \l_tmpa_fp { 1.2 } }
    {normal} { \fp_mul:Nn \l_tmpa_fp { 1.0 } }
    {high} { \fp_mul:Nn \l_tmpa_fp { 0.8 } }
    {very_high} { \fp_mul:Nn \l_tmpa_fp { 0.7 } }
  } { \fp_mul:Nn \l_tmpa_fp { 1.0 } }
  
  % Adjust for content length
  \int_compare:nNnT { #3 } > { 30 } {
    \fp_mul:Nn \l_tmpa_fp { 1.1 }
  }
  \int_compare:nNnT { #3 } < { 15 } {
    \fp_mul:Nn \l_tmpa_fp { 0.9 }
  }
  
  \fp_use:N \l_tmpa_fp
}

% Professional spacing tokens
\newlength{\ProfessionalPadding}
\newlength{\ProfessionalMargin}
\newlength{\ProfessionalGap}
\newlength{\ProfessionalMinSpacing}
\newlength{\ProfessionalMaxSpacing}

% Set professional spacing defaults
\setlength{\ProfessionalPadding}{1.5pt}
\setlength{\ProfessionalMargin}{0.8pt}
\setlength{\ProfessionalGap}{1.2pt}
\setlength{\ProfessionalMinSpacing}{0.5pt}
\setlength{\ProfessionalMaxSpacing}{3.0pt}

% Enhanced task bar spacing
\newlength{\EnhancedTaskPaddingH}
\newlength{\EnhancedTaskPaddingV}
\newlength{\EnhancedTaskGap}
\newlength{\EnhancedTaskMargin}

\setlength{\EnhancedTaskPaddingH}{1.8mm}
\setlength{\EnhancedTaskPaddingV}{0.6mm}
\setlength{\EnhancedTaskGap}{0.8pt}
\setlength{\EnhancedTaskMargin}{0.4pt}

% Typography spacing
\newlength{\EnhancedTitleSpacing}
\newlength{\EnhancedDescSpacing}
\newlength{\EnhancedOverflowSpacing}

\setlength{\EnhancedTitleSpacing}{1.2pt}
\setlength{\EnhancedDescSpacing}{0.8pt}
\setlength{\EnhancedOverflowSpacing}{0.6pt}

% Hierarchy-based spacing
\NewDocumentCommand{\HierarchySpacing}{m}{
  \str_case:nnF { #1 } {
    {CRITICAL} { 2.5pt }
    {HIGH} { 2.0pt }
    {MEDIUM} { 1.5pt }
    {LOW} { 1.0pt }
    {MINIMAL} { 0.8pt }
  } { 1.5pt }
}

% Category-based spacing
\NewDocumentCommand{\CategorySpacing}{m}{
  \str_case:nnF { #1 } {
    {PROPOSAL} { 2.2pt }
    {LASER} { 2.0pt }
    {IMAGING} { 1.8pt }
    {ADMIN} { 1.5pt }
    {DISSERTATION} { 2.5pt }
    {RESEARCH} { 2.0pt }
    {PUBLICATION} { 2.2pt }
  } { 1.5pt }
}

% Enhanced alignment macros
\NewDocumentCommand{\EnhancedAlign}{mm}{
  % Args: horizontal, vertical
  \str_case:nnF { #1 } {
    {left} { \raggedright }
    {center} { \centering }
    {right} { \raggedleft }
    {justify} { \justifying }
  } { \centering }
  
  \str_case:nnF { #2 } {
    {top} { \vspace{0pt} }
    {middle} { \vspace{\dimexpr0.5\baselineskip-0.5\height\relax} }
    {bottom} { \vspace{\dimexpr\baselineskip-\height\relax} }
    {baseline} { \vspace{0pt} }
  } { \vspace{0pt} }
}

% Professional task bar with enhanced spacing
\NewDocumentCommand{\ProfessionalTaskBar}{mmmm}{
  % Args: prominence, color, name, description
  \begin{tcolorbox}[
    enhanced,
    colback=\CategoryColor{#2},
    colbacktitle=\CategoryColor{#2},
    coltitle=white,
    boxrule=\HierarchySpacing{#1},
    arc=\dimexpr\HierarchySpacing{#1}*0.3\relax,
    left=\EnhancedTaskPaddingH,
    right=\EnhancedTaskPaddingH,
    top=\EnhancedTaskPaddingV,
    bottom=\EnhancedTaskPaddingV,
    lefttitle=0pt,
    righttitle=0pt,
    toptitle=\EnhancedTaskPaddingV,
    bottomtitle=\EnhancedTaskPaddingV,
    fonttitle=\small\bfseries,
    title=\EnhancedAlign{left}{middle}{#3},
    before skip=\EnhancedTaskGap,
    after skip=\EnhancedTaskGap,
    breakable,
    pad at break*=0pt,
    boxsep=0pt,
    left=0pt,
    right=0pt,
    top=0pt,
    bottom=0pt,
  ]
  \ifx\empty#4\empty\else
    \vspace{\EnhancedDescSpacing}
    \tiny #4
  \fi
  \end{tcolorbox}
}

% Compact task bar with optimized spacing
\NewDocumentCommand{\CompactTaskBar}{mmmm}{
  % Args: spacing, height, color, content
  \begin{tcolorbox}[
    enhanced,
    colback=\CategoryColor{#3},
    colframe=\CategoryColor{#3},
    boxrule=0.5pt,
    arc=1pt,
    left=0.5mm,
    right=0.5mm,
    top=0.3mm,
    bottom=0.3mm,
    height=#2,
    valign=center,
    halign=center,
    fonttitle=\tiny\bfseries,
    before skip=#1,
    after skip=#1,
    breakable,
    pad at break*=0pt,
    boxsep=0pt,
    left=0pt,
    right=0pt,
    top=0pt,
    bottom=0pt,
  ]
  #4
  \end{tcolorbox}
}

% Enhanced continuation chevrons with proper spacing
\NewDocumentCommand{\EnhancedContLeft}{m}{
  \begin{tikzpicture}[overlay, remember picture]
    \node[anchor=east, inner sep=0pt, outer sep=0pt] at (0,0) {
      \begin{tcolorbox}[
        enhanced,
        colback=\CategoryColor{#1},
        colframe=\CategoryColor{#1},
        boxrule=0.5pt,
        arc=0.5pt,
        left=0.3mm,
        right=0.3mm,
        top=0.2mm,
        bottom=0.2mm,
        width=2mm,
        height=4mm,
        valign=center,
        halign=center,
        fonttitle=\tiny,
        before skip=0pt,
        after skip=0pt,
      ]
      $\blacktriangleleft$
      \end{tcolorbox}
    };
  \end{tikzpicture}
}

\NewDocumentCommand{\EnhancedContRight}{m}{
  \begin{tikzpicture}[overlay, remember picture]
    \node[anchor=west, inner sep=0pt, outer sep=0pt] at (0,0) {
      \begin{tcolorbox}[
        enhanced,
        colback=\CategoryColor{#1},
        colframe=\CategoryColor{#1},
        boxrule=0.5pt,
        arc=0.5pt,
        left=0.3mm,
        right=0.3mm,
        top=0.2mm,
        bottom=0.2mm,
        width=2mm,
        height=4mm,
        valign=center,
        halign=center,
        fonttitle=\tiny,
        before skip=0pt,
        after skip=0pt,
      ]
      $\blacktriangleright$
      \end{tcolorbox}
    };
  \end{tikzpicture}
}

% Professional overflow indicator
\NewDocumentCommand{\ProfessionalOverflow}{m}{
  \begin{tcolorbox}[
    enhanced,
    colback=gray!20,
    colframe=gray!40,
    boxrule=0.3pt,
    arc=0.5pt,
    left=0.5mm,
    right=0.5mm,
    top=0.2mm,
    bottom=0.2mm,
    valign=center,
    halign=center,
    fonttitle=\tiny\bfseries,
    before skip=\EnhancedOverflowSpacing,
    after skip=\EnhancedOverflowSpacing,
  ]
  +#1 more
  \end{tcolorbox}
}

% Enhanced calendar cell spacing
\NewDocumentCommand{\EnhancedCalendarCell}{}{
  \setlength{\tabcolsep}{\EnhancedTaskGap}
  \setlength{\arrayrulewidth}{0.4pt}
  \renewcommand{\arraystretch}{1.1}
}

% Professional spacing for different view types
\NewDocumentCommand{\ViewSpecificSpacing}{m}{
  \str_case:nnF { #1 } {
    {monthly} {
      \setlength{\EnhancedTaskPaddingH}{1.8mm}
      \setlength{\EnhancedTaskPaddingV}{0.6mm}
      \setlength{\EnhancedTaskGap}{0.8pt}
    }
    {weekly} {
      \setlength{\EnhancedTaskPaddingH}{2.2mm}
      \setlength{\EnhancedTaskPaddingV}{0.8mm}
      \setlength{\EnhancedTaskGap}{1.0pt}
    }
    {daily} {
      \setlength{\EnhancedTaskPaddingH}{2.5mm}
      \setlength{\EnhancedTaskPaddingV}{1.0mm}
      \setlength{\EnhancedTaskGap}{1.2pt}
    }
  } {
    \setlength{\EnhancedTaskPaddingH}{1.8mm}
    \setlength{\EnhancedTaskPaddingV}{0.6mm}
    \setlength{\EnhancedTaskGap}{0.8pt}
  }
}

% Quality validation spacing
\NewDocumentCommand{\ValidateSpacing}{}{
  % Check minimum spacing requirements
  \ifdim\EnhancedTaskPaddingH<0.5mm
    \PackageWarning{visual_spacing}{Task padding too small for readability}
  \fi
  \ifdim\EnhancedTaskGap<0.3pt
    \PackageWarning{visual_spacing}{Task gap too small for visual separation}
  \fi
  \ifdim\EnhancedTitleSpacing<0.5pt
    \PackageWarning{visual_spacing}{Title spacing too small for readability}
  \fi
}

\ExplSyntaxOff

% Initialize professional spacing
\ViewSpecificSpacing{monthly}
\ValidateSpacing
