package ui

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/nzin/dctycoon/global"
	"github.com/nzin/sws"
	"github.com/veandco/go-sdl2/sdl"
)

const HEADER_HEIGHT = 25

type TableWithDetailsRow struct {
	detailsopened bool
	labels        []string
	details       sws.Widget
	bgcolor       uint32
}

// NewTableWithDetailsRow is about create a new row for a TableWithDetailsRow
func NewTableWithDetailsRow(bgcolor uint32, labels []string, details sws.Widget) *TableWithDetailsRow {
	row := &TableWithDetailsRow{
		detailsopened: false,
		labels:        labels,
		details:       details,
		bgcolor:       bgcolor,
	}
	return row
}

type TableWithDetailsRowBy func(l1, l2 string) bool

// TableWithDetailsRowByPriceDollar is TableWithDetailsRowBy helper implementation for string like "123.45 $"
func TableWithDetailsRowByPriceDollar(l1, l2 string) bool {
	var f1, f2 float64
	fmt.Sscanf(l1, "%f $", &f1)
	fmt.Sscanf(l2, "%f $", &f2)
	return f1 < f2
}

// TableWithDetailsRowByInteger is a TableWithDetailsRowBy helper implementation for string like "123"
func TableWithDetailsRowByInteger(l1, l2 string) bool {
	nb1, _ := strconv.Atoi(l1)
	nb2, _ := strconv.Atoi(l2)
	return nb1 < nb2
}

// TableWithDetailsRowByYearMonthDay is a TableWithDetailsRowBy helper implementation for string like "24-1-1995"
func TableWithDetailsRowByYearMonthDay(l1, l2 string) bool {
	var d1, d2, m1, m2, y1, y2 int
	fmt.Sscanf(l1, "%d-%d-%d", &d1, &m1, &y1)
	fmt.Sscanf(l2, "%d-%d-%d", &d2, &m2, &y2)
	if y1 != y2 {
		return y1 < y2
	}
	if m1 != m2 {
		return m1 < m2
	}
	return d1 < d2
}

// TableWithDetails is a data table with header + expandable details widget
type TableWithDetails struct {
	sws.CoreWidget
	vertical        *sws.ScrollbarWidget
	headertextcolor sdl.Color
	textcolor       sdl.Color
	yoffset         int32
	arrowUp         *sdl.Surface
	arrowDown       *sdl.Surface
	//data
	header          []string
	headerSize      []int32
	headerSort      []TableWithDetailsRowBy
	rows            []*TableWithDetailsRow
	currentSorter   int32
	directionSorter bool
}

func (t *TableWithDetails) Len() int      { return len(t.rows) }
func (t *TableWithDetails) Swap(i, j int) { t.rows[i], t.rows[j] = t.rows[j], t.rows[i] }
func (t *TableWithDetails) Less(i, j int) bool {
	if t.directionSorter {
		return t.headerSort[t.currentSorter](t.rows[i].labels[t.currentSorter], t.rows[j].labels[t.currentSorter])
	} else {
		return !t.headerSort[t.currentSorter](t.rows[i].labels[t.currentSorter], t.rows[j].labels[t.currentSorter])
	}
}

// NewTableWithDetails creates a TableWithDetails i.e. a data table with header + expandable details widget
func NewTableWithDetails(w, h int32) *TableWithDetails {
	corewidget := sws.NewCoreWidget(w, h)
	widget := &TableWithDetails{
		CoreWidget:      *corewidget,
		vertical:        sws.NewScrollbarWidget(15, 20, false),
		yoffset:         0,
		textcolor:       sdl.Color{0, 0, 0, 255},
		headertextcolor: sdl.Color{0, 0, 0, 255},
		header:          make([]string, 0, 0),
		headerSize:      make([]int32, 0, 0),
		rows:            make([]*TableWithDetailsRow, 0, 0),
		currentSorter:   -1,
	}

	widget.vertical.Move(w-15, HEADER_HEIGHT)
	widget.vertical.SetMinimum(0)
	widget.vertical.SetMaximum(0)
	widget.vertical.SetCallback(func(currentposition int32) {
		widget.yoffset = currentposition
		widget.PostUpdate()
	})

	if surface, err := global.LoadImageAsset("assets/ui/icon-arrowhead-pointing-to-the-right.png"); err == nil {
		widget.arrowUp = surface
	}
	if surface, err := global.LoadImageAsset("assets/ui/icon-sort-down.png"); err == nil {
		widget.arrowDown = surface
	}

	return widget
}

