% Table of Contents - Clickable Task Index
\hypertarget{task-index}{}
{\Large\textbf{Task Index}}

\vspace{0.4cm}

% Data Sources Summary Table
\noindent\begin{tabularx}{\linewidth}{@{}lX@{}}
\textbf{Data Sources:} & {{.Body.CSVFileCount}} CSV file(s) merged \\[2pt]
\textbf{Files:} & {\footnotesize {{range $i, $file := .Body.CSVFiles}}{{if $i}}, {{end}}{{$file}}{{end}}} \\[2pt]
\textbf{Total Tasks:} & {{.Body.TotalTasks}} tasks{{if .Body.MilestoneCount}} ({{.Body.MilestoneCount}} milestones){{end}}{{if .Body.CompletedCount}} | {{.Body.CompletedCount}} completed{{end}} \\
\end{tabularx}

\vspace{0.4cm}

{{- $currentSection := ""}}
{{- range .Body.PhaseOrder}}
    {{- $phase := .}}
    {{- if index $.Body.TaskIndex $phase}}
{{- $stats := index $.Body.PhaseStats $phase }}
{{- $phaseName := index $.Body.PhaseNames $phase }}
{{- $phaseColor := index $.Body.PhaseColors $phase }}
{{- $section := index $.Body.PhaseToSection $phase }}

{{- if ne $section $currentSection}}
{{- $currentSection = $section}}

% Section: {{$section}}
\vspace{0.5cm}
{\LARGE\textbf{ {{- $section -}} }}
\vspace{0.2cm}

\noindent\rule{\linewidth}{0.8pt}
\vspace{0.3cm}
{{- end}}

% Phase: {{$phaseName}}
\vspace{0.25cm}
\noindent\colorbox[RGB]{ {{- $phaseColor -}} }{\parbox{0.98\linewidth}{\vspace{2pt}\textbf{\large {{$phaseName}}}\hfill{\small {{$stats.total}} tasks{{if $stats.milestones}}, {{$stats.milestones}} milestones{{end}}{{if $stats.completed}}, {{$stats.progress}}\% complete{{end}}}\vspace{2pt}}}

\vspace{0.15cm}

\noindent\begin{tabularx}{\linewidth}{@{\hspace{0.5em}}c@{\hspace{0.8em}}>{\RaggedRight}X@{\hspace{0.8em}}l@{\hspace{0.8em}}l@{\hspace{0.5em}}}
\hline
\textbf{\#} & \textbf{Task} & \textbf{Start} & \textbf{End} \\
\hline
    {{- range $i, $task := index $.Body.TaskIndex $phase}}
        {{- $taskName := $task.Name }}
        {{- $taskIcon := "" }}
        {{- if $task.IsMilestone}}{{- $taskIcon = "$\\star$" }}{{- $taskName = printf "\\textbf{%s}" $taskName}}{{- end}}
        {{- if eq ($task.Status | lower) "completed"}}{{- $taskIcon = "$\\checkmark$" }}{{- $taskName = printf "\\textcolor{gray}{%s}" $taskName}}{{- end}}
{{plus $i 1}} & \hyperlink{ {{- $task.StartDate.Format "2006-01-02T15:04:05Z07:00" -}} }{ {{- $taskName -}} } {{$taskIcon}} & {\footnotesize {{$task.StartDate.Format "Jan 02"}}} & {\footnotesize {{$task.EndDate.Format "Jan 02"}}} \\
    {{- end}}
\hline
\end{tabularx}

    {{- end}}
{{- end}}

\pagebreak
