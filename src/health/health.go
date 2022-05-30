package health

type status string

// Health a Health status struct
type Health struct {
	IsHealthy bool
	Info      map[string]any
}

func NewHealth(isHealthy bool) Health {
	return Health{
		IsHealthy: isHealthy,
		Info:      make(map[string]any),
	}
}
