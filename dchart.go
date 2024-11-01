// Package dchart - make charts in the deck format
package dchart

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/deck/generate"
)

// ChartData defines the name,value pairs
type ChartData struct {
	label string
	value float64
	note  string
}

// Flags define chart on/off switches
type Flags struct {
	DataMinimum,
	FullDeck,
	ReadCSV,
	ShowAxis,
	ShowBar,
	ShowBowtie,
	ShowDonut,
	ShowDot,
	ShowFan,
	ShowFrame,
	ShowGrid,
	ShowHBar,
	ShowLine,
	ShowLego,
	ShowNote,
	ShowPercentage,
	ShowPGrid,
	ShowPMap,
	ShowRadial,
	ShowRegressionLine,
	ShowScatter,
	ShowSlope,
	ShowSpokes,
	ShowTitle,
	ShowValues,
	ShowVolume,
	ShowWBar,
	ShowXLast,
	ShowXstagger,
	SolidPMap bool
}

// Attributes define chart attributes
type Attributes struct {
	BackgroundColor,
	DataColor,
	FrameColor,
	LabelColor,
	RegressionLineColor,
	ValueColor,
	ChartTitle,
	CSVCols,
	DataCondition,
	DataFmt,
	HLine,
	NoteLocation,
	ValuePosition,
	YAxisR string
}

// Measures define chart measures
type Measures struct {
	CanvasWidth,
	CanvasHeight,
	TextSize,
	Left,
	Right,
	Top,
	Bottom,
	LineSpacing,
	BarWidth,
	LineWidth,
	PSize,
	PWidth,
	UserMin,
	UserMax,
	VolumeOpacity,
	XLabelRotation float64
	Boundary string
	XLabelInterval,
	PMapLength int
}

// Settings is a collection of all chart settings
type Settings struct {
	Flags
	Attributes
	Measures
}

var blue7 = []string{
	"rgb(8,69,148)",
	"rgb(33,113,181)",
	"rgb(66,146,198)",
	"rgb(107,174,214)",
	"rgb(158,202,225)",
	"rgb(198,219,239)",
	"rgb(239,243,255)",
}

var xmlmap = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;")

var nlmap = strings.NewReplacer(
	`\n`, " ",
	`\t`, " ",
)

const (
	// Titlecolor is the color of the title string
	Titlecolor = "black"
	// Dotlinecolor is the dotted line color
	Dotlinecolor = "lightgray"
	// Defaultfmt is the default number format
	Defaultfmt   = "%.1f"
	wbop         = 30.0
	largest      = math.MaxFloat64
	smallest     = -math.MaxFloat64
	topclock     = math.Pi / 2
	fullcircle   = math.Pi * 2
	transparency = 50.0
)

// xmlesc escapes XML
func xmlesc(s string) string {
	return xmlmap.Replace(s)
}

// vmap maps one range into another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// getheader returns the indicies of the comma-separated list of fields
// by default or on error, return 0, 1
// For example given this header:
// First,Second,Third,Sum
// First,Sum returns 0,3 and First,Third returns 0,2
func getheader(s []string, lv string) (int, int) {
	li := 0
	vi := 1
	cv := strings.Split(lv, ",")
	if len(cv) != 2 {
		return li, vi
	}
	for i, p := range s {
		if p == cv[0] {
			li = i
		}
		if p == cv[1] {
			vi = i
		}
	}
	return li, vi
}

// Getdata reads input from a Reader, either tab-separated or CSV
func Getdata(r io.ReadCloser, readcsv bool, cols string) ([]ChartData, float64, float64, string) {
	var min, max float64
	var title string
	var data []ChartData
	if readcsv {
		data, min, max, title = CSVdata(r, cols)
	} else {
		data, min, max, title = TSVdata(r)
	}
	return data, min, max, title
}

// CSVdata reads CSV structured name,value pairs, with optional comments,
// returning a slice with the data, allong with min, max and title
func CSVdata(r io.ReadCloser, csvcols string) ([]ChartData, float64, float64, string) {
	var (
		data []ChartData
		d    ChartData
		err  error
	)
	input := csv.NewReader(r)
	maxval := smallest
	minval := largest
	title := ""
	n := 0
	li := 0
	vi := 1
	for {
		n++
		fields, csverr := input.Read()
		if csverr == io.EOF {
			break
		}
		if csverr != nil {
			fmt.Fprintf(os.Stderr, "%v %v\n", csverr, fields)
			continue
		}

		if len(fields) < 2 {
			continue
		}
		if fields[0] == "#" {
			title = fields[1]
			continue
		}
		if len(fields) == 3 {
			d.note = xmlesc(fields[2])
		} else {
			d.note = ""
		}
		if n == 1 && len(csvcols) > 0 { // column header is assumed to be the first row
			li, vi = getheader(fields, csvcols)
			title = fields[vi]
			continue
		}

		d.label = xmlesc(fields[li])
		d.value, err = strconv.ParseFloat(fields[vi], 64)
		if err != nil {
			d.value = 0
		}
		if d.value > maxval {
			maxval = d.value
		}
		if d.value < minval {
			minval = d.value
		}
		data = append(data, d)
	}
	r.Close()
	return data, minval, maxval, xmlesc(title)
}

// TSVdata reads tab-delimited name,value pairs, with optional comments,
// returning a slice with the data, allong with min, max and title
func TSVdata(r io.ReadCloser) ([]ChartData, float64, float64, string) {
	var (
		data []ChartData
		d    ChartData
		err  error
	)

	maxval := smallest
	minval := largest
	title := ""
	scanner := bufio.NewScanner(r)
	// read a line, parse into name, value pairs
	// compute min and max values
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if t[0] == '#' && len(t) > 2 { // process titles
			title = strings.TrimSpace(t[1:])
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 2 {
			continue
		}
		if len(fields) == 3 {
			d.note = xmlesc(fields[2])
		} else {
			d.note = ""
		}
		d.label = xmlesc(fields[0])
		d.value, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			d.value = 0
		}
		if d.value > maxval {
			maxval = d.value
		}
		if d.value < minval {
			minval = d.value
		}
		data = append(data, d)
	}
	r.Close()
	return data, minval, maxval, xmlesc(title)
}

