package handler

import (
	"encoding/json"
	"gosignaler-cluster/signalerconst"
	"gosignaler-cluster/util"
	"github.com/lexkong/log"
	"sync"
	"time"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	//clients map[*Client]bool
	Clients sync.Map

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	//count of client
	ClientNum uint16
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// send json object to a client with peerId
func (this *Hub) SendJsonToClient(peerId string, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		//logrus.Errorf("[Client.jsonResponse] Marshal err: %s", err.Error())
		return
	}
	client, ok := this.Clients.Load(peerId)
	if !ok {
		log.Error("sendJsonToClient error",err)
		return
	}
	if err := client.(*Client).sendMessage(b); err != nil {
		//logrus.Errorf("[Client.jsonResponse] sendMessage err: %s", err.Error())
	}
	//if err := client.(*Client).conn.WriteJSON(value); err != nil {
	//	//logrus.Errorf("[Client.jsonResponse] sendMessage err: %s", err.Error())
	//}
}

func (this *Hub) Run() {
	for {
		select {
		case client := <-this.register:
			this.doRegister(client)
		case client := <-this.unregister:
			this.doUnregister(client)
		}
	}
}

func (this *Hub) doRegister(client *Client) {
	//	logrus.Debugf("[Hub.doRegister] %s", client.id)
	if client.PeerId != "" {
		this.Clients.Store(client.PeerId, client)
		this.ClientNum ++
		util.Redis.Set(client.PeerId, util.LOCAL_IP, time.Second*signalerconst.REDIS_EXPIRE)
	}
}

func (this *Hub) doUnregister(client *Client) {
	//	logrus.Debugf("[Hub.doUnregister] %s", client.id)

	if client.PeerId == "" {
		return
	}

	_, ok := this.Clients.Load(client.PeerId)

	if ok {
		//delRecordCh <- client.id
		this.Clients.Delete(client.PeerId)
		close(client.send)
		this.ClientNum --
		util.Redis.Del(client.PeerId)
	}

}
