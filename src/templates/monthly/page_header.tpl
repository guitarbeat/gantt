{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
\begin{center}%
{{- range $color, $category := $taskColors -}}
\ColorCircle{ {{- $color -}} } \small{ {{- $category -}} }%
\hspace{1.5em}%
{{- end -}}
\end{center}%
\begin{center}%
\textcolor{gray!60}{\rule{0.6\textwidth}{0.8pt}}%
\end{center}%
{{- else -}}
\ColorLegend
{{- end -}}
{\noindent\Large\renewcommand{\arraystretch}{\myNumArrayStretch}
{{- .Body.Breadcrumb -}}
\hfill%
{{ .Body.Extra.Table false -}}
}
\myLineThick