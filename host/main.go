package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/loopholelabs/scale"
	"github.com/loopholelabs/scale/scalefunc"
	signature "go.dev.scale.sh/v1/shivanshvij/kafkatransforms/latest/host"
)

//go:embed local-kafkatransforms-latest.scale
var functionData []byte

const inputRecord = "ewogICAgIm9yZGVyX2lkIjogMSwKICAgICJwcm9kdWN0X2lkIjogMTIsCiAgICAicXVhbnRpdHkiOiAyLAogICAgImFtb3VudCI6IDU2LjAsCiAgICAic2hpcHBpbmciOiAxNS4wLAogICAgInRheCI6IDIuMAp9"

func main() {
	sf := new(scalefunc.Schema)
	err := sf.Decode(functionData)
	if err != nil {
		panic(err)
	}

	s, err := scale.New(scale.NewConfig(signature.New).WithFunction(sf))
	if err != nil {
		panic(err)
	}

	i, err := s.Instance()
	if err != nil {
		panic(err)
	}

	decodedRecord, err := base64.StdEncoding.DecodeString(inputRecord)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded Input Record: %s\n", decodedRecord)

	sig := signature.New()
	sig.Context.Input.Topic = "test-topic"
	sig.Context.Input.Partition = 32
	sig.Context.Input.Offset = 64
	sig.Context.Input.Record = decodedRecord

	err = i.Run(context.Background(), sig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded Output Record: %s\n", sig.Context.Output.Record)
}
