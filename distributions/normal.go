package distributions

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// Normal represents a Normal (Gaussian) distribution
type Normal struct {
	Mu    float64
	Sigma float64
	dist  distuv.Normal
}

// NewNormal creates a new Normal distribution
func NewNormal(mu, sigma float64) *Normal {
	return &Normal{
		Mu:    mu,
		Sigma: sigma,
		dist:  distuv.Normal{Mu: mu, Sigma: sigma},
	}
}

// PDF returns the probability density function at x
func (n *Normal) PDF(x float64) float64 {
	return n.dist.Prob(x)
}

// LogPDF returns the log probability density function at x
func (n *Normal) LogPDF(x float64) float64 {
	return n.dist.LogProb(x)
}

// CDF returns the cumulative distribution function at x
func (n *Normal) CDF(x float64) float64 {
	return n.dist.CDF(x)
}

// Quantile returns the inverse CDF at probability p
func (n *Normal) Quantile(p float64) float64 {
	return n.dist.Quantile(p)
}

// Sample generates a random sample
func (n *Normal) Sample() float64 {
	return n.dist.Rand()
}

// SampleN generates n random samples
func (n *Normal) SampleN(nSamples int) []float64 {
	samples := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		samples[i] = n.Sample()
	}
	return samples
}

// Mean returns the expected value
func (n *Normal) Mean() float64 {
	return n.dist.Mean()
}

// Variance returns the variance
func (n *Normal) Variance() float64 {
	return n.dist.Variance()
}

// StdDev returns the standard deviation
func (n *Normal) StdDev() float64 {
	return n.dist.StdDev()
}

// Mode returns the mode
func (n *Normal) Mode() []float64 {
	return []float64{n.Mu}
}

// Median returns the median
func (n *Normal) Median() float64 {
	return n.Mu
}

// Entropy returns the differential entropy
func (n *Normal) Entropy() float64 {
	return n.dist.Entropy()
}

// NormalConjugate implements conjugate update for Normal likelihood with known variance
type NormalConjugate struct {
	*Normal
	KnownVariance float64
}

// NewNormalConjugate creates a conjugate prior for Normal likelihood
func NewNormalConjugate(mu, sigma, knownVariance float64) *NormalConjugate {
	return &NormalConjugate{
		Normal:        NewNormal(mu, sigma),
		KnownVariance: knownVariance,
	}
}

// Update performs conjugate update with Normal likelihood
func (nc *NormalConjugate) Update(data []float64) Posterior {
	n := float64(len(data))
	sumX := 0.0
	for _, x := range data {
		sumX += x
	}
	xBar := sumX / n

	// Conjugate update formulas
	tau0 := 1.0 / (nc.Sigma * nc.Sigma)
	tau := 1.0 / nc.KnownVariance

	tauNew := tau0 + n*tau
	muNew := (tau0*nc.Mu + n*tau*xBar) / tauNew
	sigmaNew := math.Sqrt(1.0 / tauNew)

	return &NormalPosterior{
		Normal: NewNormal(muNew, sigmaNew),
	}
}

// UpdateSingle updates with a single observation
func (nc *NormalConjugate) UpdateSingle(observation float64) Posterior {
	return nc.Update([]float64{observation})
}

// NormalPosterior represents a Normal posterior distribution
type NormalPosterior struct {
	*Normal
}

// CredibleInterval returns the credible interval
func (np *NormalPosterior) CredibleInterval(confidence float64) (lower, upper float64) {
	alpha := (1 - confidence) / 2
	return np.Quantile(alpha), np.Quantile(1 - alpha)
}

// MAP returns the maximum a posteriori estimate
func (np *NormalPosterior) MAP() float64 {
	return np.Mean()
}

// HPD returns the highest posterior density interval
func (np *NormalPosterior) HPD(confidence float64) (lower, upper float64) {
	// For Normal distribution, HPD equals credible interval
	return np.CredibleInterval(confidence)
}
