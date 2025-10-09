% Table of Contents - Clickable Task Index
\hypertarget{task-index}{}
{\Large\textbf{Task Index}}

\vspace{0.1cm}

{\small\textit{Total: {{.Body.TotalTasks}} tasks
{{- if .Body.MilestoneCount}} ({{.Body.MilestoneCount}} milestones){{- end -}}
{{- if .Body.CompletedCount}} | {{.Body.CompletedCount}} completed{{- end -}}
}}

\vspace{0.1cm}

{{- range .Body.PhaseOrder}}
    {{- $phase := .}}
    {{- if index $.Body.TaskIndex $phase}}
        {{- if not (eq $phase "1") }}
\vspace{0.2cm}
\hrule height 0.3pt
\vspace{0.1cm}
        {{- end}}

\vspace{0.2cm}
{{- $stats := index $.Body.PhaseStats $phase }}
{{- $phaseName := index $.Body.PhaseNames $phase }}
\textbf{\large {{$phaseName}}} ({{$stats.total}} tasks
{{- if $stats.milestones}}, {{$stats.milestones}} milestones{{end}}
{{- if $stats.completed}}, {{$stats.progress}}\% complete{{end}})

\vspace{0.2cm}

\begin{tabularx}{\linewidth}{
  >{\RaggedRight}X
  >{\RaggedRight}X
}
    {{- range $i, $task := index $.Body.TaskIndex $phase}}
        {{- if mod $i 2 | eq 0}}
{{- $taskName := $task.Name }}
            {{- if $task.IsMilestone}}{{- $taskName = printf "\\textbf{%s} $\\star$" $taskName}}{{- end}}
            {{- if eq ($task.Status | lower) "completed"}}{{- $taskName = printf "$\\checkmark$ \\textcolor{gray}{%s}" $taskName}}{{- end}}
\hyperlink{ {{- $task.StartDate.Format "2006-01-02T15:04:05Z07:00" -}} }{ {{- $taskName -}} }
            {{- if eq (plus $i 1) (len (index $.Body.TaskIndex $phase)) }} & \ {{end}}
        {{- else}}
 &          {{- $taskName := $task.Name }}
            {{- if $task.IsMilestone}}{{- $taskName = printf "\\textbf{%s} $\\star$" $taskName}}{{- end}}
            {{- if eq ($task.Status | lower) "completed"}}{{- $taskName = printf "$\\checkmark$ \\textcolor{gray}{%s}" $taskName}}{{- end}}
\hyperlink{ {{- $task.StartDate.Format "2006-01-02T15:04:05Z07:00" -}} }{ {{- $taskName -}} } \
        {{- end}}
    {{- end}}
\end{tabularx}
\vspace{0.2cm}

    {{- end}}
{{- end}}

\pagebreak
