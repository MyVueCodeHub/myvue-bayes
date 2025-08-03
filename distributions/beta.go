package distributions

import (
	"gonum.org/v1/gonum/stat/distuv"
)

// Beta represents a Beta distribution
type Beta struct {
	Alpha float64
	Beta  float64
	dist  distuv.Beta
}

// NewBeta creates a new Beta distribution
func NewBeta(alpha, beta float64) *Beta {
	return &Beta{
		Alpha: alpha,
		Beta:  beta,
		dist:  distuv.Beta{Alpha: alpha, Beta: beta},
	}
}

// PDF returns the probability density function at x
func (b *Beta) PDF(x float64) float64 {
	return b.dist.Prob(x)
}

// LogPDF returns the log probability density function at x
func (b *Beta) LogPDF(x float64) float64 {
	return b.dist.LogProb(x)
}

// CDF returns the cumulative distribution function at x
func (b *Beta) CDF(x float64) float64 {
	return b.dist.CDF(x)
}

// Quantile returns the inverse CDF at probability p
func (b *Beta) Quantile(p float64) float64 {
	return b.dist.Quantile(p)
}

// Sample generates a random sample
func (b *Beta) Sample() float64 {
	return b.dist.Rand()
}

// SampleN generates n random samples
func (b *Beta) SampleN(n int) []float64 {
	samples := make([]float64, n)
	for i := 0; i < n; i++ {
		samples[i] = b.Sample()
	}
	return samples
}

// Mean returns the expected value
func (b *Beta) Mean() float64 {
	return b.dist.Mean()
}

// Variance returns the variance
func (b *Beta) Variance() float64 {
	return b.dist.Variance()
}

// StdDev returns the standard deviation
func (b *Beta) StdDev() float64 {
	return b.dist.StdDev()
}

// Mode returns the mode(s) of the distribution
func (b *Beta) Mode() []float64 {
	if b.Alpha > 1 && b.Beta > 1 {
		return []float64{(b.Alpha - 1) / (b.Alpha + b.Beta - 2)}
	} else if b.Alpha < 1 && b.Beta < 1 {
		return []float64{0, 1} // Bimodal
	} else if b.Alpha < 1 && b.Beta >= 1 {
		return []float64{0}
	} else {
		return []float64{1}
	}
}

// Median returns the median
func (b *Beta) Median() float64 {
	return b.Quantile(0.5)
}

// Entropy returns the differential entropy
func (b *Beta) Entropy() float64 {
	return b.dist.Entropy()
}

// Update performs Bayesian update with binomial likelihood
func (b *Beta) Update(data []float64) Posterior {
	successes := 0.0
	trials := float64(len(data))
	for _, x := range data {
		if x > 0 {
			successes++
		}
	}

	return &BetaPosterior{
		Beta: NewBeta(b.Alpha+successes, b.Beta+trials-successes),
	}
}

// UpdateSingle updates with a single observation
func (b *Beta) UpdateSingle(observation float64) Posterior {
	if observation > 0 {
		return &BetaPosterior{
			Beta: NewBeta(b.Alpha+1, b.Beta),
		}
	}
	return &BetaPosterior{
		Beta: NewBeta(b.Alpha, b.Beta+1),
	}
}

// BetaPosterior represents a Beta posterior distribution
type BetaPosterior struct {
	*Beta
}

// CredibleInterval returns the credible interval
func (bp *BetaPosterior) CredibleInterval(confidence float64) (lower, upper float64) {
	alpha := (1 - confidence) / 2
	return bp.Quantile(alpha), bp.Quantile(1 - alpha)
}

// MAP returns the maximum a posteriori estimate
func (bp *BetaPosterior) MAP() float64 {
	modes := bp.Mode()
	if len(modes) > 0 {
		return modes[0]
	}
	return bp.Mean()
}

// HPD returns the highest posterior density interval
func (bp *BetaPosterior) HPD(confidence float64) (lower, upper float64) {
	// Simplified implementation - for Beta, often similar to credible interval
	// In production, use numerical optimization for true HPD
	return bp.CredibleInterval(confidence)
}
