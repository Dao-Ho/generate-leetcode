package handler

import "strings"

type Command int

const (
	CommandUnknown Command = iota
	CommandRandom
	CommandHelp
)

var commandStrings = map[string]Command{
	"random": CommandRandom,
	"help":   CommandHelp,
}

var commandNames = map[Command]string{
	CommandRandom: "random",
	CommandHelp:   "help",
}

func ParseCommand(text string) Command {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		return CommandUnknown
	}
	cmd, ok := commandStrings[strings.ToLower(parts[1])]
	if !ok {
		return CommandUnknown
	}
	return cmd
}

func (c Command) String() string {
	if name, ok := commandNames[c]; ok {
		return name
	}
	return "unknown"
}

func ListCommands() []string {
	cmds := make([]string, 0, len(commandStrings))
	for name := range commandStrings {
		cmds = append(cmds, name)
	}
	return cmds
}
