package latex

import "fmt"

const nl = "\n"

// CellColor creates a colored cell
func CellColor(color, text string) string {
	return fmt.Sprintf(`\cellcolor{%s}{%s}`, color, text)
}

// TextColor creates colored text
func TextColor(color, text string) string {
	return fmt.Sprintf(`\textcolor{%s}{%s}`, color, text)
}

// Hyperlink creates a hyperlink
func Hyperlink(ref, text string) string {
	return fmt.Sprintf(`\hyperlink{%s}{%s}`, ref, text)
}

// Hypertarget creates a hypertarget
func Hypertarget(ref, text string) string {
	return fmt.Sprintf(`\hypertarget{%s}{%s}`, ref, text)
}

// Tabular creates a tabular environment
func Tabular(format, text string) string {
	return `\begin{tabular}{` + format + `}` + nl + text + nl + `\end{tabular}`
}

// ResizeBoxW creates a resized box with specified width
func ResizeBoxW(width, text string) string {
	return fmt.Sprintf(`\resizebox{!}{%s}{%s}`, width, text)
}

// Multirow creates a multirow cell
func Multirow(rows int, text string) string {
	return fmt.Sprintf(`\multirow{%d}{*}{%s}`, rows, text)
}

// Bold creates bold text
func Bold(text string) string {
	return fmt.Sprintf(`\textbf{%s}`, text)
}

// Target creates a target for hyperlinks
func Target(ref, text string) string {
	return "\\hypertarget{" + ref + "}{" + text + "}"
}

// Link creates a hyperlink
func Link(ref, text string) string {
	return "\\hyperlink{" + ref + "}{" + text + "}"
}

// EmphCell creates an emphasized cell with black background and white text
func EmphCell(text string) string {
	return CellColor("black", TextColor("white", text))
}

