package aggregate

// Summary 汇总多目标成功率。
type Summary struct {
	IDs   []string
	Rates map[string]float64
}

func (b *Book) Summary() Summary {
	s := Summary{Rates: make(map[string]float64, len(b.m))}
	for id := range b.m {
		s.IDs = append(s.IDs, id)
		s.Rates[id] = b.Rate(id)
	}
	return s
}

// Mean 平均成功率。
func Mean(rates map[string]float64) float64 {
	if len(rates) == 0 {
		return 0
	}
	var sum float64
	for _, r := range rates {
		sum += r
	}
	return sum / float64(len(rates))
}
