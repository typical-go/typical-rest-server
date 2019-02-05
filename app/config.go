package app

// FIXME: Application Name, Description, Config Prefix
const (
	appName         = "[Service Name]"
	appUsage        = "API Server for [Service Description]"
	appConfigPrefix = "APP"
)

type config struct {
	Address string `envconfig:"ADDRESS"`
}
