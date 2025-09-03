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

\NewDocumentCommand{\myDotGrid}{mm}{\leavevmode\multido{\dC=0mm+5mm}{#1}{\multido{\dR=0mm+5mm}{#2}{\put(\dR,\dC){\circle*{0.1}}}}}

\NewDocumentCommand{\myMash}{O{}mm}{
  {{- if $.Cfg.Dotted -}} \vskip\myLenLineHeightButLine#1\myDotGrid{#2}{#3} {{- else -}} \Repeat{#2}{\myLineGrayVskipTop} {{- end -}}
}

\NewDocumentCommand{\remainingHeight}{}{%
  \ifdim\pagegoal=\maxdimen
  \dimexpr\textheight-9.4pt\relax
  \else
  \dimexpr\pagegoal-\pagetotal-\lineskip-9.4pt\relax
  \fi%
}