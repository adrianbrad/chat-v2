package configs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadDBconfig(t *testing.T) {
	err := ioutil.WriteFile("mock-config.yaml", []byte("host: test"), 0444)
	if err != nil {
		t.Fatalf("Could not create mock config file, error: %v", err)
	}
	defer func() {
		err := os.Remove("mock-config.yaml")
		if err != nil {
			t.Fatalf("Could not delete mock-config.yaml, error: %s", err.Error())
		}
	}()

	dbConfig, err := LoadDBconfig("mock-config.yaml")
	assert.Equal(t, DBconfig{Host: "test"}, dbConfig)
	assert.Nil(t, err)
}
