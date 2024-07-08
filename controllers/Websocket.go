package controllers

import (
	"fmt"
	"gelio/m/middleware"
	"io"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type Server struct {
	Conn map[string]*websocket.Conn
	Mu   sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Conn: make(map[string]*websocket.Conn),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	defer func() {
		s.removeConn(ws)
		ws.Close()
	}()

	var idMessage string
	if err := websocket.Message.Receive(ws, &idMessage); err != nil {
		fmt.Println("Read error:", err)
		return
	}

	ids := strings.Split(idMessage, "-")
	if len(ids) != 2 {
		fmt.Println("Invalid ID message format. Expected format: SenderId-ReceiverId")
		return
	}

	SenderId := ids[0]
	ReceiverId := ids[1]

	s.addConn(SenderId, ws)

	fmt.Println("New incoming connection from Sender ID:", SenderId, "for Receiver ID:", ReceiverId)

	s.readLoop(ws, ReceiverId)

}

func (s *Server) addConn(senderId string, ws *websocket.Conn) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Conn[senderId] = ws
}

func (s *Server) removeConn(ws *websocket.Conn) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for id, conn := range s.Conn {
		if conn == ws {
			delete(s.Conn, id)
			break
		}
	}
}
func (s *Server) readLoop(ws *websocket.Conn, receiverId string) {
	var msg string

	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println(err)
			continue
		}

		Err := s.sendMessageToUser(msg, receiverId)
		if Err != nil {
			fmt.Println(Err)
		}
	}

}

func (s *Server) sendMessageToUser(message string, receiverId string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if value, exists := s.Conn[receiverId]; exists {
		return websocket.Message.Send(value, message)
	}

	return fmt.Errorf("User ID %s not connected", receiverId)
}

func (s *Server) InitializeRoutes(r *gin.Engine) {

	r.GET("/ws", middleware.RequireAuth, func(c *gin.Context) {
		websocket.Handler(s.handleWS).ServeHTTP(c.Writer, c.Request)
	})
}
