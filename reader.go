package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	viperDefaultDelimiter = "."
	defaultTagName        = "def"
	squashTagValue        = ",squash"
	structureTagName      = "mapstructure"
	defaultEnvFileName    = ".env"
)

func Read(config any, opts ...viper.DecoderConfigOption) error {
	viperInstance := viper.New()
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(viperDefaultDelimiter, "_"))

	if _, err := os.Stat(defaultEnvFileName); !os.IsNotExist(err) {
		err := godotenv.Load(defaultEnvFileName)
		if err != nil {
			return err
		}
	}

	viperInstance.AutomaticEnv()
	viperInstance.SetTypeByDefaultValue(true)
	elem := reflect.ValueOf(config).Elem()
	err := setDefaultValues("", viperInstance, &reflect.StructField{}, &elem)
	if err != nil {
		return errors.WithMessage(err, "failed to apply default values")
	}
	err = viperInstance.Unmarshal(config, opts...)
	if err != nil {
		return errors.WithMessage(err, "failed to parse configuration object")
	}
	return nil
}

func setDefaultValues(parentName string, viperInstance *viper.Viper,
	reflectField *reflect.StructField, reflectValue *reflect.Value) error {
	if reflectValue.Kind() == reflect.Struct {
		envValue, ok := reflectField.Tag.Lookup(structureTagName)
		if ok && envValue != squashTagValue {
			if parentName != "" {
				parentName += viperDefaultDelimiter
			}
			parentName += strings.ToUpper(envValue)
		}
		for i := 0; i < reflectValue.NumField(); i++ {
			typeField := reflectValue.Type().Field(i)
			field := reflectValue.Field(i)
			if err := setDefaultValues(parentName, viperInstance,
				&typeField, &field); err != nil {
				return err
			}
		}
		return nil
	}
	defaultValue, _ := reflectField.Tag.Lookup(defaultTagName)
	fieldName, ok := reflectField.Tag.Lookup(structureTagName)
	if ok && fieldName != squashTagValue {
		if parentName != "" {
			fieldName = parentName + viperDefaultDelimiter + strings.ToUpper(fieldName)
		}
		viperInstance.SetDefault(strings.ToUpper(fieldName), defaultValue)
	}
	return nil
}
