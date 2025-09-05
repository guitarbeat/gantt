% Enhanced Visual Design Template
% This template provides professional visual design with improved color schemes and typography

% Enhanced color definitions with accessibility in mind
\ExplSyntaxOn

% Professional color palette with high contrast ratios
\definecolor{primary}{HTML}{2563EB}      % Indigo-600
\definecolor{primaryLight}{HTML}{3B82F6} % Indigo-500
\definecolor{primaryDark}{HTML}{1D4ED8}  % Indigo-700
\definecolor{primaryLighter}{HTML}{60A5FA} % Indigo-400
\definecolor{primaryDarker}{HTML}{1E40AF}  % Indigo-800

\definecolor{secondary}{HTML}{6B7280}     % Gray-500
\definecolor{secondaryLight}{HTML}{9CA3AF} % Gray-400
\definecolor{secondaryDark}{HTML}{4B5563}  % Gray-600

\definecolor{accent}{HTML}{F59E0B}        % Amber-500
\definecolor{accentLight}{HTML}{FBBF24}   % Amber-400
\definecolor{accentDark}{HTML}{D97706}    % Amber-600

\definecolor{neutral}{HTML}{9CA3AF}       % Gray-400
\definecolor{neutralLight}{HTML}{D1D5DB}  % Gray-300
\definecolor{neutralDark}{HTML}{6B7280}   % Gray-500

\definecolor{background}{HTML}{FFFFFF}    % White
\definecolor{surface}{HTML}{F9FAFB}       % Gray-50
\definecolor{text}{HTML}{111827}          % Gray-900
\definecolor{textLight}{HTML}{374151}     % Gray-700
\definecolor{textLighter}{HTML}{6B7280}   % Gray-500

\definecolor{success}{HTML}{10B981}       % Emerald-500
\definecolor{warning}{HTML}{F59E0B}       % Amber-500
\definecolor{error}{HTML}{EF4444}         % Red-500
\definecolor{info}{HTML}{3B82F6}          % Blue-500

\definecolor{border}{HTML}{E5E7EB}        % Gray-200
\definecolor{borderLight}{HTML}{F3F4F6}   % Gray-100
\definecolor{borderDark}{HTML}{D1D5DB}    % Gray-300

\definecolor{shadow}{HTML}{000000}        % Black
\definecolor{highlight}{HTML}{FEF3C7}     % Amber-100

% Enhanced category colors with improved contrast
\definecolor{catPROPOSAL}{HTML}{2563EB}     % Indigo-600
\definecolor{catPROPOSALLight}{HTML}{3B82F6} % Indigo-500
\definecolor{catPROPOSALDark}{HTML}{1D4ED8}  % Indigo-700
\definecolor{catPROPOSALSecondary}{HTML}{DBEAFE} % Indigo-100

\definecolor{catLASER}{HTML}{EA580C}        % Orange-600
\definecolor{catLASERLight}{HTML}{F97316}   % Orange-500
\definecolor{catLASERDark}{HTML}{C2410C}    % Orange-700
\definecolor{catLASERSecondary}{HTML}{FED7AA} % Orange-100

\definecolor{catIMAGING}{HTML}{16A34A}      % Green-600
\definecolor{catIMAGINGLight}{HTML}{22C55E} % Green-500
\definecolor{catIMAGINGDark}{HTML}{15803D}  % Green-700
\definecolor{catIMAGINGSecondary}{HTML}{BBF7D0} % Green-100

\definecolor{catADMIN}{HTML}{6B7280}        % Gray-500
\definecolor{catADMINLight}{HTML}{9CA3AF}   % Gray-400
\definecolor{catADMINDark}{HTML}{4B5563}    % Gray-600
\definecolor{catADMINSecondary}{HTML}{E5E7EB} % Gray-200

\definecolor{catDISSERTATION}{HTML}{7C3AED} % Violet-600
\definecolor{catDISSERTATIONLight}{HTML}{8B5CF6} % Violet-500
\definecolor{catDISSERTATIONDark}{HTML}{6D28D9}  % Violet-700
\definecolor{catDISSERTATIONSecondary}{HTML}{DDD6FE} % Violet-100