// dottedvline makes dotted vertical line, using circles,
// with specified step
func dottedvline(deck *generate.Deck, x, y1, y2, dotsize, step float64, color string) {

	if y1 < y2 { // positive
		for y := y1; y <= y2; y += step {
			deck.Circle(x, y, dotsize, color)
		}
	} else { // negative
		for y := y2; y <= y1; y += step {
			deck.Circle(x, y, dotsize, color)
		}
	}
}

// dottedhline makes a dotted horizontal line, using circles,
// with specified step and separation
func dottedhline(d *generate.Deck, x, y, width, height, step, space float64, color string) {
	for xp := x; xp < x+width; xp += step {
		d.Circle(xp, y, height, color)
		xp += space
	}
}

// yrange parses the min, max, step for axis labels
func yrange(s string) (float64, float64, float64) {
	var min, max, step float64
	n, err := fmt.Sscanf(s, "%f,%f,%f", &min, &max, &step)
	if n != 3 || err != nil {
		return 0, 0, 0
	}
	return min, max, step
}

// cyrange computes "optimal" min, max, step for axis labels
// rounding the max to the appropriate number, given the number of labels
func cyrange(min, max float64, n int) (float64, float64, float64) {
	l := math.Log10(max)
	p := math.Pow10(int(l))
	pl := math.Ceil(max / p)
	ymax := pl * p
	return min, ymax, ymax / float64(n)
}

// yaxis constructs y axis labels
func (s *Settings) yaxis(deck *generate.Deck, x, dmin, dmax float64) {
	var axismin, axismax, step float64
	if s.Attributes.YAxisR == "" {
		axismin, axismax, step = cyrange(dmin, dmax, 5)
	} else {
		axismin, axismax, step = yrange(s.Attributes.YAxisR)
	}
	if step <= 0 {
		return
	}
	var axisfmt = "%0.f"
	if step < 1 {
		axisfmt = "%3.2f"
	}
	left := s.Measures.Left
	if left < 0 {
		left = 10.0
	}
	for y := axismin; y <= axismax; y += step {
		yp := vmap(y, dmin, dmax, s.Measures.Bottom, s.Measures.Top)
		deck.TextEnd(x, yp, fmt.Sprintf(axisfmt, y), "sans", s.Measures.TextSize*0.75, s.Attributes.LabelColor)
		if s.Flags.ShowGrid {
			deck.Line(left, yp, s.Measures.Right, yp, 0.1, "lightgray")
		}
	}
}

