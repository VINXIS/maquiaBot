package pokemontools

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// APICall makes an API call to the pokemon API
func APICall(endpoint, param string, structure interface{}) (interface{}, error) {
	res, err := http.Get("https://pokeapi.co/api/v2/" + endpoint + "/" + param)
	if err != nil {
		return nil, err
	}

	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(string(byteArray)) == "not found" || strings.HasPrefix(string(byteArray), "<") {
		return nil, errors.New(strings.Title(endpoint) + " **" + param + "** does not exist!")
	}

	marshaller := reflect.New(reflect.TypeOf(structure))
	err = json.Unmarshal(byteArray, marshaller.Interface())
	if err != nil {
		return nil, err
	}

	return marshaller.Interface(), nil
}
