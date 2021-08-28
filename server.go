package main

import (
	"os"
)

func main() {
	a := App{}
	a.Initialize(os.Getenv("HISTORY_SERVER_LISTEN_ADDR"))
	a.Run()
}
