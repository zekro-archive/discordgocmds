package main

type TestDatabase struct{}

func (db *TestDatabase) GetUserPermissionLevel(userID string, roles []string) (int, error) {
	return 5, nil
}

func (db *TestDatabase) GetGuildPrefix(guildID string) (string, error) {
	return "lel!", nil
}
