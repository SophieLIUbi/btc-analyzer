package visualizer

import (
	"btc-analyzer/internal/types"
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// ChartConfig holds configuration for chart generation
type ChartConfig struct {
	Width       int
	Height      int
	Title       string
	XLabel      string
	YLabel      string
	ShowGrid    bool
	ShowLegend  bool
	LineWidth   vg.Length
	FontSize    vg.Length
	Theme       string
}

// DefaultChartConfig returns default chart configuration
func DefaultChartConfig() ChartConfig {
	return ChartConfig{
		Width:      1000,
		Height:     600,
		Title:      "Bitcoin Technical Indicators",
		XLabel:     "Time",
		YLabel:     "Value",
		ShowGrid:   true,
		ShowLegend: true,
		LineWidth:  vg.Points(2),
		FontSize:   vg.Points(12),
		Theme:      "default",
	}
}

// writeBuffer implements io.Writer for byte slice
type writeBuffer struct {
	buf *[]byte
}

func (wb *writeBuffer) Write(p []byte) (n int, err error) {
	*wb.buf = append(*wb.buf, p...)
	return len(p), nil
}

// DrawTechnicalIndicatorsChart creates a chart with RSI and MACD indicators
func DrawTechnicalIndicatorsChart(bts *types.BTCTimeSeries, analytics types.BTCAnalytics, config ChartConfig) ([]byte, error) {
	if len(bts.Data) == 0 {
		return nil, fmt.Errorf("no data to plot")
	}

	p := plot.New()
	p.Title.Text = config.Title
	p.X.Label.Text = config.XLabel
	p.Y.Label.Text = config.YLabel

	// Add grid
	if config.ShowGrid {
		p.Add(plotter.NewGrid())
	}

	// Plot RSI if available (0-100 scale)
	if len(analytics.RSI) > 0 {
		rsiLine, err := plotter.NewLine(makeSimpleXYs(analytics.RSI))
		if err == nil {
			rsiLine.LineStyle.Color = color.RGBA{R: 150, G: 0, B: 150, A: 255}
			rsiLine.LineStyle.Width = config.LineWidth
			p.Add(rsiLine)
			
			if config.ShowLegend {
				p.Legend.Add("RSI", rsiLine)
			}
		}
	}

	// Plot MACD if available (scaled to fit with RSI)
	if len(analytics.MACD.MACD) > 0 {
		// Scale MACD to 0-100 range to match RSI
		scaledMACD := make([]float64, len(analytics.MACD.MACD))
		for i, val := range analytics.MACD.MACD {
			scaledMACD[i] = (val * 10) + 50 // Scale and shift to 0-100 range
		}
		
		macdLine, err := plotter.NewLine(makeSimpleXYs(scaledMACD))
		if err == nil {
			macdLine.LineStyle.Color = color.RGBA{R: 0, G: 100, B: 200, A: 255}
			macdLine.LineStyle.Width = config.LineWidth
			p.Add(macdLine)
			
			if config.ShowLegend {
				p.Legend.Add("MACD (scaled)", macdLine)
			}
		}
	}

	// Add RSI reference lines at 30 and 70
	if len(analytics.RSI) > 0 {
		// Oversold line at 30
		oversoldLine, _ := plotter.NewLine(plotter.XYs{
			{X: 0, Y: 30},
			{X: float64(len(analytics.RSI)), Y: 30},
		})
		oversoldLine.LineStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 100}
		oversoldLine.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		oversoldLine.LineStyle.Width = vg.Points(1)
		p.Add(oversoldLine)

		// Overbought line at 70
		overboughtLine, _ := plotter.NewLine(plotter.XYs{
			{X: 0, Y: 70},
			{X: float64(len(analytics.RSI)), Y: 70},
		})
		overboughtLine.LineStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 100}
		overboughtLine.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		overboughtLine.LineStyle.Width = vg.Points(1)
		p.Add(overboughtLine)

		if config.ShowLegend {
			p.Legend.Add("RSI 30/70", oversoldLine)
		}
	}

	return renderPlot(p, config)
}

// Helper function to create simple XY points
func makeSimpleXYs(values []float64) plotter.XYs {
	points := make(plotter.XYs, len(values))
	for i, v := range values {
		points[i].X = float64(i)
		points[i].Y = v
	}
	return points
}

// Helper function to render plot to bytes
func renderPlot(p *plot.Plot, config ChartConfig) ([]byte, error) {
	w, err := p.WriterTo(vg.Length(config.Width), vg.Length(config.Height), "png")
	if err != nil {
		return nil, err
	}

	var buf []byte
	buf = make([]byte, 0)
	_, err = w.WriteTo(&writeBuffer{buf: &buf})
	return buf, err
}

// GenerateIndicatorChart creates just the technical indicators chart
func GenerateIndicatorChart(bts *types.BTCTimeSeries, analytics types.BTCAnalytics) ([]byte, error) {
	config := DefaultChartConfig()
	config.Title = "Bitcoin Technical Indicators (RSI & MACD)"
	
	return DrawTechnicalIndicatorsChart(bts, analytics, config)
}