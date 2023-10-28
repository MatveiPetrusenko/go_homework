package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Money struct {
	Value    float64
	Currency Currency
}

type Currency string

var ExchangeRates = make(map[Currency]Money)

var formatLineCSV = 2

func (m Money) Equals(other Money) string {
	if m.Currency == other.Currency {
		if m.Value > other.Value {
			return ">"
		} else if m.Value < other.Value {
			return "<"
		} else {
			return "="
		}
	}

	convertedValue := other.ConvertTo(m.Currency)
	if m.Value > convertedValue.Value {
		return ">"
	} else if m.Value < convertedValue.Value {
		return "<"
	} else {
		return "="
	}
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
	if m.Currency != newCurrency && newCurrency == "USD" {
		convertedValue := m.Value * ExchangeRates[m.Currency].Value
		return Money{Value: convertedValue, Currency: newCurrency}
	} else if m.Currency != newCurrency && newCurrency != "USD" {
		convertedValue := (m.Value / ExchangeRates[m.Currency].Value) * ExchangeRates[newCurrency].Value
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
			fmt.Println("Error Line", i, line)
			continue
		}

		currency := Currency(line[0])
		rate, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ExchangeRates[currency] = Money{Currency: currency, Value: rate}
	}

	return nil
}

func main() {
	err := loadFile("exchange_rates.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	money1 := Money{Value: 200.0, Currency: "RUR"}
	money2 := Money{Value: 2.0, Currency: "EUR"}

	//Equals
	result := money1.Equals(money2)
	fmt.Printf("%v %v %v\n", money1.Currency, result, money2.Currency)

	//Add
	total := money1.Add(money2)
	fmt.Println(total)

	//Subtract
	diff := money1.Subtract(money2)
	fmt.Println(diff)

	//ConvertTo
	convertedValue := money1.ConvertTo("EUR")
	fmt.Println(convertedValue)

}
