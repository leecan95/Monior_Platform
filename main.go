package main

import "Monitor_Platform/routes"

func main() {

	r := routes.SetupRouter()

	err := r.Run(":8933")
	if err != nil {
		panic(err)
	}

}
