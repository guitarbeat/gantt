package templates

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const nl = "\n"

// LaTeX Functions

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

// Header Items

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
		out = Target(s, s)
	} else {
		out = Link(s, s)
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
		return Target(ref, text)
	}

	return Link(ref, text)
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

func (c *CellItem) Select() *CellItem {
	c.selected = true
	return c
}

func (c *CellItem) Selected(selected bool) *CellItem {
	c.selected = selected
	return c
}

func (c *CellItem) Refer(ref string) *CellItem {
	c.Ref = ref
	return c
}

func (c *CellItem) Display() string {
	if len(c.Ref) == 0 {
		c.Ref = c.Text
	}

	link := c.Text

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
		return Target(ref, out)
	}

	return Link(ref, out)
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

func (i Items) WithTopRightCorner(flag bool, kernSpacing string) Items {
	if !flag {
		return i
	}

	return append(i, plainItem(`\kern `+kernSpacing))
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
