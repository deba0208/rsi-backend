package rsi

func Calculate(
	closes []float64,
	period int,
) float64 {

	if len(closes) < period+1 {
		return 0.0
	}

	var gain float64
	var loss float64

	for i := 1; i <= period; i++ {

		change := closes[i] - closes[i-1]

		if change > 0 {
			gain += change
		} else {
			loss += -change
		}
	}

	avgGain := gain / float64(period)
	avgLoss := loss / float64(period)

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss

	return 100 - (100 / (1 + rs))
}