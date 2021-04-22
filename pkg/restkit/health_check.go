package restkit

// HealthStatusOK is health status ok
var HealthStatusOK = "OK"

type (
	// HealthMap to handle health-check
	HealthMap map[string]error
)

// Status of HealthCheck
func (m HealthMap) Status() (map[string]string, bool) {
	ok := true
	status := make(map[string]string)

	for name, err := range m {
		if err != nil {
			status[name] = err.Error()
			ok = false
		} else {
			status[name] = HealthStatusOK
		}
	}
	return status, ok
}
