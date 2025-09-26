package plotcalc

import (
	"fmt"
	"math"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/plotutil"
)

// Mean calculates the mean of durations
func Mean(durations []time.Duration) time.Duration {
	var sum int64
	for _, d := range durations {
		sum += d.Nanoseconds()
	}
	return time.Duration(sum / int64(len(durations)))
}

// StdDev calculates the standard deviation of durations
func StdDev(durations []time.Duration) time.Duration {
	m := Mean(durations).Nanoseconds()
	var variance float64
	for _, d := range durations {
		diff := float64(d.Nanoseconds() - m)
		variance += diff * diff
	}
	variance /= float64(len(durations))
	return time.Duration(math.Sqrt(variance))
}

// PlotDurations plots the data, mean, and mean ± stddev
func PlotDurations(durations []time.Duration, filename string) error {
	meanVal := Mean(durations)
	stdVal := StdDev(durations)

	p := plot.New()
	p.Title.Text = "Durations with Mean ± StdDev"
	p.X.Label.Text = "Sample"
	p.Y.Label.Text = "Time (ms)"

	// Scatter points
	pts := make(plotter.XYs, len(durations))
	for i, d := range durations {
		pts[i].X = float64(i)
		pts[i].Y = float64(d.Milliseconds())
	}
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		return err
	}
	p.Add(scatter)

	// Mean line
	meanLine := plotter.NewFunction(func(x float64) float64 {
		return float64(meanVal.Milliseconds())
	})
	meanLine.Color = plotutil.Color(1)
	p.Add(meanLine)

	// Mean ± stddev
	stdLineUp := plotter.NewFunction(func(x float64) float64 {
		return float64((meanVal + stdVal).Milliseconds())
	})
	stdLineDown := plotter.NewFunction(func(x float64) float64 {
		return float64((meanVal - stdVal).Milliseconds())
	})
	stdLineUp.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	stdLineDown.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	p.Add(stdLineUp, stdLineDown)

	// Save file
	if err := p.Save(6*vg.Inch, 4*vg.Inch, filename); err != nil {
		return fmt.Errorf("failed to save plot: %w", err)
	}
	return nil
}
