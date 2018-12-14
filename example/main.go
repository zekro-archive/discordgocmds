package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cmds ".."
	"./commands/testcmd"
	"github.com/bwmarrin/discordgo"
)

func main() {
	session, _ := discordgo.New("Bot NDE5ODM3NDcyMDQ2ODQxODY2.Duyx9A.i0FjveuEJhC2lIzC90KPIV_jwDc")

	cmdHandler := cmds.New(session, new(TestDatabase), cmds.NewCmdHandlerOptions())
	cmdHandler.RegisterCommand(new(testcmd.TestCmd))

	session.Open()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	session.Close()
}
