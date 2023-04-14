package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vvnguyen00/funtemps/conv"
)

func StringToFloat(s string) (float64, error) {
	// Sjekk om strengen er tom
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

func main() {
	// Åpne filen for lesing
	filepath := "kjevik-temp-celsius-20220318-20230318.csv"
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Åpne filen for skriving
	filename := "kjevik-temp-fahr-20220318-20230318.csv"
	newFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// Les og konverter temperaturer, skriv til ny fil
	scanner := bufio.NewScanner(file)
	firstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		if firstLine {
			// Skriv første linje uendret
			fmt.Fprintln(newFile, line)
			firstLine = false
		} else {
			// Konverter temperatur og skriv til ny fil
			fields := strings.Split(line, ";")
			tempC, err := StringToFloat(fields[3])
			if err != nil {
				log.Fatal(err)
			}
			tempF := conv.CelsiusToFahrenheit(tempC)
			fmt.Fprintf(newFile, "%s;%s;%s;%.2f\n", fields[0], fields[1], fields[2], tempF)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Skriv siste linje
	fmt.Fprintf(newFile, "Data er basert på gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET); endringen er gjort av %s\n", "STUDENTENS_NAVN")

	// Les filen og skriv til terminalen
	outputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	fmt.Println("The output file contains the following data:")
	scanner = bufio.NewScanner(outputFile)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
