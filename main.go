package main

import (
    "fmt"
    "os/exec"
    "log"
	"net/http"
	"github.com/gin-gonic/gin"
)


func wake(c *gin.Context){
	cmd := exec.Command("wakeonlan f0:de:f1:86:3a:3d")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
		c.JSON(500)
	} else {
		c.JSON(200)
	}


 
}

const (
	Addr = "127.0.0.1:8089"
)

func main() {
	r := gin.Default()
	r.GET("/wake", wake)
	if err := r.Run(Addr); err != nil {
		log.Printf("Error: %v", err)
	}
}
