package base

import (
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/memrepo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
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

type UpdateExampleElement struct {
	Name,
	Response,
	Media *string
}

func (u UpdateExampleElement) ID() string {
	if u.Name == nil {
		return ""
	}
	return *u.Name
}

type ExampleModule struct {
	Module[ExampleElement]
}

func newDefault() *ExampleModule {
	repo := memrepo.New()
	collection, err := repos.GetCollection(&repo, "Example")
	if err != nil {
		log.Fatalf("Something went horribly wrong and we could not create a new collection in the memrepo: %v", err)
	}
	return &ExampleModule{*New[ExampleElement](collection)}
}

var defaultElement = ExampleElement{
	Name:     "aya",
	Response: "bop",
}

func newModuleWithDefaultCommand(t *testing.T, commands []ExampleElement) *ExampleModule {
	module := newModuleWithCommands(t, commands)
	createDefaultCommandChecked(t, module)
	return module
}

func createDefaultCommandChecked(t *testing.T, mod *ExampleModule) {
	checkedCreate(t, mod, defaultElement)
}

func newModuleWithCommands(t *testing.T, commands []ExampleElement) *ExampleModule {
	module := newDefault()

	for _, cmd := range commands {
		checkedCreate(t, module, cmd)
	}

	return module
}

func checkedCreate(t *testing.T, mod *ExampleModule, cmd ExampleElement) {
	err := mod.Create(cmd)
	require.Nil(t, err, "Error when trying to create an element")
}

func checkedDelete(t *testing.T, mod *ExampleModule, name string) {
	err := mod.Delete(name)
	require.Nil(t, err, "Error when trying to delete an element")
}

func checkedGet(t *testing.T, mod *ExampleModule, name string) ExampleElement {
	cmd, err := mod.Get(name)
	require.Nil(t, err, "Error when trying to get an element")
	return cmd
}

func checkedList(t *testing.T, mod *ExampleModule) []ExampleElement {
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
	module := newDefault()

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
		module := newDefault()

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
	module := newDefault()

	createDefaultCommandChecked(t, module)

	checkedDelete(t, module, defaultElement.Name)

	_, err := module.Get(defaultElement.Name)

	assert.NotNil(t, err, "Successfully retrieved an element that should have been deleted")
}

func TestUpdate(t *testing.T) {
	module := newDefault()

	expected := ExampleElement{Name: defaultElement.Name, Response: "newResponse"}

	newResponse := "newResponse"
	updateElement := UpdateExampleElement{Name: &newResponse}

	createDefaultCommandChecked(t, module)

	err := module.Update(defaultElement.Name, updateElement)

	assert.Nil(t, err, "Error when trying to update an element")

	cmd := checkedGet(t, module, newResponse)

	assert.Equal(t, expected, cmd, "ExampleElement was not updated")
}
