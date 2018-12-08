package testcmd

import (
	cmds "../../.."
)

type TestCmd struct {
}

func (t *TestCmd) GetInvokes() []string {
	return []string{"test"}
}

func (t *TestCmd) GetDescription() string {
	return ""
}

func (t *TestCmd) GetHelp() string {
	return ""
}

func (t *TestCmd) GetPermission() int {
	return 4
}

func (t *TestCmd) Exec(args *cmds.CommandArgs) error {
	args.Session.ChannelMessageSend(args.Channel.ID, "test123")
	return nil
}
