package grizzly

import "fmt"

func (series *Series) CountWord(word string) float64 {
	if series.DataType == "float" {
		return 0
	} else {
		return arrayStringCountWord(series.String, word)
	}
}

func (series *Series) GetMax() (float64, error) {
	if series.DataType == "string" {
		return 0, fmt.Errorf("to get max select a float column")
	} else if series.GetLength() == 0 {
		return 0, fmt.Errorf("GetMax requires a non-empty array")
	}
	return arrayMax(series.Float), nil
}

func (series *Series) GetMin() (float64, error) {
	if series.DataType == "string" {
		return 0, fmt.Errorf("to get min select a float column")
	} else if series.GetLength() == 0 {
		return 0, fmt.Errorf("GetMin requires a non-empty array")
	}
	return arrayMin(series.Float), nil
}

func (series *Series) GetMean() (float64, error) {
	if series.DataType == "string" {
		return 0, fmt.Errorf("to get mean select a float column")
	} else if series.GetLength() == 0 {
		return 0, fmt.Errorf("GetMean requires a non-empty array")
	}
	return arrayMean(series.Float), nil
}

func (series *Series) GetMedian() (float64, error) {
	if series.DataType == "string" {
		return 0, fmt.Errorf("to get median select a float column")
	} else if series.GetLength() == 0 {
		return 0, fmt.Errorf("GetMedian requires a non-empty array")
	}
	return arrayMedian(series.Float), nil
}

func (series *Series) GetProduct() (float64, error) {
	if series.DataType == "string" {
		return 0, fmt.Errorf("to get product select a float column")
	} else if series.GetLength() == 0 {
		return 0, fmt.Errorf("GetProduct requires a non-empty array")
	}
	return arrayProduct(series.Float), nil
}

func (series *Series) GetSum() (float64, error) {
	if series.DataType == "string" {
		return 0, fmt.Errorf("to get sum select a float column")
	} else if series.GetLength() == 0 {
		return 0, fmt.Errorf("GetSum requires a non-empty array")
	}
	return arraySum(series.Float), nil
}

func (series *Series) GetVariance() (float64, error) {
	if series.DataType == "string" {
		return 0, fmt.Errorf("to get variance select a float column")
	} else if series.GetLength() == 0 {
		return 0, fmt.Errorf("GetVariance requires a non-empty array")
	}
	return arrayVariance(series.Float), nil
}

func (series *Series) GetNonFloatValues() []string {
	if series.DataType == "float" {
		return []string{}
	}
	return arrayGetNonFloatValues(series.String)
}
