// dchart - make charts in the deck format
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ajstarks/dchart"
	"github.com/ajstarks/deckgen"
)

var usageMsg = `
dchart [options] file..

Options categories with defaults and descriptions:

Chart Types
.......................................................................
-bar        true                      bar chart
-wbar       false                     word bar chart
-hbar       false                     horizontal bar chart
-donut      false                     donut chart
-dot        false                     dot chart
-lego       false                     lego chart
-line       false                     line chart
-pgrid      false                     proportional grid
-pmap       false                     proportional map
-bowtie     false                     bowtie chart
-fan        false                     fan chart
-radial     false                     radial chart
-scatter    false                     scatter chart
-slope      false                     slope chart
-vol        false                     volume (area) chart


Chart Elements
.......................................................................
-csv        false                     read CSV files
-frame      false                     show a colored frame
-fulldeck   true                      generate full deck markup
-grid       false                     show gridlines on the y axis
-note       true                      show annotations
-pct        false                     show computed percentage
-rline      false                     show a regression line
-solidpmap  false                     show solid pmap colors
-spokes     false                     show spokes in radial chart
-title      true                      show the title
-val        true                      show values
-xlast      false                     show the last x label
-xstagger   false                     stagger x axis labels
-yaxis      false                     show a y axis
-chartitle  override title in data    specify the title
-datacond   low,high,colors           conditional data colors
-hline      value,label2              label horizontal line at value
-valpos     t=top, b=bottom, m=middle value position
-xlabel     default=1, 0 to suppress  x axis label interval
-yrange     min,max.step              specify the y axis label range


Position and Scaling
.......................................................................
-left       20                        left margin
-right      80                        right margin
-top        80                        top of the chart
-bottom     30                        bottom of the chart
-min        data min                  set the minimum data value
-max        data max                  set the maximum data value
-bounds     ""                        set left,right,top,bottom


Measures and Attributes
.......................................................................
-bgcolor    white                     background color
-barwidth   computed from data size   barwidth
-color      lightsteelblue            data color
-csvcol     labe1,label2              specify csv columns
-datafmt    %.1f                      format for values (%f or %,)
-dmin       false                     use data minimum, not zero
-framecolor rgb(127,127,127)          frame color
-lcolor     rgb(75,75,75)             label color
-linewidth  0.20                      linewidth
-ls         2.4                       linespacing
-noteloc    c=center, r=right, l=left annotation location
-pmlen      20                        pmap label length
-psize      30                        diameter of the donut
-pwidth     30                        width of the donut or pmap
-rlcolor    rgb(127,0,0)              regression line color
-textsize   1.50                      text size
-xlabrot    0                         xlabel rotation (deg.)
-vcolor     rgb(127,0,0)              value color
-volop      50                        volume opacity %
`

func printusage() {
	fmt.Fprintln(flag.CommandLine.Output(), usageMsg)
}

func cmdflags() dchart.Settings {
	var chart dchart.Settings

	// Measures
	flag.Float64Var(&chart.TextSize, "textsize", 1.5, "text size")
	flag.Float64Var(&chart.CanvasWidth, "cw", 792, "canvas width")
	flag.Float64Var(&chart.CanvasHeight, "ch", 612, "canvas height")
	flag.Float64Var(&chart.Left, "left", -1, "left margin")
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
	flag.StringVar(&chart.Boundary, "bounds", "", "chart boundary (left,right,top,bottom)")

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
	flag.BoolVar(&chart.ShowLego, "lego", false, "show lego chart")
	flag.BoolVar(&chart.ShowBowtie, "bowtie", false, "show bowtie chart")
	flag.BoolVar(&chart.ShowFan, "fan", false, "show fan chart")
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
	flag.Usage = printusage
	flag.Parse()
	if len(chart.Boundary) > 0 {
		chart.Left, chart.Right, chart.Top, chart.Bottom = dchart.Parsebounds(chart.Boundary)
	}

	return chart
}
func main() {
	chart := cmdflags()
	fulldeck := chart.Flags.FullDeck

	deck := deckgen.NewSlides(os.Stdout, 0, 0)
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
