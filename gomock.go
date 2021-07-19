package gomock

//go:generate mockgen -source=./internal/tcsrdrserver/portfolio_response_struct.go -destination=./test/mocks/tcsrdrserver/portfolio_response_mock.go -package=tcsrdrserver
//go:generate mockgen -source=./internal/tcsrdrserver/orderbook_response_struct.go -destination=./test/mocks/tcsrdrserver/orderbook_response_mock.go -package=tcsrdrserver
//go:generate mockgen -source=./internal/tcsrdrconfig/config_struct.go -destination=./test/mocks/tcsrdrconfig/config_mock.go -package=tcsrgrconfig
