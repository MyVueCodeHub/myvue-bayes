package metrics

import (
	"math"

	"github.com/MyVueCodeHub/myvue-bayes/distributions"
	"gonum.org/v1/gonum/stat"
)

// MetricEstimate represents a business metric with uncertainty
type MetricEstimate struct {
	Mean    float64
	Median  float64
	Mode    float64
	CI95    [2]float64
	CI99    [2]float64
	Samples []float64
	Summary distributions.Summary
}

// BusinessMetrics provides Bayesian estimates for common business metrics
type BusinessMetrics struct {
	DefaultPriors map[string]distributions.Prior
}

// NewBusinessMetrics creates a new BusinessMetrics instance with sensible defaults
func NewBusinessMetrics() *BusinessMetrics {
	return &BusinessMetrics{
		DefaultPriors: map[string]distributions.Prior{
			"conversion": distributions.NewBeta(1, 1),
			"revenue":    distributions.NewNormalConjugate(100, 50, 100),
			"retention":  distributions.NewBeta(1, 1),
			"churn":      distributions.NewBeta(1, 1),
		},
	}
}

// ConversionRate estimates conversion rate with uncertainty
func (bm *BusinessMetrics) ConversionRate(successes, trials int) MetricEstimate {
	prior := bm.DefaultPriors["conversion"]

	// Create binary data
	data := make([]float64, trials)
	for i := 0; i < successes; i++ {
		data[i] = 1.0
	}

	posterior := prior.Update(data)
	samples := posterior.SampleN(10000)
	summary := distributions.ComputeSummary(samples)

	return MetricEstimate{
		Mean:    summary.Mean,
		Median:  summary.Median,
		Mode:    summary.Mode,
		CI95:    summary.CI95,
		CI99:    summary.CI99,
		Samples: samples,
		Summary: summary,
	}
}

// AverageOrderValue estimates AOV with uncertainty
func (bm *BusinessMetrics) AverageOrderValue(orders []float64) MetricEstimate {
	// Use log-normal for revenue data
	logOrders := make([]float64, len(orders))
	for i, order := range orders {
		if order > 0 {
			logOrders[i] = math.Log(order)
		}
	}

	// Estimate parameters
	mean := stat.Mean(logOrders, nil)
	stdDev := stat.StdDev(logOrders, nil)

	// Create posterior samples
	normalPost := distributions.NewNormal(mean, stdDev/math.Sqrt(float64(len(orders))))
	samples := make([]float64, 10000)
	for i := range samples {
		samples[i] = math.Exp(normalPost.Sample())
	}

	summary := distributions.ComputeSummary(samples)

	return MetricEstimate{
		Mean:    summary.Mean,
		Median:  summary.Median,
		Mode:    summary.Mode,
		CI95:    summary.CI95,
		CI99:    summary.CI99,
		Samples: samples,
		Summary: summary,
	}
}

// RetentionRate estimates retention rate with cohort data
func (bm *BusinessMetrics) RetentionRate(cohortData [][]int) []MetricEstimate {
	// cohortData[i][j] = number of users from cohort i active in period j
	results := make([]MetricEstimate, len(cohortData[0]))

	for period := 0; period < len(cohortData[0]); period++ {
		totalUsers := 0
		activeUsers := 0

		for cohort := 0; cohort < len(cohortData); cohort++ {
			if period < len(cohortData[cohort]) {
				if period == 0 {
					totalUsers += cohortData[cohort][0]
				} else if cohort+period < len(cohortData) {
					totalUsers += cohortData[cohort][0]
					activeUsers += cohortData[cohort][period]
				}
			}
		}

		if totalUsers > 0 {
			results[period] = bm.ConversionRate(activeUsers, totalUsers)
		}
	}

	return results
}

// ChurnProbability estimates the probability of churn
func (bm *BusinessMetrics) ChurnProbability(activeCustomers, churnedCustomers int) MetricEstimate {
	return bm.ConversionRate(churnedCustomers, activeCustomers+churnedCustomers)
}

// CustomerLifetimeValue estimates CLV with uncertainty
func (bm *BusinessMetrics) CustomerLifetimeValue(
	avgOrderValue MetricEstimate,
	purchaseFrequency MetricEstimate,
	churnRate MetricEstimate,
) MetricEstimate {
	// CLV = AOV × Purchase Frequency × (1 / Churn Rate)
	// Using Monte Carlo to propagate uncertainty

	nSamples := 10000
	clvSamples := make([]float64, nSamples)

	for i := 0; i < nSamples; i++ {
		aov := avgOrderValue.Samples[i%len(avgOrderValue.Samples)]
		freq := purchaseFrequency.Samples[i%len(purchaseFrequency.Samples)]
		churn := churnRate.Samples[i%len(churnRate.Samples)]

		if churn > 0 {
			clvSamples[i] = aov * freq * (1.0 / churn)
		} else {
			clvSamples[i] = aov * freq * 100 // Cap at 100 periods
		}
	}

	summary := distributions.ComputeSummary(clvSamples)

	return MetricEstimate{
		Mean:    summary.Mean,
		Median:  summary.Median,
		Mode:    summary.Mode,
		CI95:    summary.CI95,
		CI99:    summary.CI99,
		Samples: clvSamples,
		Summary: summary,
	}
}

// RevenueProjection projects future revenue with uncertainty
func (bm *BusinessMetrics) RevenueProjection(
	historicalRevenue []float64,
	periods int,
) []MetricEstimate {
	// Simple Bayesian linear regression for trend
	n := float64(len(historicalRevenue))
	x := make([]float64, len(historicalRevenue))
	for i := range x {
		x[i] = float64(i)
	}

	// Calculate regression parameters
	meanX := stat.Mean(x, nil)
	meanY := stat.Mean(historicalRevenue, nil)

	var num, den float64
	for i := range x {
		num += (x[i] - meanX) * (historicalRevenue[i] - meanY)
		den += (x[i] - meanX) * (x[i] - meanX)
	}

	slope := num / den
	intercept := meanY - slope*meanX

	// Calculate residual variance
	var residualSS float64
	for i := range x {
		pred := intercept + slope*x[i]
		residualSS += math.Pow(historicalRevenue[i]-pred, 2)
	}
	sigma := math.Sqrt(residualSS / (n - 2))

	// Project forward with increasing uncertainty
	projections := make([]MetricEstimate, periods)

	for t := 0; t < periods; t++ {
		futureX := n + float64(t)
		meanPred := intercept + slope*futureX

		// Prediction interval widens with distance from data
		predSE := sigma * math.Sqrt(1+1/n+math.Pow(futureX-meanX, 2)/den)

		// Generate samples
		samples := make([]float64, 10000)
		dist := distributions.NewNormal(meanPred, predSE)
		for i := range samples {
			samples[i] = math.Max(0, dist.Sample()) // Revenue can't be negative
		}

		summary := distributions.ComputeSummary(samples)
		projections[t] = MetricEstimate{
			Mean:    summary.Mean,
			Median:  summary.Median,
			Mode:    summary.Mode,
			CI95:    summary.CI95,
			CI99:    summary.CI99,
			Samples: samples,
			Summary: summary,
		}
	}

	return projections
}
