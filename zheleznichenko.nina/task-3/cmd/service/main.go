package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
	"gopkg.in/yaml.v3"
)

const (
	outputDirPermissions  = 0o750
	outputFilePermissions = 0o600
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode   int     `xml:"NumCode"`
	CharCode  string  `xml:"CharCode"`
	ValueStr  string  `xml:"Value"`
	ParsedVal float64 `xml:"-"`
}

type CurrencyResult struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func readConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config yaml: %w", err)
	}

	return config, nil
}

func readCurrencies(path string) ([]Valute, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input xml file: %w", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		encoding, err := htmlindex.Get(charset)
		if err != nil {
			return nil, fmt.Errorf("failed to detect XML encoding: %w", err)
		}

		return transform.NewReader(input, encoding.NewDecoder()), nil
	}

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to parse xml data: %w", err)
	}

	return valCurs.Valutes, nil
}

func parseCurrencyValues(currencies []Valute) error {
	for index := range currencies {
		value := strings.ReplaceAll(currencies[index].ValueStr, ",", ".")

		parsedValue, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return fmt.Errorf("failed to convert currency value: %w", err)
		}

		currencies[index].ParsedVal = parsedValue
	}

	return nil
}

func convertCurrencies(currencies []Valute) []CurrencyResult {
	sort.Slice(currencies, func(first, second int) bool {
		return currencies[first].ParsedVal > currencies[second].ParsedVal
	})

	results := make([]CurrencyResult, 0, len(currencies))
	for _, currency := range currencies {
		results = append(results, CurrencyResult{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    currency.ParsedVal,
		})
	}

	return results
}

func writeResults(path string, results []CurrencyResult) error {
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	outputDir := filepath.Dir(path)
	if err := os.MkdirAll(outputDir, outputDirPermissions); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(path, output, outputFilePermissions); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := readConfig(*configPath)
	if err != nil {
		panic(err)
	}

	currencies, err := readCurrencies(strings.TrimSpace(config.InputFile))
	if err != nil {
		panic(err)
	}

	if err := parseCurrencyValues(currencies); err != nil {
		panic(err)
	}

	results := convertCurrencies(currencies)
	if err := writeResults(config.OutputFile, results); err != nil {
		panic(err)
	}
}
