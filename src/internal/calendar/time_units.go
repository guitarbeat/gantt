package calendar

import (
	"strconv"
	"time"

	"phd-dissertation-planner/internal/rendering"
)

// Years and Year

type Years []*Year

type Year struct {
	Number   int
	Quarters Quarters
	Weeks    Weeks
}

func NewYear(wd time.Weekday, year int) *Year {
	out := &Year{Number: year}
	out.Weeks = NewWeeksForYear(wd, out)
	for q := 1; q <= 4; q++ {
		out.Quarters = append(out.Quarters, NewQuarter(wd, out, q))
	}
	return out
}

func (y Year) Breadcrumb() string {
	return rendering.Items{
		rendering.NewIntItem(y.Number).Ref(),
		rendering.NewItemsGroup(
			rendering.NewTextItem("Q1"),
			rendering.NewTextItem("Q2"),
			rendering.NewTextItem("Q3"),
			rendering.NewTextItem("Q4"),
		),
	}.Table(true)
}

func (y Year) SideQuarters(sel ...int) []rendering.CellItem {
	out := make([]rendering.CellItem, 0, len(y.Quarters))
	for i := len(y.Quarters) - 1; i >= 0; i-- {
		mark := false
		for _, oneof := range sel {
			if oneof == y.Quarters[i].Number {
				mark = true
				break
			}
		}
		out = append(out, rendering.NewCellItem(y.Quarters[i].Name()).Selected(mark))
	}
	return out
}

func (y Year) SideMonths(sel ...time.Month) []rendering.CellItem {
	out := make([]rendering.CellItem, 0, 12)
	for i := len(y.Quarters) - 1; i >= 0; i-- {
		for j := len(y.Quarters[i].Months) - 1; j >= 0; j-- {
			mon := y.Quarters[i].Months[j]
			mark := false
			for _, month := range sel {
				if month == mon.Month {
					mark = true
					break
				}
			}
			cell := rendering.NewCellItem(mon.ShortName()).Refer(mon.Month.String()).Selected(mark)
			out = append(out, cell)
		}
	}
	return out
}

func (y Year) HeadingMOS() string {
	return rendering.ResizeBoxW(`\myLenHeaderResizeBox`, rendering.Hypertarget("Calendar", strconv.Itoa(y.Number)))
}

// Quarters and Quarter

type Quarters []*Quarter

type Quarter struct {
	Year   *Year
	Number int
	Months Months
}

func (q Quarters) Numbers() []int {
	if len(q) == 0 {
		return nil
	}
	out := make([]int, 0, len(q))
	for _, quarter := range q {
		out = append(out, quarter.Number)
	}
	return out
}

func NewQuarter(wd time.Weekday, year *Year, qrtr int) *Quarter {
	out := &Quarter{Year: year, Number: qrtr}
	start := time.Month(qrtr*3 - 2)
	end := start + 2
	for month := start; month <= end; month++ {
		out.Months = append(out.Months, NewMonth(wd, year, out, month))
	}
	return out
}

func (q *Quarter) Breadcrumb() string {
	return rendering.Items{rendering.NewIntItem(q.Year.Number), rendering.NewItemsGroup(
		rendering.NewTextItem("Q1").Bold(q.Number == 1).Ref(q.Number == 1),
		rendering.NewTextItem("Q2").Bold(q.Number == 2).Ref(q.Number == 2),
		rendering.NewTextItem("Q3").Bold(q.Number == 3).Ref(q.Number == 3),
		rendering.NewTextItem("Q4").Bold(q.Number == 4).Ref(q.Number == 4),
	)}.Table(true)
}

func (q *Quarter) Name() string { return "Q" + strconv.Itoa(q.Number) }

func (q *Quarter) HeadingMOS() string {
	return ` \begin{tabular}{@{}l}
  \resizebox{!}{\myLenHeaderResizeBox}{` + rendering.Target(q.Name(), q.Name()) + `}
\end{tabular}`
}