// commaf returns a string from a floating point value using
// commas to separate thousands.
// (from https://github.com/dustin/go-humanize/blob/master/comma.go)
func commaf(v float64, prec int) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}

	comma := []byte{','}

	parts := strings.Split(strconv.FormatFloat(v, 'f', prec, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

// dformat returns the string representation of a float64
// according to the datafmt flag value.
// if there is no fractional portion of the float64, override the flag and
// return the string with no decimals.
func dformat(datafmt string, x float64) string {

	if datafmt != Defaultfmt {
		if strings.HasPrefix(datafmt, "%,") {
			if len(datafmt) > 2 {
				prec, _ := strconv.Atoi(datafmt[2:])
				return commaf(x, prec)
			}
			return commaf(x, -1)
		}
		return fmt.Sprintf(datafmt, x)
	}

	frac := x - float64(int(x))
	if frac == 0 {
		return fmt.Sprintf("%0.f", x)
	}
	return fmt.Sprintf(datafmt, x)
}

// datasum computes the sum of the chart data
func datasum(data []ChartData) float64 {
	sum := 0.0
	for _, d := range data {
		sum += d.value
	}
	return sum
}

// pct computs the percentage of a range of values
func pct(data []ChartData) []float64 {
	sum := 0.0
	for _, d := range data {
		sum += d.value
	}

	p := make([]float64, len(data))
	for i, d := range data {
		p[i] = (d.value / sum) * 100
	}
	return p
}

// parsecondition parses the expression low,high,color. For example "0,10,red"
// means color the data red if the value is between 0 and 10.
func parsecondition(s string) (float64, float64, string, error) {
	cs := strings.Split(s, ",")
	if len(cs) != 3 {
		return smallest, largest, "", fmt.Errorf("%s bad condition", s)
	}
	low, err := strconv.ParseFloat(cs[0], 64)
	if err != nil {
		return smallest, largest, "", err
	}
	high, err := strconv.ParseFloat(cs[1], 64)
	if err != nil {
		return smallest, largest, "", err
	}
	return low, high, cs[2], nil
}

// parsebounds returns the boundary (left, right, top, bottom) of a
// comma-separated list
func Parsebounds(s string) (float64, float64, float64, float64) {

	var err error
	var left, right, top, bottom float64

	bounds := strings.Split(s, ",")
	if len(bounds) != 4 {
		return 0, 0, 0, 0
	}
	left, err = strconv.ParseFloat(strings.TrimSpace(bounds[0]), 64)
	if err != nil {
		left = 0
	}
	right, err = strconv.ParseFloat(strings.TrimSpace(bounds[1]), 64)
	if err != nil {
		right = 0
	}
	top, err = strconv.ParseFloat(strings.TrimSpace(bounds[2]), 64)
	if err != nil {
		top = 0
	}
	bottom, err = strconv.ParseFloat(strings.TrimSpace(bounds[3]), 64)
	if err != nil {
		bottom = 0
	}
	return left, right, top, bottom
}

// pgrid makes a proportional grid with the specified rows and columns
func (s *Settings) pgrid(deck *generate.Deck, data []ChartData, title string, rows, cols int) {

	ls := s.Measures.LineSpacing
	ts := s.Measures.TextSize
	top := s.Measures.Top
	left := s.Measures.Left
	valuecolor := s.Attributes.ValueColor

	// sanity checks
	if left < 0 {
		left = 30.0
	}

	if rows*cols != 100 {
		return
	}

	sum := 0.0
	for _, d := range data {
		sum += d.value
	}
	pct := make([]float64, len(data))
	for i, d := range data {
		pct[i] = math.Floor((d.value / sum) * 100)
	}

	// encode the data in a string vector
	chars := make([]string, 100)
	cb := 0
	for k := 0; k < len(data); k++ {
		for l := 0; l < int(pct[k]); l++ {
			chars[cb] = data[k].note
			cb++
		}
	}

	// make rows and cols
	n := 0
	y := s.Measures.Top
	for i := 0; i < rows; i++ {
		x := s.Measures.Left
		for j := 0; j < cols; j++ {
			if n >= 100 {
				break
			}
			deck.Circle(x, y, ts, chars[n])
			n++
			x += ls
		}
		y -= ls
	}

	// title and legend
	if len(title) > 0 && s.Flags.ShowTitle {
		deck.Text(s.Measures.Left-ts/2, top+ts*2, title, "sans", ts*1.5, Titlecolor)
	}
	cx := (float64(cols-1) * ls) + ls/2
	df := s.Attributes.DataFmt
	for i, d := range data {
		y -= ls * 1.2
		deck.Circle(left, y, ts, d.note)
		deck.Text(left+ts, y-(ts/2), d.label+" ("+dformat(df, pct[i])+"%)", "sans", ts, "")
		if s.Flags.ShowValues {
			deck.TextEnd(left+cx, y-(ts/2), dformat(df, d.value), "sans", ts, valuecolor)
		}
	}
}

// dotgrid makes a grid 10x10 grid of dots colored by value
func dotgrid(deck *generate.Deck, x, y, left, step float64, n int, fillcolor string) (float64, float64) {
	edge := (((step * 0.3) + step) * 7) + left
	for i := 0; i < n; i++ {
		if x > edge {
			x = left
			y -= step
		}
		deck.Circle(x, y, 2*step*0.3, fillcolor)
		deck.Rect(x, y, step*0.9, step*0.9, fillcolor, 30)
		x += step
	}
	return x, y
}

// lego makes lego charts (a variation of pgrid)
func (s *Settings) lego(deck *generate.Deck, data []ChartData, title string) {
	left := s.Measures.Left
	x := left
	y := s.Measures.Top
	step := s.Measures.TextSize

	if len(title) > 0 && s.Flags.ShowTitle {
		deck.Text(left-step/2, y+step*2, title, "sans", step, Titlecolor)
	}
	sum := 0.0
	for _, d := range data {
		sum += d.value
	}
	for _, d := range data {
		pct := (d.value / sum) * 100
		v := int(math.Round(pct))
		px, py := dotgrid(deck, x, y, left, step, v, d.note)
		x = px
		y = py
	}
	y -= step * 2
	for _, d := range data {
		pct := (d.value / sum) * 100
		v := int(math.Round(pct))
		deck.Circle(left, y, 2*step*0.3, d.note)
		deck.Text(left+step, y-step*0.2, fmt.Sprintf("%s (%.d%%)", d.label, v), "sans", step*0.5, "")
		y -= step
	}
}

// polar converts polar to Cartesian coordinates
func polar(x, y, r, t float64) (float64, float64) {
	px := x + r*math.Cos(t)
	py := y + r*math.Sin(t)
	return px, py
}

// cpolar converts polar to Cartesion coordinates, compensating for the canvas aspect ratio
func cpolar(x, y, r, t, w, h float64) (float64, float64) {
	px := x + (r * math.Cos(t))
	ry := r * (w / h)
	py := y + (ry * math.Sin(t))
	return px, py
}

// spokes draws the points and lines like spokes on a wheel
func spokes(deck *generate.Deck, cx, cy, r, spokesize, w, h float64, n int, color string) {
	t := topclock
	step := fullcircle / float64(n)
	for i := 0; i < n; i++ {
		px, py := cpolar(cx, cy, r, t, w, h)
		deck.Line(cx, cy, px, py, spokesize, "lightgray")
		deck.Circle(px, py, 0.5, color)
		t -= step
	}
}

// radial draws a radial plot
func (s *Settings) radial(deck *generate.Deck, data []ChartData, title string, maxd float64) {
	top := s.Measures.Top
	left := s.Measures.Left
	pwidth := s.Measures.PWidth
	psize := s.Measures.PSize
	datacolor := s.Attributes.DataColor
	umax := s.Measures.UserMax
	ts := s.Measures.TextSize

	if s.Measures.CanvasHeight == 0 || s.Measures.CanvasWidth == 0 {
		s.Measures.CanvasWidth, s.Measures.CanvasHeight = 792.9, 612.0
	}

	if left < 0 {
		left = 50.0
	}

	rw, rh := s.Measures.CanvasWidth, s.Measures.CanvasHeight

	dx := left
	dy := top
	if len(title) > 0 && s.Flags.ShowTitle {
		deck.TextMid(dx, dy, title, "sans", ts*1.5, Titlecolor)
	}
	if umax > 0 {
		maxd = umax
	}
	t := topclock
	deck.Circle(dx, dy, pwidth*2, "silver", 10)
	step := fullcircle / float64(len(data))
	var color string

	for _, d := range data {
		cv := vmap(d.value, 0, maxd, 2, psize)
		px, py := cpolar(dx, dy, pwidth, t, rw, rh)
		tx, ty := cpolar(dx, dy, pwidth+(psize/2)+(ts*2), t, rw, rh)

		if len(d.note) > 0 {
			color = d.note
		} else {
			color = datacolor
		}

		deck.TextMid(tx, ty, d.label, "sans", ts/2, "black")
		if s.Flags.ShowValues {
			deck.TextMid(px, py-ts/3, dformat(s.Attributes.DataFmt, d.value), "mono", ts, s.Attributes.ValueColor)
		}
		if s.Flags.ShowSpokes {
			spokes(deck, px, py, psize/2, 0.05, rw, rh, int(d.value), color)
		} else {
			deck.Circle(px, py, cv, color, transparency)
			deck.Line(tx, ty, px, py, 0.05, "gray", 50)
		}
		t -= step
	}
}

// Slopechart draws a slope chart
func (s *Settings) Slopechart(deck *generate.Deck, r io.ReadCloser) {
	data, mindata, maxdata, title := Getdata(r, s.Flags.ReadCSV, s.Attributes.CSVCols)
	if len(data) < 2 {
		fmt.Fprintf(os.Stderr, "slope graphs need at least two data points")
		return
	}

	datamin := s.Flags.DataMinimum
	umin := s.Measures.UserMin
	umax := s.Measures.UserMax
	top := s.Measures.Top
	bottom := s.Measures.Bottom
	left := s.Measures.Left
	right := s.Measures.Right

	if !datamin {
		mindata = 0
	}

	if umin >= 0 {
		mindata = umin
	}

	var Showslopemax bool
	if umax >= 0 && umax > mindata {
		maxdata = umax
		Showslopemax = true
	}

	chartitle := s.Attributes.ChartTitle
	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}

	bgcolor := s.Attributes.BackgroundColor
	labelcolor := s.Attributes.LabelColor
	datacolor := s.Attributes.DataColor
	valuecolor := s.Attributes.ValueColor
	linewidth := s.Measures.LineWidth
	ts := s.Measures.TextSize

	if s.Flags.FullDeck {
		deck.StartSlide(bgcolor)
	}

	lw := linewidth / 2
	lsize := ts * 0.75
	tsize := ts * 1.5
	w := right - left
	h := top - bottom
	if len(title) > 0 && s.Flags.ShowTitle {
		hsize := ts * 2
		deck.Text(left, top+10, title, "sans", hsize, Titlecolor)
	}

	// these are magical
	hskip := w * .60
	vskip := h * 1.4

	x1 := left
	x2 := right
	// Process the data in pairs
	for i := 0; i < len(data)-1; i += 2 {
		if len(data[i].label) > 0 {
			deck.TextMid(x1+(w/2), top+3, data[i].note, "sans", tsize, labelcolor)
		}
		v1 := data[i].value
		v2 := data[i+1].value
		v1y := vmap(v1, mindata, maxdata, bottom, top)
		v2y := vmap(v2, mindata, maxdata, bottom, top)
		deck.Line(x1, bottom, x1, top, lw, "black")
		deck.Line(x2, bottom, x2, top, lw, "black")
		deck.Circle(x1, v1y, ts, datacolor)
		deck.Circle(x2, v2y, ts, datacolor)
		deck.Line(x1, v1y, x2, v2y, linewidth, datacolor)
		deck.TextMid(x1, bottom-2, data[i].label, "sans", ts, labelcolor)
		deck.TextMid(x2, bottom-2, data[i+1].label, "sans", ts, labelcolor)

		// only Show max value id user-specified
		df := s.Attributes.DataFmt
		if Showslopemax {
			deck.TextEnd(x1-1, top, dformat(df, maxdata), "sans", lsize, labelcolor)
		}
		deck.TextEnd(x1-1, v1y, dformat(df, v1), "sans", lsize, valuecolor)
		deck.Text(x2+1, v2y, dformat(df, v2), "sans", lsize, valuecolor)
		x1 += w + hskip
		x2 += w + hskip
		if x2 > 100 {
			x1 = left
			x2 = right
			top -= vskip
			bottom -= vskip
		}
	}
	if s.Flags.FullDeck {
		deck.EndSlide()
	}
}

// pmap draws a porpotional map
func (s *Settings) pmap(deck *generate.Deck, data []ChartData, title string) {
	top := s.Measures.Top
	left := s.Measures.Left
	right := s.Measures.Right
	pwidth := s.Measures.PWidth
	pmlen := s.Measures.PMapLength
	datacolor := s.Attributes.DataColor
	ts := s.Measures.TextSize

	if left < 0 {
		left = 20.0
	}
	x := left
	pl := (right - left)
	bl := pl / 100.0
	hspace := 0.10
	var ty float64
	var textcolor string
	if len(title) > 0 && s.Flags.ShowTitle {
		deck.TextMid(x+pl/2, top+(pwidth*2), title, "sans", ts*1.5, Titlecolor)
	}
	for i, p := range pct(data) {
		bx := (p * bl)
		if p < 3 || len(data[i].label) > pmlen {
			ty = top - pwidth*1.2
			deck.Line(x+(bx/2), ty+(ts*1.5), x+(bx/2), top, 0.1, Dotlinecolor)
		} else {
			ty = top
		}
		linecolor, lineop := stdcolor(i, data[i].note, datacolor, p, s.Flags.SolidPMap)
		deck.Line(x, top, bx+x, top, pwidth, linecolor, lineop)
		if lineop == 100 {
			textcolor = "white"
		} else {
			textcolor = "black"
		}

		df := s.Attributes.DataFmt
		if s.Flags.ShowValues {
			deck.TextMid(x+(bx/2), ty-pwidth, dformat(df, data[i].value), "mono", ts/2, s.Attributes.ValueColor)
		}
		deck.TextMid(x+(bx/2), ty+(pwidth), data[i].label, "sans", ts*0.75, s.LabelColor)
		deck.TextMid(x+(bx/2), ty-(ts/2), fmt.Sprintf(df+"%%", p), "sans", ts, textcolor)

		x += bx - hspace
	}
}

// stdcolor uses either the standard color (cycling through a list) or specified color and opacity
func stdcolor(i int, dcolor, color string, op float64, solid bool) (string, float64) {
	if color == "std" {
		return blue7[i%len(blue7)], 100
	}
	if len(dcolor) > 0 {
		if solid {
			return dcolor, 100
		}
		return dcolor, 40
	}
	return color, op
}

// donut makes a donut chart
func (s *Settings) donut(deck *generate.Deck, data []ChartData, title string) {
	top := s.Measures.Top
	left := s.Measures.Left
	psize := s.Measures.PSize
	pwidth := s.Measures.PWidth
	ts := s.Measures.TextSize

	if left < 0 {
		left = 50.0
	}
	a1 := 0.0
	dx := left // + (psize / 2)
	dy := top - (psize / 2)
	if len(title) > 0 && s.Flags.ShowTitle {
		deck.TextMid(dx, dy+(psize*1.2), title, "sans", s.Measures.TextSize*1.5, Titlecolor)
	}
	for i, p := range pct(data) {
		angle := (p / 100) * 360 // fullcircle
		a2 := a1 + angle
		mid := (a1 + a2) / 2

		bcolor, op := stdcolor(i, data[i].note, s.Attributes.DataColor, p, s.Flags.SolidPMap)
		deck.Arc(dx, dy, psize, psize, pwidth, a1, a2, bcolor, op)
		tx, ty := polar(dx, dy, psize*.85, mid*(math.Pi/180))
		if s.Flags.ShowValues {
			deck.TextMid(tx, ty, fmt.Sprintf("%s "+s.Attributes.DataFmt+"%%", data[i].label, p), "sans", ts, "")
		}
		a1 = a2
	}
}

const (
	topbegAngle   = 145.0 // top beginning angle
	botbegAngle   = 215.0 // bottom beginning angle
	fanspan       = 110.0 // span size of the top and bottom of the fan
	leftbegAngle  = 135.0 // left beginning angle
	rightbegAngle = 315.0 // right beginning angle
	wingspan      = 90.0  // span size of the left and right wings
	lpadding      = 10.0  // label padding
)

// data split divides a data set into two sections
func datasplit(data []ChartData) ([]ChartData, []ChartData) {
	half := len(data) / 2
	top := make([]ChartData, half)
	bottom := make([]ChartData, half)
	copy(top, data[0:half])
	copy(bottom, data[half:])
	return top, bottom
}

// fpolar to Cartesian coordinates, corrected for aspect ratio
func fpolar(cx, cy, r, theta, cw, ch float64) (float64, float64) {
	ry := r * (cw / ch)
	t := theta * (math.Pi / 180)
	return cx + (r * math.Cos(t)), cy + (ry * math.Sin(t))
}

// legend makes a balanced left and right hand legend
func legend(deck *generate.Deck, data []ChartData, orientation string, rows int, cx, cy, asize, ts float64) {
	var x, y, xoffset float64
	var alignment string
	right := len(data) % rows
	left := len(data) - right
	r := ts + 1.0
	leading := ts * 5.5

	alignment = "l"
	switch orientation {
	case "tb":
		x = cx - asize - lpadding
		y = cy + lpadding
	case "lr":
		x = cx - r
		y = cy + asize + (lpadding)
	}
	// left/top legend
	xoffset = 3
	for i := 0; i < left; i++ {
		label := data[i].label
		deck.Circle(x, y, r, data[i].note)
		legendlabel(deck, label, alignment, x+xoffset, y, ts)
		y -= leading
	}
	// right/bottom legend
	switch orientation {
	case "tb":
		x = cx + asize + lpadding
		y = cy + lpadding
		xoffset = -3
		alignment = "e"
	case "lr":
		x = cx - r
		y = cy - (asize * 0.6)
	}
	for i := left; i < len(data); i++ {
		label := data[i].label
		deck.Circle(x, y, r, data[i].note)
		legendlabel(deck, label, alignment, x+xoffset, y, ts)
		y -= leading
	}
}

// legendlabel lays out the legend labels for fan and bowtie charts
func legendlabel(deck *generate.Deck, s, alignment string, x, y, ts float64) {
	w := strings.Split(s, `\n`)
	lw := len(w)
	if lw == 1 {
		showtext(deck, x, y-(ts/3), ts, s, alignment)
	} else {
		y = y + (ts * (float64(lw / 3)))
		for i := 0; i < lw; i++ {
			showtext(deck, x, y, ts, w[i], alignment)
			y -= (ts * 1.8)
		}
	}
}

// showtext places text beginning center, or end
func showtext(deck *generate.Deck, x, y, ts float64, s, align string) {
	switch align {
	case "l", "b":
		deck.Text(x, y, s, "sans", ts, "")
	case "r", "e":
		deck.TextEnd(x, y, s, "sans", ts, "")
	case "c", "m":
		deck.TextEnd(x, y, s, "sans", ts, "")
	default:
		deck.Text(x, y, s, "sans", ts, "")
	}
}

// arclabel labels the data items
func arclabel(deck *generate.Deck, cx, cy, a1, a2, asize, value, cw, ch, ts float64) {
	v := strconv.FormatFloat(value, 'f', 1, 64)
	diff := a2 - a1
	lx, ly := fpolar(cx, cy, asize*0.9, a1+(diff*0.5), cw, ch)
	deck.TextMid(lx, ly, v+"%", "sans", ts, "")
}

// wedge makes data wedges
func wedge(deck *generate.Deck, data []ChartData, cx, cy, begAngle, asize, cw, ch, ts float64) {
	start := begAngle
	for _, d := range data {
		m := (d.value / 100) * wingspan
		a1 := start
		a2 := start + m
		deck.Arc(cx, cy, asize, asize, asize, a1, a2, d.note)
		arclabel(deck, cx, cy, a1, a2, asize, d.value, cw, ch, ts)
		start = a2
	}
}

// bowtie makes a bowtie chart
func (s *Settings) bowtie(deck *generate.Deck, data []ChartData, title string) {
	top := s.Measures.Top
	left := s.Measures.Left
	asize := s.Measures.PSize
	ts := s.TextSize

	if left < 0 {
		left = 50.0
	}
	cx := left
	cy := top - (asize / 2)
	cw := s.CanvasWidth
	ch := s.CanvasHeight

	topdata, botdata := datasplit(data)

	//var lx, ly float64
	//lx, ly = cpolar(cx, cy, asize+1, 180, cw, ch)
	//deck.TextEnd(lx, ly, "", "sans", ts, s.LabelColor)
	wedge(deck, topdata, cx, cy, leftbegAngle, asize, cw, ch, ts)
	//lx, ly = cpolar(cx, cy, asize+1, 0, cw, ch)
	//deck.TextEnd(lx, ly, "", "sans", ts, s.LabelColor)
	wedge(deck, botdata, cx, cy, rightbegAngle, asize, cw, ch, ts)

	ty := cy + (asize * 1.2)
	if s.Flags.ShowValues {
		legend(deck, topdata, "lr", 3, cx, cy, asize, ts)
		ty = 92.0
	}
	if len(title) > 0 && s.Flags.ShowTitle {
		deck.TextMid(cx, ty, title, "sans", s.Measures.TextSize*1.5, Titlecolor)
	}
}

// fan makes a fan chart
func (s *Settings) fan(deck *generate.Deck, data []ChartData, title string) {
	top := s.Measures.Top
	left := s.Measures.Left
	asize := s.Measures.PSize
	ts := s.Measures.TextSize

	if left < 0 {
		left = 50.0
	}
	cx := left
	cy := top - (asize / 2)
	cw := s.CanvasWidth
	ch := s.CanvasHeight

	topdata, botdata := datasplit(data)
	if len(title) > 0 && s.Flags.ShowTitle {
		deck.TextMid(cx, cy+(asize*1.5), title, "sans", ts*1.5, Titlecolor)
	}
	var start float64
	// the top of the fan chart
	start = topbegAngle
	for _, d := range topdata {
		m := (d.value / 100) * fanspan
		a1 := start - m
		a2 := start
		deck.Arc(cx, cy, asize, asize, asize, a1, a2, d.note)
		arclabel(deck, cx, cy, a1, a2, asize, d.value, cw, ch, ts)
		start = a1
	}
	// bottom of the fan chart
	start = botbegAngle
	for i := len(botdata) - 1; i >= 0; i-- {
		d := botdata[i]
		m := (d.value / 100) * fanspan
		a1 := start + m
		a2 := start
		deck.Arc(cx, cy, asize, asize, asize, a2, a1, d.note)
		arclabel(deck, cx, cy, a1, a2, asize, d.value, cw, ch, ts)
		start = a1
	}

	if s.Flags.ShowValues {
		legend(deck, topdata, "tb", 3, cx, cy, asize, ts)
	}
}

// Pchart draws proportional data, either a pmap, pgrid, radial or donut using input from a Reader
func (s *Settings) Pchart(deck *generate.Deck, r io.ReadCloser) {
	f := s.Flags
	data, _, maxdata, title := Getdata(r, f.ReadCSV, s.Attributes.CSVCols)
	chartitle := s.Attributes.ChartTitle
	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}
	if f.FullDeck {
		deck.StartSlide(s.Attributes.BackgroundColor)
	}
	switch {
	case f.ShowDonut:
		s.donut(deck, data, title)
	case f.ShowPMap:
		s.pmap(deck, data, title)
	case f.ShowBowtie:
		s.bowtie(deck, data, title)
	case f.ShowFan:
		s.fan(deck, data, title)
	case f.ShowPGrid:
		s.pgrid(deck, data, title, 10, 10)
	case f.ShowLego:
		s.lego(deck, data, title)
	case f.ShowRadial:
		s.radial(deck, data, title, maxdata)
	}
	if f.FullDeck {
		deck.EndSlide()
	}
}

