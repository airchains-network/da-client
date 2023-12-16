package main

import (
    "github.com/airchains-network/da-client/routes"	
)

func main() {
    r := routes.SetupRouter()
    r.Run(":5050")
}
