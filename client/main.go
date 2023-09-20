package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	address := "http://localhost:8080/api"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}
	fmt.Println(ts(), "fetching:", address)
	resp, err := http.Get(address)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("http request failed: %+v", resp))
	}
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("%s http stream failed: %v", ts(), err))
		}
		fmt.Print(ts(), " got: ", line)
	}
}

func ts() string {
	return time.Now().Format(time.DateTime)
}
