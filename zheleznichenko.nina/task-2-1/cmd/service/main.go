package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readInt(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, fmt.Errorf("failed to scan token")
	}

	parsedValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, err
	}

	return parsedValue, nil
}

func readString(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		return "", fmt.Errorf("failed to scan token")
	}

	return scanner.Text(), nil
}

func processCondition(operation string, threshold int, minTemp *int, maxTemp *int) {
	switch operation {
	case ">=":
		if threshold > *minTemp {
			*minTemp = threshold
		}
	case "<=":
		if threshold < *maxTemp {
			*maxTemp = threshold
		}
	case "=":
		if threshold > *minTemp {
			*minTemp = threshold
		}
		if threshold < *maxTemp {
			*maxTemp = threshold
		}
	}
}

func processTestCases(scanner *bufio.Scanner, testCasesCount int) error {
	for range testCasesCount {
		operationsCount, err := readInt(scanner)
		if err != nil {
			return err
		}

		minTemp := 15
		maxTemp := 30

		for range operationsCount {
			operation, err := readString(scanner)
			if err != nil {
				return err
			}

			threshold, err := readInt(scanner)
			if err != nil {
				return err
			}

			processCondition(operation, threshold, &minTemp, &maxTemp)

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	testCasesCount, err := readInt(scanner)
	if err != nil {
		return
	}

	err = processTestCases(scanner, testCasesCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
	}
}
