package metrics

// Counter is metrics counter.
type Counter interface {
	Inc()
	Add(delta float64)
}

// Gauge is metrics gauge.
type Gauge interface {
	Set(value float64)
	Inc()
	Dec()
	Add(delta float64)
	Sub(delta float64)
	Value() float64
}
