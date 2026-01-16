package handler

import (
	"strings"
)

type Command int

const (
	CommandUnknown Command = iota
	CommandRandom
	CommandHelp
)

type ParsedCommand struct {
	Command    Command
	Difficulty string
}

var commandStrings = map[string]Command{
	"random": CommandRandom,
	"help":   CommandHelp,
}

var commandNames = map[Command]string{
	CommandRandom: "random",
	CommandHelp:   "help",
}

var validDifficulties = map[string]string{
	"--easy":  "EASY",
	"-medium": "MEDIUM",
	"-hard":   "HARD",
}

func ParseCommand(text string) ParsedCommand {
	parts := strings.Fields(text)
	result := ParsedCommand{Command: CommandUnknown}

	for _, part := range parts {
		lower := strings.ToLower(part)

		// Skip bot mentions
		if strings.HasPrefix(part, "<@") && strings.HasSuffix(part, ">") {
			continue
		}

		// Check if it's a command
		if cmd, ok := commandStrings[lower]; ok {
			result.Command = cmd
			continue
		}

		// Check if it's a difficulty
		if diff, ok := validDifficulties[lower]; ok {
			result.Difficulty = diff
		}
	}

	return result
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
