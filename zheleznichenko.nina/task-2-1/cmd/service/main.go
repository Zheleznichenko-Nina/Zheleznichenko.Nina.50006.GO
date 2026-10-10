package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var errFailedScan = errors.New("failed to scan token")

func readInt(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return 0, fmt.Errorf("scan integer: %w", err)
		}

		return 0, errFailedScan
	}

	parsedValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, fmt.Errorf("parse integer: %w", err)
	}

	return parsedValue, nil
}

func readString(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("scan string: %w", err)
		}

		return "", errFailedScan
	}

	return scanner.Text(), nil
}

func processCondition(operation string, threshold int, minTemp *int, maxTemp *int) {
	switch operation {
	case ">=":
		*minTemp = max(*minTemp, threshold)
	case "<=":
		*maxTemp = min(*maxTemp, threshold)
	case "=":
		*minTemp = max(*minTemp, threshold)
		*maxTemp = min(*maxTemp, threshold)
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

	if err := processTestCases(scanner, testCasesCount); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
	}
}
