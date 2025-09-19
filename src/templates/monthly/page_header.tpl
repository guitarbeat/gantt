{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
{{- range $color, $category := $taskColors -}}
\ColorCircle{ {{- $color -}} } \small{ {{- $category -}} }%
\hspace{1.5em}%
{{- end -}}
\vspace*{-0.5ex}%
\textcolor{gray!60}{\rule{0.6\textwidth}{0.8pt}}%
{{- else -}}
\ColorLegend
{{- end -}}
{\noindent\Large\renewcommand{\arraystretch}{\myNumArrayStretch}
{{- .Body.Breadcrumb -}}
\hfill%
{{ .Body.Extra.Table false -}}
}
\myLineThick