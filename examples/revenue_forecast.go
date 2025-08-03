package examples

// import (
// 	"fmt"

// 	"github.com/MyVueCodeHub/myvue-bayes/metrics"
// )

// func ForecastExample() {
// 	// Historical monthly revenue data
// 	historicalRevenue := []float64{
// 		45000, 48000, 47500, 52000, 54000,
// 		53500, 58000, 61000, 59500, 63000,
// 		65000, 68000,
// 	}

// 	// Create business metrics analyzer
// 	bm := metrics.NewBusinessMetrics()

// 	// Project revenue for next 6 months
// 	projections := bm.RevenueProjection(historicalRevenue, 6)

// 	fmt.Println("Revenue Projections for Next 6 Months:")
// 	fmt.Println("======================================")

// 	for i, proj := range projections {
// 		fmt.Printf("Month %d:\n", i+1)
// 		fmt.Printf("  Expected: $%.2f\n", proj.Mean)
// 		fmt.Printf("  95%% CI: [$%.2f, $%.2f]\n", proj.CI95[0], proj.CI95[1])
// 		fmt.Printf("  99%% CI: [$%.2f, $%.2f]\n\n", proj.CI99[0], proj.CI99[1])
// 	}
// }
