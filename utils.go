package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	// "time"

	"github.com/joho/godotenv"
	"gonum.org/v1/gonum/stat"
)

var processBar = []string{
	"00%: [                                          ]",
	"05%: [##                                        ]",
	"10%: [####                                      ]",
	"15%: [######                                    ]",
	"20%: [########                                  ]",
	"25%: [##########                                ]",
	"30%: [############                              ]",
	"35%: [##############                            ]",
	"40%: [################                          ]",
	"45%: [##################                        ]",
	"50%: [####################                      ]",
	"55%: [######################                    ]",
	"60%: [########################                  ]",
	"65%: [##########################                ]",
	"70%: [############################              ]",
	"75%: [##############################            ]",
	"80%: [################################          ]",
	"85%: [##################################        ]",
	"90%: [####################################      ]",
	"95%: [######################################    ]",
	"100%:[##########################################]\n",
}

var processBar2 = []string{
	"00%: [                                          ]",
	"10%: [####                                      ]",
	"30%: [############                              ]",
	"40%: [################                          ]",
	"60%: [########################                  ]",
	"70%: [############################              ]",
	"90%: [####################################      ]",
	"100%:[##########################################]\n",
}

func printCustomProgressBar(progress int) {
	if progress < 0 || progress > 100 {
		return
	}
	index := progress / 10
	fmt.Printf("\r%s", processBar2[index])
}

/*
Function to get back the var in .env file
*/
func getDotEnvVar(key string) string {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// Convert The price from currency to EUR, (Values September 13)

func convertToEUR(price float64, currency string) (float64, error) {
	// Define exchange rates for some currencies (you can extend this list)
	exchangeRates := map[string]float64{
		"EUR": 1.00,         // Euro
		"USD": 0.93,         // US Dollar
		"BRL": 0.19,         // Brazilian Real
		"MYR": 0.20,         // Malaysian Ringgit
		"IDR": 0.000061,     // Indonesian Rupiah
		"CNY": 0.13,         // Chinese Yuan
		"PLN": 0.22,         // Polish Zloty
		"RSD": 0.0085,       // Serbian Dinar
		"PKR": 0.0032,       // Pakistani Rupee
		"GTQ": 0.12,         // Guatemalan Quetzal
		"PHP": 0.016,        // Philippines Pes
		"BAM": 0.51,         // Bosnian Marks
		"ALL": 0.0094108537, // Albanian Lek
		"RUB": 0.0097,       // Russian Ruble
		"HNL": 0.038,        // Honduran Lempira
		"JPY": 0.0063,       // Japanese Yen
		"TND": 0.30,         // Tunisian Dinar
		"NOK": 0.087,        // Norwegian Krone
		"KRW": 0.00070,      // Won south korean
		"COP": 0.00023,      //Colombian Peso
		"MXN": 0.054,        // Mexico peso
		"GMD": 0.015,        // Gambian Dalasi
		"THB": 0.026,        // Thai Baht
		"CAD": 0.69,         //Canadian Dollar
		"AMD": 0.0024,       // Armenian Dram
		"NGN": 0.0012,       // Nigerian Naira
		"SEK": 0.084,        // Swedish Krone
		"VND": 0.000039,     // Vietnamese Dong
		"DKK": 0.13,         // Danish Krone
		"IQD": 0.00071,      // Iraqi Dinar
		"HRK": 0.13272281,   // croatian Kunas
		"CZK": 0.041,        // Czech Republic Koruna
		"BGN": 0.51,         // bulgarian leva
		"ETB": 0.017,        // ethiopian birs
		// Add more currencies and their exchange rates here
	}

	// Check if the currency is in the exchangeRates map
	rate, found := exchangeRates[currency]
	if !found {
		return 0.0, fmt.Errorf("Exchange rate not found for currency: %s", currency)
	}

	// Convert the price to euros
	priceInEUR := price * rate
	return priceInEUR, nil
}

/*
Function to print 10 random entries of the map
*/

func printRandomEntriesMap(customersMoneySpent map[int64]float64) {

	fmt.Printf("------------------------------------------------------------\n")
	fmt.Printf("Printing 10 random entries of the map:\n")
	// Shuffle the keys of the map
	keys := make([]int64, 0, len(customersMoneySpent))
	for customerID := range customersMoneySpent {
		keys = append(keys, customerID)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// Log 10 random entries
	numEntriesToLog := 10
	for i := 0; i < numEntriesToLog && i < len(keys); i++ {
		customerID := keys[i]
		spent := customersMoneySpent[customerID]
		log.Printf("CustomerID: %d, Spent: %.2f\n", customerID, spent)
	}
	fmt.Printf("------------------------------------------------------------\n")
}

/*

*/

func calculateQuantiles(data []float64, quantiles []float64) []float64 {
    result := make([]float64, len(quantiles))
    for i, q := range quantiles {
        result[i] = stat.Quantile(q, stat.Empirical, data, nil)
    }
    return result
}