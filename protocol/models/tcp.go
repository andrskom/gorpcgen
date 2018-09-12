package models

// TCPRequest is struct of tcp request
//   Method is method of rpc
//   Data is field with request
type TCPRequest struct {
	Method string  `json:"method"`
	Data   Request `json:"data"`
}
