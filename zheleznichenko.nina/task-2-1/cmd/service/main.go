package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	Valute []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	ValueStr string `xml:"Value"`
}

type CurrencyResult struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if strings.TrimSpace(*configPath) == "" {
		return errors.New("config file flag -config is required")
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		return err
	}

	results, err := processCurrencyData(cfg.InputFile)
	if err != nil {
		return err
	}

	return saveResults(cfg.OutputFile, results)
}

func loadConfig(path string) (Config, error) {
	var cfg Config

	configData, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal yaml: %w", err)
	}

	cfg.InputFile = strings.TrimSpace(cfg.InputFile)
	cfg.OutputFile = strings.TrimSpace(cfg.OutputFile)

	if cfg.InputFile == "" {
		return cfg, errors.New("input-file is required in config")
	}

	if cfg.OutputFile == "" {
		return cfg, errors.New("output-file is required in config")
	}

	return cfg, nil
}

func processCurrencyData(inputFile string) ([]CurrencyResult, error) {
	xmlData, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("read xml file: %w", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	results := make([]CurrencyResult, 0, len(valCurs.Valute))

	for _, item := range valCurs.Valute {
		cleanValue := strings.ReplaceAll(
			strings.TrimSpace(item.ValueStr),
			",",
			".",
		)

		value, err := strconv.ParseFloat(cleanValue, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"parse currency value %q: %w",
				item.ValueStr,
				err,
			)
		}

		results = append(results, CurrencyResult{
			NumCode:  item.NumCode,
			CharCode: item.CharCode,
			Value:    value,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	return results, nil
}

func saveResults(outputFile string, results []CurrencyResult) error {
	outDir := filepath.Dir(outputFile)

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	jsonData = append(jsonData, '\n')

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}
