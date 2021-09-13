package main

func (app *Application) InitializeRoutes() {
	app.Router.HandleFunc("/publish", app.PublishEvent)
	app.Router.HandleFunc("/", app.RootHandler)
}

