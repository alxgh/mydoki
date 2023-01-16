package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"mydoki/master"
)

func main() {
	totalWorkers := 18

	mu := master.NewUsecase(totalWorkers)

	r := gin.Default()

	master.NewHTTP(r, mu)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}
}
