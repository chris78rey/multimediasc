package main

import (
	"log"

	"multimediasc/internal/app"
)

func main() {
	if err := app.RunDesktop(); err != nil {
		log.Fatal(err)
	}
}
