package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
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
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("Config file flag -config is required")
	}

	configData, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		panic(err)
	}

	xmlData, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	// Создаем XML-декодер с поддержкой кодировки windows-1251
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		panic(err)
	}

	var results []CurrencyResult
	for _, v := range valCurs.Valute {
		cleanValue := strings.ReplaceAll(v.ValueStr, ",", ".")
		valFloat, err := strconv.ParseFloat(cleanValue, 64)
		if err != nil {
			panic(err)
		}

		results = append(results, CurrencyResult{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    valFloat,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	outDir := filepath.Dir(cfg.OutputFile)
	if outDir != "" {
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(cfg.OutputFile, jsonData, 0644)
	if err != nil {
		panic(err)
	}
}
