package commands

import (
	"fmt"
)

type Command struct {
	Name,
	Response,
	Media string
}

type CommandsModule struct {
	Commands []Command
}

func NewModule() CommandsModule {
	return CommandsModule{
		Commands: make([]Command, 0, 10),
	}
}

func (mod *CommandsModule) Create(cmd Command) error {
	_, err := mod.Get(cmd.Name)
	if err == nil {
		return fmt.Errorf("command with name '%v' already exists", cmd.Name)
	}
	mod.Commands = append(mod.Commands, cmd)
	return nil
}

func (mod *CommandsModule) Get(name string) (Command, error) {
	cmd, _, err := findCommand(mod.Commands, name)
	if err != nil {
		return Command{}, fmt.Errorf("could not get command: %v", err)
	}
	return cmd, nil
}

func (mod *CommandsModule) List() ([]Command, error) {
	r := make([]Command, len(mod.Commands))
	copy(r, mod.Commands)
	return r, nil
}

func (mod *CommandsModule) Delete(name string) error {
	_, i, err := findCommand(mod.Commands, name)
	if err != nil {
		return fmt.Errorf("could not delete command: %v", err)
	}
	mod.Commands = removeCommand(mod.Commands, i)
	return nil
}

func (mod *CommandsModule) Update(name string, newcmd Command) error {
	_, i, err := findCommand(mod.Commands, name)
	if err != nil {
		return fmt.Errorf("could not update command: %v", err)
	}
	mod.Commands[i] = newcmd
	return nil
}

func removeCommand(s []Command, i int) []Command {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findCommand(s []Command, name string) (Command, int, error) {
	for i, cmd := range s {
		if cmd.Name == name {
			return cmd, i, nil
		}
	}
	return Command{}, 0, fmt.Errorf("no command with name '%v' found", name)
}
