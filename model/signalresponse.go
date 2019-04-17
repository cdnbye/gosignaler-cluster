package model

/**
   connections's singal message
 */
type SignalResponse struct {
	Action     string      `json:"action"`       // "signal" or "others"
	FromPeerId string      `json:"from_peer_id"` //
	Data       interface{} `json:"data"`         //message data
}
