package main

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"

	"io/ioutil"
	"os"
)

func main(){
	g := generator.New()

	data,err := ioutil.ReadAll(os.Stdin)
    if err != nil  {
    	g.Error(err,"reading input")
	}

	if err:= proto.Unmarshal(data,g.Request);err!=nil {
		g.Error(err,"parsing  input proto")
	}

	if len(g.Request.FileToGenerate)==0 {
		g.Fail("no files t generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())
    g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	g.GenerateAllFiles()

	data,err = proto.Marchal(g.Response)
	if err !=nil {
		g.Error(err,"failed  to marshal output proto")
	}

	_,err = os.Stdout.Write(data)
	if err !=nil {
		g.Error(err."fail to write output proto")
	}
}