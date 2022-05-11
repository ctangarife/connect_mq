package main

import (
	rabbit "github.com/ctangarife/connect_mq/consume"
	utl "github.com/ctangarife/connect_mq/utils"
)

func main() {

	conf, err := utl.ReadConf("config/config.yml")
	if err != nil {
		panic(err)
	}
	rabbit.Worker(conf.Monitor.Queue)

}
