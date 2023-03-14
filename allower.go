package main

import (
	"fmt"
	"github.com/RestartFU/list"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"net"
)

type Allower struct {
	list *list.List
}

func (a *Allower) Allow(addr net.Addr, d login.IdentityData, c login.ClientData) (string, bool) {
	if a.list.Listed(d.DisplayName) {
		fmt.Printf("%s[%v] couldn't join the server: player is banned\n", d.DisplayName, addr.String())
		return "You are banned", false
	}

	fmt.Printf("%s[%v] has joined the server\n", d.DisplayName, addr.String())
	return "", true
}
