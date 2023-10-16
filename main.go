package main

import "Monitor_Platform/routes"

func main() {

	r := routes.SetupRouter()

	err := r.Run("localhost:8933")
	if err != nil {
		panic(err)
	}

}
