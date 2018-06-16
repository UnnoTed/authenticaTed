package generator

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/ungerik/go-dry"
	"gopkg.in/yaml.v1"
)

var wd string

type Config struct {
	Schema  map[string]interface{}
	Dest    string
	API     string
	Package string

	Database string
	File     string
	Debug    bool

	Logging bool
	Fmt     bool
}

// NewConfig creates and load a new config
func NewConfig() (*Config, error) {
	c := &Config{}
	err := c.Load()

	return c, err
}

// Find the config file in the current folder
func (c *Config) Find() (string, bool) {
	return dry.FileFind([]string{"."}, "authenticaTed.yaml", "authenticaTed.yml")
}

// Load reads the config file if found
// then inserts it into the config
func (c *Config) Load() error {
	log.Println("Reading config...")

	// finds the config file
	filePath, found := c.Find()
	if !found {
		return errors.New("Error: authenticaTed config file not found.")
	}

	// reads the config file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// parse it into c
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	// get absolute path from filePath
	d, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// get working dir + destination
	wd = filepath.Dir(d)
	wd = filepath.Join(wd, c.Dest)

	// set package
	c.Package = filepath.Base(c.Dest)

	return nil
}
