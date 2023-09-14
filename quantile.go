package main

import (
	"fmt"
	// "math"
	"sort"
)

// "fmt"

type QuantileInfo struct {
	NumClients int
	MaxRevenue float64
	MinRevenue float64
}

/*
	Function to create a map with the best clients, last quartile (0.9725)
*/

func createBestClientMap(spentValues []float64, moneySpentSlice []CustomerSpent) (map[int64]float64, error) {
	// Calculate the quantile value for the first quantile (2.5%)
	lastQuantile := 0.975

	// Calculate the index of the client that represents the first quantile
	index := int(float64(len(spentValues)) * lastQuantile)

	// Create a map to store the N clients of the first quantile
	lastQuantileClients := make(map[int64]float64)

	// Iterate through the sorted spentValues and add N clients to the map
	for i := index; i < len(spentValues); i++ {
		lastQuantileClients[moneySpentSlice[i].CustomerID] = spentValues[i]
	}

	// Print the number of clients in the last quantile
	fmt.Printf("Number of clients in the last quantile (Best clients): %d\n", len(lastQuantileClients))
	// Iterate over the map keys and values
	for customerID, spent := range lastQuantileClients {
		fmt.Printf("CustomerID: %d, Spent: %.2f\n", customerID, spent)
	}
	// Simulate a database error (you should replace this with your actual database operation)
	if len(lastQuantileClients) == 0 {
		return nil, fmt.Errorf("no best clients found")
	}
	return lastQuantileClients, nil
}

/*
	Function who create all the quantiles for all the customers, 2.5% by 2.5 until the end (Because we put 40 in num of quartiles)
*/

func CalculateQuantilesNearestRank(spentValues []float64, numQuantiles int) (map[float64]QuantileInfo, []float64) {
	// Sort the spentValues array
	sort.Float64s(spentValues)
	// Create a map to store quantile information
	quantileInfoMap := make(map[float64]QuantileInfo)
	// Calculate the number of clients per quantile
	clientsPerQuantile := len(spentValues) / numQuantiles
	// Create a slice to store quantiles for sorting
	quantiles := make([]float64, numQuantiles)

	for i := 0; i < numQuantiles; i++ {
		// Calculate the index range for the current quantile
		startIndex := i * clientsPerQuantile
		endIndex := (i + 1) * clientsPerQuantile
		if i == numQuantiles-1 {
			endIndex = len(spentValues)
		}
		// Initialize minRevenue and maxRevenue with the first value in the current quantile
		minRevenue := spentValues[startIndex]
		maxRevenue := spentValues[startIndex]
		numClients := endIndex - startIndex
		// Calculate max and min revenue in the current quantile
		for j := startIndex; j < endIndex; j++ {
			if spentValues[j] > maxRevenue {
				maxRevenue = spentValues[j]
			}
			if spentValues[j] < minRevenue {
				minRevenue = spentValues[j]
			}
		}
		// Calculate the quantile value
		quantile := float64(i+1) / float64(numQuantiles)
		quantiles[i] = quantile
		// Store quantile information in the map
		quantileInfoMap[quantile] = QuantileInfo{
			NumClients: numClients,
			MaxRevenue: maxRevenue,
			MinRevenue: minRevenue,
		}
	}
	printQuantileInfo(quantileInfoMap, quantiles)
	return quantileInfoMap, quantiles
}

/*
	Function who print the quantile information
*/

func printQuantileInfo(quantileInfoMap map[float64]QuantileInfo, quantiles []float64) {
	for _, q := range quantiles {
		info := quantileInfoMap[q]
		fmt.Printf("%.2f%% Quantile Info:\n", q*100-2.5)
		fmt.Printf("Number of clients: %d\n", info.NumClients)
		fmt.Printf("Max Revenue: %.2f\n", info.MaxRevenue)
		fmt.Printf("Min Revenue: %.2f\n\n", info.MinRevenue)
	}
}

/*
	Function to generate all of the quantile number automatically, with a start a end and step
	Here : start = 0, end = 1, step = 0.025
*/

func GenerateQuantiles(start, end, step float64) []float64 {
	var quantiles []float64

	for q := start; q <= end; q += step {
		quantiles = append(quantiles, q)
	}

	return quantiles
}
