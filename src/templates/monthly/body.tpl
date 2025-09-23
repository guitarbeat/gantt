% Setup category palette for this month
\SetupDefaultCategoryPalette

{{- template "calendar_table.tpl" dict "Cfg" .Cfg "Body" .Body -}}

% Legend at bottom of page
\vfill
{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
{\centering
{{- range $color, $category := $taskColors -}}
\ColorCircle{ {{- $color -}} } \small{ {{- $category -}} }%
\hspace{1.5em}%
{{- end -}}
\par}
{{- else -}}
\ColorLegend
{{- end -}}
