package main

import "Monitor_Platform/routes"

func main() {

	r := routes.SetupRouter()

	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}

}
