package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gosignaler-cluster/util"
	"github.com/lexkong/log"
	"net/http"
	"time"
	//"github.com/sirupsen/logrus"

)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 5120
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { //跨域(Cross domain)
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	PeerId string //唯一标识
}

func (c *Client) sendMessage(msg []byte) error {
	select {
	case c.send <- msg:
		return nil
	default:
		close(c.send)

		//		delete(c.hub.clients, c.id)
		c.hub.Clients.Delete(c.PeerId)
		//      delete from redis
		util.Redis.Del(c.PeerId)

		//		logrus.Errorf("[Client.sendMessage] send to chan error")
		err:=errors.New("send to chan error")
		log.Error("[Client.sendMessage] send to chan error",err)
		return fmt.Errorf("send to chan error")
	}
}

func (this *Client) jsonResponse(value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		//logrus.Errorf("[Client.jsonResponse] Marshal err: %s", err.Error())
		return
	}
	if err := this.sendMessage(b); err != nil {
		//logrus.Errorf("[Client.jsonResponse] sendMessage err: %s", err.Error())
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		//fmt.Println(c.PeerId + "read message" + string(message))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
		c.handle(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			//test
			//fmt.Println(c.PeerId + "write message" + string(message))
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	temp := <-c.send
			//	fmt.Printf("temp %s", message)
			//	w.Write(newline)
			//	w.Write(temp)
			//}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, id string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("[serverWs]",err)
		return
	}
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		PeerId: id,
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

}
