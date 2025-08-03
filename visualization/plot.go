package visualization

import (
	"fmt"
	"image/color"

	"github.com/MyVueCodeHub/myvue-bayes/distributions"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// PlotType represents the type of plot
type PlotType int

const (
	LinePlot PlotType = iota
	ScatterPlot
	HistogramPlot
	DensityPlot
	BoxPlot
)

// PlotData represents data for plotting
type PlotData struct {
	X     []float64
	Y     []float64
	Label string
	Type  PlotType
}

// BayesianPlotter provides plotting functionality for Bayesian analyses
type BayesianPlotter struct {
	plot *plot.Plot
}

// NewBayesianPlotter creates a new plotter
func NewBayesianPlotter(title string) (*BayesianPlotter, error) {
	p := plot.New()

	p.Title.Text = title
	return &BayesianPlotter{plot: p}, nil
}

// PriorPosteriorPlot creates a plot comparing prior and posterior
func (bp *BayesianPlotter) PriorPosteriorPlot(
	prior distributions.Distribution,
	posterior distributions.Distribution,
	xMin, xMax float64,
	nPoints int,
) error {
	x := make([]float64, nPoints)
	priorY := make([]float64, nPoints)
	postY := make([]float64, nPoints)

	step := (xMax - xMin) / float64(nPoints-1)
	for i := 0; i < nPoints; i++ {
		x[i] = xMin + float64(i)*step
		priorY[i] = prior.PDF(x[i])
		postY[i] = posterior.PDF(x[i])
	}

	// Prior line
	priorLine, err := plotter.NewLine(plotter.XYs{})
	if err != nil {
		return err
	}
	for i := range x {
		priorLine.XYs = append(priorLine.XYs, plotter.XY{X: x[i], Y: priorY[i]})
	}
	priorLine.Color = color.RGBA{0, 0, 255, 50}
	priorLine.Width = vg.Points(2)

	// Posterior line
	postLine, err := plotter.NewLine(plotter.XYs{})
	if err != nil {
		return err
	}
	for i := range x {
		postLine.XYs = append(postLine.XYs, plotter.XY{X: x[i], Y: postY[i]})
	}
	postLine.Color = color.RGBA{255, 0, 0, 50} //red
	postLine.Width = vg.Points(2)

	bp.plot.Add(priorLine, postLine)
	bp.plot.Legend.Add("Prior", priorLine)
	bp.plot.Legend.Add("Posterior", postLine)
	bp.plot.X.Label.Text = "Value"
	bp.plot.Y.Label.Text = "Density"

	return nil
}

// CredibleIntervalPlot creates a plot with credible intervals
func (bp *BayesianPlotter) CredibleIntervalPlot(
	samples []float64,
	credibleLevel float64,
) error {
	h, err := plotter.NewHist(plotter.Values(samples), 50)
	if err != nil {
		return err
	}
	h.Normalize(1)

	// Calculate credible interval
	summary := distributions.ComputeSummary(samples)
	var lower, upper float64
	if credibleLevel == 0.95 {
		lower, upper = summary.CI95[0], summary.CI95[1]
	} else {
		lower, upper = summary.CI99[0], summary.CI99[1]
	}

	// Add vertical lines for credible interval
	lowerLine, err := plotter.NewLine(plotter.XYs{
		{X: lower, Y: 0},
		{X: lower, Y: 1},
	})
	if err != nil {
		return err
	}
	lowerLine.LineStyle.Color = color.RGBA{255, 0, 0, 50} //red
	lowerLine.LineStyle.Width = vg.Points(2)
	lowerLine.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	upperLine, err := plotter.NewLine(plotter.XYs{
		{X: upper, Y: 0},
		{X: upper, Y: 1},
	})
	if err != nil {
		return err
	}
	upperLine.LineStyle = lowerLine.LineStyle

	bp.plot.Add(h, lowerLine, upperLine)
	bp.plot.X.Label.Text = "Value"
	bp.plot.Y.Label.Text = "Density"

	return nil
}

// TracePlot creates trace plots for MCMC diagnostics
func (bp *BayesianPlotter) TracePlot(chains [][]float64) error {
	for i, chain := range chains {
		line, err := plotter.NewLine(plotter.XYs{})
		if err != nil {
			return err
		}

		for j, value := range chain {
			line.XYs = append(line.XYs, plotter.XY{X: float64(j), Y: value})
		}

		bp.plot.Add(line)
		bp.plot.Legend.Add(fmt.Sprintf("Chain %d", i+1), line)
	}

	bp.plot.X.Label.Text = "Iteration"
	bp.plot.Y.Label.Text = "Value"

	return nil
}

// Save saves the plot to a file
func (bp *BayesianPlotter) Save(filename string, width, height vg.Length) error {
	return bp.plot.Save(width, height, filename)
}

// PlotABTestResults creates a comprehensive visualization of A/B test results
func PlotABTestResults(
	controlSamples, treatmentSamples []float64,
	filename string,
) error {
	p := plot.New()

	p.Title.Text = "A/B Test Posterior Distributions"

	// Create histograms
	controlHist, err := plotter.NewHist(plotter.Values(controlSamples), 50)
	if err != nil {
		return err
	}
	controlHist.FillColor = color.RGBA{0, 0, 255, 50} //blue
	controlHist.Color = color.RGBA{255, 0, 0, 50}     //red
	controlHist.Normalize(1)

	treatmentHist, err := plotter.NewHist(plotter.Values(treatmentSamples), 50)
	if err != nil {
		return err
	}
	treatmentHist.FillColor = color.RGBA{255, 0, 0, 50} //red
	treatmentHist.Color = color.RGBA{255, 0, 0, 50}     //red
	treatmentHist.Normalize(1)

	p.Add(controlHist, treatmentHist)
	p.Legend.Add("Control", controlHist)
	p.Legend.Add("Treatment", treatmentHist)
	p.X.Label.Text = "Conversion Rate"
	p.Y.Label.Text = "Density"

	return p.Save(8*vg.Inch, 6*vg.Inch, filename)
}
