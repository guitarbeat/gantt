% Table of Contents - Clickable Task Index
{\Large\textbf{Task Index}}

\vspace{0.5cm}

% Simple table of all tasks
\begin{tabularx}{\linewidth}{|>{\raggedright\arraybackslash}p{0.5\linewidth}|>{\centering\arraybackslash}p{0.15\linewidth}|>{\centering\arraybackslash}p{0.15\linewidth}|>{\centering\arraybackslash}p{0.1\linewidth}|}
\hline
\textbf{Task} & \textbf{Category} & \textbf{Start Date} & \textbf{Type} \\
\hline
{{- range .Body.TaskIndex -}}
\hyperlink{ {{- .DateRef -}} }{ {{- .Name -}} }
{{- if .IsMilestone -}}
\ (*)
{{- end -}}
& {{- .Category -}} & {{- .StartDate -}} &
{{- if .IsMilestone -}}
Milestone
{{- else -}}
Task
{{- end -}} \\
\hline
{{- end -}}
\end{tabularx}

\vspace{0.3cm}

% Legend
{\small
\textbf{Legend:} (*) Milestone tasks with enhanced borders \\
Click on any task name to jump to its location in the timeline.
}

\pagebreak
