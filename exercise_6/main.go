package main

import (
	"cmp"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const formatLineCSV = 2
const divisorValue = 1000
const defaultCurrency = "USD"

type Money struct {
	Value    int
	Currency Currency
}

type Currency string

var ExchangeRates = make(map[Currency]Money)

func (m Money) Equals(other Money) int {
	if m.Currency == other.Currency {
		var result int
		compared := []int{m.Value, other.Value}

		slices.SortFunc(compared, func(a, b int) int {
			result = cmp.Compare(a, b)
			return result
		})

		return result
	}

	var result int
	convertedValue := other.ConvertTo(m.Currency)
	compared := []int{m.Value, convertedValue.Value}

	slices.SortFunc(compared, func(a, b int) int {
		result = cmp.Compare(a, b)
		return result
	})

	return result

}

func (m Money) Add(other Money) Money {
	if m.Currency != other.Currency {
		convertedValue := other.ConvertTo(m.Currency)
		return Money{Value: m.Value + convertedValue.Value, Currency: m.Currency}
	}
	return Money{Value: m.Value + other.Value, Currency: m.Currency}
}

func (m Money) Subtract(other Money) Money {
	if m.Currency != other.Currency {
		convertedValue := other.ConvertTo(m.Currency)
		return Money{Value: m.Value - convertedValue.Value, Currency: m.Currency}
	}
	return Money{Value: m.Value - other.Value, Currency: m.Currency}
}

func (m Money) ConvertTo(newCurrency Currency) Money {
	if m.Currency != newCurrency && newCurrency == defaultCurrency {
		convertedValue := (m.Value * ExchangeRates[m.Currency].Value) / divisorValue
		return Money{Value: convertedValue, Currency: newCurrency}
	} else if m.Currency != newCurrency && newCurrency != defaultCurrency {
		convertedValue := ((m.Value / ExchangeRates[m.Currency].Value) * ExchangeRates[newCurrency].Value) / divisorValue
		return Money{Value: convertedValue, Currency: newCurrency}
	}
	return m
}

func loadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, line := range records {
		if len(line) != formatLineCSV {
			fmt.Println("line skipped", i, line)
			continue
		}

		currency := Currency(line[0])

		rate, err := strconv.Atoi(strings.Replace(line[1], ".", "", -1))
		if err != nil {
			fmt.Println("rate skipped", i, line)
			continue
		}

		ExchangeRates[currency] = Money{Currency: currency, Value: rate}
	}

	return nil
}

func main() {
	err := loadFile("exchange_rates.csv")
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
