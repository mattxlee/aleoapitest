package nodes

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadNodesFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var nodes []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Check if the trimmedLine has the correct format of IP:Port
		if !strings.Contains(trimmedLine, ":") {
			continue
		}
		parts := strings.Split(trimmedLine, ":")
		if len(parts) != 2 {
			continue
		}
		ip := net.ParseIP(parts[0])
		if ip == nil {
			continue
		}
		port, err := strconv.Atoi(parts[1])
		if err != nil || port < 1 || port > 65535 {
			continue
		}
		if trimmedLine != "" {
			nodes = append(nodes, trimmedLine)
		}
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found in the file")
	}

	return nodes, nil
}

func SaveNodesToFile(nodes []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	for _, node := range nodes {
		_, err := file.WriteString(node + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}

func TestNodesFromUrl(url string) ([]string, error) {
	nodeList, err := retrieveNodeList(url)
	if err != nil {
		fmt.Println("failed to get node list, reason: ", err)
	}
	const replacedPort = "3030"
	var successfulNodes []string
	for _, node := range nodeList {
		ip, err := testNode(node, replacedPort)
		if err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", *ip, err)
			continue
		}
		newNode := net.JoinHostPort(*ip, replacedPort)
		successfulNodes = append(successfulNodes, newNode)
		fmt.Printf("Successfully connected to %s\n", *ip)
	}
	return successfulNodes, nil
}

func retrieveNodeList(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the response code is %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body, reason: %v", err)
	}

	var nodeList []string
	err = json.Unmarshal(body, &nodeList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return nodeList, nil
}

func testNode(entry string, replacedPort string) (*string, error) {
	parts := strings.Split(entry, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid address format: %s", entry)
	}
	ip := parts[0]
	port := replacedPort

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), 5*time.Second)
	if err != nil {
		return &ip, fmt.Errorf("failed to connect to %s: %v", entry, err)
	}
	defer conn.Close()

	return &ip, nil
}
