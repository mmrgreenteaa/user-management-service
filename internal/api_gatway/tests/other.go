package tests

import(
		"github.com/mmrgreenteaa/user-management-service/internal/api_gatway/handlers"
)

const (
	loginTst = "testLogin"
	passTst  = "testPass"
)

var cfg = handlers.ApiGatwayConfig{
	AuthIp:     "27.0.0.1:62480",
	UserMengIP: "127.0.0.1:62380",
	
}
