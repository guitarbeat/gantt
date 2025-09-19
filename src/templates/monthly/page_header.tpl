{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
{\centering
{{- range $color, $category := $taskColors -}}
\ColorCircle{ {{- $color -}} } \small{ {{- $category -}} }%
\hspace{1.5em}%
{{- end -}}
\par}
{\centering
\textcolor{gray!60}{\rule{0.6\textwidth}{0.8pt}}%
\par}
{{- else -}}
\ColorLegend
{{- end -}}
{\noindent\Large\renewcommand{\arraystretch}{\myNumArrayStretch}
{{- .Body.Breadcrumb -}}
\hfill%
{{ .Body.Extra.Table false -}}
}
\myLineThick