package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/discordgocmds"
	"github.com/zekroTJA/discordgocmds/example/commands"
)

var (
	flagToken = flag.String("t", "", "bot token")
)

func main() {
	flag.Parse()

	session, _ := discordgo.New("Bot " + *flagToken)

	cmdHandler := discordgocmds.New(session, new(TestDatabase), discordgocmds.NewCmdHandlerOptions())
	cmdHandler.RegisterCommand(new(commands.TestCmd))

	session.Open()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	session.Close()
}
