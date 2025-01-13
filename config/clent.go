package config

import "tihai/global"

func initClient() {
	global.UserClients = make(map[uint]*global.Client)
}
