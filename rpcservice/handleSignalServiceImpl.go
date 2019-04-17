package rpcservice

import (
	"encoding/json"
	"gosignaler-cluster/handler"
	"gosignaler-cluster/model"
	"gosignaler-cluster/signalerconst"
	"log"
)

/****
   the implement fo the rpc service
 */
type HandleSignalService struct {
	Hub *handler.Hub
}

func (p *HandleSignalService) HandleSignal(request []byte, reply *model.RpcResponse) error {

	rpcRequest := model.Rpcrequest{}
	err := json.Unmarshal(request, &rpcRequest)

	if err != nil {
		log.Println("rpc HandleSignal serial the signal message fail", err)
	}

	to_peer_id := rpcRequest.ToPeerId
	_, ok := p.Hub.Clients.Load(to_peer_id)

	signalresponse := model.SignalResponse{
		Action:     "signal",
		FromPeerId: rpcRequest.PeerId,
		Data:       rpcRequest.Data,   //data
	}

	//send message
	if ok {

		p.Hub.SendJsonToClient(to_peer_id, signalresponse)
		*reply = model.RpcResponse{
			Code:   signalerconst.RPC_SUCCESS,
			Messge: "dial the rpc service success",
			FromPeerId: rpcRequest.PeerId,
			ToPeerId:   rpcRequest.ToPeerId,
		}

		// if peer not exists ,maybe you have to tell the original something
	} else {
		*reply = model.RpcResponse{
			Code:       signalerconst.RPC_FAIL,
			Messge:     "dial the rpc service fail",
			FromPeerId: rpcRequest.PeerId,
			ToPeerId:   rpcRequest.ToPeerId,
		}
	}
	return err
}
