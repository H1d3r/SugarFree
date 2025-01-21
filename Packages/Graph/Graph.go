package Graph

import (
	"SugarFree/Packages/Manager"
	"fmt"
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// barChart represents the data needed for a bar in the chart
type barChart struct {
	label       string
	originalVal float64
	reducedVal  float64
}

// GenerateEntropyGraph creates a visualization comparing original and reduced entropy
func GenerateEntropyGraph(originalSections, reducedSections []Manager.SectionEntropy, outputPath string) error {
	// Create a new plot
	p := plot.New()

	// Set plot title and labels
	p.Title.Text = "PE Section Entropy Comparison"
	p.X.Label.Text = "Sections"
	p.Y.Label.Text = "Entropy"

	// Create bars for the chart
	var bars []barChart
	for i := 0; i < len(originalSections); i++ {
		bars = append(bars, barChart{
			label:       originalSections[i].Name,
			originalVal: originalSections[i].Entropy,
			reducedVal:  reducedSections[i].Entropy,
		})
	}

	// Create grouped bars
	groupWidth := vg.Points(20)
	w := groupWidth / 2

	// Create bars for original entropy
	originalBars, err := createBars(bars, true, w)
	if err != nil {
		return fmt.Errorf("error creating original bars: %w", err)
	}
	originalBars.Color = color.RGBA{R: 255, B: 0, A: 255}
	p.Add(originalBars)

	// Create bars for reduced entropy
	reducedBars, err := createBars(bars, false, w)
	if err != nil {
		return fmt.Errorf("error creating reduced bars: %w", err)
	}
	reducedBars.Color = color.RGBA{G: 255, B: 0, A: 255}
	p.Add(reducedBars)

	// Add legend
	p.Legend.Add("Original", originalBars)
	p.Legend.Add("Reduced", reducedBars)
	p.Legend.Top = true

	// Set axis ranges
	p.Y.Min = 0
	p.Y.Max = 8 // Maximum possible entropy value
	p.X.Min = -0.5
	p.X.Max = float64(len(bars)) - 0.5

	// Customize the axis labels
	p.X.Label.TextStyle.Font.Size = vg.Points(12)
	p.Y.Label.TextStyle.Font.Size = vg.Points(12)
	p.Title.TextStyle.Font.Size = vg.Points(14)

	// Add grid lines
	p.Add(plotter.NewGrid())

	// Save the plot
	if err := p.Save(8*vg.Inch, 6*vg.Inch, outputPath); err != nil {
		return fmt.Errorf("error saving plot: %w", err)
	}

	return nil
}

// createBars creates a bar plotter for either original or reduced values
func createBars(data []barChart, isOriginal bool, width vg.Length) (*plotter.BarChart, error) {
	var values plotter.Values
	var labels []string

	for _, bar := range data {
		if isOriginal {
			values = append(values, bar.originalVal)
		} else {
			values = append(values, bar.reducedVal)
		}
		labels = append(labels, bar.label)
	}

	bars, err := plotter.NewBarChart(values, width)
	if err != nil {
		return nil, err
	}

	// Customize bar appearance
	bars.LineStyle.Width = vg.Points(1)
	bars.Offset = width
	if !isOriginal {
		bars.Offset = -width
	}

	return bars, nil
}

// calculateYAxisRange calculates appropriate Y-axis range based on data
func calculateYAxisRange(data []barChart) (float64, float64) {
	minVal := math.MaxFloat64
	maxVal := -math.MaxFloat64

	for _, bar := range data {
		minVal = math.Min(minVal, math.Min(bar.originalVal, bar.reducedVal))
		maxVal = math.Max(maxVal, math.Max(bar.originalVal, bar.reducedVal))
	}

	// Add some padding
	padding := (maxVal - minVal) * 0.1
	return math.Max(0, minVal-padding), maxVal + padding
}
