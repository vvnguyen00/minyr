package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/vvnguyen00/funtemps/conv"
)

// TestNumberOfLines tester antall linjer i den konverterte filen
func TestNumberOfLines(t *testing.T) {
	filename := "kjevik-temp-fahr-20220318-20230318.csv"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	want := 16757
	if lines != want {
		t.Errorf("got %d lines, want %d", lines, want)
	}
}

// TestConvertTemperature tester konvertering av temperaturer fra Celsius til Fahrenheit
func TestConvertTemperature(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"Kjevik;SN39040;18.03.2022 01:50;6", "Kjevik;SN39040;18.03.2022 01:50;42.80"},
		{"Kjevik;SN39040;07.03.2023 18:20;0", "Kjevik;SN39040;07.03.2023 18:20;32.00"},
		{"Kjevik;SN39040;08.03.2023 02:20;-11", "Kjevik;SN39040;08.03.2023 02:20;12.20"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Convert %s", tc.input), func(t *testing.T) {
			fields := strings.Split(tc.input, ";")
			tempC, err := StringToFloat(fields[3])
			if err != nil {
				t.Fatal(err)
			}
			tempF := conv.CelsiusToFahrenheit(tempC)
			output := fmt.Sprintf("%s;%s;%s;%.2f", fields[0], fields[1], fields[2], tempF)

			if output != tc.want {
				t.Errorf("got %q, want %q", output, tc.want)
			}
		})
	}
}

// TestFinalLine tester at siste linje i den konverterte filen inneholder riktig informasjon
func TestFinalLine(t *testing.T) {
	filename := "kjevik-temp-fahr-20220318-20230318.csv"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	want := "Data er basert på gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET); endringen er gjort av STUDENTENS_NAVN"

	if !strings.HasPrefix(lastLine, want) {
		t.Errorf("got %q, want %q", lastLine, want)
	}
}
// TestAverageTemperature tester om gjennomsnittstemperaturen er riktig
func TestAverageTemperature(t *testing.T) {
	filepath := "kjevik-temp-celsius-20220318-20230318.csv"
	file, err := os.Open(filepath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	firstLine := true
	var sum float64
	var count int

	for scanner.Scan() {
		line := scanner.Text()
		if firstLine {
			firstLine = false
		} else {
			fields := strings.Split(line, ";")
			tempC, err := StringToFloat(fields[3])
			if err != nil {
				t.Fatal(err)
			}
			sum += tempC
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	average := sum / float64(count)

	want := 8.56
	allowedDifference := 0.01

	if average < want-allowedDifference || average > want+allowedDifference {
		t.Errorf("got average temperature of %.2f, want %.2f", average, want)
	}
}


