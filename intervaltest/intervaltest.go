package main

import "log"

func main() {
	// start := 1478791360469
	// end := 1478813411126
	start := 1478791360469000
	end := 1478813411126000

	// write the duration for easier calculation
	value := end - start
	if value < 0 {
		value = 0
	}

	// At every 1hr interval
	ms_in_hr := 3600000

	interval_start := start
	interval_end := interval_start + (ms_in_hr - (interval_start % ms_in_hr))

	// For cases within an hour
	if end <= interval_end {
		interval_end = end
	}

	for interval_end <= end {

		log.Printf("interval [%d:%d], value = %d", interval_start, interval_end, int64((interval_end-interval_start)/1000))

		if interval_end == end {
			break
		} // break out if the last interval

		interval_start = interval_end
		interval_end = interval_end + ms_in_hr

		// i.e. last interval
		if interval_end > end {
			interval_end = end
		}
	}

}
