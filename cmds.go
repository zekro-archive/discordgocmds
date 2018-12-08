package discordgocmds

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CmdHandler is the main controller of the command
// handler and contains the bot session,  the
// registered commands and registers the message
// handler checking messages for commands.
type CmdHandler struct {
	discordSession         *discordgo.Session
	options                *CmdHandlerOptions
	permHandler            PermissionHandler
	databaseMiddleware     DatabaseMiddleware
	registeredCmds         map[string]Command
	registeredCmdInstances []Command
}

// New creates a new instance of CmdHandler by passing
// the discordgo session, the database middleware instance
// and the command handler options as argument.
func New(session *discordgo.Session, dbMiddleware DatabaseMiddleware, options *CmdHandlerOptions) *CmdHandler {
	c := &CmdHandler{
		discordSession:         session,
		options:                options,
		databaseMiddleware:     dbMiddleware,
		registeredCmds:         make(map[string]Command),
		registeredCmdInstances: make([]Command, 0),
	}
	c.discordSession.AddHandler(c.messageHandler)
	c.discordSession.AddHandler(c.readyHandler)
	return c
}

// RegisterCommand registers a Command class in the
// command handler and will be available for execution.
func (c *CmdHandler) RegisterCommand(cmd Command) {
	c.registeredCmdInstances = append(c.registeredCmdInstances, cmd)
	for _, invoke := range cmd.GetInvokes() {
		c.registeredCmds[invoke] = cmd
	}
}

//////// discordgo event handlers ////////

func (c *CmdHandler) messageHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Message.Author.ID == s.State.User.ID {
		return
	}
	channel, err := s.Channel(e.ChannelID)
	if err != nil {
		// TODO: Do some logging stuff
		return
	}
	if channel.Type != discordgo.ChannelTypeGuildText {
		return
	}
	guildPrefix, err := c.databaseMiddleware.GetGuildPrefix(e.GuildID)
	if err != nil {
		// TODO: Do some logging stuff
	}

	var pre string
	if strings.HasPrefix(e.Message.Content, c.options.Prefix) {
		pre = c.options.Prefix
	} else if guildPrefix != "" && strings.HasPrefix(e.Message.Content, guildPrefix) {
		pre = guildPrefix
	} else {
		return
	}

	contSplit := strings.Fields(e.Message.Content)
	invoke := contSplit[0][len(pre):]
	if c.options.InvokeToLower {
		invoke = strings.ToLower(invoke)
	}

	if cmdInstance, ok := c.registeredCmds[invoke]; ok {
		guild, _ := s.Guild(e.GuildID)
		cmdArgs := &CommandArgs{
			Args:       contSplit[1:],
			Channel:    channel,
			CmdHandler: c,
			Guild:      guild,
			Message:    e.Message,
			Session:    s,
			User:       e.Author,
		}
		hasPerm, err := c.permHandler.CheckUserPermission(cmdArgs, cmdInstance)
		if err != nil {
			// TODO: Print error message for failed permission check
			return
		}
		if !hasPerm {
			// TODO: Send missing permission message to discord channel
			return
		}
		err = cmdInstance.Exec(cmdArgs)
		if err != nil {
			// TODO: Send error message to discord channel
		}
	}
}

func (c *CmdHandler) readyHandler(s *discordgo.Session, e *discordgo.Ready) {
	if c.databaseMiddleware == nil {
		panic("Database middleware must be registered")
	}
	if c.permHandler == nil {
		c.permHandler = NewDefaultPermissionHandler(c.databaseMiddleware)
	}
}
