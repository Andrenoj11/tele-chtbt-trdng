package strategy

import "math"

func EMA(values []float64, period int) []float64 {
	if len(values) == 0 {
		return nil
	}
	k := 2.0 / (float64(period) + 1.0)
	out := make([]float64, len(values))
	prev := values[0]
	out[0] = prev

	for i := 1; i < len(values); i++ {
		next := values[i]*k + prev*(1.0-k)
		out[i] = next
		prev = next
	}
	return out
}

func RSI(values []float64, period int) []float64 {
	if len(values) < period+1 {
		return nil
	}

	var gains, losses float64
	for i := 1; i <= period; i++ {
		diff := values[i] - values[i-1]
		if diff >= 0 {
			gains += diff
		} else {
			losses += -diff
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	rsis := make([]float64, 0, len(values)-period)

	rs := rsRatio(avgGain, avgLoss)
	rsis = append(rsis, 100.0-100.0/(1.0+rs))

	for i := period + 1; i < len(values); i++ {
		diff := values[i] - values[i-1]
		gain := math.Max(0, diff)
		loss := math.Max(0, -diff)

		avgGain = (avgGain*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)

		rs = rsRatio(avgGain, avgLoss)
		rsis = append(rsis, 100.0-100.0/(1.0+rs))
	}

	return rsis
}

func rsRatio(avgGain, avgLoss float64) float64 {
	if avgLoss == 0 {
		return 100
	}
	return avgGain / avgLoss
}
