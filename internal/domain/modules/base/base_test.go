package base

import (
	"FICSIT-Ordis/internal/domain/modules/auth"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"FICSIT-Ordis/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

func TestExampleModuleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleModuleTestSuite))
}

type ExampleModuleTestSuite struct {
	suite.Suite
	rep repo.Repository[ExampleElement]
	mod *ExampleModule
}

func (s *ExampleModuleTestSuite) SetupSuite() {
	rep, err := test.GetRepo()
	s.Require().NoError(err)
	rep, err = repos.Retype[ExampleElement, id.IDer](rep)
	s.Require().NoError(err)
	s.rep = rep
}

func (s *ExampleModuleTestSuite) SetupTest() {
	collection, err := repos.CreateCollection[ExampleElement](s.rep, "Example")
	s.Require().NoError(err)
	authModule, err := auth.New(auth.Config{Secret: "test-secret"})
	s.Require().NoError(err)
	s.mod = &ExampleModule{*New[ExampleElement](NewDefaultConfigNoPerm(authModule), collection)}
}

func (s *ExampleModuleTestSuite) TearDownTest() {
	err := s.rep.DeleteCollection("Example")
	s.Require().NoError(err)
}

var defaultElement = ExampleElement{
	Name:     "aya",
	Response: "bop",
}

//func newModuleWithDefaultCommand(t *testing.T, commands []ExampleElement) *ExampleModule {
//	module := newModuleWithCommands(t, commands)
//	createDefaultCommandChecked(t, module)
//	return module
//}

func createDefaultCommandChecked(t *testing.T, mod *ExampleModule) {
	checkedCreate(t, mod, defaultElement)
}

//func newModuleWithCommands(t *testing.T, commands []ExampleElement) *ExampleModule {
//	module, err := newDefault()
//	require.NoError(t, err)
//
//	for _, cmd := range commands {
//		checkedCreate(t, module, cmd)
//	}
//
//	return module
//}

func checkedCreate(t *testing.T, mod *ExampleModule, cmd ExampleElement) {
	err := mod.Create(nil, cmd)
	require.NoError(t, err)
}

func checkedDelete(t *testing.T, mod *ExampleModule, name string) {
	err := mod.Delete(nil, name)
	require.NoError(t, err)
}

func checkedGet(t *testing.T, mod *ExampleModule, name string) ExampleElement {
	cmd, err := mod.Get(nil, name)
	require.NoError(t, err)
	return cmd
}

func checkedList(t *testing.T, mod *ExampleModule) []ExampleElement {
	list, err := mod.List(nil)
	require.NoError(t, err)
	return list
}

func (s *ExampleModuleTestSuite) TestCreate() {
	t := s.T()
	createDefaultCommandChecked(t, s.mod)

	err := s.mod.Create(nil, defaultElement)

	assert.Error(t, err, "Was able to create already existing element")
}

func (s *ExampleModuleTestSuite) TestGet() {
	t := s.T()
	createDefaultCommandChecked(t, s.mod)

	cmd := checkedGet(t, s.mod, defaultElement.Name)

	assert.Equal(t, defaultElement, cmd, "Retrieved element is not the inserted element")
}

func (s *ExampleModuleTestSuite) TestList() {
	t := s.T()
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
		for _, cmd := range test {
			checkedCreate(t, s.mod, cmd)
		}

		list := checkedList(t, s.mod)

		assert.Equal(t, test, list, "List of elements isn't the same as the list of created elements")

		if len(list) > 0 {
			list[0].Name = "ThisShouldNotBeAlreadyAName"

			newlist := checkedList(t, s.mod)

			assert.NotEqual(t, newlist, list,
				"Modifying the return value of Elements.List modified the internal value of the Elements module")
		}

		for _, cmd := range test {
			checkedDelete(t, s.mod, cmd.Name)
		}
	}
}

func (s *ExampleModuleTestSuite) TestDelete() {
	t := s.T()
	createDefaultCommandChecked(t, s.mod)

	checkedDelete(t, s.mod, defaultElement.Name)

	_, err := s.mod.Get(nil, defaultElement.Name)

	assert.NotNil(t, err, "Successfully retrieved an element that should have been deleted")
}

func (s *ExampleModuleTestSuite) TestUpdate() {
	t := s.T()
	expected := ExampleElement{Name: defaultElement.Name, Response: "newResponse"}

	newResponse := "newResponse"
	updateElement := UpdateExampleElement{Response: &newResponse}

	createDefaultCommandChecked(t, s.mod)

	_, r, err := s.mod.Update(nil, defaultElement.Name, updateElement)

	require.NoError(t, err, "Error when trying to update an element")
	assert.Equal(t, expected, r, "Returned element is not updated")

	cmd := checkedGet(t, s.mod, defaultElement.Name)

	assert.Equal(t, expected, cmd, "ExampleElement was not updated")
}

func (s *ExampleModuleTestSuite) TestDeleteNonExistentElemErrors() {
	err := s.mod.Delete(nil, "doesntexist")
	s.Require().Error(err)
}