// Wbchart makes a word bar chart
func (s *Settings) Wbchart(deck *generate.Deck, r io.ReadCloser) {
	left := s.Measures.Left
	right := s.Measures.Right
	top := s.Measures.Top
	ls := s.Measures.LineSpacing
	ts := s.Measures.TextSize
	datamin := s.Flags.DataMinimum
	chartitle := s.Attributes.ChartTitle
	Showtitle := s.Flags.ShowTitle

	if left < 0 {
		left = 20.0
	}
	hts := ts / 2
	mts := ts * 0.75
	linespacing := ts * ls

	bardata, mindata, maxdata, title := Getdata(r, s.Flags.ReadCSV, s.Attributes.CSVCols) // getdata(r)
	if !datamin {
		mindata = 0
	}
	if s.Flags.FullDeck {
		deck.StartSlide(s.Attributes.BackgroundColor)
	}

	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}

	if len(title) > 0 && Showtitle {
		deck.Text(left, top+(linespacing*1.5), title, "sans", ts*1.5, Titlecolor)
	}

	var sum float64
	if s.Flags.ShowPercentage {
		sum = datasum(bardata)
	}

	// for every name, value pair, make the chart
	y := top
	labelcolor, datacolor, valuecolor := s.Attributes.LabelColor, s.Attributes.DataColor, s.Attributes.ValueColor
	for _, data := range bardata {
		deck.Text(left+hts, y, data.label, "sans", ts, labelcolor)
		bv := vmap(data.value, mindata, maxdata, left, right)
		deck.Line(left+hts, y+hts, bv, y+hts, ts*1.5, datacolor, wbop)
		if s.Flags.ShowValues {
			df := s.Attributes.DataFmt
			if s.Flags.ShowPercentage {
				avgs := fmt.Sprintf(" ("+df+"%%)", 100*(data.value/sum))
				deck.TextEnd(left, y+(hts/2), dformat(df, data.value)+avgs, "mono", mts, valuecolor)
			} else {
				deck.TextEnd(left, y+(hts/2), dformat(df, data.value), "mono", mts, valuecolor)
			}
		}
		y -= linespacing
	}
	if s.Flags.FullDeck {
		deck.EndSlide()
	}
}

