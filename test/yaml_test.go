package test

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"nat-pernetration/conf"
	"testing"
)

func TestUnMarshalYaml(t *testing.T) {
	s := new(conf.Server)
	b, err := ioutil.ReadFile("../conf/server.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(b, s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}
