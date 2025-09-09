package header

import (
	"strconv"
	"strings"
	"time"

	"latex-yearly-planner/internal/latex"
)

// Item interface defines the contract for all header items
type Item interface {
	Display() string
}

// Items is a slice of Item interfaces
type Items []Item

// plainItem creates a simple item that just displays the given text
type plainItem string

func (p plainItem) Display() string {
	return string(p)
}

// ItemsGroup groups multiple items with a delimiter
type ItemsGroup struct {
	Items Items
	delim string
}

func NewItemsGroup(items ...Item) ItemsGroup {
	return ItemsGroup{
		Items: items,
		delim: "\\quad{}",
	}
}

func (i ItemsGroup) Display() string {
	list := make([]string, 0, len(i.Items))

	for _, item := range i.Items {
		list = append(list, item.Display())
	}

	return strings.Join(list, i.delim)
}

func (i ItemsGroup) Delim(delim string) ItemsGroup {
	i.delim = delim
	return i
}

// IntItem represents an integer item with optional reference
type IntItem struct {
	Val int
	ref bool
}

func (i IntItem) Display() string {
	var out string
	s := strconv.Itoa(i.Val)

	if i.ref {
		out = latex.Target(s, s)
	} else {
		out = latex.Link(s, s)
	}

	return out
}

func (i IntItem) Ref() IntItem {
	i.ref = true
	return i
}

func NewIntItem(val int) IntItem {
	return IntItem{Val: val}
}

// MonthItem represents a month item with optional reference and shortening
type MonthItem struct {
	Val     time.Month
	ref     bool
	shorten bool
}

func (m MonthItem) Display() string {
	ref := m.Val.String()
	text := ref

	if m.shorten {
		text = text[:3]
	}

	if m.ref {
		return latex.Target(ref, text)
	}

	return latex.Link(ref, text)
}

func (m MonthItem) Ref() MonthItem {
	m.ref = true
	return m
}

func (m MonthItem) Shorten(f bool) MonthItem {
	m.shorten = f
	return m
}

func NewMonthItem(mo time.Month) MonthItem { 
	return MonthItem{Val: mo} 
}

// CellItem represents a cell item with text, reference, and selection state
type CellItem struct {
	Text     string
	Ref      string
	selected bool
}

func NewCellItem(text string) CellItem {
	return CellItem{Text: text}
}

func (c CellItem) Select() CellItem {
	c.selected = true
	return c
}

func (c CellItem) Selected(selected bool) CellItem {
	c.selected = selected
	return c
}

func (c CellItem) Refer(ref string) CellItem {
	c.Ref = ref
	return c
}

func (c CellItem) Display() string {
	if len(c.Ref) == 0 {
		c.Ref = c.Text
	}

	link := `\hyperlink{` + c.Ref + `}{` + c.Text + `}`

	if c.selected {
		return `\cellcolor{black}{\textcolor{white}{` + link + `}}`
	}

	return link
}

// TextItem represents a text item with formatting and reference options
type TextItem struct {
	Name      string
	bold      bool
	ref       bool
	refPrefix string
	refText   string
}

func NewTextItem(name string) TextItem {
	return TextItem{
		Name: name,
	}
}

func (t TextItem) Display() string {
	var (
		out string
		ref string
	)
	if t.bold {
		out = "\\textbf{" + t.Name + "}"
	} else {
		out = t.Name
	}

	if len(t.refText) > 0 {
		ref = t.refText
	} else {
		ref = t.refPrefix + t.Name
	}

	if t.ref {
		return latex.Target(ref, out)
	}

	return latex.Link(ref, out)
}

func (t TextItem) Ref(ref bool) TextItem {
	t.ref = ref
	return t
}

func (t TextItem) Bold(f bool) TextItem {
	t.bold = f
	return t
}

func (t TextItem) RefPrefix(refPrefix string) TextItem {
	t.refPrefix = refPrefix
	return t
}

func (t TextItem) RefText(refText string) TextItem {
	t.refText = refText
	return t
}

// Items methods for working with collections of items

func (i Items) WithTopRightCorner(flag bool) Items {
	if !flag {
		return i
	}

	return append(i, plainItem(`\kern 5mm`))
}

func (i Items) Length() int {
	return len(i)
}

func (i Items) ColSetup(left bool) string {
	if left {
		return "|" + strings.Join(strings.Split(strings.Repeat("l", len(i)), ""), "|")
	}

	return strings.Join(strings.Split(strings.Repeat("r", len(i)), ""), "|") + "@{}"
}

func (i Items) Row() string {
	out := make([]string, 0, len(i))

	for _, item := range i {
		out = append(out, item.Display())
	}

	return strings.Join(out, " & ")
}

func (i Items) Table(left bool) string {
	if len(i) == 0 {
		return ""
	}

	return `\begin{tabular}{` + i.ColSetup(left) + `}
` + i.Row() + `
\end{tabular}`
}

