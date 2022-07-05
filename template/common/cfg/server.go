package cfg

type Server struct {
	HTTPAddr string `json:"http_addr" yaml:"http_addr"`
	GRPCAddr string `json:"grpc_addr" yaml:"grpc_addr"`
	AddrName string `json:"addr_name" yaml:"addr_name"`
}
