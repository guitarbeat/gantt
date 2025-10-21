% Setup category palette for this month
\SetupDefaultCategoryPalette{}

{{- template "calendar.tpl" dict "Cfg" .Cfg "Body" .Body -}}

% Legend at bottom of page - just colors and categories
\vfill
{{- $phaseGroups := .Body.Month.GetTaskColorsByPhase -}}
{{- if $phaseGroups -}}
{\small{{- range $idx, $phase := $phaseGroups -}}
% Phase header with subtle background
{\colorbox[RGB]{245,245,245}{\makebox[\linewidth][l]{\textbf{ {{- $phase.PhaseName -}} }}}\\
\vspace{1pt}
{{- range $subIdx, $subPhase := $phase.SubPhases -}}\ColorCircle{ {{- $subPhase.Color -}} }{ {{- $subPhase.Name -}} }\quad{{- end -}}\\

{{- end -}}}%
}
{{- else -}}
% Fallback to simple legend if no phase data
{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
\noindent{\small
% Legend header
{\colorbox[RGB]{250,250,250]{\makebox[\linewidth][l]{\textbf{Task Categories}}}}\\
\vspace{1pt}
{{- range $color, $category := $taskColors -}}{{- if ne $color "" -}}\ColorCircle{ {{- $color -}} }{ {{- $category -}} }\quad{{- end -}}{{- end -}}}
}
{{- else -}}
% Ultimate fallback to algorithmic legend
\noindent{\small
{\colorbox[RGB]{248,248,248]{\makebox[\linewidth][l]{\textbf{Default Categories}}}}\\
\vspace{1pt}
\ColorLegend}
{{- end -}}
{{- end -}}
