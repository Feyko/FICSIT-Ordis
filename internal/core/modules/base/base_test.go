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

func (elem ExampleElement) SearchFields() []string {
	return []string{"Name", "Response"}
}

var defaultElement = ExampleElement{
	Name:     "aya",
	Response: "bop",
}

func newModuleWithDefaultCommand(t *testing.T, commands []ExampleElement) *Module[ExampleElement] {
	module := newModuleWithCommands(t, commands)
	createDefaultCommandChecked(t, module)
	return module
}

func createDefaultCommandChecked(t *testing.T, mod *Module[ExampleElement]) {
	checkedCreate(t, mod, defaultElement)
}

func newModuleWithCommands(t *testing.T, commands []ExampleElement) *Module[ExampleElement] {
	module := newDefault[ExampleElement]()

	for _, cmd := range commands {
		checkedCreate(t, module, cmd)
	}

	return module
}

func checkedCreate(t *testing.T, mod *Module[ExampleElement], cmd ExampleElement) {
	err := mod.Create(cmd)
	require.Nil(t, err, "Error when trying to create an element")
}

func checkedDelete(t *testing.T, mod *Module[ExampleElement], name string) {
	err := mod.Delete(name)
	require.Nil(t, err, "Error when trying to delete an element")
}

func checkedGet(t *testing.T, mod *Module[ExampleElement], name string) ExampleElement {
	cmd, err := mod.Get(name)
	require.Nil(t, err, "Error when trying to get an element")
	return cmd
}

func checkedList(t *testing.T, mod *Module[ExampleElement]) []ExampleElement {
	list, err := mod.List()
	require.Nil(t, err, "Error when trying to list elements")
	return list
}

func TestCreate(t *testing.T) {
	module := newModuleWithDefaultCommand(t, []ExampleElement{})

	err := module.Create(defaultElement)

	assert.NotNil(t, err, "Was able to create already existing element")
}

func TestGet(t *testing.T) {
	module := newDefault[ExampleElement]()

	createDefaultCommandChecked(t, module)

	cmd := checkedGet(t, module, defaultElement.Name)

	assert.Equal(t, defaultElement, cmd, "Retrieved element is not the inserted element")
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
		module := newDefault[ExampleElement]()

		for _, cmd := range test {
			checkedCreate(t, module, cmd)
		}

		list := checkedList(t, module)

		assert.Equal(t, test, list, "List of elements isn't the same as the list of created elements")

		if len(list) > 0 {
			list[0].Name = "ThisShouldNotBeAlreadyAName"

			newlist := checkedList(t, module)

			assert.NotEqual(t, newlist, list,
				"Modifying the return value of Elements.List modified the internal value of the Elements module")
		}
	}
}

func TestDelete(t *testing.T) {
	module := newDefault[ExampleElement]()

	createDefaultCommandChecked(t, module)

	checkedDelete(t, module, defaultElement.Name)

	_, err := module.Get(defaultElement.Name)

	assert.NotNil(t, err, "Successfully retrieved an element that should have been deleted")
}

func TestUpdate(t *testing.T) {
	module := newDefault[ExampleElement]()

	newcmd := defaultElement
	newcmd.Name = "new"

	createDefaultCommandChecked(t, module)

	err := module.Update(defaultElement.Name, newcmd)

	assert.Nil(t, err, "Error when trying to update an element")

	cmd := checkedGet(t, module, newcmd.Name)

	assert.Equal(t, newcmd, cmd, "ExampleElement was not updated")
}
