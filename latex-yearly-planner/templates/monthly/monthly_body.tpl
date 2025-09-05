% Setup category palette for this month
\SetupDefaultCategoryPalette

{{- template "calendar_table.tpl" dict "Cfg" .Cfg "Body" .Body -}}

% Enhanced layout visualization if available
{{ if hasLayoutData .Body }}
\smallskip
% Layout Statistics Section
{{ $stats := getLayoutStats .Body }}
{{ if $stats }}
\parbox{\textwidth}{
  \myUnderline{Layout Statistics}
  \footnotesize{
    Total Tasks: {{ $stats.TotalTasks }} | 
    Processed Bars: {{ $stats.ProcessedBars }} | 
    Space Efficiency: {{ printf "%.1f" $stats.SpaceEfficiency }} | 
    Visual Quality: {{ printf "%.1f" $stats.VisualQuality }}
  }
}
\smallskip
{{ end }}

% Task Bars Visualization
{{ $taskBars := getTaskBars .Body }}
{{ if $taskBars }}
\parbox{\textwidth}{
  \myUnderline{Task Visualization}
  \begin{tikzpicture}[overlay, remember picture]
    {{ range $taskBars }}
    {{ formatTaskBar . }}
    {{ end }}
  \end{tikzpicture}
}
\smallskip
{{ end }}
{{ end }}

% Single full-width Notes area (replaces previous two-column layout)
\parbox{\textwidth}{
  \myUnderline{Notes}
  \vbox to \dimexpr\textheight-\pagetotal-\myLenLineHeightButLine\relax {%
    \leaders\vbox to \myLenLineHeightButLine{\vfil\hrule width \linewidth height \myLenLineThicknessDefault}\vskip \dimexpr\textheight-\pagetotal-\myLenLineHeightButLine\relax
  }%
}