// Hchart makes horizontal bar charts using input from a Reader
func (s *Settings) Hchart(deck *generate.Deck, r io.ReadCloser) {
	ts := s.Measures.TextSize
	ls := s.Measures.LineSpacing
	left := s.Measures.Left
	right := s.Measures.Right
	top := s.Measures.Top

	hts := ts / 2
	mts := ts * 0.75
	linespacing := ts * ls

	bardata, mindata, maxdata, title := Getdata(r, s.Flags.ReadCSV, s.Attributes.CSVCols) // getdata(r)

	if left < 0 {
		left = 30.0
	}

	f := s.Flags
	if !f.DataMinimum {
		mindata = 0
	}

	bgcolor := s.Attributes.BackgroundColor
	valuecolor := s.Attributes.ValueColor
	datacolor := s.Attributes.DataColor
	labelcolor := s.Attributes.LabelColor

	if f.FullDeck {
		deck.StartSlide(bgcolor)
	}

	chartitle := s.Attributes.ChartTitle
	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}

	if len(title) > 0 && f.ShowTitle {
		deck.TextMid(50, top+(linespacing*1.5), title, "sans", ts*1.5, Titlecolor)
	}

	var sum float64
	if f.ShowPercentage {
		sum = datasum(bardata)
	}

	// for every name, value pair, make the chart
	y := top

	for _, data := range bardata {
		label := nlmap.Replace(data.label) // replace '\n' with spaces
		deck.TextEnd(left-hts, y+(hts/2), label, "sans", ts, labelcolor)
		bv := vmap(data.value, mindata, maxdata, left, right)
		if f.ShowDot {
			dottedhline(deck, left, y+hts, bv-left, ts/5, 1, 0.25, Dotlinecolor)
			deck.Circle(bv, y+hts, mts, datacolor)
		} else {
			bw := ts
			barw := s.Measures.BarWidth
			if barw > 0 {
				bw = barw
			}
			deck.Line(left, y+hts, bv, y+hts, bw, datacolor)
		}
		if f.ShowValues {
			df := s.Attributes.DataFmt
			if f.ShowPercentage {
				avgs := fmt.Sprintf(" ("+df+"%%)", 100*(data.value/sum))
				deck.Text(bv+hts, y+(hts/2), dformat(df, data.value)+avgs, "mono", mts, valuecolor)
			} else {
				deck.Text(bv+hts, y+(hts/2), dformat(df, data.value), "mono", mts, valuecolor)
			}
		}
		y -= linespacing
	}
	if f.FullDeck {
		deck.EndSlide()
	}
}

