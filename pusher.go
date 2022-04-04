package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Connection struct {
	data chan string
}

type Config struct {
	Host string
	Port int
}

type Pusher struct {
	config  Config
	clients map[string]Connection
	server  *echo.Echo
}

func NewPusher(cfg Config) Pusher {
	pusher := Pusher{
		config:  cfg,
		clients: make(map[string]Connection),
	}
	e := echo.New()
	e.Group("")
	e.GET("/listen/:id", pusher.ConnectionHandler)
	e.POST("/push/:id", pusher.Push)
	pusher.server = e
	return pusher
}

func (p *Pusher) Serve() {
	address := fmt.Sprintf("%s:%d", p.config.Host, p.config.Port)

	if err := p.server.Start(address); !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("ok...")
	}
}

func (p *Pusher) ConnectionHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	enc := json.NewEncoder(c.Response())

	id := c.Param("id")
	dataChannel := make(chan string)
	p.clients[id] = Connection{
		data: dataChannel,
	}
	for {
		select {
		case s := <-dataChannel:
			if err := enc.Encode(s); err != nil {
				return err
			}
			c.Response().Flush()
		}
	}

}

func (p *Pusher) Push(c echo.Context) error {
	id := c.Param("id")
	client, ok := p.clients[id]

	if !ok {
		return c.String(http.StatusNotFound, "client not found")
	}
	client.data <- id
	return c.String(http.StatusOK, "")
}
