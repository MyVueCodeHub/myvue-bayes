package distributions

import (
	"slices"

	"gonum.org/v1/gonum/stat"
)

// Distribution defines the interface for all probability distributions
type Distribution interface {
	// PDF returns the probability density function at x
	PDF(x float64) float64

	// LogPDF returns the log probability density function at x
	LogPDF(x float64) float64

	// CDF returns the cumulative distribution function at x
	CDF(x float64) float64

	// Quantile returns the inverse CDF at probability p
	Quantile(p float64) float64

	// Sample generates a random sample from the distribution
	Sample() float64

	// SampleN generates n random samples from the distribution
	SampleN(n int) []float64

	// Mean returns the expected value
	Mean() float64

	// Variance returns the variance
	Variance() float64

	// StdDev returns the standard deviation
	StdDev() float64
}

// ContinuousDistribution extends Distribution for continuous random variables
type ContinuousDistribution interface {
	Distribution

	// Mode returns the mode(s) of the distribution
	Mode() []float64

	// Median returns the median
	Median() float64

	// Entropy returns the differential entropy
	Entropy() float64
}

// DiscreteDistribution extends Distribution for discrete random variables
type DiscreteDistribution interface {
	Distribution

	// PMF returns the probability mass function at k
	PMF(k int) float64

	// LogPMF returns the log probability mass function at k
	LogPMF(k int) float64
}

// Prior represents a prior distribution that can be updated
type Prior interface {
	Distribution

	// Update returns the posterior distribution given observed data
	Update(data []float64) Posterior

	// UpdateSingle returns the posterior distribution given a single observation
	UpdateSingle(observation float64) Posterior
}

// Posterior represents a posterior distribution
type Posterior interface {
	Distribution

	// CredibleInterval returns the credible interval at the given confidence level
	CredibleInterval(confidence float64) (lower, upper float64)

	// MAP returns the maximum a posteriori estimate
	MAP() float64

	// HPD returns the highest posterior density interval
	HPD(confidence float64) (lower, upper float64)
}

// Summary provides a statistical summary of a distribution
type Summary struct {
	Mean     float64
	Median   float64
	Mode     float64
	Variance float64
	StdDev   float64
	CI95     [2]float64
	CI99     [2]float64
	Samples  []float64
}

// ComputeSummary generates a statistical summary from samples
func ComputeSummary(samples []float64) Summary {
	sort := make([]float64, len(samples))
	copy(sort, samples)

	slices.Sort(sort)

	return Summary{
		Mean:     stat.Mean(samples, nil),
		Median:   stat.Quantile(0.5, stat.Empirical, sort, nil),
		Mode:     estimateMode(samples),
		Variance: stat.Variance(samples, nil),
		StdDev:   stat.StdDev(samples, nil),
		CI95: [2]float64{
			stat.Quantile(0.025, stat.Empirical, sort, nil),
			stat.Quantile(0.975, stat.Empirical, sort, nil),
		},
		CI99: [2]float64{
			stat.Quantile(0.005, stat.Empirical, sort, nil),
			stat.Quantile(0.995, stat.Empirical, sort, nil),
		},
		Samples: samples,
	}
}

func estimateMode(samples []float64) float64 {
	// Simple kernel density estimation for mode
	//bandwidth := 1.06 * stat.StdDev(samples, nil) * math.Pow(float64(len(samples)), -0.2)
	// Simplified implementation - in production use more sophisticated KDE
	return stat.Mean(samples, nil) // Placeholder
}
