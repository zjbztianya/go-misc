package metrics

// Gauge is metrics gauge.
type Gauge interface {
	Set(value float64)
	Inc()
	Dec()
	Add(delta float64)
	Sub(delta float64)
	Value() float64
}