// Vchart makes charts using input from a Reader
// the types of charts are bar (column), dot, line, and volume
func (s *Settings) Vchart(deck *generate.Deck, r io.ReadCloser) {
	chartdata, mindata, maxdata, title := Getdata(r, s.Flags.ReadCSV, s.Attributes.CSVCols) // getdata(r)

	left := s.Measures.Left
	right := s.Measures.Right
	top := s.Measures.Top
	bottom := s.Measures.Bottom
	umin := s.Measures.UserMin
	umax := s.Measures.UserMax
	barw := s.Measures.BarWidth
	ls := s.Measures.LineSpacing
	ts := s.Measures.TextSize

	datamin := s.Flags.DataMinimum
	showvolume := s.Flags.ShowVolume
	showrline := s.Flags.ShowRegressionLine
	showframe := s.Flags.ShowFrame
	showaxis := s.Flags.ShowAxis
	showline := s.Flags.ShowLine
	showdot := s.Flags.ShowDot
	showtitle := s.Flags.ShowTitle
	showpct := s.Flags.ShowPercentage
	showscatter := s.Flags.ShowScatter
	showbar := s.Flags.ShowBar
	showval := s.Flags.ShowValues
	shownote := s.Flags.ShowNote

	valuecolor := s.Attributes.ValueColor
	valpos := s.Attributes.ValuePosition
	noteloc := s.Attributes.NoteLocation
	hline := s.Attributes.HLine
	datacolor := s.Attributes.DataColor
	labelcolor := s.Attributes.LabelColor
	framecolor := s.Attributes.FrameColor
	linewidth := s.Measures.LineWidth

	if left < 0 {
		left = 10.0
	}

	if !datamin {
		mindata = 0
	}

	if umin >= 0 {
		mindata = umin
	}

	if umax >= 0 && umax > mindata {
		maxdata = umax
	}

	l := len(chartdata)
	dlen := float64(l - 1)

	// define the width of bars
	var dw = (right-left)/dlen - 1
	if barw > 0 && barw <= dw {
		dw = barw
	}

	// for volume plots, allocate, fill in the extrema
	var xvol, yvol []float64
	if showvolume {
		xvol = make([]float64, l+2)
		yvol = make([]float64, l+2)
		xvol[0] = left
		yvol[0] = bottom
		xvol[l+1] = right
		yvol[l+1] = bottom
	}

	var xreg, yreg []float64
	if showrline {
		xreg = make([]float64, l)
		yreg = make([]float64, l)
	}

	linespacing := ts * ls
	spacing := ts * 1.5

	bgcolor := s.Attributes.BackgroundColor
	if s.Flags.FullDeck {
		deck.StartSlide(bgcolor)
	}

	// Show a frame if specified
	if showframe {
		fw := right - left
		fh := top - bottom
		deck.Rect(left+(fw/2), bottom+(fh/2), fw, fh, framecolor, 5)
	}

	chartitle := s.Attributes.ChartTitle
	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}

	if len(title) > 0 && showtitle {
		deck.TextMid(left+((right-left)/2), top+(linespacing*1.5), title, "sans", spacing, Titlecolor)
	}

	if showaxis {
		// yaxis(deck, left-spacing-(dw*0.5), mindata, maxdata)
		s.yaxis(deck, left-spacing, mindata, maxdata)
	}

	if len(hline) > 0 {
		var hl float64
		var hs string
		fmt.Sscanf(hline, "%f,%s", &hl, &hs)
		hy := vmap(hl, mindata, maxdata, bottom, top)
		deck.Line(left, hy, right, hy, 0.1, valuecolor, 50)
		if len(hs) > 0 {
			deck.Text(right+ts/2, hy-ts/4, hs, "serif", ts*0.75, labelcolor)
		}
	}

	var clow, chigh float64
	var cerr error
	var condcolor string
	datacond := s.Attributes.DataCondition
	if len(datacond) > 0 {
		clow, chigh, condcolor, cerr = parsecondition(datacond)
		if cerr != nil {
			fmt.Fprintf(os.Stderr, "%v\n", cerr)
			return
		}
	}

	var sum float64
	if showpct {
		sum = datasum(chartdata)
	}

	// for every name, value pair, make the chart elements
	var px, py float64
	var defcolor = datacolor
	for i, data := range chartdata {
		x := vmap(float64(i), 0, dlen, left, right)
		y := vmap(data.value, mindata, maxdata, bottom, top)

		if showvolume {
			xvol[i+1] = x
			yvol[i+1] = y
		}

		if showrline {
			xreg[i] = float64(i)
			yreg[i] = data.value
		}

		if len(datacond) > 0 {
			if data.value <= chigh && data.value >= clow {
				datacolor = condcolor
			} else {
				datacolor = defcolor
			}
		}
		if showline && i > 0 {
			deck.Line(px, py, x, y, linewidth, datacolor)
		}

		if showdot {
			dottedvline(deck, x, bottom, y, ts/6, 1, Dotlinecolor)
			deck.Circle(x, y, ts*.6, datacolor)
		}

		if showscatter {
			deck.Circle(x, y, ts*.6, datacolor)
		}

		if showbar {
			deck.Line(x, bottom, x, y, dw, datacolor)
		}

		if showval {
			yv := y + ts
			switch valpos {
			case "t":
				if data.value < 0 {
					yv = y - ts
				} else {
					yv = y + ts
				}
			case "b":
				yv = bottom + ts
			case "m":
				yv = y - ((y - bottom) / 2)
			}
			df := s.Attributes.DataFmt
			if showpct {
				avgs := fmt.Sprintf(" ("+df+"%%)", 100*(data.value/sum))
				deck.TextMid(x, yv, dformat(df, data.value)+avgs, "sans", ts*0.75, valuecolor)
			} else {
				deck.TextMid(x, yv, dformat(df, data.value), "sans", ts*0.75, valuecolor)
			}
		}
		if len(data.note) > 0 && shownote {
			xoffset := ts / 2
			yoffset := ts / 2
			notesize := ts * 0.75
			switch noteloc {
			case "l", "b":
				deck.Text(x+xoffset, y, data.note, "serif", notesize, labelcolor)
			case "r", "e":
				deck.TextEnd(x-xoffset, y, data.note, "serif", notesize, labelcolor)
			case "c":
				deck.TextMid(x, y+yoffset, data.note, "serif", notesize, labelcolor)
			default:
				deck.TextMid(x, y+yoffset, data.note, "serif", notesize, labelcolor)
			}
		}
		// Show x label every xinit times, Show the last, if specified
		xint := s.Measures.XLabelInterval
		if xint > 0 && (i%xint == 0 || (s.Flags.ShowXLast && i == l-1)) {
			xlabels := strings.Split(data.label, `\n`)
			xly := bottom - (ts * 2)
			for _, xl := range xlabels {
				if s.Flags.ShowXstagger && (i+1)%2 == 0 {
					xly -= (ts * 2)
				}
				if s.Measures.XLabelRotation == 0 {
					deck.TextMid(x, xly, xl, "sans", ts*0.8, labelcolor)
				} else {
					deck.TextRotate(x, xly, xl, "", "sans", s.Measures.XLabelRotation, ts*0.8, labelcolor)
				}
				xly -= ts * 1.2
			}
		}
		px = x
		py = y
	}
	if showvolume {
		deck.Polygon(xvol, yvol, datacolor, s.Measures.VolumeOpacity)
	}

	if showrline {
		s.Measures.rline(deck, xreg, yreg, mindata, maxdata, s.Attributes.RegressionLineColor)
	}

	if s.Flags.FullDeck {
		deck.EndSlide()
	}
}

