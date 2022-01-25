package api

type Command string

var ValidCommands = map[string]Command{"suscribe": "suscribirme", "unsuscribe": "baja"}

func (c Command) IsValid() bool {
	for _, cmd := range ValidCommands {
		if c == cmd {
			return true
		}
	}
	return false
}

func GetValidCommands() []Command {
	commands := []Command{}
	for _, cmd := range ValidCommands {
		commands = append(commands, cmd)
	}
	return commands
}

func GetValidCommandsString() []string {
	commands := []string{}
	for _, cmd := range ValidCommands {
		commands = append(commands, string(cmd))
	}
	return commands
}

func ToCommand(cmd string) Command {
	return Command(cmd)
}
