package env

import (
	"flag"
	"os"
)

var (
	ZipkinEndpoint string
	GRPC string
	GRPC_REG_ADDR string
)

func Init() {
	flag.StringVar(&ZipkinEndpoint, "zipkin.endpoint", os.Getenv("ZIPKIN_ENDPOINT"), "zipkin.endpoint is zipkin  endpoint. value: localhost:9411 etc.")
	GRPC = flag.Lookup("grpc").Value.String()
	flag.StringVar(&GRPC_REG_ADDR, "grpc.reg.addr", os.Getenv("GRPC_REG_ADDR"), "usage: -grpc.reg.addr=grpc://127.0.0.1:9000")
}
