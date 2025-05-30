package utils

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func BindFromJSON(dest any, fileName, path string) error {
	v := viper.New()

	v.SetConfigName(strings.TrimSuffix(fileName, ".json"))
	v.SetConfigType("json") // Pastikan tipe file diatur ke JSON
	v.AddConfigPath(path)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&dest)

	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
	}

	return nil
}

func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)

	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
	}

	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.Itoa(int(valOf.Float()))
		}

		err = os.Setenv(k, val)

		if err != nil {
			logrus.Errorf("failed to unmarshal: %v", err)
			return err
		}

	}

	return nil
}

func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		logrus.Errorf("failed to add remote provider: %v", err)
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("failed to add remote config: %v", err)
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}

	err = SetEnvFromConsulKV(v)

	if err != nil {
		logrus.Errorf("failed to set env from consul kv:: %v", err)
		return err
	}

	return nil
}
