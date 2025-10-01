% Setup category palette for this month
\SetupDefaultCategoryPalette{}

{{- template "calendar.tpl" dict "Cfg" .Cfg "Body" .Body -}}

% Legend at bottom of page - grouped by phase
\vfill
{{- $phaseGroups := .Body.Month.GetTaskColorsByPhase -}}
{{- if $phaseGroups -}}
{\small
{{- range $idx, $phase := $phaseGroups -}}
{{- if $idx -}}\vspace{2pt}{{- end -}}
\textbf{ {{- $phase.PhaseName -}} }\\
{{- range $subPhase := $phase.SubPhases -}}
\ColorCircle{ {{- $subPhase.Color -}} } \small{ {{- $subPhase.Name -}} }~%
{{- end -}}
\\
{{- end -}}
\par}
{{- else -}}
% Fallback to simple legend if no phase data
{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
{\centering
{{- range $color, $category := $taskColors -}}\ColorCircle{ {{- $color -}} } \small{ {{- $category -}} }\hspace{ {{$.Cfg.Layout.Spacing.ColorLegendSep}} }{{- end -}}
\par}
{{- else -}}
\ColorLegend
{{- end -}}
{{- end -}}
