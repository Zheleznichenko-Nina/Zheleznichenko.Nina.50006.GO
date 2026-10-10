package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return
	}

	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			break
		}
		k, err := strconv.Atoi(scanner.Text())
		if err != nil {
			break
		}

		minTemp := 15
		maxTemp := 30

		for j := 0; j < k; j++ {
			if !scanner.Scan() {
				break
			}
			op := scanner.Text()

			if !scanner.Scan() {
				break
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				break
			}

			switch op {
			case ">=":
				if val > minTemp {
					minTemp = val
				}
			case "<=":
				if val < maxTemp {
					maxTemp = val
				}
			case "=":
				if val > minTemp {
					minTemp = val
				}
				if val < maxTemp {
					maxTemp = val
				}
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
	}
}
