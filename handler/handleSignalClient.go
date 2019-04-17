package handler

import (
	"gosignaler-cluster/model"
	"gosignaler-cluster/signalerconst"
	"net/rpc"
)

type HandleSignalServiceClient struct {
	*rpc.Client
}

func DialHandleSignalService(network, address string) (*HandleSignalServiceClient, error) {

	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HandleSignalServiceClient{Client: c}, nil
}

func (p *HandleSignalServiceClient) HandleSignal(request []byte, reply *model.RpcResponse) error {
	return p.Client.Call(signalerconst.HANDLE_SIGNAL_SERVICE_NAME+".HandleSignal", request, reply)
}
