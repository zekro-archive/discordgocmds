package discordgocmds

// DefaultPermissionHandler is the default handler
// for level based permission handling
type DefaultPermissionHandler struct {
	db DatabaseMiddleware
}

// NewDefaultPermissionHandler creates an instance of DefaultPermissionHandler
// getting passed the instance of the database middileware
func NewDefaultPermissionHandler(db DatabaseMiddleware) *DefaultPermissionHandler {
	return &DefaultPermissionHandler{
		db: db,
	}
}

// CheckUserPermission compares the command executing users permission level to the
// required permission level of the command and returns if the user matches the
// required permission.
func (p *DefaultPermissionHandler) CheckUserPermission(cmdArgs *CommandArgs, cmdInstance Command) (bool, error) {
	lvl, err := p.db.GetUserPermissionLevel(cmdArgs.User.ID)
	if err != nil {
		return false, err
	}
	return (cmdInstance.GetPermission() <= lvl), nil
}
