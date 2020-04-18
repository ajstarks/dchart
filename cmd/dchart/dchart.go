// dchart - make charts in the deck format
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ajstarks/dchart"
	"github.com/ajstarks/deck/generate"
)

func cmdflags() dchart.Settings {
	var chart dchart.Settings

	// Measures
	flag.Float64Var(&chart.Measures.TextSize, "textsize", 1.5, "text size")
	flag.Float64Var(&chart.Left, "left", -1.0, "left margin") // default set to out of bounds because different charts need individual defaults
	flag.Float64Var(&chart.Right, "right", 90.0, "right margin")
	flag.Float64Var(&chart.Top, "top", 80.0, "top of the plot")
	flag.Float64Var(&chart.Bottom, "bottom", 30.0, "bottom of the plot")
	flag.Float64Var(&chart.LineSpacing, "ls", 2.4, "ls")
	flag.Float64Var(&chart.BarWidth, "barwidth", 0, "barwidth")
	flag.Float64Var(&chart.UserMin, "min", -1, "minimum")
	flag.Float64Var(&chart.UserMax, "max", -1, "maximum")
	flag.Float64Var(&chart.PSize, "psize", 40.0, "size of the donut")
	flag.Float64Var(&chart.PWidth, "pwidth", chart.Measures.TextSize*3, "width of the pmap/donut/radial")
	flag.Float64Var(&chart.LineWidth, "linewidth", 0.2, "width of line for line charts")
	flag.Float64Var(&chart.VolumeOpacity, "volop", 50, "volume opacity")
	flag.Float64Var(&chart.XLabelRotation, "xlabrot", 0, "xlabel rotation (degrees)")
	flag.IntVar(&chart.XLabelInterval, "xlabel", 1, "x axis label interval (show every n labels, 0 to show no labels)")
	flag.IntVar(&chart.PMapLength, "pmlen", 20, "pmap label length")

	// Flags (On/Off)
	flag.BoolVar(&chart.ShowBar, "bar", true, "show a bar chart")
	flag.BoolVar(&chart.ShowDot, "dot", false, "show a dot chart")
	flag.BoolVar(&chart.ShowVolume, "vol", false, "show a volume chart")
	flag.BoolVar(&chart.ShowDonut, "donut", false, "show a donut chart")
	flag.BoolVar(&chart.ShowPMap, "pmap", false, "show a proportional map")
	flag.BoolVar(&chart.ShowLine, "line", false, "show a line chart")
	flag.BoolVar(&chart.ShowHBar, "hbar", false, "show a horizontal bar chart")
	flag.BoolVar(&chart.ShowValues, "val", true, "show data values")
	flag.BoolVar(&chart.ShowAxis, "yaxis", false, "show y axis")
	flag.BoolVar(&chart.ShowSlope, "slope", false, "show a slope graph")
	flag.BoolVar(&chart.ShowTitle, "title", true, "show title")
	flag.BoolVar(&chart.ShowGrid, "grid", false, "show y axis grid")
	flag.BoolVar(&chart.ShowScatter, "scatter", false, "show scatter chart")
	flag.BoolVar(&chart.ShowRadial, "radial", false, "show a radial chart")
	flag.BoolVar(&chart.ShowSpokes, "spokes", false, "show spokes on radial charts")
	flag.BoolVar(&chart.ShowPGrid, "pgrid", false, "show proportional grid")
	flag.BoolVar(&chart.ShowNote, "note", true, "show annotations")
	flag.BoolVar(&chart.ShowFrame, "frame", false, "show frame")
	flag.BoolVar(&chart.ShowRegressionLine, "rline", false, "show regression line")
	flag.BoolVar(&chart.ShowXLast, "xlast", false, "show the last label")
	flag.BoolVar(&chart.ShowXstagger, "xstagger", false, "stagger x axis labels")
	flag.BoolVar(&chart.FullDeck, "fulldeck", true, "generate full markup")
	flag.BoolVar(&chart.DataMinimum, "dmin", false, "zero minimum")
	flag.BoolVar(&chart.ReadCSV, "csv", false, "read CSV data")
	flag.BoolVar(&chart.ShowWBar, "wbar", false, "show word bar chart")
	flag.BoolVar(&chart.ShowPercentage, "pct", false, "show computed percentages with values")
	flag.BoolVar(&chart.SolidPMap, "solidpmap", false, "solid pmap colors")

	// Attributes
	flag.StringVar(&chart.ChartTitle, "chartitle", "", "specify the title (overiding title in the data)")
	flag.StringVar(&chart.CSVCols, "csvcol", "", "label,value from the CSV header")
	flag.StringVar(&chart.ValuePosition, "valpos", "t", "value position (t=top, b=bottom, m=middle)")
	flag.StringVar(&chart.LabelColor, "lcolor", "rgb(75,75,75)", "label color")
	flag.StringVar(&chart.DataColor, "color", "lightsteelblue", "data color")
	flag.StringVar(&chart.ValueColor, "vcolor", "rgb(127,0,0)", "value color")
	flag.StringVar(&chart.RegressionLineColor, "rlcolor", "rgb(127,0,0)", "regression line color")
	flag.StringVar(&chart.FrameColor, "framecolor", "rgb(127,127,127)", "framecolor")
	flag.StringVar(&chart.BackgroundColor, "bgcolor", "white", "background color")
	flag.StringVar(&chart.DataFmt, "datafmt", dchart.Defaultfmt, "data format")
	flag.StringVar(&chart.YAxisR, "yrange", "", "y-axis range (min,max,step)")
	flag.StringVar(&chart.HLine, "hline", "", "horizontal line value,label")
	flag.StringVar(&chart.NoteLocation, "noteloc", "c", "note location (c-center, r-right aligned, l-left aligned)")
	flag.StringVar(&chart.DataCondition, "datacond", "", "data condition: low,high,color")
	flag.Parse()

	return chart
}
func main() {
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
			r.Close()
		}
	} else {
		chart.GenerateChart(deck, os.Stdin)
	}
	if fulldeck {
		deck.EndDeck()
	}
}
