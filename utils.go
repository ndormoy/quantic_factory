package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
		"EUR": 1.07,     // Euro
		"USD": 1.18,     // US Dollar
		"BRL": 5.30,     // Brazilian Real
		"MYR": 5.03,     // Malaysian Ringgit
		"IDR": 16500.95, // Indonesian Rupiah
		"CNY": 7.72,     // Chinese Yuan
		"PLN": 4.62,     // Polish Zloty
		"RSD": 117.28,   // Serbian Dinar
		"PKR": 316.89,   // Pakistani Rupee
		"GTQ": 8.45,     // Guatemalan Quetzal
		"PHP": 61.02,    // Philippines Pes
		"BAM": 1.95,     // Bosnian Marks
		"ALL": 106.37,   // Albanian Lek
		"RUB": 103.55,   // Russian Ruble
		"HNL": 26.45,    // Honduran Lempira
		"JPY": 158.33,   // Japanese Yen
		"TND": 3.36,     // Tunisian Dinar
		"NOK": 11.49,    // Norwegian Krone
		"KRW": 1425.36,  // Won south korean
		"COP": 4258.73,  //Colombian Peso
		"MXN": 18.39,    // Mexico peso
		// Add more currencies and their exchange rates here
	}

	// Check if the currency is in the exchangeRates map
	rate, found := exchangeRates[currency]
	if !found {
		return 0.0, fmt.Errorf("Exchange rate not found for currency: %s", currency)
	}

	// Convert the price to euros
	priceInEUR := price / rate
	return priceInEUR, nil
}
