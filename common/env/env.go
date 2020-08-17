package env

import (
	"flag"
	"os"
)

var (
	ZipkinEndpoint string
)

func init() {
	flag.StringVar(&ZipkinEndpoint, "zipkin.endpoint", os.Getenv("ZIPKIN_ENDPOINT"), "zipkin.endpoint is zipkin  endpoint. value: localhost:9411 etc.")
}
