{%
{{ if not .Body.Large }} \renewcommand{\arraystretch}{\myNumArrayStretch} {{ end }}
\setlength{\tabcolsep}{\myLenTabColSep}
{{ $tbl := .Body.Month.DefineTable .Body.TableType .Body.Large }}
{{ if $tbl }}
{{$tbl}}
{{ end }}
  {{ $mname := .Body.Month.MaybeName .Body.Large }}
  {{ if $mname }}{{$mname}}{{ end }}
  {{ if $.Body.Large }} \hline {{ end }}
  {{ $wh := .Body.Month.WeekHeader .Body.Large }}
  {{ if $wh }}{{$wh}} \\ {{ if .Body.Large }} \noalign{\hrule height \myLenLineThicknessThick} {{ else }} \hline {{ end}}{{ end }}
  {{ range $i, $week := .Body.Month.Weeks }}
  {{ if $week.HasDays }}
  {{$week.WeekNumber $.Body.Large}} &
    {{ range $j, $day := $week.Days }}
      {{ $cell := $day.Day $.Body.Today $.Body.Large }}
      {{ if $cell }}
        {{$cell}}
      {{ end }}
      {{ if eq $j 6 }}
        \\[\myLenMonthlyCellHeight] {{ if $.Body.Large }} \hline {{ end }}
      {{ else }} & {{ end }}
    {{ end }}
  {{ end }}
  {{ end }}
  {{ .Body.Month.EndTable .Body.TableType }}
}
