package aleorpc

import (
	"fmt"
	"net/rpc"
)

type Node struct {
	Node string
}

func (node *Node) GetHeight() (int, error) {
	client, err := rpc.DialHTTP("tcp", node.Node)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to RPC server: %v", err)
	}
	defer client.Close()

	var height int
	err = client.Call("latest/height", struct{}{}, &height)
	if err != nil {
		return 0, fmt.Errorf("RPC call failed: %v", err)
	}
	return height, nil
}
