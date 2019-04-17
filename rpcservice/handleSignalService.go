package rpcservice

import (
	"gosignaler-cluster/model"
	"gosignaler-cluster/signalerconst"
	"net/rpc"
)

/***
  Defining remote services
 */

type HandleSignalServiceInterface = interface {
	HandleSignal(request []byte, reply *model.RpcResponse) error
}

func RegisterHandleSignalService(service HandleSignalServiceInterface) error {
	return rpc.RegisterName(signalerconst.HANDLE_SIGNAL_SERVICE_NAME, service)
}
