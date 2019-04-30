package promutil

import (
	"fmt"
	"math"
	"strconv"
)

func BucketsToStrings(floats []float64) []string {
	strings := make([]string, len(floats))

	for i, v := range floats {
		strings[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}

	return strings
}

func StringsToBuckets(strings []string) ([]float64, error) {
	buckets := make([]float64, len(strings))

	var err error

	for i, s := range strings {
		buckets[i], err = strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
	}

	for i, bucket := range buckets {
		if i < len(buckets)-1 {
			if bucket >= buckets[i+1] {
				return nil, fmt.Errorf(
					"histogram buckets must be in increasing order: %f >= %f",
					bucket, buckets[i+1],
				)
			}
		} else {
			if math.IsInf(bucket, +1) {
				// The +Inf bucket is implicit. Remove it here.
				buckets = buckets[:i]
			}
		}
	}

	return buckets, nil
}
