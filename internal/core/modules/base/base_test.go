package base

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type ExampleElement struct {
	Name,
	Response,
	Media string
}

func (elem ExampleElement) ID() string {
	return elem.Name
}

var defaultCommand = ExampleElement{
	Name:     "aya",
	Response: "bop",
}

func newModuleWithDefaultCommand(t *testing.T, commands []ExampleElement) *BasicModule[ExampleElement] {
	module := newModuleWithCommands(t, commands)
	createDefaultCommandChecked(t, module)
	return module
}

func createDefaultCommandChecked(t *testing.T, mod *BasicModule[ExampleElement]) {
	checkedCreate(t, mod, defaultCommand)
}

func newModuleWithCommands(t *testing.T, commands []ExampleElement) *BasicModule[ExampleElement] {
	module := New[ExampleElement]()

	for _, cmd := range commands {
		checkedCreate(t, module, cmd)
	}

	return module
}

func checkedCreate(t *testing.T, mod *BasicModule[ExampleElement], cmd ExampleElement) {
	err := mod.Create(cmd)
	require.Nil(t, err, "Error when trying to create a command")
}

func checkedDelete(t *testing.T, mod *BasicModule[ExampleElement], name string) {
	err := mod.Delete(name)
	require.Nil(t, err, "Error when trying to delete a command")
}

func checkedGet(t *testing.T, mod *BasicModule[ExampleElement], name string) ExampleElement {
	cmd, err := mod.Get(name)
	require.Nil(t, err, "Error when trying to get a command")
	return cmd
}

func checkedList(t *testing.T, mod *BasicModule[ExampleElement]) []ExampleElement {
	list, err := mod.List()
	require.Nil(t, err, "Error when trying to list base")
	return list
}

func TestCreate(t *testing.T) {
	module := newModuleWithDefaultCommand(t, []ExampleElement{})

	err := module.Create(defaultCommand)

	assert.NotNil(t, err, "Was able to create already existing command")
}

func TestGet(t *testing.T) {
	module := New[ExampleElement]()

	createDefaultCommandChecked(t, module)

	cmd := checkedGet(t, module, defaultCommand.Name)

	assert.Equal(t, defaultCommand, cmd, "Retrieved command is not the inserted command")
}

func TestList(t *testing.T) {
	tests := [][]ExampleElement{
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
		module := New[ExampleElement]()

		for _, cmd := range test {
			checkedCreate(t, module, cmd)
		}

		list := checkedList(t, module)

		assert.Equal(t, test, list, "List of base isn't the same as the list of created base")

		if len(list) > 0 {
			list[0].Name = "ThisShouldNotBeAlreadyAName"

			newlist := checkedList(t, module)

			assert.NotEqual(t, newlist, list,
				"Modifying the return value of Elements.List modified the internal value of the Elements module")
		}
	}
}

func TestDelete(t *testing.T) {
	module := New[ExampleElement]()

	createDefaultCommandChecked(t, module)

	checkedDelete(t, module, defaultCommand.Name)

	_, err := module.Get(defaultCommand.Name)

	assert.NotNil(t, err, "Successfully retrieved a command that should have been deleted")
}

func TestUpdate(t *testing.T) {
	module := New[ExampleElement]()

	newcmd := defaultCommand
	newcmd.Name = "new"

	createDefaultCommandChecked(t, module)

	err := module.Update(defaultCommand.Name, newcmd)

	assert.Nil(t, err, "Error when trying to update a command")

	cmd := checkedGet(t, module, newcmd.Name)

	assert.Equal(t, newcmd, cmd, "ExampleElement was not updated")
}
