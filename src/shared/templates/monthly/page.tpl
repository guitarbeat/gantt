{{ template "macros.tpl" . }}

{{ if .Body.TOCContent }}
% Table of Contents Page
{{ .Body.TOCContent }}
{{ else }}
{{ template "header.tpl" dict "Cfg" .Cfg "Body" .Body }}
{{ template "body.tpl" dict "Cfg" .Cfg "Body" .Body }}

\pagebreak
{{ end }}
