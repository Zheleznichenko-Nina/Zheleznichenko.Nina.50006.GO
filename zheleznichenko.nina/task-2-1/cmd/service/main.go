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

var errConfigRequired = errors.New("config file flag -config is required")

type Config struct {
	InputFile  string `yaml:"inputFile"`
	OutputFile string `yaml:"outputFile"`
}

func (c *Config) UnmarshalYAML(node *yaml.Node) error {
	var m map[string]string
	if err := node.Decode(&m); err != nil {
		return fmt.Errorf("decode yaml to map: %w", err)
	}

	if val, ok := m["inputFile"]; ok {
		c.InputFile = val
	} else if val, ok := m["input-file"]; ok {
		c.InputFile = val
	}

	if val, ok := m["outputFile"]; ok {
		c.OutputFile = val
	} else if val, ok := m["output-file"]; ok {
		c.OutputFile = val
	}

	return nil
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
	NumCode  int     `json:"numCode"`
	CharCode string  `json:"charCode"`
	Value    float64 `json:"value"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		return errConfigRequired
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
		cleanValue := strings.ReplaceAll(item.ValueStr, ",", ".")

		valFloat, err := strconv.ParseFloat(cleanValue, 64)
		if err != nil {
			return nil, fmt.Errorf("parse valute float: %w", err)
		}

		results = append(results, CurrencyResult{
			NumCode:  item.NumCode,
			CharCode: item.CharCode,
			Value:    valFloat,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	return results, nil
}

func saveResults(outputFile string, results []CurrencyResult) error {
	outDir := filepath.Dir(outputFile)

	if outDir != "" {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err := os.WriteFile(outputFile, jsonData, 0600); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}
