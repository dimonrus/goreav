package main

import (
	"flag"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"goreav/gen"
)

type CliArgs struct {
	Generator string
	File      string
}

var args CliArgs

func init() {
	flag.StringVar(&args.File, "f", "", " config")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	//Load template file
	data, err := ioutil.ReadFile(args.File)
	check(err)

	//Template struct
	var template = make(gen.AppTemplate)

	//Unmarshal yaml
	err = yaml.Unmarshal([]byte(data), template)
	check(err)

	//Parse template
	err = gen.ParseTemplate(template)
	check(err)
}
