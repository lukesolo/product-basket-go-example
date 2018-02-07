package main

import (
	"math"
	"runtime"
	"sort"
)

type basket struct {
	products productSlice
	totalSum int64
}

func NewBasket(products []product, totalSum int64) basket {
	cheapProducts := excludeExpensive(products, totalSum)
	return basket{cheapProducts, totalSum}
}

func (b basket) FindBest() []product {
	// При определённом количестве товаров
	// количество операций для полного перебора стоновится очень большим
	if len(b.products) > 24 {
		return b.greedy()
	}
	return b.bruteForce()
}

// Жадный алгоритм даст приблизительный результат
func (b basket) greedy() []product {
	sort.Sort(b.products)
	var best []product
	var sum int64

	for _, p := range b.products {
		if p.price+sum <= b.totalSum {
			best = append(best, p)
			sum += p.price

			if sum == b.totalSum {
				break
			}
		}
	}
	return best
}

// Метод полного перебора даст точный результат
func (b basket) bruteForce() []product {
	maxBitmap := int64(math.Pow(2.0, float64(len(b.products))))
	cores := int64(runtime.NumCPU())

	results := make(chan int64)
	var from, to, routine int64
	// Выполняем полный перебор на всех ядрах
	for routine = 0; routine < cores; routine++ {
		from = to
		to = from + maxBitmap/cores
		if routine == cores-1 {
			to = maxBitmap
		}
		go b.bruteForceFromTo(from, to, results)
	}

	var bestSum, bestBitmap int64
	for routine = 0; routine < cores; routine++ {
		bitmap := <-results
		sum := b.bitmapSum(bitmap)
		if sum == b.totalSum {
			bestBitmap = bitmap
			break
		}
		if sum > bestSum {
			bestBitmap = bitmap
			bestSum = sum
		}
	}

	best := b.bitmapFilter(bestBitmap)
	sort.Sort(best)
	return best
}

func (b basket) bruteForceFromTo(from, to int64, resultChan chan<- int64) {
	var best, bestBitmap int64

	for bitmap := from; bitmap < to; bitmap++ {
		sum := b.bitmapSum(bitmap)
		if sum == b.totalSum {
			bestBitmap = bitmap
			break
		}
		if sum > best && sum <= b.totalSum {
			best = sum
			bestBitmap = bitmap
		}
	}
	resultChan <- bestBitmap
}

// Сумма по товарам, которые соответствуют битовой карте
func (b basket) bitmapSum(bitmap int64) int64 {
	var sum int64
	var bit int64 = 1

	for _, p := range b.products {
		if bitmap&bit != 0 {
			sum += p.price
		}
		bit <<= 1
	}
	return sum
}

// Выбрать товары по позициям, соответствующим битовой карте
func (b basket) bitmapFilter(bitmap int64) productSlice {
	var filtered productSlice
	var bit int64 = 1

	for _, p := range b.products {
		if bitmap&bit != 0 {
			filtered = append(filtered, p)
		}
		bit <<= 1
	}
	return filtered
}

func excludeExpensive(products []product, totalSum int64) productSlice {
	var cheap productSlice
	for _, p := range products {
		if p.price <= totalSum {
			cheap = append(cheap, p)
		}
	}
	return cheap
}

type productSlice []product

func (ps productSlice) Len() int           { return len(ps) }
func (ps productSlice) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps productSlice) Less(i, j int) bool { return ps[i].price > ps[j].price }