// mean computes the arithmetic mean of a set of data
func mean(x []float64) float64 {
	sum := 0.0
	n := len(x)
	for i := 0; i < n; i++ {
		sum += x[i]
	}
	return sum / float64(n)
}

// slope computes the slope (m, b) of a set of x, y points
func slope(x, y []float64) (float64, float64) {
	n := len(x) // assume x and y have the same length
	xy := make([]float64, n)
	for i := 0; i < n; i++ {
		xy[i] = x[i] * y[i]
	}
	sqx := make([]float64, n)
	for i := 0; i < n; i++ {
		sqx[i] = x[i] * x[i]
	}
	meanxy := mean(xy)
	meanx := mean(x)
	meany := mean(y)
	meanxsq := mean(sqx)

	rise := (meanxy - (meanx * meany))
	run := (meanxsq - (meanx * meanx))
	m := rise / run
	b := meany - (m * meanx)
	return m, b
}

// rline makes a regression line
func (measures *Measures) rline(deck *generate.Deck, x, y []float64, mindata, maxdata float64, color string) {
	top := measures.Top
	left := measures.Left
	if left < 0 {
		left = 10.0
	}
	bottom := measures.Bottom
	right := measures.Right
	lw := measures.LineWidth
	m, b := slope(x, y)
	dl := len(x) - 1
	l := float64(dl)
	x1 := x[0]
	x2 := x[dl]
	y1 := m*x1 + b
	y2 := m*x2 + b
	rx1 := vmap(x1, 0, l, left, right)
	rx2 := vmap(x2, 0, l, left, right)
	ry1 := vmap(y1, mindata, maxdata, bottom, top)
	ry2 := vmap(y2, mindata, maxdata, bottom, top)
	deck.Line(rx1, ry1, rx2, ry2, lw, color)
}