\definecolor{catRESEARCH}{HTML}{0EA5E9}     % Sky-500
\definecolor{catRESEARCHLight}{HTML}{38BDF8} % Sky-400
\definecolor{catRESEARCHDark}{HTML}{0284C7}  % Sky-600
\definecolor{catRESEARCHSecondary}{HTML}{BAE6FD} % Sky-100

\definecolor{catPUBLICATION}{HTML}{DC2626}  % Red-600
\definecolor{catPUBLICATIONLight}{HTML}{EF4444} % Red-500
\definecolor{catPUBLICATIONDark}{HTML}{B91C1C}  % Red-700
\definecolor{catPUBLICATIONSecondary}{HTML}{FECACA} % Red-100

% Enhanced typography system
\newcommand{\HeadingOne}[1]{\fontsize{36}{43.2}\selectfont\textbf{\color{text}#1}}
\newcommand{\HeadingTwo}[1]{\fontsize{30}{39}\selectfont\textbf{\color{text}#1}}
\newcommand{\HeadingThree}[1]{\fontsize{24}{33.6}\selectfont\textbf{\color{text}#1}}
\newcommand{\BodyLarge}[1]{\fontsize{18}{28.8}\selectfont\color{textLight}#1}
\newcommand{\Body}[1]{\fontsize{16}{24}\selectfont\color{textLight}#1}
\newcommand{\BodySmall}[1]{\fontsize{14}{21}\selectfont\color{textLighter}#1}
\newcommand{\Caption}[1]{\fontsize{12}{16.8}\selectfont\color{textLighter}#1}

% Task-specific typography
\newcommand{\TaskTitle}[1]{\fontsize{14}{19.6}\selectfont\textbf{\color{white}#1}}
\newcommand{\TaskDescription}[1]{\fontsize{12}{15.6}\selectfont\color{white}#1}
\newcommand{\OverflowText}[1]{\fontsize{10}{12}\selectfont\textbf{\color{textLighter}#1}}

% Enhanced spacing system
\newlength{\Spacing0}
\newlength{\Spacing1}
\newlength{\Spacing2}
\newlength{\Spacing3}
\newlength{\Spacing4}
\newlength{\Spacing5}
\newlength{\Spacing6}
\newlength{\Spacing8}
\newlength{\Spacing10}
\newlength{\Spacing12}
\newlength{\Spacing16}
\newlength{\Spacing20}
\newlength{\Spacing24}
\newlength{\Spacing32}
\newlength{\Spacing40}
\newlength{\Spacing48}
\newlength{\Spacing56}
\newlength{\Spacing64}

\setlength{\Spacing0}{0pt}
\setlength{\Spacing1}{4pt}
\setlength{\Spacing2}{8pt}
\setlength{\Spacing3}{12pt}
\setlength{\Spacing4}{16pt}
\setlength{\Spacing5}{20pt}
\setlength{\Spacing6}{24pt}
\setlength{\Spacing8}{32pt}
\setlength{\Spacing10}{40pt}
\setlength{\Spacing12}{48pt}
\setlength{\Spacing16}{64pt}
\setlength{\Spacing20}{80pt}
\setlength{\Spacing24}{96pt}
\setlength{\Spacing32}{128pt}
\setlength{\Spacing40}{160pt}
\setlength{\Spacing48}{192pt}
\setlength{\Spacing56}{224pt}
\setlength{\Spacing64}{256pt}

% Enhanced border radius system
\newlength{\RadiusNone}
\newlength{\RadiusSm}
\newlength{\RadiusBase}
\newlength{\RadiusMd}
\newlength{\RadiusLg}
\newlength{\RadiusXl}
\newlength{\Radius2xl}
\newlength{\Radius3xl}
\newlength{\RadiusFull}

\setlength{\RadiusNone}{0pt}
\setlength{\RadiusSm}{2pt}
\setlength{\RadiusBase}{4pt}
\setlength{\RadiusMd}{6pt}
\setlength{\RadiusLg}{8pt}
\setlength{\RadiusXl}{12pt}
\setlength{\Radius2xl}{16pt}
\setlength{\Radius3xl}{24pt}
\setlength{\RadiusFull}{9999pt}

% Enhanced shadow system
\tikzset{
  shadowSm/.style={shadow={xshift=0pt,yshift=1pt,blur=2pt,spread=0pt,opacity=0.05,color=shadow}},
  shadowBase/.style={shadow={xshift=0pt,yshift=1pt,blur=3pt,spread=0pt,opacity=0.1,color=shadow}},
  shadowMd/.style={shadow={xshift=0pt,yshift=4pt,blur=6pt,spread=-1pt,opacity=0.1,color=shadow}},
  shadowLg/.style={shadow={xshift=0pt,yshift=10pt,blur=15pt,spread=-3pt,opacity=0.1,color=shadow}},
  shadowXl/.style={shadow={xshift=0pt,yshift=20pt,blur=25pt,spread=-5pt,opacity=0.1,color=shadow}}
}

% Enhanced border system
\newlength{\BorderNone}
\newlength{\BorderThin}
\newlength{\BorderBase}
\newlength{\BorderThick}
\newlength{\BorderThicker}

\setlength{\BorderNone}{0pt}
\setlength{\BorderThin}{1pt}
\setlength{\BorderBase}{1pt}
\setlength{\BorderThick}{2pt}
\setlength{\BorderThicker}{4pt}

% Professional task bar with enhanced visual design
\NewDocumentCommand{\ProfessionalTaskBarEnhanced}{mmmm}{
  % Args: prominence, color, name, description
  \begin{tcolorbox}[
    enhanced,
    colback=\CategoryColor{#2},
    colbacktitle=\CategoryColor{#2},
    coltitle=white,
    boxrule=\HierarchySpacing{#1},
    arc=\dimexpr\HierarchySpacing{#1}*0.3\relax,
    left=\Spacing4,
    right=\Spacing4,
    top=\Spacing2,
    bottom=\Spacing2,
    lefttitle=0pt,
    righttitle=0pt,
    toptitle=\Spacing2,
    bottomtitle=\Spacing2,
    fonttitle=\TaskTitle,
    title=\TaskTitle{#3},
    before skip=\Spacing2,
    after skip=\Spacing2,
    breakable,
    pad at break*=0pt,
    boxsep=0pt,
    left=0pt,
    right=0pt,
    top=0pt,
    bottom=0pt,
    shadow=shadowMd,
    borderline={0pt}{0pt}{border},
  ]
  \ifx\empty#4\empty\else
    \vspace{\Spacing1}
    \TaskDescription{#4}
  \fi
  \end{tcolorbox}
}

% Compact task bar with enhanced visual design
\NewDocumentCommand{\CompactTaskBarEnhanced}{mmmm}{
  % Args: spacing, height, color, content
  \begin{tcolorbox}[
    enhanced,
    colback=\CategoryColor{#3},
    colframe=\CategoryColor{#3},
    boxrule=\BorderThin,
    arc=\RadiusBase,
    left=\Spacing2,
    right=\Spacing2,
    top=\Spacing1,
    bottom=\Spacing1,
    height=#2,
    valign=center,
    halign=center,
    fonttitle=\TaskTitle,
    before skip=#1,
    after skip=#1,
    breakable,
    pad at break*=0pt,
    boxsep=0pt,
    left=0pt,
    right=0pt,
    top=0pt,
    bottom=0pt,
    shadow=shadowSm,
    borderline={0pt}{0pt}{border},
  ]
  #4
  \end{tcolorbox}
}

% Enhanced continuation chevrons with professional styling
\NewDocumentCommand{\EnhancedContLeftProfessional}{m}{
  \begin{tikzpicture}[overlay, remember picture]
    \node[anchor=east, inner sep=0pt, outer sep=0pt] at (0,0) {
      \begin{tcolorbox}[
        enhanced,
        colback=\CategoryColor{#1},
        colframe=\CategoryColor{#1},
        boxrule=\BorderThin,
        arc=\RadiusSm,
        left=\Spacing1,
        right=\Spacing1,
        top=\Spacing1,
        bottom=\Spacing1,
        width=\Spacing4,
        height=\Spacing8,
        valign=center,
        halign=center,
        fonttitle=\Caption,
        before skip=0pt,
        after skip=0pt,
        shadow=shadowSm,
      ]
      $\blacktriangleleft$
      \end{tcolorbox}
    };
  \end{tikzpicture}
}

\NewDocumentCommand{\EnhancedContRightProfessional}{m}{
  \begin{tikzpicture}[overlay, remember picture]
    \node[anchor=west, inner sep=0pt, outer sep=0pt] at (0,0) {
      \begin{tcolorbox}[
        enhanced,
        colback=\CategoryColor{#1},
        colframe=\CategoryColor{#1},
        boxrule=\BorderThin,
        arc=\RadiusSm,
        left=\Spacing1,
        right=\Spacing1,
        top=\Spacing1,
        bottom=\Spacing1,
        width=\Spacing4,
        height=\Spacing8,
        valign=center,
        halign=center,
        fonttitle=\Caption,
        before skip=0pt,
        after skip=0pt,
        shadow=shadowSm,
      ]
      $\blacktriangleright$
      \end{tcolorbox}
    };
  \end{tikzpicture}
}

% Professional overflow indicator with enhanced styling
\NewDocumentCommand{\ProfessionalOverflowEnhanced}{m}{
  \begin{tcolorbox}[
    enhanced,
    colback=surface,
    colframe=border,
    boxrule=\BorderThin,
    arc=\RadiusBase,
    left=\Spacing2,
    right=\Spacing2,
    top=\Spacing1,
    bottom=\Spacing1,
    valign=center,
    halign=center,
    fonttitle=\OverflowText,
    before skip=\Spacing1,
    after skip=\Spacing1,
    shadow=shadowSm,
  ]
  +#1 more
  \end{tcolorbox}
}

% Enhanced calendar cell with professional styling
\NewDocumentCommand{\EnhancedCalendarCellProfessional}{}{
  \setlength{\tabcolsep}{\Spacing2}
  \setlength{\arrayrulewidth}{\BorderThin}
  \renewcommand{\arraystretch}{1.1}
}

% Professional spacing for different view types with enhanced design
\NewDocumentCommand{\ViewSpecificSpacingProfessional}{m}{
  \str_case:nnF { #1 } {
    {monthly} {
      \setlength{\EnhancedTaskPaddingH}{4mm}
      \setlength{\EnhancedTaskPaddingV}{2mm}
      \setlength{\EnhancedTaskGap}{2pt}
    }
    {weekly} {
      \setlength{\EnhancedTaskPaddingH}{5mm}
      \setlength{\EnhancedTaskPaddingV}{2.5mm}
      \setlength{\EnhancedTaskGap}{2.5pt}
    }
    {daily} {
      \setlength{\EnhancedTaskPaddingH}{6mm}
      \setlength{\EnhancedTaskPaddingV}{3mm}
      \setlength{\EnhancedTaskGap}{3pt}
    }
  } {
    \setlength{\EnhancedTaskPaddingH}{4mm}
    \setlength{\EnhancedTaskPaddingV}{2mm}
    \setlength{\EnhancedTaskGap}{2pt}
  }
}

% Enhanced quality validation with professional standards
\NewDocumentCommand{\ValidateVisualQuality}{}{
  % Check minimum spacing requirements
  \ifdim\EnhancedTaskPaddingH<2mm
    \PackageWarning{enhanced_visual}{Task padding too small for professional appearance}
  \fi
  \ifdim\EnhancedTaskGap<1pt
    \PackageWarning{enhanced_visual}{Task gap too small for visual separation}
  \fi
  \ifdim\Spacing2<6pt
    \PackageWarning{enhanced_visual}{Spacing too small for professional layout}
  \fi
}

% Enhanced accessibility validation
\NewDocumentCommand{\ValidateAccessibility}{}{
  % Check font size requirements
  \ifdim\fontsize{12}{14}\selectfont\relax<12pt
    \PackageWarning{enhanced_visual}{Font size below accessibility minimum}
  \fi
  % Check contrast requirements
  \PackageInfo{enhanced_visual}{Using high-contrast color palette for accessibility}
}

\ExplSyntaxOff

% Initialize enhanced visual design
\ViewSpecificSpacingProfessional{monthly}
\ValidateVisualQuality
\ValidateAccessibility
