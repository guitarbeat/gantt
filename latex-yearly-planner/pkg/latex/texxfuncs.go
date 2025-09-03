package latex



func EmphCell(text string) string {
	return CellColor("black", TextColor("white", text))
}