func (self *TableWithDetails) AddHeader(label string, size int32, sorterfunction TableWithDetailsRowBy) {
	self.header = append(self.header, label)
	self.headerSize = append(self.headerSize, size)
	self.headerSort = append(self.headerSort, sorterfunction)
}

func (self *TableWithDetails) AddRowTop(row *TableWithDetailsRow) {
	self.rows = append([]*TableWithDetailsRow{row}, self.rows...)
	if self.currentSorter >= 0 && self.headerSort[self.currentSorter] != nil {
		sort.Sort(self)
	}
	self.Resize(self.Width(), self.Height())
	self.PostUpdate()
}

func (self *TableWithDetails) AddRow(row *TableWithDetailsRow) {
	self.rows = append(self.rows, row)
	if self.currentSorter >= 0 && self.headerSort[self.currentSorter] != nil {
		sort.Sort(self)
	}
	self.Resize(self.Width(), self.Height())
	self.PostUpdate()
}

func (self *TableWithDetails) SetHeaderTextColor(color sdl.Color) {
	self.headertextcolor = color
	self.PostUpdate()
}

func (self *TableWithDetails) SetTextColor(color sdl.Color) {
	self.textcolor = color
	self.PostUpdate()
}

func (self *TableWithDetails) Resize(width, height int32) {
	self.CoreWidget.Resize(width, height)

	// sanity
	if height < HEADER_HEIGHT {
		height = HEADER_HEIGHT
	}

	self.vertical.Resize(15, height-HEADER_HEIGHT)
	// compute table height
	tableheight := int32(HEADER_HEIGHT)
	for _, r := range self.rows {
		tableheight += 25
		if r.detailsopened == true {
			tableheight += r.details.Height()
		}
	}
	if height < tableheight {
		self.vertical.Move(width-15, HEADER_HEIGHT)
		self.AddChild(self.vertical)
		self.vertical.SetMaximum(tableheight - self.Height())
	} else {
		self.RemoveChild(self.vertical)
	}
	self.PostUpdate()
}

func (self *TableWithDetails) renderText(x, y, width, height int32, label string, textcolor sdl.Color) {
	text, err := self.Font().RenderUTF8Blended(label, textcolor)
	if err != nil {
		return
	}
	wGap := width - text.W
	hGap := 25 - text.H
	if wGap < 0 {
		wGap = 0
	}
	if hGap < 0 {
		hGap = 0
	}
	maxwidth := text.W
	maxheight := text.H
	if maxwidth > width {
		maxwidth = width
	}
	if maxheight > height {
		maxheight = height
	}
	rectSrc := sdl.Rect{0, 0, maxwidth, maxheight}
	rectDst := sdl.Rect{x + (wGap / 2), y + (hGap / 2), width - (wGap / 2), height - (hGap / 2)}
	if err = text.Blit(&rectSrc, self.Surface(), &rectDst); err != nil {
	}
}

