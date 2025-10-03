% Setup category palette for this month
\SetupDefaultCategoryPalette{}

{{- template "calendar.tpl" dict "Cfg" .Cfg "Body" .Body -}}

% Legend at bottom of page - just colors and categories
\vfill
{{- $phaseGroups := .Body.Month.GetTaskColorsByPhase -}}
{{- if $phaseGroups -}}
\noindent{\small%
{{- range $idx, $phase := $phaseGroups -}}
{{- range $subIdx, $subPhase := $phase.SubPhases -}}
{{- if or $idx $subIdx -}}\quad{{- end -}}\ColorCircle{ {{- $subPhase.Color -}} }{ {{- $subPhase.Name -}} }%
{{- end -}}
{{- end -}}
}
{{- else -}}
% Fallback to simple legend if no phase data
{{- $taskColors := .Body.Month.GetTaskColors -}}
{{- if $taskColors -}}
\noindent{\small%
{{- range $color, $category := $taskColors -}}{{- if ne $color "" -}}\ColorCircle{ {{- $color -}} }{ {{- $category -}} }\quad{{- end -}}{{- end -}}
}
{{- else -}}
\ColorLegend
{{- end -}}
{{- end -}}
