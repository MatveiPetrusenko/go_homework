package main

import (
	"testing"
)

func TestMoney_Equals(t *testing.T) {
	money1 := Money{Value: 100, Currency: "RUR"}
	money2 := Money{Value: 100, Currency: "RUR"}

	result := money1.Equals(money2)
	if result != 0 {
		t.Errorf("Equals test failed: expected 0, got %d", result)
	}
}

func TestMoney_Add(t *testing.T) {
	money1 := Money{Value: 100, Currency: "RUR"}
	money2 := Money{Value: 100, Currency: "USD"}

	total := money1.Add(money2)

	expectedTotal := Money{Value: 20100, Currency: "RUR"}
	if total != expectedTotal {
		t.Errorf("Add test failed: expected %v, got %v", expectedTotal, total)
	}
}

func TestMoney_Subtract(t *testing.T) {
	money1 := Money{Value: 2000000, Currency: "RUR"}
	money2 := Money{Value: 100, Currency: "EUR"}

	diff := money1.Subtract(money2)
	expectedDiff := Money{Value: 1999900, Currency: "RUR"}
	if diff != expectedDiff {
		t.Errorf("Subtract test failed: expected %v, got %v", expectedDiff, diff)
	}
}

func TestMoney_ConvertTo(t *testing.T) {
	money1 := Money{Value: 2000000, Currency: "RUR"}

	convertedValue := money1.ConvertTo("EUR")
	expectedConvertedValue := Money{Value: 2000, Currency: "EUR"}
	if convertedValue != expectedConvertedValue {
		t.Errorf("ConvertTo test failed: expected %v, got %v", expectedConvertedValue, convertedValue)
	}
}
