package health

type status string

// Health a Health status struct
type Health struct {
	IsHealthy bool
	info      map[string]any
}

func NewHealth(isHealthy bool) Health {
	return Health{
		IsHealthy: isHealthy,
		info:      make(map[string]any),
	}
}
