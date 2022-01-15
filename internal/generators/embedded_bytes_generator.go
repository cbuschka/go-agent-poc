package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var dataSlice []string

	outfile, err := os.Create("./data.go")
	if err != nil {
		panic(err.Error())
	}
	defer outfile.Close()

	infile, err := ioutil.ReadFile("../../dist/agent")
	if err != nil {
		panic(err.Error())
	}

	outfile.Write([]byte("package master\n\nvar (\n\tagentBytes = []byte{"))

	for _, b := range infile {
		bString := fmt.Sprintf("%v", b)
		dataSlice = append(dataSlice, bString)
	}

	dataString := strings.Join(dataSlice, ", ")
	outfile.WriteString(dataString)
	outfile.Write([]byte("}\n)"))
}
