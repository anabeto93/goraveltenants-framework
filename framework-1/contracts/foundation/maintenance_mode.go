package foundation

type MaintenanceMode interface {
	Activate(payload map[string]interface{})
	Deactivate()
	Active() bool
	Data() map[string]interface{}
}
