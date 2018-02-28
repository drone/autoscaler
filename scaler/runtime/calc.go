// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package runtime

import "math"

// helper function returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		x = x * -1
	}
	return x
}

// helper function returns the larger of x or y.
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// helper function calculates the different between the existing
// server count and required server count to handle queue volume.
func serverDiff(pending, available, concurrency int) int {
	return int(
		math.Ceil(
			float64(pending-available) /
				float64(concurrency),
		),
	)
}

// helper function adjusts the number of servers to provision
// to ensure it does not exceed the max server count.
func serverCeil(count, additions, ceiling int) int {
	if count+additions >= ceiling {
		additions = ceiling - count
	}
	return additions
}

// helper function adjusts the number of servers to provision
// to ensure the minimum server count is maintained.
func serverFloor(count, deletions, floor int) int {
	if deletions == 0 {
		return 0
	}
	if floor > count-deletions {
		deletions = count - floor
	}
	return max(deletions, 0)
}
