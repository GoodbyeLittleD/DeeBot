package main

import (
	"dvabot/pkg/army_of_the_hell"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	// arena60.Handle()
	// army_of_the_hell.PlayConsole()
	army_of_the_hell.Handle()

	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "#",
		SuperUsers:    []int64{2446629225},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("ws://127.0.0.1:6700", ""),
		},
	}, nil)
}
