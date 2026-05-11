package main

import (
	"fmt"
	"maths/utils"
)

func main() {
	myFile := utils.CheckFile()
	file, err := utils.OpenFile(myFile)
	if err != nil {
		panic("Cannot open the file. File path,name or type is incorrect")
	}
	numLines := utils.ReadConvStoreLines(file)
	roundedAverage, unRoundedAverage := utils.CalculateAverage(numLines)
	utils.DataValidation(unRoundedAverage, numLines)
	roundedMedian := utils.CalculateMedian(numLines)
	roundedVariance, unRoundedVariance := utils.CalculateVariance(numLines, unRoundedAverage)
	roundedStandardDeviation := utils.CalculateStandardDeviation(unRoundedVariance)

	fmt.Printf("Average: %v\n", roundedAverage)
	fmt.Printf("Median: %v\n", roundedMedian)
	fmt.Printf("Variance: %v\n", roundedVariance)
	fmt.Printf("Standard Deviation: %v\n", roundedStandardDeviation)

}
