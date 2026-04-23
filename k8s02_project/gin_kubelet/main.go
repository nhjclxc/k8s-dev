package main

import "gin_kubelet/cmd"

func main() {
	cmd.Execute()
}

/*

go run main.go cron -c ./config/config.yaml
go run main.go http -c ./config/config.yaml
go run main.go start -c ./config/config.yaml

*/
