package model

/**
   the response of  dialing rpc service
 */
type RpcResponse struct {
	Code       int    `json:"code"`         //Status code
	Messge     string `json:"message"`      //return message
	FromPeerId string `json:"from_peer_id"` //from peerId
	ToPeerId   string `json:"to_peer_id"`   // to peerId
}
