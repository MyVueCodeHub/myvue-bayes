package examples

// import (
// 	"fmt"
// 	"log"

// 	"github.com/MyVueCodeHub/myvue-bayes/models"
// 	"github.com/MyVueCodeHub/myvue-bayes/visualization"
// )

// func AbTestExample() {
// 	// Create a new A/B test
// 	test := models.NewABTest()

// 	// Simulate some data
// 	// Control: 120 conversions out of 1000 visitors (12% conversion rate)
// 	controlData := make([]float64, 1000)
// 	for i := 0; i < 120; i++ {
// 		controlData[i] = 1.0
// 	}

// 	// Treatment: 150 conversions out of 1000 visitors (15% conversion rate)
// 	treatmentData := make([]float64, 1000)
// 	for i := 0; i < 150; i++ {
// 		treatmentData[i] = 1.0
// 	}

// 	// Add data to the test
// 	test.AddControlData(controlData)
// 	test.AddTreatmentData(treatmentData)

// 	// Print results
// 	fmt.Println(test.Summary())

// 	// Create visualization
// 	if test.ControlPost != nil && test.TreatmentPost != nil {
// 		controlSamples := test.ControlPost.SampleN(10000)
// 		treatmentSamples := test.TreatmentPost.SampleN(10000)

// 		err := visualization.PlotABTestResults(
// 			controlSamples,
// 			treatmentSamples,
// 			"ab_test_results.png",
// 		)
// 		if err != nil {
// 			log.Printf("Error creating plot: %v", err)
// 		} else {
// 			fmt.Println("Plot saved to ab_test_results.png")
// 		}
// 	}
// }
