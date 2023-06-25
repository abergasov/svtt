package routes

import (
	"encoding/json"
	"svtt/internal/entities"

	"github.com/gofiber/contrib/websocket"
)

func (s *Server) dutiesHandler(c *websocket.Conn) {
	var (
		msg []byte
		err error
	)
	go func() {
		for resp := range s.service.GetDutyResponser() {
			if err = c.WriteJSON(resp); err != nil {
				s.log.Error("error while sending message", err)
				break
			}
		}
	}()
	for {
		if _, msg, err = c.ReadMessage(); err != nil {
			s.log.Error("error while reading message", err)
			break
		}
		// handle message, try parse it and send to duty processor
		var message entities.DutyRequest
		if err = json.Unmarshal(msg, &message); err != nil {
			s.log.Error("error while unmarshalling message", err)
			continue
		}
		s.service.HandleDutyRequest(message)
	}
}
