package utility

// StatusOK is status when no error
const StatusOK = "OK"

// HealthCheck is key-value store that contain health check information
type HealthCheck map[string]string

// NewHealthCheck return new instance of HealthCheck
func NewHealthCheck() HealthCheck {
	return HealthCheck(make(map[string]string))
}

// Add name and error to register as heath check
func (c HealthCheck) Add(name string, err error) HealthCheck {
	status := StatusOK
	if err != nil {
		status = err.Error()
	}

	c[name] = status
	return c
}

// NotOK return true is some error registered
func (c HealthCheck) NotOK() bool {
	for _, value := range c {
		if value != StatusOK {
			return true
		}
	}

	return false
}
