package node

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)
	defer cancel()
	cmd := exec.CommandContext(
		ctx,
		"geth",
		"--dev",
		"--http",
		"--http.api",
		"eth,web3,personal,net",
		"--http.corsdomain",
		"http://remix.ethereum.org",
	)
	go func() {
		cmd.Stdout = os.Stdout
		fmt.Printf("running node")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}()
	return <-errChan
}
