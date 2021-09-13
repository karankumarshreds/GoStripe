package main 

import (
	"os"
	// "log"
)


func main() {
	
	STRIPE_KEY     :=  os.Getenv("STRIPE_KEY")
	STRIPE_SECRET  :=  os.Getenv("STRIPE_SECRET")
	DB_NAME        :=  os.Getenv("DB_NAME")

	// if len(STRIPE_KEY) < 1 || len(STRIPE_SECRET) < 1 || len(DB_NAME) < 1 {
	// 	log.Fatal("STRIPE_KEY|STRIPE_SECRET|DB_NAME missing from environment variables")
	// }

	app := Application{}
	app.Initialize(
		STRIPE_KEY,
		STRIPE_SECRET,
		DB_NAME,
	)
	app.Run()
	// defer app.Channel.Close()
}
