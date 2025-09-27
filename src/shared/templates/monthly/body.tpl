% Setup category palette for this month
\SetupDefaultCategoryPalette{}

{{- template "calendar.tpl" dict "Cfg" .Cfg "Body" .Body -}}

% Legend at bottom of page
\vfill
{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
{\centering
{{- range $color, $category := $taskColors -}}\ColorCircle{ {{- $color -}} } \small{ {{- $category -}} }\hspace{ {{$.Cfg.Layout.Spacing.ColorLegendSep}} }{{- end -}}
\par}
{{- else -}}
\ColorLegend
{{- end -}}
