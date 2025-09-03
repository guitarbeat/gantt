package header

import (
	"strconv"

	"github.com/kudrykv/latex-yearly-planner/internal/latex"
)

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