// GenerateChart makes charts according to the orientation:
// horizontal bar or line, bar, dot, or donut volume charts
func (s *Settings) GenerateChart(deck *generate.Deck, r io.ReadCloser) {
	f := s.Flags
	switch {
	case f.ShowHBar:
		s.Hchart(deck, r)
	case f.ShowWBar:
		s.Wbchart(deck, r)
	case f.ShowDonut, f.ShowPMap, f.ShowPGrid, f.ShowRadial, f.ShowLego, f.ShowFan, f.ShowBowtie:
		s.Pchart(deck, r)
	case f.ShowSlope:
		s.Slopechart(deck, r)
	default:
		s.Vchart(deck, r)
	}
}

// Write performs chart I/O
func (s *Settings) Write(w io.Writer, r io.ReadCloser) {
	s.GenerateChart(generate.NewSlides(w, 0, 0), r)
}

// NewFullChart initializes the settings required to make a chart
// that includes the enclosing deck markup
func NewFullChart(chartType string, top, bottom, left, right float64) Settings {
	s := NewChart(chartType, top, bottom, left, right)
	s.Flags.FullDeck = true
	return s
}

// NewChart initializes the settings required to make a chart
// chartType may be one of: "line", "slope", "bar", "wbar", "hbar",
// "volume, "scatter", "donut", "pmap", "pgrid", "lego", "radial", "bowtie", "fan"
func NewChart(chartType string, top, bottom, left, right float64) Settings {
	var s Settings

	switch chartType {
	case "bar":
		s.Flags.ShowBar = true
	case "wbar":
		s.Flags.ShowWBar = true
	case "hbar":
		s.Flags.ShowHBar = true
	case "donut":
		s.Flags.ShowDonut = true
	case "bowtie":
		s.Flags.ShowBowtie = true
	case "fan":
		s.Flags.ShowFan = true
	case "pmap":
		s.Flags.ShowPMap = true
	case "pgrid":
		s.Flags.ShowPGrid = true
	case "lego":
		s.Flags.ShowLego = true
	case "radial":
		s.Flags.ShowRadial = true
	case "line":
		s.Flags.ShowLine = true
	case "scatter":
		s.Flags.ShowScatter = true
	case "volume", "area":
		s.Flags.ShowVolume = true
	case "slope":
		s.Flags.ShowSlope = true
	}
	if left <= 0 {
		left = 10
	}
	if right <= 0 {
		right = 90
	}
	if top <= 0 {
		top = 90
	}
	if bottom <= 0 {
		bottom = 30
	}
	s.Measures.Left = left
	s.Measures.Right = right
	s.Measures.Top = top
	s.Measures.Bottom = bottom
	s.Measures.XLabelInterval = 1
	s.Measures.TextSize = 1.5
	s.Measures.LineSpacing = 2.4
	s.Measures.CanvasWidth = 792
	s.Measures.CanvasHeight = 612

	s.Attributes.BackgroundColor = "white"
	s.Attributes.DataColor = "lightsteelblue"
	s.Attributes.LabelColor = "rgb(75,75,75)"

	return s
}
