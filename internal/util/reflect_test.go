package util

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type value struct {
	S  string
	I  int
	P  *string
	SL []string
}

type patch struct {
	S  *string
	I  *int
	P  **string
	SL *[]string
}

type invalidPatch struct {
	S string
	I *string
	P *string
}

type ReflectTestSuite struct {
	suite.Suite
	defaultValue    value
	defaultPatch    patch
	defaultExpected value
}

func (s *ReflectTestSuite) SetupTest() {
	defaultString := "Default"
	s.defaultValue = value{
		S:  defaultString,
		I:  42,
		P:  &defaultString,
		SL: []string{defaultString, "AnotherString"},
	}
	newString := "Patched"
	newStringPtr := &newString
	s.defaultPatch = patch{
		S:  &newString,
		I:  Ptr(420),
		P:  Ptr(newStringPtr),
		SL: Ptr([]string{newString}),
	}
	s.defaultExpected = value{
		S:  newString,
		I:  420,
		P:  newStringPtr,
		SL: []string{newString},
	}
}

func (s *ReflectTestSuite) TestPatchStructFull() {
	err := PatchStruct(&s.defaultValue, s.defaultPatch)
	s.NoError(err)
	s.Equal(s.defaultExpected, s.defaultValue)
}

func (s *ReflectTestSuite) TestPatchStructInterface() {
	var v any
	v = &s.defaultValue
	err := PatchStruct(v, s.defaultPatch)
	s.NoError(err)
	s.Equal(s.defaultExpected, s.defaultValue)
}

func (s *ReflectTestSuite) TestPatchStructValue() {
	err := PatchStruct(s.defaultValue, s.defaultPatch)
	s.Error(err)
}

func (s *ReflectTestSuite) TestPatchStructPtrPatch() {
	err := PatchStruct(&s.defaultValue, &s.defaultPatch)
	s.NoError(err)
	s.Equal(s.defaultExpected, s.defaultValue)
}

func (s *ReflectTestSuite) TestPatchStructInterfacePatch() {
	var patch any
	patch = s.defaultPatch
	err := PatchStruct(&s.defaultValue, patch)
	s.NoError(err)
	s.Equal(s.defaultExpected, s.defaultValue)
}

func (s *ReflectTestSuite) TestPatchStructInterfacePtrPatch() {
	var patch any
	patch = &s.defaultPatch
	err := PatchStruct(&s.defaultValue, patch)
	s.NoError(err)
	s.Equal(s.defaultExpected, s.defaultValue)
}

func (s *ReflectTestSuite) TestPatchStructPartial() {
	newString := "Patched"
	patch := patch{
		S:  &newString,
		SL: Ptr([]string{newString}),
	}
	expected := value{
		S:  newString,
		I:  s.defaultValue.I,
		P:  s.defaultValue.P,
		SL: []string{newString},
	}
	err := PatchStruct(&s.defaultValue, &patch)
	s.NoError(err)
	s.Equal(expected, s.defaultValue)
}

func (s *ReflectTestSuite) TestPatchStructInvalid() {
	newString := "Patched"
	newStringPtr := &newString
	patch := invalidPatch{
		S: newString,
		I: newStringPtr,
		P: newStringPtr,
	}
	expected := s.defaultValue
	err := PatchStruct(&s.defaultValue, &patch)
	s.NoError(err)
	s.Equal(expected, s.defaultValue)
}

func TestReflectTestSuite(t *testing.T) {
	suite.Run(t, new(ReflectTestSuite))
}
