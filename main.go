package main

import (
	"log"
	"os"

	"github.com/mattxlee/aleoapitest/aleorpc"
	"github.com/mattxlee/aleoapitest/nodes"
	"golang.org/x/exp/rand"
)

const (
	CACHE_FILENAME = "nodes.cache"
	NODELIST_URL   = "http://34.16.24.72:3030/testnet/peers/all"
)

func main() {
	_, err := os.Stat(CACHE_FILENAME)
	var nodeList []string
	if os.IsNotExist(err) {
		log.Printf("Cache file %s does not exist.\n", CACHE_FILENAME)
		nodeList, err = nodes.TestNodesFromUrl(NODELIST_URL)
		if err != nil {
			log.Printf("failed to retrieve nodes from internet, reason: %v\n", err)
			return
		}
		err = nodes.SaveNodesToFile(nodeList, CACHE_FILENAME)
		if err != nil {
			log.Printf("Error saving nodes to file: %v\n", err)
			return
		}
		log.Printf("Successfully saved %d nodes to %s\n", len(nodeList), CACHE_FILENAME)
	} else if err != nil {
		log.Printf("Error checking cache file: %v\n", err)
		return
	} else {
		log.Printf("Cache file %s exists.\n", CACHE_FILENAME)
		nodeList, err = nodes.LoadNodesFromFile(CACHE_FILENAME)
	}

	if len(nodeList) > 0 {
		randomIndex := rand.Intn(len(nodeList))
		chosenNode := nodeList[randomIndex]
		log.Printf("Randomly chosen node: %s\n", chosenNode)

		// create a node
		node := aleorpc.Node{Node: chosenNode}
		height, err := node.GetHeight()
		if err != nil {
			log.Printf("failed to get height from aleo node: %v\n", err)
			return
		}
		log.Printf("aleo height: %d\n", height)
	} else {
		log.Println("No successful nodes to choose from.")
	}
}
