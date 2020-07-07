package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func main() {
	fp, err := os.Open("./helloworld/FILE")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}

	var fds protobuf.FileDescriptorSet
	if err := proto.Unmarshal(b, &fds); err != nil {
		panic(err)
	}

	locs := map[string]*protobuf.SourceCodeInfo_Location{}

	for _, file := range fds.GetFile() {
		for _, loc := range file.GetSourceCodeInfo().GetLocation() {
			var keys []string
			for _, p := range loc.GetPath() {
				keys = append(keys, strconv.FormatInt(int64(p), 10))
			}
			locs[strings.Join(keys, ".")] = loc
		}

		for si, svc := range file.GetService() {
			path := fmt.Sprintf("6.%d", si)
			for mi, method := range svc.GetMethod() {
				path := fmt.Sprintf("%s.2.%d", path, mi)
				fmt.Printf("%s ===\n", method.GetName())
				loc := locs[path]
				fmt.Printf("l: %v\n", loc.GetLeadingComments())
				fmt.Printf("d: %v\n", loc.GetLeadingDetachedComments())
				fmt.Printf("t: %v\n", loc.GetTrailingComments())
			}
		}

	}
}
