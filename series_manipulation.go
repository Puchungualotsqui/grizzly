package grizzly

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func (series *Series) RemoveIndexes(indexes []int) {
	if series.DataType == "float" {
		filteredFloats := make([]float64, len(indexes))
		for i, idx := range indexes {
			filteredFloats[i] = series.Float[idx]
		}
		series.Float = filteredFloats
	} else {
		filteredStrings := make([]string, len(indexes))
		for i, idx := range indexes {
			filteredStrings[i] = series.String[idx]
		}
		series.String = filteredStrings
	}
}

func (series *Series) FilterFloatSeries(condition func(float64) bool) []int {
	if series.DataType != "float" {
		panic("FilterFloatSeries only works with float series")
	}

	// Number of Goroutines
	numGoroutines := runtime.NumCPU()
	length := series.GetLength()
	if length == 0 {
		return nil
	}

	var wg sync.WaitGroup
	ch := make(chan []int, numGoroutines)

	// Splitting the work across goroutines
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > length {
			end = length
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			var localFiltered []int
			for j := start; j < end; j++ {
				if condition(series.Float[j]) {
					localFiltered = append(localFiltered, j)
				}
			}
			ch <- localFiltered
		}(start, end)
	}

	// Closing channel after all goroutines finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results
	var filteredIndexes []int
	for indexes := range ch {
		filteredIndexes = append(filteredIndexes, indexes...)
	}

	series.RemoveIndexes(filteredIndexes)
	return filteredIndexes
}

func (series *Series) FilterStringSeries(condition func(string) bool) []int {
	if series.DataType != "string" {
		panic("FilterStringSeries only works with string series")
	}

	// Number of Goroutines
	numGoroutines := runtime.NumCPU()
	length := series.GetLength()
	if length == 0 {
		return nil
	}

	var wg sync.WaitGroup
	ch := make(chan []int, numGoroutines)

	// Splitting the work across goroutines
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > length {
			end = length
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			var localFiltered []int
			for j := start; j < end; j++ {
				if condition(series.String[j]) {
					localFiltered = append(localFiltered, j)
				}
			}
			ch <- localFiltered
		}(start, end)
	}

	// Closing channel after all goroutines finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results
	var filteredIndexes []int
	for indexes := range ch {
		filteredIndexes = append(filteredIndexes, indexes...)
	}

	series.RemoveIndexes(filteredIndexes)
	return filteredIndexes
}

func (series *Series) ConvertStringToFloat() {
	if series.DataType == "float" {
		return
	}
	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	length := len(series.String)
	floatArray := make([]float64, length)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	// Calculate chunk size
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	// Launch multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}

		// Increment the WaitGroup counter
		wg.Add(1)

		// Process the chunk in a goroutine
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if firstErr != nil {
					// Stop if there is an error
					return
				}
				val, err := strconv.ParseFloat(series.String[j], 64)
				if err != nil {
					once.Do(func() {
						firstErr = err
					})
					return
				}
				mu.Lock()
				floatArray[j] = val
				mu.Unlock()
			}
		}(start, end)
	}
	wg.Wait()

	if firstErr != nil {
		fmt.Println("Processing stopped due to error: ", firstErr)
	} else {
		series.Float = floatArray
		series.String = []string{}
	}
}

func (series *Series) ConvertFloatToString() {
	if series.DataType == "string" {
		return
	}
	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	length := len(series.Float)
	stringArray := make([]string, length)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	// Calculate chunk size
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	// Launch multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}

		// Increment the WaitGroup counter
		wg.Add(1)

		// Process the chunk in a goroutine
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if firstErr != nil {
					// Stop if there is an error
					return
				}
				// Convert float to string with desired formatting
				val := strconv.FormatFloat(series.Float[j], 'f', -1, 64)
				mu.Lock()
				stringArray[j] = val
				mu.Unlock()
			}
		}(start, end)
	}
	wg.Wait()

	if firstErr != nil {
		fmt.Println("Processing stopped due to error: ", firstErr)
	} else {
		series.String = stringArray
		series.Float = []float64{}
	}
}

func (series *Series) ReplaceWholeWord(old, new string) {
	if series.DataType == "float" || series.GetLength() == 0 {
		return
	}

	numGoroutines := runtime.NumCPU()
	length := series.GetLength()
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	// Compile the regular expression once
	pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(old))
	re := regexp.MustCompile(pattern)

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				// No mutex needed; each goroutine works on separate slice elements
				series.String[j] = re.ReplaceAllString(series.String[j], new)
			}
		}(start, end)
	}
	wg.Wait() // Wait for all goroutines to complete
}

func (series *Series) Replace(old, new string) {
	if series.DataType == "float" || series.GetLength() == 0 {
		return
	}
	numGoroutines := runtime.NumCPU()
	length := series.GetLength()
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				series.String[j] = strings.ReplaceAll(series.String[j], old, new)
			}
		}(start, end)
	}
	wg.Wait() // Wait for all goroutines to complete
}
