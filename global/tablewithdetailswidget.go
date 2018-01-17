package global

import (
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

func NewTableWithDetailsRow(bgcolor uint32, labels []string, details sws.Widget) *TableWithDetailsRow {
	row := &TableWithDetailsRow{
		detailsopened: false,
		labels:        labels,
		details:       details,
		bgcolor:       bgcolor,
	}
	return row
}

type TableWithDetails struct {
	sws.CoreWidget
	vertical        *sws.ScrollbarWidget
	headertextcolor sdl.Color
	textcolor       sdl.Color
	yoffset         int32
	arrowUp         *sdl.Surface
	arrowDown       *sdl.Surface
	//data
	header     []string
	headerSize []int32
	rows       []*TableWithDetailsRow
}

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
	}

	widget.vertical.Move(w-15, HEADER_HEIGHT)
	widget.vertical.SetMinimum(0)
	widget.vertical.SetMaximum(0)
	widget.vertical.SetCallback(func(currentposition int32) {
		widget.yoffset = currentposition
		widget.PostUpdate()
	})

	if surface, err := LoadImageAsset("assets/ui/icon-arrowhead-pointing-to-the-right.png"); err == nil {
		widget.arrowUp = surface
	}
	if surface, err := LoadImageAsset("assets/ui/icon-sort-down.png"); err == nil {
		widget.arrowDown = surface
	}

	return widget
}

func (self *TableWithDetails) AddHeader(label string, size int32) {
	self.header = append(self.header, label)
	self.headerSize = append(self.headerSize, size)
}

func (self *TableWithDetails) AddRowTop(row *TableWithDetailsRow) {
	self.rows = append([]*TableWithDetailsRow{row}, self.rows...)
	self.PostUpdate()
}

func (self *TableWithDetails) AddRow(row *TableWithDetailsRow) {
	self.rows = append(self.rows, row)
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
	self.FillRect(0, 0, 25, self.Height(), 0xffdddddd)
	self.FillRect(25, 0, self.Width()-25, self.Height(), 0xffffffff)

	nbcolumns := int32(len(self.header))

	y := HEADER_HEIGHT - self.yoffset
	// print the cells
	for j := int32(0); j < int32(len(self.rows)); j++ {
		row := self.rows[j]
		self.FillRect(25, y, self.Width()-25, 25, row.bgcolor)

		// show the details button
		var arrow *sdl.Surface
		if row.detailsopened {
			arrow = self.arrowDown
		} else {
			arrow = self.arrowUp
		}
		if arrow != nil {
			rectSrc := sdl.Rect{0, 0, arrow.W, arrow.H}
			rectDst := sdl.Rect{0, y, arrow.W, arrow.H}
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
			rectDst := sdl.Rect{25, y, row.details.Width(), row.details.Height()}
			row.details.Surface().Blit(&rectSrc, self.Surface(), &rectDst)

			y += row.details.Height()
		}
	}

	// headers
	self.FillRect(0, 0, self.Width(), HEADER_HEIGHT, 0xffdddddd)
	xoffset := int32(25)
	for i := int32(0); i < nbcolumns; i++ {
		label := self.header[i]
		size := self.headerSize[i]
		self.renderText(xoffset, 0, size, HEADER_HEIGHT, label, self.headertextcolor)

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
	self.SetDirtyFalse()
}

func (self *TableWithDetails) MousePressDown(x, y int32, button uint8) {
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
