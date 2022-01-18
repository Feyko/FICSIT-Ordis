package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var defaultCommand = Command{
	Name:     "aya",
	Response: "ayo",
	Media:    "",
}

func createDefaultCommandChecked(t *testing.T, mod *CommandsModule) {
	checkedCreate(t, mod, defaultCommand)
}

func checkedCreate(t *testing.T, mod *CommandsModule, cmd Command) {
	err := mod.Create(cmd)
	assert.Nil(t, err, "Error when trying to create a command")
}

func checkedDelete(t *testing.T, mod *CommandsModule, name string) {
	err := mod.Delete(name)
	assert.Nil(t, err, "Error when trying to delete a command")
}

func checkedGet(t *testing.T, mod *CommandsModule, name string) Command {
	cmd, err := mod.Get(name)
	assert.Nil(t, err, "Error when trying to get a command")
	return cmd
}

func checkedList(t *testing.T, mod *CommandsModule) []Command {
	list, err := mod.List()
	assert.Nil(t, err, "Error when trying to list commands")
	return list
}

func TestCreate(t *testing.T) {
	module := NewModule()

	createDefaultCommandChecked(t, &module)

	err := module.Create(defaultCommand)

	assert.NotNil(t, err, "Was able to create already existing command")
}

func TestGet(t *testing.T) {
	module := NewModule()

	createDefaultCommandChecked(t, &module)

	cmd := checkedGet(t, &module, defaultCommand.Name)

	assert.Equal(t, defaultCommand, cmd, "Retrieved command is not the inserted command")
}

func TestList(t *testing.T) {
	tests := [][]Command{
		{
			{Name: "1"},
			{Name: "2"},
		},
		{
			{Name: "1"},
		},
		{}, // Empty array
	}

	for _, test := range tests {
		module := NewModule()

		for _, cmd := range test {
			checkedCreate(t, &module, cmd)
		}

		list := checkedList(t, &module)

		assert.Equal(t, test, list, "List of commands isn't the same as the list of created commands")

		if len(list) > 0 {
			list[0].Name = "ThisShouldNotBeAlreadyAName"

			newlist := checkedList(t, &module)

			assert.NotEqual(t, newlist, list,
				"Modifying the return value of Commands.List modified the internal value of the Commands module")
		}
	}
}

func TestDelete(t *testing.T) {
	module := NewModule()

	createDefaultCommandChecked(t, &module)

	checkedDelete(t, &module, defaultCommand.Name)

	_, err := module.Get(defaultCommand.Name)

	assert.NotNil(t, err, "Successfully retrieved a command that should have been deleted")
}

func TestUpdate(t *testing.T) {
	module := NewModule()

	newcmd := defaultCommand
	newcmd.Name = "new"

	createDefaultCommandChecked(t, &module)

	err := module.Update(defaultCommand.Name, newcmd)

	assert.Nil(t, err, "Error when trying to update a command")

	cmd := checkedGet(t, &module, newcmd.Name)

	assert.Equal(t, newcmd, cmd, "Command was not updated")
}
