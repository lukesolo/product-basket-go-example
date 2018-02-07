package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func getProducts(filename string) ([]product, error) {
	records, err := readData(filename)
	if err != nil {
		return nil, err
	}
	return convertToProducts(records)
}

func readData(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func convertToProducts(table [][]string) ([]product, error) {
	var products []product
	for _, row := range table {
		price, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			return nil, err
		}
		products = append(products, product{row[0], price})
	}
	return products, nil
}
