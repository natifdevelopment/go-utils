package utils

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	ErrInvalidReflectValue = "invalid reflect.Value"
)

func GetStructValue(instance interface{}, field string) (interface{}, error) {
	val := reflect.ValueOf(instance)
	if !val.IsValid() {
		return nil, errors.New(ErrInvalidReflectValue)
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct but got %s", val.Kind())
	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("no such field: %s", field)
	}

	if fieldVal.Kind() == reflect.Ptr || fieldVal.Kind() == reflect.Interface {
		if fieldVal.IsNil() {
			return nil, nil
		}
	}

	return fieldVal.Interface(), nil
}

func SetStructValue(instance interface{}, field string, value interface{}) error {
	val := reflect.ValueOf(instance)
	if !val.IsValid() {
		return errors.New(ErrInvalidReflectValue)
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("nil pointer dereference")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct but got %s", val.Kind())
	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.IsValid() {
		return fmt.Errorf("no such field: %s", field)
	}
	if !fieldVal.CanSet() {
		return fmt.Errorf("field %s cannot be set", field)
	}

	if value == nil {
		fieldType := fieldVal.Type()
		if fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Map || fieldType.Kind() == reflect.Chan || fieldType.Kind() == reflect.Interface {
			fieldVal.Set(reflect.Zero(fieldType))
			return nil
		}
		return fmt.Errorf("cannot set nil to non-nilable type %s", fieldType)
	}

	newValueVal := reflect.ValueOf(value)
	if fieldVal.Type() != newValueVal.Type() {
		return fmt.Errorf("provided value type %s doesn't match field type %s", newValueVal.Type(), fieldVal.Type())
	}

	fieldVal.Set(newValueVal)
	return nil
}

func ReInitiateStruct(instance interface{}) interface{} {
	instanceType := reflect.TypeOf(instance)

	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	newInstance := reflect.New(instanceType).Interface()
	return newInstance
}

func MakeSliceOfType(modelType interface{}) interface{} {
	modelTypeValue := reflect.ValueOf(modelType).Elem()
	sliceType := reflect.SliceOf(modelTypeValue.Type())
	return reflect.New(sliceType).Interface()
}

func GetModelFields(instance interface{}) ([]string, error) {
	val := reflect.ValueOf(instance)
	val = dereferencePointer(val)

	if err := validateStructValue(val); err != nil {
		return nil, err
	}

	fields := []string{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		if isEmbeddedField(field) {
			embeddedFields, err := getEmbeddedFields(fieldValue)
			if err != nil {
				return nil, err
			}
			fields = append(fields, embeddedFields...)
			continue
		}

		fields = append(fields, field.Name)
	}

	return fields, nil
}

func dereferencePointer(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

func validateStructValue(val reflect.Value) error {
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("instance must be a struct or a pointer to a struct")
	}
	if !val.IsValid() {
		return errors.New(ErrInvalidReflectValue)
	}
	return nil
}

func isEmbeddedField(field reflect.StructField) bool {
	return field.Name == "BaseModel" || field.Anonymous
}

func getEmbeddedFields(fieldValue reflect.Value) ([]string, error) {
	if fieldValue.Kind() == reflect.Struct {
		return GetModelFields(fieldValue.Interface())
	}
	if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
		return GetModelFields(fieldValue.Elem().Interface())
	}
	return []string{}, nil
}
