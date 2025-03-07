package cli

type Context struct {
	Commands       map[string]Command
	Call           string
	Input          string
	SanitizedInput []string
}
