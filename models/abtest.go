package models

import (
	"fmt"
	"sort"

	"github.com/MyVueCodeHub/myvue-bayes/distributions"
)

// ABTest represents a Bayesian A/B test
type ABTest struct {
	ControlPrior   distributions.Prior
	TreatmentPrior distributions.Prior
	ControlData    []float64
	TreatmentData  []float64
	ControlPost    distributions.Posterior
	TreatmentPost  distributions.Posterior
}

// NewABTest creates a new A/B test with default Beta(1,1) priors
func NewABTest() *ABTest {
	return &ABTest{
		ControlPrior:   distributions.NewBeta(1, 1),
		TreatmentPrior: distributions.NewBeta(1, 1),
	}
}

// NewABTestWithPriors creates a new A/B test with custom priors
func NewABTestWithPriors(controlPrior, treatmentPrior distributions.Prior) *ABTest {
	return &ABTest{
		ControlPrior:   controlPrior,
		TreatmentPrior: treatmentPrior,
	}
}

// AddControlData adds data for the control group
func (ab *ABTest) AddControlData(data []float64) {
	ab.ControlData = append(ab.ControlData, data...)
	ab.updatePosteriors()
}

// AddTreatmentData adds data for the treatment group
func (ab *ABTest) AddTreatmentData(data []float64) {
	ab.TreatmentData = append(ab.TreatmentData, data...)
	ab.updatePosteriors()
}

// updatePosteriors updates the posterior distributions
func (ab *ABTest) updatePosteriors() {
	if len(ab.ControlData) > 0 {
		ab.ControlPost = ab.ControlPrior.Update(ab.ControlData)
	}
	if len(ab.TreatmentData) > 0 {
		ab.TreatmentPost = ab.TreatmentPrior.Update(ab.TreatmentData)
	}
}

// ProbabilityOfImprovement calculates P(treatment > control)
func (ab *ABTest) ProbabilityOfImprovement() float64 {
	if ab.ControlPost == nil || ab.TreatmentPost == nil {
		return 0.5
	}

	// Monte Carlo estimation
	nSamples := 10000
	controlSamples := ab.ControlPost.SampleN(nSamples)
	treatmentSamples := ab.TreatmentPost.SampleN(nSamples)

	wins := 0
	for i := 0; i < nSamples; i++ {
		if treatmentSamples[i] > controlSamples[i] {
			wins++
		}
	}

	return float64(wins) / float64(nSamples)
}

// ExpectedLoss calculates the expected loss for each variant
func (ab *ABTest) ExpectedLoss() (controlLoss, treatmentLoss float64) {
	if ab.ControlPost == nil || ab.TreatmentPost == nil {
		return 0, 0
	}

	nSamples := 10000
	controlSamples := ab.ControlPost.SampleN(nSamples)
	treatmentSamples := ab.TreatmentPost.SampleN(nSamples)

	for i := 0; i < nSamples; i++ {
		diff := treatmentSamples[i] - controlSamples[i]
		if diff > 0 {
			controlLoss += diff
		} else {
			treatmentLoss -= diff
		}
	}

	controlLoss /= float64(nSamples)
	treatmentLoss /= float64(nSamples)
	return
}

// CredibleIntervalDifference returns the credible interval for treatment - control
func (ab *ABTest) CredibleIntervalDifference(confidence float64) (lower, upper float64) {
	if ab.ControlPost == nil || ab.TreatmentPost == nil {
		return 0, 0
	}

	nSamples := 10000
	controlSamples := ab.ControlPost.SampleN(nSamples)
	treatmentSamples := ab.TreatmentPost.SampleN(nSamples)

	differences := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		differences[i] = treatmentSamples[i] - controlSamples[i]
	}

	summary := distributions.ComputeSummary(differences)
	alpha := (1 - confidence) / 2

	if confidence == 0.95 {
		return summary.CI95[0], summary.CI95[1]
	} else if confidence == 0.99 {
		return summary.CI99[0], summary.CI99[1]
	}

	// For other confidence levels, compute from samples
	sortedDiff := make([]float64, len(differences))
	copy(sortedDiff, differences)
	sort.Float64s(sortedDiff)

	lowerIdx := int(alpha * float64(len(sortedDiff)))
	upperIdx := int((1 - alpha) * float64(len(sortedDiff)))

	return sortedDiff[lowerIdx], sortedDiff[upperIdx]
}

// RelativeUplift calculates the relative uplift of treatment over control
func (ab *ABTest) RelativeUplift() (mean, lower, upper float64) {
	if ab.ControlPost == nil || ab.TreatmentPost == nil {
		return 0, 0, 0
	}

	nSamples := 10000
	controlSamples := ab.ControlPost.SampleN(nSamples)
	treatmentSamples := ab.TreatmentPost.SampleN(nSamples)

	uplifts := make([]float64, 0, nSamples)
	for i := 0; i < nSamples; i++ {
		if controlSamples[i] > 0 {
			uplift := (treatmentSamples[i] - controlSamples[i]) / controlSamples[i]
			uplifts = append(uplifts, uplift)
		}
	}

	summary := distributions.ComputeSummary(uplifts)
	return summary.Mean, summary.CI95[0], summary.CI95[1]
}

// Summary returns a human-readable summary of the A/B test results
func (ab *ABTest) Summary() string {
	if ab.ControlPost == nil || ab.TreatmentPost == nil {
		return "Insufficient data for analysis"
	}

	prob := ab.ProbabilityOfImprovement()
	controlLoss, treatmentLoss := ab.ExpectedLoss()
	lower, upper := ab.CredibleIntervalDifference(0.95)
	upliftMean, upliftLower, upliftUpper := ab.RelativeUplift()

	return fmt.Sprintf(`
A/B Test Results:
=================
Control:    n=%d, mean=%.4f
Treatment:  n=%d, mean=%.4f

Probability of Improvement: %.2f%%
Expected Loss:
  - Control:   %.4f
  - Treatment: %.4f

95%% Credible Interval for Difference: [%.4f, %.4f]
Relative Uplift: %.2f%% [%.2f%%, %.2f%%]

Recommendation: %s
`,
		len(ab.ControlData),
		ab.ControlPost.Mean(),
		len(ab.TreatmentData),
		ab.TreatmentPost.Mean(),
		prob*100,
		controlLoss,
		treatmentLoss,
		lower, upper,
		upliftMean*100, upliftLower*100, upliftUpper*100,
		ab.getRecommendation(prob, treatmentLoss),
	)
}

func (ab *ABTest) getRecommendation(prob, treatmentLoss float64) string {
	if prob > 0.95 && treatmentLoss < 0.01 {
		return "Strong evidence favors treatment. Recommend implementation."
	} else if prob > 0.80 {
		return "Moderate evidence favors treatment. Consider implementation or continue testing."
	} else if prob < 0.20 {
		return "Evidence favors control. Treatment likely inferior."
	} else {
		return "Insufficient evidence to make a recommendation. Continue testing."
	}
}
