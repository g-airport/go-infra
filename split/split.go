package split

// Range one limit
type Range struct {
	Begin int64
	End   int64
}


func (r Range) Length(interval int64) int64 {
	if interval <= 0 {
		return 0
	}

	diff := r.End - r.Begin
	out := diff/interval + 1
	if out < 0 {
		return 0
	}

	if diff%interval == 0 {
		out--
	}

	return out
}

func Split(r Range, interval int64) []Range {
	var (
		length = r.Length(interval)
		out    = make([]Range, 0, length)
	)

	if length == 0 {
		return []Range{r}
	}

	for i := int64(0); i < length; i++ {
		nr := Range{
			Begin: r.Begin + i*interval,
			End:   r.Begin + (i+1)*interval,
		}

		if nr.End > r.End {
			nr.End = r.End
		}

		out = append(out, nr)
	}

	return out
}
