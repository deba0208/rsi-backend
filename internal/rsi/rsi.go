package rsi

func Calculate(
	closes []float64,
	period int,
) float64 {

	if len(closes) < period+1 {
		return 0
	}

	gains := make([]float64, len(closes)-1)
	losses := make([]float64, len(closes)-1)
	for i := 1; i < len(closes); i++ {

		change := closes[i] - closes[i-1]

		if change > 0 {
			gains[i-1] = change
			losses[i-1] = 0
		} else {
			gains[i-1] = 0
			losses[i-1] = -change
		}
	}

	var avgGain float64
	var avgLoss float64

	for i := 0; i < period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}

	avgGain /= float64(period)
	avgLoss /= float64(period)

	for i := period; i < len(gains); i++ {

		avgGain =
			((avgGain * float64(period-1)) + gains[i]) /
				float64(period)

		avgLoss =
			((avgLoss * float64(period-1)) + losses[i]) /
				float64(period)
	}

	if avgLoss == 0 {
		if avgGain == 0 {
			return 50
		}
		return 100
	}

	rs := avgGain / avgLoss

	return 100 - (100 / (1 + rs))
}
