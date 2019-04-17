package model

/**
  peer 之间传递的信息
 */
type Rpcrequest struct {
	PeerId   string      `json:"peer_id"`    //original peer
	Action   string      `json:"action"`     // "signal" or "others"
	ToPeerId string      `json:"to_peer_id"` // to peerId
	Data     interface{} `json:"data"`       //message data
}
