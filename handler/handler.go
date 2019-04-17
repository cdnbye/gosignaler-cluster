package handler

import (
	"encoding/json"
	"fmt"
	"gosignaler-cluster/model"
	"gosignaler-cluster/signalerconst"
	"gosignaler-cluster/util"
	"github.com/lexkong/log"
	"strings"
)

type Handler interface {
	Handle()
}

type SignalMsg struct {
	To_peer_id string      `json:"to_peer_id"`
	Data       interface{} `json:"data"`
}

func (this *Client) handle(message []byte) {
	//	logrus.Debugf("[Client.handle] %s", string(message))
	action := struct {
		Action string `json:"action"`
	}{}
	if err := json.Unmarshal(message, &action); err != nil {
		//logrus.Errorf("[Client.handle] json.Unmarshal %s", err.Error())
		log.Fatal("[Client.handle] json.Unmarshal %s", err)
		return
	}

	this.CreateHandler(action.Action, message).Handle()
}

func (this *Client) CreateHandler(action string, payload []byte) Handler {
	switch action {
	case "signal":
		msg := SignalMsg{}
		if err := json.Unmarshal(payload, &msg); err != nil {
			//logrus.Errorf("[PullHandler.Handle] json.Unmarshal %s", err.Error())

			return &ExceptionHandler{err.Error()}
		}
		return &SignalHandler{msg, this}
	}

	return &ExceptionHandler{message: fmt.Sprintf("unregnized action %s", action)}
}

type ExceptionHandler struct {
	message string
}

func (this *ExceptionHandler) Handle() {
	log.Warnf("[ExceptionHandler] err %s", this.message)
}

type SignalHandler struct {
	message SignalMsg
	client  *Client
}

func (this *SignalHandler) Handle() {
	//log.Printf("SignalHandler Handle %v", this.message)

	signalResponse := model.SignalResponse{
		Action : "signal",
		FromPeerId : this.client.PeerId,
		Data : this.message.Data,                    //data
	}

	/*b, err := json.Marshal(response)*/
	//log.Printf("sendJsonToClient %v", this.message.To_peer_id)
	//1. get the serverurl from the redis server
	//2.if the serverurl is equal to the localserver ip
	//3.else have to dial the rpc service

	to_peer_id_serverurl := util.Redis.Get(this.message.To_peer_id).Val()

	// if to_peer_id_serverurl  is ""
	if(strings.EqualFold(to_peer_id_serverurl,"")){

		signalResponseFail:=model.SignalResponse{
			Action:"signal",
			FromPeerId: this.message.To_peer_id,//to_peer_id
		}
		log.Warnf("Peer not found:%s",this.message.To_peer_id)
		this.client.hub.SendJsonToClient(this.client.PeerId, signalResponseFail)
		return
	}

	//if to_peer_id is on the local server
	if(strings.EqualFold(to_peer_id_serverurl,util.LOCAL_IP)){
		//send json to client
		this.HandleJsonToClient(signalResponse)
		//dial the rpc service
	}else{
		this.DialRpcService(to_peer_id_serverurl)
	}


}

//send json to client
func (this *SignalHandler) HandleJsonToClient(value interface{}) {

	_, ok := this.client.hub.Clients.Load(this.message.To_peer_id) //Determine if the node is still online
	if ok {
		this.client.hub.SendJsonToClient(this.message.To_peer_id, value)
	} else {
		log.Warnf("Peer not found:%s",this.message.To_peer_id)

		signalResponseFail := model.SignalResponse{
			Action:     "signal",
			FromPeerId: this.message.To_peer_id, //to_peer_id
		}
		this.client.hub.SendJsonToClient(this.client.PeerId, signalResponseFail)
	}
}

//Dial rpc service
func (this *SignalHandler) DialRpcService(to_peer_id_serverurl string) {

	//rpc request param
	request := model.Rpcrequest{
		PeerId:   this.client.PeerId,
		Action:   "signal",
		ToPeerId: this.message.To_peer_id,
		Data:     this.message.Data,
	}

	b, err := json.Marshal(request)
	if (err != nil) {
		log.Error("dialing DialRpcService  serial fail", err)
	}

	client, err := DialHandleSignalService("tcp", fmt.Sprintf(signalerconst.REMOTE_ADDRESS, to_peer_id_serverurl))
	if err != nil {
		log.Error("dialing DialRpcService fail:", err)
	}
	var reply = model.RpcResponse{}

	err = client.HandleSignal(b, &reply)

	if err != nil {
		log.Error("rpc service dial fail",err)
	}

	if(reply.Code==signalerconst.RPC_SUCCESS){
		log.Infof("rpc service success:%s---->%s",reply.FromPeerId,reply.ToPeerId)
	}

	//if rpc service fail,if peer not exists ,maybe you have to tell the original something
	if (reply.Code == signalerconst.RPC_FAIL) {

		log.Warnf("rpc service fail:%s--->%s",reply.FromPeerId,reply.ToPeerId)
		fromPeerId := reply.FromPeerId
		toPeerId := reply.ToPeerId

		signalResponseFail := model.SignalResponse{
			Action:     "signal",
			FromPeerId: toPeerId, //to_peer_id
		}
		this.client.hub.SendJsonToClient(fromPeerId, signalResponseFail)
	}
}
