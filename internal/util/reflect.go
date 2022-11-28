package util

import (
	"github.com/pkg/errors"
	"reflect"
)

// Takes a pointer to a struct to patch and a pointer to a patch struct.
// Any exported pointer field of the patch struct with a matching exported field (name + type) in the patched struct will have its value copied over
func PatchStruct(v, patch any) error {
	structV, isPtr, ok := getStructValue(reflect.ValueOf(v))
	if !ok || !isPtr {
		return errors.New("The patched value must be a pointer to a struct")
	}

	patchV, _, ok := getStructValue(reflect.ValueOf(patch))
	if !ok {
		return errors.New("The patch must be a struct")
	}

	return patchStructNoCheck(structV, patchV)
}

func patchStructNoCheck(structV, patchV reflect.Value) error {
	structT := structV.Type()
	patchT := patchV.Type()

	numPatchFields := patchT.NumField()
	for i := 0; i < numPatchFields; i++ {
		patchFieldT := patchT.Field(i)
		if !patchFieldT.IsExported() {
			continue
		}
		structFieldT, ok := structT.FieldByName(patchFieldT.Name)
		if !ok || !structFieldT.IsExported() {
			continue
		}
		patchField := patchV.Field(i)
		if !IsNilable(patchField) || patchField.IsNil() {
			continue
		}
		if patchField.Kind() == reflect.Pointer {
			patchField = patchField.Elem()
		}
		structField := structV.FieldByIndex(structFieldT.Index)

		if structField.Kind() == reflect.Struct && patchField.Kind() == reflect.Struct {
			if structField.Type() == patchField.Type() {
				structField.Set(patchField)
				continue
			}

			err := patchStructNoCheck(structField, patchField)
			if err != nil {
				return errors.Wrap(err, "error updating struct field")
			}
			continue
		}

		if structField.Type() != patchField.Type() {
			return errors.Errorf("Incompatible types between value and patch on fields %s.%s.%s and %s.%s.%s",
				structT.PkgPath(), structT.Name(), structFieldT.Name, patchT.PkgPath(), patchT.Name(), patchFieldT.Name)
		}

		structField.Set(patchField)
	}
	return nil
}

func getStructValue(v reflect.Value) (r reflect.Value, isPtr, ok bool) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Pointer {
		if v.Kind() != reflect.Struct {
			return reflect.Value{}, false, false
		}
		return v, false, true
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, false, false
	}
	return v, true, true
}

func IsNil(v any) bool {
	if v == nil {
		return true
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(v).IsNil()
	}
	return false
}

func IsNilable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return true
	default:
		return false
	}
}
