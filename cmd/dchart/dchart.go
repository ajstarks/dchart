// dchart - make charts in the deck format
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ajstarks/dchart"
	"github.com/ajstarks/deck/generate"
)

// command line options
func cmdflags() dchart.Settings {

	var s dchart.Settings

	// Measures
	m := s.Measures
	flag.Float64Var(&m.TextSize, "textsize", 1.5, "text size")
	flag.Float64Var(&m.Left, "left", -1.0, "left margin") // default set to out of bounds because different charts need individual defaults
	flag.Float64Var(&m.Right, "right", 90.0, "right margin")
	flag.Float64Var(&m.Top, "top", 80.0, "top of the plot")
	flag.Float64Var(&m.Bottom, "bottom", 30.0, "bottom of the plot")
	flag.Float64Var(&m.LineSpacing, "ls", 2.4, "ls")
	flag.Float64Var(&m.BarWidth, "barwidth", 0, "barwidth")
	flag.Float64Var(&m.UserMin, "min", -1, "minimum")
	flag.Float64Var(&m.UserMax, "max", -1, "maximum")
	flag.Float64Var(&m.PSize, "psize", 40.0, "size of the donut")
	flag.Float64Var(&m.PWidth, "pwidth", m.TextSize*3, "width of the pmap/donut/radial")
	flag.Float64Var(&m.LineWidth, "linewidth", 0.2, "width of line for line charts")
	flag.Float64Var(&m.VolumeOpacity, "volop", 50, "volume opacity")
	flag.Float64Var(&m.XLabelRotation, "xlabrot", 0, "xlabel rotation (degrees)")
	flag.IntVar(&m.XLabelInterval, "xlabel", 1, "x axis label interval (show every n labels, 0 to show no labels)")
	flag.IntVar(&m.PMapLength, "pmlen", 20, "pmap label length")

	// Flags (On/Off)
	f := s.Flags
	flag.BoolVar(&f.ShowBar, "bar", true, "show a bar chart")
	flag.BoolVar(&f.ShowDot, "dot", false, "show a dot chart")
	flag.BoolVar(&f.ShowVolume, "vol", false, "show a volume chart")
	flag.BoolVar(&f.ShowDonut, "donut", false, "show a donut chart")
	flag.BoolVar(&f.ShowPMap, "pmap", false, "show a proportional map")
	flag.BoolVar(&f.ShowLine, "line", false, "show a line chart")
	flag.BoolVar(&f.ShowHBar, "hbar", false, "show a horizontal bar chart")
	flag.BoolVar(&f.ShowValues, "val", true, "show data values")
	flag.BoolVar(&f.ShowAxis, "yaxis", false, "show y axis")
	flag.BoolVar(&f.ShowSlope, "slope", false, "show a slope graph")
	flag.BoolVar(&f.ShowTitle, "title", true, "show title")
	flag.BoolVar(&f.ShowGrid, "grid", false, "show y axis grid")
	flag.BoolVar(&f.ShowScatter, "scatter", false, "show scatter chart")
	flag.BoolVar(&f.ShowRadial, "radial", false, "show a radial chart")
	flag.BoolVar(&f.ShowSpokes, "spokes", false, "show spokes on radial charts")
	flag.BoolVar(&f.ShowPGrid, "pgrid", false, "show proportional grid")
	flag.BoolVar(&f.ShowNote, "note", true, "show annotations")
	flag.BoolVar(&f.ShowFrame, "frame", false, "show frame")
	flag.BoolVar(&f.ShowRegressionLine, "rline", false, "show regression line")
	flag.BoolVar(&f.ShowXLast, "xlast", false, "show the last label")
	flag.BoolVar(&f.ShowXstagger, "xstagger", false, "stagger x axis labels")
	flag.BoolVar(&f.FullDeck, "fulldeck", true, "generate full markup")
	flag.BoolVar(&f.DataMinimum, "dmin", false, "zero minimum")
	flag.BoolVar(&f.ReadCSV, "csv", false, "read CSV data")
	flag.BoolVar(&f.ShowWBar, "wbar", false, "show word bar chart")
	flag.BoolVar(&f.ShowPercentage, "pct", false, "show computed percentages with values")
	flag.BoolVar(&f.SolidPMap, "solidpmap", false, "solid pmap colors")

	// Attributes
	a := s.Attributes
	flag.StringVar(&a.ChartTitle, "chartitle", "", "specify the title (overiding title in the data)")
	flag.StringVar(&a.CSVCols, "csvcol", "", "label,value from the CSV header")
	flag.StringVar(&a.ValuePosition, "valpos", "t", "value position (t=top, b=bottom, m=middle)")
	flag.StringVar(&a.LabelColor, "lcolor", "rgb(75,75,75)", "label color")
	flag.StringVar(&a.DataColor, "color", "lightsteelblue", "data color")
	flag.StringVar(&a.ValueColor, "vcolor", "rgb(127,0,0)", "value color")
	flag.StringVar(&a.RegressionLineColor, "rlcolor", "rgb(127,0,0)", "regression line color")
	flag.StringVar(&a.FrameColor, "framecolor", "rgb(127,127,127)", "framecolor")
	flag.StringVar(&a.BackgroundColor, "bgcolor", "white", "background color")
	flag.StringVar(&a.DataFmt, "datafmt", dchart.Defaultfmt, "data format")
	flag.StringVar(&a.YAxisR, "yrange", "", "y-axis range (min,max,step)")
	flag.StringVar(&a.HLine, "hline", "", "horizontal line value,label")
	flag.StringVar(&a.NoteLocation, "noteloc", "c", "note location (c-center, r-right aligned, l-left aligned)")
	flag.StringVar(&a.DataCondition, "datacond", "", "data condition: low,high,color")

	flag.Parse()
	return s
}

func main() {
	// process command line options,
	// start the deck, for every file name make a slide.
	// Read from standard input, if no files are specified.
	chart := cmdflags()
	fulldeck := chart.Flags.FullDeck
	deck := generate.NewSlides(os.Stdout, 0, 0)
	if fulldeck {
		deck.StartDeck()
	}
	if len(flag.Args()) > 0 {
		for _, file := range flag.Args() {
			r, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			chart.GenerateChart(deck, r)
		}
	} else {
		chart.GenerateChart(deck, os.Stdin)
	}
	if fulldeck {
		deck.EndDeck()
	}
}
