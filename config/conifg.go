package config

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v3"
	"gs/fs"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"package": "config"})

type GSPathConfig struct {
	Gen    string `yaml:"gen"`
	Config string `yaml:"config"`
}

type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	MaxAge           int      `yaml:"max_age"`
}

type SSTConfig struct {
	StacksPath string `yaml:"stacks_path"`
}

type GSConfig struct {
	Module          string       `yaml:"-"`
	Paths           GSPathConfig `yaml:"paths"`
	SST             *SSTConfig   `yaml:"sst,omitempty"`
	WatchExtensions []string     `yaml:"watch_extensions,omitempty"`
	CORS            *CORSConfig  `yaml:"cors,omitempty"`
}

var config *GSConfig

func Get() *GSConfig {
	if config == nil {
		_cnf, err := read()
		if err != nil {
			log.Errorf("Could not read module name: %s", err)
			panic(err)
		}
		config = _cnf
	}
	return config
}

func (c *GSConfig) Reload() {
	_cnf, err := read()
	if err != nil {
		log.Errorf("Could not read module name: %s", err)
		panic(err)
	}
	config = _cnf
}

func Exists() bool {
	exists, _ := fs.Exists("gs.yaml")
	return exists
}

func (c *GSConfig) Write() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return fs.WriteFile("gs.yaml", string(data))
}

func readConfig() *GSConfig {
	cnf := &GSConfig{
		Paths: GSPathConfig{
			Gen:    "gen",
			Config: "config",
		},
	}
	if exists, _ := fs.Exists("gs.yaml"); !exists {
		return cnf
	}

	configData, err := fs.ReadFile("gs.yaml")
	if err != nil {
		return cnf
	}

	_ = yaml.Unmarshal([]byte(configData), cnf)
	cnf.Paths.Gen = strings.Replace(cnf.Paths.Gen, "\\", "/", -1)
	cnf.Paths.Config = strings.Replace(cnf.Paths.Config, "\\", "/", -1)
	log.Debugf("Read config: %+v", cnf)
	return cnf
}

func read() (*GSConfig, error) {
	log.Debugf("Reading config...")
	cnf := readConfig()
	module, err := readModule()
	if err != nil {
		return nil, err
	}
	cnf.Module = module
	return cnf, nil
}

func readModule() (string, error) {
	mod, err := fs.ReadFile("go.mod")
	if err != nil {
		return "", err
	}
	module := modulePath([]byte(mod))
	if module == "" {
		return "", errors.New("could not read the module name")
	}
	return module, nil
}

// Copied from https://github.com/golang/mod/blob/master/modfile/read.go#L882
var (
	slashSlash = []byte("//")
	moduleStr  = []byte("module")
)

// ModulePath returns the module path from the gomod file text.
// If it cannot find a module path, it returns an empty string.
// It is tolerant of unrelated problems in the go.mod file.
func modulePath(mod []byte) string {
	for len(mod) > 0 {
		line := mod
		mod = nil
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, mod = line[:i], line[i+1:]
		}
		if i := bytes.Index(line, slashSlash); i >= 0 {
			line = line[:i]
		}
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, moduleStr) {
			continue
		}
		line = line[len(moduleStr):]
		n := len(line)
		line = bytes.TrimSpace(line)
		if len(line) == n || len(line) == 0 {
			continue
		}

		if line[0] == '"' || line[0] == '`' {
			p, err := strconv.Unquote(string(line))
			if err != nil {
				return "" // malformed quoted string or multiline module path
			}
			return p
		}

		return string(line)
	}
	return "" // missing module path
}
