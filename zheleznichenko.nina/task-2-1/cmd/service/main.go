package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}

	nStr := strings.TrimSpace(scanner.Text())
	if nStr == "" {
		return
	}
	n, err := strconv.Atoi(nStr)
	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			break
		}
		kStr := strings.TrimSpace(scanner.Text())
		if kStr == "" {
			i--
			continue
		}
		k, err := strconv.Atoi(kStr)
		if err != nil {
			continue
		}

		minTemp := 15
		maxTemp := 30

		for j := 0; j < k; j++ {
			if !scanner.Scan() {
				break
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				j--
				continue
			}

			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			op := parts[0]
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
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
			}

			if minTemp <= maxTemp {
				fmt.Println(maxTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
	}
}
