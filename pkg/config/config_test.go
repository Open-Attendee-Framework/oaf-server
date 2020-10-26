package config

import (
	"io/ioutil"
	"testing"
)

const testConf = `{
    "host":"localhost",
    "port":"2020",
    "databaseconnection": {
        "driver":"sqlite3",
        "connection":"test.db"
    }
}`

const brokenConf = `{
    "host":"localhost"
    "port":"2020",
    "databaseconnection": 
        "driver":"sqlite3",
        "connection":"test.db"
    }
}`

func TestLoadConfig(t *testing.T) {
	d := []byte(testConf)
	err := ioutil.WriteFile("testConf.json", d, 0644)
	if err != nil {
		t.Fatalf("Error creating demo config: %v", err)
	}
	d = []byte(brokenConf)
	err = ioutil.WriteFile("brokenConf.json", d, 0644)
	if err != nil {
		t.Fatalf("Error creating broken config: %v", err)
	}
	_, err = LoadConfig("thisshouldnotexist.json")
	if err == nil {
		t.Errorf("File not exist Error expected but go non")
	}
	_, err = LoadConfig("brokenConf.json")
	if err == nil {
		t.Errorf("Parsing Error expected but go non")
	}
	c, err := LoadConfig("testConf.json")
	if err != nil {
		t.Errorf("No error expected but got %v", err)
	}
	if c.Host != "localhost" {
		t.Errorf("Expected host to be %v but got %v", "localhost", c.Host)
	}
	if c.Port != "2020" {
		t.Errorf("Expected port to be %v but got %v", "2020", c.Port)
	}
	if c.DatabaseConnection.Connection != "test.db" {
		t.Errorf("Expected connection to be %v but got %v", "test.db", c.DatabaseConnection.Connection)
	}
	if c.DatabaseConnection.Driver != "sqlite3" {
		t.Errorf("Expected driver to be %v but got %v", "sqlite3", c.DatabaseConnection.Driver)
	}
}
