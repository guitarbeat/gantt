{%
{{ if not .Body.Large -}} \renewcommand{\arraystretch}{\myNumArrayStretch}% {{- end}}
\setlength{\tabcolsep}{\myLenTabColSep}%
%
{{ .Body.Month.DefineTable .Body.TableType .Body.Large }}
  {{ .Body.Month.MaybeName .Body.Large }}
  {{ if $.Body.Large -}} \hline {{- end }}
  {{ .Body.Month.WeekHeader .Body.Large }} \\ {{ if .Body.Large -}} \noalign{\hrule height \myLenLineThicknessThick} {{- else -}} \hline {{- end}}
  {{- range $i, $week := .Body.Month.Weeks }}
  {{$week.WeekNumber $.Body.Large}} &
    {{- range $j, $day := $week.Days -}}
      {{- $day.Day $.Body.Today $.Body.Large -}}
      {{- if eq $j 6 -}}
        \\ {{ if $.Body.Large -}} \hline {{- end -}}
      {{- else -}} & {{- end -}}
    {{- end -}}
  {{ end }}
  {{ .Body.Month.EndTable .Body.TableType -}}
}