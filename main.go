package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"wol-proxy/websocket"

	"github.com/gin-gonic/gin"
)

func get(c *gin.Context) {

	if websocket.IsWebSocketRequest(c.Request) {
		websocket.Websocket(c.Writer, c.Request, c)
		return
	}

	err := start()
	if err != nil {
		log.Fatal(err)

		c.String(500, " error")
	}

	send(c)
}

func post(c *gin.Context) {

	path := c.Param("path")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
		c.String(500, " error")
		return
	}

	response, err := http.Post("http://192.168.0.106:2342"+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)

		c.Status(http.StatusServiceUnavailable)

	} else {

		extraHeaders := map[string]string{}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		for k, v := range response.Header {
			extraHeaders[k] = strings.Join(v, ",")
			fmt.Println("%s: %s ", k, strings.Join(v, ","))
		}
		c.DataFromReader(response.StatusCode, contentLength, contentType, reader, extraHeaders)
	}
}

func send(c *gin.Context) {
	path := c.Param("path")

	response, err := http.Get("http://192.168.0.106:2342" + path)
	if err != nil {
		log.Fatal(err)

		c.Status(http.StatusServiceUnavailable)

	} else {

		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{}
		for k, v := range response.Header {
			extraHeaders[k] = strings.Join(v, ",")
			fmt.Println("%s: %s ", k, strings.Join(v, ","))
		}

		xValue := response.Header.Get("X-Session-Id")
		fmt.Printf("session id : %s ", xValue)
		c.DataFromReader(response.StatusCode, response.ContentLength, contentType, response.Body, extraHeaders)
	}
}

const (
	Addr = "0.0.0.0:8089"
)

func main() {
	r := gin.Default()
	r.GET("/*path", get)
	r.POST("/*path", post)
	if err := r.Run(Addr); err != nil {
		log.Printf("Error: %v", err)
	}
}

func start() error {
	macAddress, err := net.ParseMAC("f0-de-f1-86-3a-3d")
	if err != nil {
		fmt.Println("Invalid MAC address:", err)
		return err
	}

	//broadcastAddr := os.Args[2]

	magicPacket := buildMagicPacket(macAddress)

	udpAddr, err := net.ResolveUDPAddr("udp", "192.168.0.255:9")
	if err != nil {
		fmt.Println("Error resolving broadcast address:", err)
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error connecting to UDP:", err)
		return err
	}
	defer conn.Close()

	_, err = conn.Write(magicPacket)
	if err != nil {
		fmt.Println("Error sending magic packet:", err)
		return err
	}
	fmt.Println("Magic packet sent successfully to", macAddress.String())
	return nil

}

func buildMagicPacket(mac net.HardwareAddr) []byte {
	payload := make([]byte, 102)

	for i := 0; i < 6; i++ {
		payload[i] = 0xFF
	}

	for i := 1; i <= 16; i++ {
		copy(payload[i*6:], mac)
	}

	magicPacket, err := hex.DecodeString(hex.EncodeToString(payload))
	if err != nil {
		fmt.Println("Error encoding Magic Packet:", err)
		os.Exit(1)
	}

	return magicPacket
}
