package limiter

type Limiter interface {
	Next() chan struct{}
	Ready()
}
