package grizzly

import (
	"math"
)

func arrayMean(data []float64) float64 {
	chain := arrayFloatBase(0, data, func(info float64, result float64) float64 {
		result = result + info
		return result
	})
	var result float64

	for val := range chain {
		result += val
	}
	result /= float64(len(data))
	return result
}

func arrayProduct(data []float64) float64 {
	chain := arrayFloatBase(1, data, func(info float64, result float64) float64 {
		result = result * info
		return result
	})
	result := 1.0

	for val := range chain {
		result *= val
	}
	return result
}

func arraySum(data []float64) float64 {
	chain := arrayFloatBase(0, data, func(info float64, result float64) float64 {
		result = result + info
		return result
	})
	var result float64

	for val := range chain {
		result += val
	}
	return result
}

func arrayVariance(data []float64, meanP ...float64) float64 {
	var mean float64
	if len(meanP) != 0 {
		mean = meanP[0]
	} else {
		mean = arrayMean(data)
	}
	chain := arrayFloatBase(0, data, func(info float64, result float64) float64 {
		diff := info - mean
		return result + diff*diff // Accumulate the squared difference
	})

	var sumOfSquaredDiffs float64
	for val := range chain {
		sumOfSquaredDiffs += val
	}

	// Step 3: Calculate the variance (sum of squared differences divided by the number of elements)
	return sumOfSquaredDiffs / float64(len(data))
}

func arrayMin(data []float64) float64 {
	maxChan := arrayFloatBase(math.MaxFloat64, data, func(info float64, result float64) float64 {
		if info < result {
			result = info
		}
		return result
	})

	minVal := <-maxChan // Initialize minVal with the first value received from the channel
	for val := range maxChan {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

func arrayMax(data []float64) float64 {
	maxChan := arrayFloatBase(math.MaxFloat64*-1, data, func(info float64, result float64) float64 {
		if info > result {
			result = info
		}
		return result
	})

	minVal := <-maxChan // Initialize minVal with the first value received from the channel
	for val := range maxChan {
		if val > minVal {
			minVal = val
		}
	}
	return minVal
}

func arrayMedian(nums []float64) float64 {
	nums = ParallelSortFloat(nums)
	n := len(nums)

	if n%2 == 1 {
		// Odd length, return the middle element
		return nums[n/2]
	} else {
		// Even length, return the average of the two middle elements
		return (nums[n/2-1] + nums[n/2]) / 2.0
	}
}

func arrayCalculatePercentile(nums []float64, percentile float64) float64 {
	size := len(nums)
	index := (percentile / 100.0) * float64(size)
	lower := int(index)
	upper := lower + 1
	weight := index - float64(lower)

	if upper >= size {
		return nums[lower]
	}
	return nums[lower]*(1-weight) + nums[upper]*weight
}
