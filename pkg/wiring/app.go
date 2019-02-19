package wiring

import "context"

const source = "api-cab-data"

// App defines the application
type App struct {
	Config    Config
	GitCommit string
	Version   string
}

// Run starts the application
func (a App) Run() {
	StartServer(context.Background(),"appname")
}
