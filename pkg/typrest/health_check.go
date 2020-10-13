package typrest

// HealthStatusOK is health status ok
var HealthStatusOK = "OK"

type (
	// HealthMap to handle health-check
	HealthMap map[string]error
)

// HealthStatus of HealthCheck
func HealthStatus(m HealthMap) (healthy bool, details map[string]string) {
	healthy = true
	details = make(map[string]string)

	for name, err := range m {
		if err != nil {
			details[name] = err.Error()
			healthy = false
		} else {
			details[name] = HealthStatusOK
		}
	}
	return
}
