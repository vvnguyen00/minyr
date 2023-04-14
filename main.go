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

func convertTemperatures() {
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
	fmt.Fprintf(newFile, "Data er basert på gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET); endringen er gjort av %s\n", "Vinh Nguyen")

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

func calculateAverageTemperature(unit string) {
	filepath := "kjevik-temp-celsius-20220318-20230318.csv"
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
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
				log.Fatal(err)
			}

			if unit == "c" {
				sum += tempC
			} else if unit == "f" {
				tempF := conv.CelsiusToFahrenheit(tempC)
				sum += tempF
			}
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	average := sum / float64(count)
	if unit == "c" {
		fmt.Printf("Gjennomsnittstemperaturen er %.2f grader Celsius.\n", average)
	} else if unit == "f" {
		fmt.Printf("Gjennomsnittstemperaturen er %.2f grader Fahrenheit.\n", average)
	}
}

func userInteraction() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Skriv 'minyr' for å starte programmet:")
	cmd, _ := reader.ReadString('\n')
	cmd = strings.TrimSpace(cmd)

	if cmd == "minyr" {
		for {
			fmt.Println("Velg en handling:")
			fmt.Println("  'convert' for å konvertere temperaturer")
			fmt.Println("  'average' for å beregne gjennomsnittstemperatur")
			action, _ := reader.ReadString('\n')
			action = strings.TrimSpace(action)

			if action == "convert" {
				if _, err := os.Stat("kjevik-temp-fahr-20220318-20230318.csv"); err == nil {
					fmt.Println("Filen eksisterer allerede. Vil du generere den på nytt? (j/n)")
					choice, _ := reader.ReadString('\n')
					choice = strings.TrimSpace(choice)

					if choice == "j" {
						convertTemperatures()
					} else if choice == "n" {
						fmt.Println("Konvertering avbrutt.")
					} else {
						fmt.Println("Ugyldig valg. Prøv igjen.")
						continue
					}
				} else {
					convertTemperatures()
				}
			} else if action == "average" {
				fmt.Println("Velg enhet for gjennomsnittstemperatur:")
				fmt.Println("  'c' for grader Celsius")
				fmt.Println("  'f' for grader Fahrenheit")
				unit, _ := reader.ReadString('\n')
				unit = strings.TrimSpace(unit)

				if unit == "c" || unit == "f" {
					calculateAverageTemperature(unit)
				} else {
					fmt.Println("Ugyldig valg. Prøv igjen.")
					continue
				}
			} else {
				fmt.Println("Ugyldig kommando. Prøv igjen.")
				continue
			}

			fmt.Println("Vil du fortsette? (j/n)")
			cont, _ := reader.ReadString('\n')
			cont = strings.TrimSpace(cont)
			if cont != "j" {
				break
			}
		}
	} else {
		fmt.Println("Ugyldig kommando. Skriv 'minyr' for å starte programmet.")
	}
}

func main() {
	userInteraction()
}
