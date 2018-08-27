package main

import (
	"flag"
	nm "consensus_layer/node" // node manager
	"fmt"
)

func main() {
	var (
		address = flag.String("address", "", "address of your node")
		target = flag.String("target", "", "address of target peer")
	)
	flag.Parse()
	fmt.Println("address, target: ", *address, *target)
	if *address == "" || *target == "" {
		*address = "0.0.0.0:2000"
		*target = "localhost:2001"
	}
	node := nm.NewNode(*address, []string{*target})
	done := make(chan struct{})
	node.Start()
	<- done
}