package plc

// Common PLC counter maximums for rollover detection.
const (
	max16 int64 = 1<<16 - 1 // 65535
	max32 int64 = 1<<32 - 1 // 4294967295
)

// CalculateDelta computes the production delta between two counter readings.
// Returns the delta and any anomaly type ("reset", "jump", or "").
// Detects 16-bit and 32-bit counter rollover (wrap-around) vs genuine PLC resets.
func CalculateDelta(lastCount, newCount, jumpThreshold int64) (delta int64, anomaly string) {
	if newCount == lastCount {
		return 0, ""
	}
	if newCount < lastCount {
		// Check for 16-bit or 32-bit counter rollover (wrap-around).
		// If the rollover delta is small (below jump threshold), treat as normal production.
		if rollover := tryRollover(lastCount, newCount, jumpThreshold); rollover > 0 {
			return rollover, ""
		}
		// Genuine backward count = PLC restore/reset.
		return newCount, "reset"
	}
	delta = newCount - lastCount
	if delta > jumpThreshold {
		return delta, "jump"
	}
	return delta, ""
}

// tryRollover checks if the counter wrapped around a known PLC bit width.
// Returns the rollover delta if plausible (below jump threshold), or 0.
func tryRollover(lastCount, newCount, jumpThreshold int64) int64 {
	for _, max := range []int64{max16, max32} {
		if lastCount <= max {
			d := (max - lastCount) + newCount + 1
			if d > 0 && d <= jumpThreshold {
				return d
			}
		}
	}
	return 0
}
