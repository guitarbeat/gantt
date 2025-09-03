package texext

import tex "github.com/kudrykv/latex-yearly-planner/components/latex"

func EmphCell(text string) string {
	return tex.CellColor("black", tex.TextColor("white", text))
}
