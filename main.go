package main

import (
	"flag"
	"fmt"
)

var sum int

func init() {
	flag.IntVar(&sum, "sum", 0, "maximum sum")
	flag.Parse()
	if sum < 1 {
		panic("sum cannot be lower than 1")
	}
}

func main() {
	products, err := getProducts("products.csv")
	if err != nil {
		panic(err)
	}
	basket := NewBasket(products, int64(sum))
	best := basket.FindBest()
	printResult(best)
}

func printResult(products []product) {
	var sum int64
	for _, p := range products {
		fmt.Printf("%s - %d %s\n", p.name, p.price, formatRuble(p.price))
		sum += p.price
	}
	fmt.Printf("\nИтого: %d %s", sum, formatRuble(sum))
}

func formatRuble(amount int64) string {
	lastNumber := amount % 10
	if lastNumber == 0 || lastNumber >= 5 || int((amount%100)/10) == 1 {
		return "рублей"
	} else if lastNumber == 1 {
		return "рубль"
	} else {
		return "рубля"
	}
}

type product struct {
	name  string
	price int64
}