func (self *TableWithDetails) Repaint() {
	self.FillRect(0, 0, self.Width(), self.Height(), 0xffffffff)

	nbcolumns := int32(len(self.header))

	y := HEADER_HEIGHT - self.yoffset
	// print the cells
	for j := int32(0); j < int32(len(self.rows)); j++ {
		row := self.rows[j]
		self.FillRect(0, y, self.Width(), 25, row.bgcolor)

		// show the details button
		var arrow *sdl.Surface
		if row.detailsopened {
			arrow = self.arrowDown
		} else {
			arrow = self.arrowUp
		}
		if arrow != nil {
			rectSrc := sdl.Rect{0, 0, arrow.W, arrow.H}
			rectDst := sdl.Rect{(25 - arrow.W) / 2, y + (25-arrow.H)/2, arrow.W, arrow.H}
			arrow.Blit(&rectSrc, self.Surface(), &rectDst)
		}

		// show the row
		xoffset := int32(25)
		for i := int32(0); i < nbcolumns; i++ {
			label := row.labels[i]
			self.renderText(xoffset, y, self.headerSize[i], 25, label, self.textcolor)
			xoffset += self.headerSize[i]
		}
		y += 25
		// show the details
		if row.detailsopened {
			if row.details.IsDirty() {
				row.details.Repaint()
			}
			rectSrc := sdl.Rect{0, 0, row.details.Width(), row.details.Height()}
			rectDst := sdl.Rect{0, y, row.details.Width(), row.details.Height()}
			row.details.Surface().Blit(&rectSrc, self.Surface(), &rectDst)

			y += row.details.Height()
		}
		self.SetDrawColorHex(0xffeeeeee)
		self.DrawLine(1, y-1, self.Width()-2, y-1)
	}

	// headers
	self.FillRect(0, 0, self.Width(), HEADER_HEIGHT, 0xffdddddd)

	//bezel (+)
	self.SetDrawColorHex(0xffffffff)
	self.DrawLine(0, 0, HEADER_HEIGHT-1, 0)
	self.DrawLine(0, 0, 0, HEADER_HEIGHT-1)
	self.SetDrawColor(50, 50, 50, 255)
	self.DrawLine(HEADER_HEIGHT-1, 1, HEADER_HEIGHT-1, HEADER_HEIGHT-1)
	self.DrawLine(1, HEADER_HEIGHT-1, HEADER_HEIGHT-1, HEADER_HEIGHT-1)

	xoffset := int32(25)
	for i := int32(0); i < nbcolumns; i++ {
		label := self.header[i]
		size := self.headerSize[i]
		self.renderText(xoffset, 0, size, HEADER_HEIGHT, label, self.headertextcolor)

		if self.currentSorter == i {
			self.SetDrawColorHex(0xff000000)
			if self.directionSorter {
				for j := int32(0); j < 8; j++ {
					self.DrawLine(xoffset+size-23+j, 8+j, xoffset+size-7-j, 8+j)
				}
			} else {
				for j := int32(0); j < 8; j++ {
					self.DrawLine(xoffset+size-15-j, 8+j, xoffset+size-15+j, 8+j)
				}
			}
		}

		//bezel
		self.SetDrawColorHex(0xffffffff)
		self.DrawLine(xoffset, 0, xoffset+size-1, 0)
		self.DrawLine(xoffset, 0, xoffset, HEADER_HEIGHT-1)
		self.SetDrawColor(50, 50, 50, 255)
		self.DrawLine(xoffset+size-1, 1, xoffset+size-1, HEADER_HEIGHT-1)
		self.DrawLine(xoffset+1, HEADER_HEIGHT-1, xoffset+size-1, HEADER_HEIGHT-1)

		xoffset += size
	}

	for _, child := range self.GetChildren() {
		// adjust the clipping to the current child
		if child.IsDirty() {
			child.Repaint()
		}
		rectSrc := sdl.Rect{0, 0, child.Width(), child.Height()}
		rectDst := sdl.Rect{child.X(), child.Y(), child.Width(), child.Height()}
		child.Surface().Blit(&rectSrc, self.Surface(), &rectDst)
	}

	// global bezel
	self.SetDrawColorHex(0xffffffff)
	self.DrawLine(self.Width()-1, HEADER_HEIGHT, self.Width()-1, self.Height()-1)
	self.DrawLine(0, self.Height()-1, self.Width()-1, self.Height()-1)
	self.SetDrawColor(50, 50, 50, 255)
	self.DrawLine(0, HEADER_HEIGHT, 0, self.Height()-1)

	self.SetDirtyFalse()
}

func (self *TableWithDetails) MousePressDown(x, y int32, button uint8) {
	if button == sdl.BUTTON_LEFT && y <= HEADER_HEIGHT {
		xoffset := int32(25)
		for i := int32(0); i < int32(len(self.header)); i++ {
			size := self.headerSize[i]
			if x >= xoffset && x < xoffset+size {
				if self.currentSorter == i {
					self.directionSorter = !self.directionSorter
				} else {
					self.currentSorter = i
				}
				if self.headerSort[self.currentSorter] != nil {
					sort.Sort(self)
					self.PostUpdate()
				}
			}
			xoffset += size
		}
	}
	if button == sdl.BUTTON_LEFT && y > HEADER_HEIGHT {
		yoffset := HEADER_HEIGHT - self.yoffset
		for j := int32(0); j < int32(len(self.rows)); j++ {
			row := self.rows[j]
			if y >= yoffset && y < yoffset+25 {
				row.detailsopened = !row.detailsopened
				self.Resize(self.Width(), self.Height())
				self.PostUpdate()
			}
			yoffset += 25
			if row.detailsopened {
				yoffset += row.details.Height()
			}
		}
	}
}

func (self *TableWithDetails) ClearRows() {
	self.rows = make([]*TableWithDetailsRow, 0, 0)
	self.PostUpdate()
}
