// Copyright 2019 ChainSafe Systems (ON) Corp.
// This file is part of gossamer.
//
// The gossamer library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The gossamer library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the gossamer library. If not, see <http://www.gnu.org/licenses/>.

package dot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"unicode"

	"github.com/ChainSafe/gossamer/node/gssmr"
	"github.com/ChainSafe/gossamer/node/ksmcc"

	log "github.com/ChainSafe/log15"
	"github.com/naoina/toml"
)

// Config is a collection of configurations throughout the system
type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Account AccountConfig `toml:"account"`
	Core    CoreConfig    `toml:"core"`
	Network NetworkConfig `toml:"network"`
	RPC     RPCConfig     `toml:"rpc"`
}

// GlobalConfig is to marshal/unmarshal toml global config vars
type GlobalConfig struct {
	Name    string `toml:"name"`
	ID      string `toml:"id"`
	Config  string `toml:"config"`
	Genesis string `toml:"genesis"`
	DataDir string `toml:"datadir"`
}

// AccountConfig is to marshal/unmarshal account config vars
type AccountConfig struct {
	Key    string `toml:"key"`
	Unlock string `toml:"unlock"`
}

// NetworkConfig is to marshal/unmarshal toml network config vars
type NetworkConfig struct {
	Port        uint32   `toml:"port"`
	Bootnodes   []string `toml:"bootnodes"`
	ProtocolID  string   `toml:"protocol"`
	Roles       byte     `toml:"roles"`
	NoBootstrap bool     `toml:"nobootstrap"`
	NoMDNS      bool     `toml:"nomdns"`
}

// CoreConfig is to marshal/unmarshal toml core config vars
type CoreConfig struct {
	Authority bool `toml:"authority"`
}

// RPCConfig is to marshal/unmarshal toml RPC config vars
type RPCConfig struct {
	Enabled bool     `toml:"enabled"`
	Port    uint32   `toml:"port"`
	Host    string   `toml:"host"`
	Modules []string `toml:"modules"`
}

// String will return the json representation for a Config
func (c *Config) String() string {
	out, _ := json.MarshalIndent(c, "", "\t")
	return string(out)
}

// GssmrConfig returns a new test configuration using the provided datadir
func GssmrConfig() *Config {
	return &Config{
		Global: GlobalConfig{
			Name:    gssmr.DefaultName,
			ID:      gssmr.DefaultID,
			Config:  gssmr.DefaultConfig,
			Genesis: gssmr.DefaultGenesis,
			DataDir: gssmr.DefaultDataDir,
		},
		Account: AccountConfig{
			Key:    gssmr.DefaultKey,
			Unlock: gssmr.DefaultUnlock,
		},
		Core: CoreConfig{
			Authority: gssmr.DefaultAuthority,
		},
		Network: NetworkConfig{
			Port:        gssmr.DefaultNetworkPort,
			Bootnodes:   gssmr.DefaultNetworkBootnodes,
			ProtocolID:  gssmr.DefaultNetworkProtocolID,
			Roles:       gssmr.DefaultRoles,
			NoBootstrap: gssmr.DefaultNoBootstrap,
			NoMDNS:      gssmr.DefaultNoMDNS,
		},
		RPC: RPCConfig{
			Port:    gssmr.DefaultRPCHTTPPort,
			Host:    gssmr.DefaultRPCHTTPHost,
			Modules: gssmr.DefaultRPCModules,
		},
	}
}

// KsmccConfig returns a "ksmcc" node configuration
func KsmccConfig() *Config {
	return &Config{
		Global: GlobalConfig{
			Name:    ksmcc.DefaultName,
			ID:      ksmcc.DefaultID,
			Config:  ksmcc.DefaultConfig,
			Genesis: ksmcc.DefaultGenesis,
			DataDir: ksmcc.DefaultDataDir,
		},
		Account: AccountConfig{
			Key:    ksmcc.DefaultKey,
			Unlock: ksmcc.DefaultUnlock,
		},
		Core: CoreConfig{
			Authority: ksmcc.DefaultAuthority,
		},
		Network: NetworkConfig{
			Port:        ksmcc.DefaultNetworkPort,
			Bootnodes:   ksmcc.DefaultNetworkBootnodes,
			ProtocolID:  ksmcc.DefaultNetworkProtocolID,
			Roles:       ksmcc.DefaultRoles,
			NoBootstrap: ksmcc.DefaultNoBootstrap,
			NoMDNS:      ksmcc.DefaultNoMDNS,
		},
		RPC: RPCConfig{
			Port:    ksmcc.DefaultRPCHTTPPort,
			Host:    ksmcc.DefaultRPCHTTPHost,
			Modules: ksmcc.DefaultRPCModules,
		},
	}
}

// LoadConfig loads the values from the toml configuration file into the provided configuration
func LoadConfig(cfg *Config, fp string) error {
	fp, err := filepath.Abs(fp)
	if err != nil {
		log.Error("[dot] Failed to create absolute path for toml configuration file", "error", err)
		return err
	}

	file, err := os.Open(filepath.Clean(fp))
	if err != nil {
		log.Error("[dot] Failed to open toml configuration file", "error", err)
		return err
	}

	var tomlSettings = toml.Config{
		NormFieldName: func(rt reflect.Type, key string) string {
			return key
		},
		FieldToKey: func(rt reflect.Type, field string) string {
			return field
		},
		MissingField: func(rt reflect.Type, field string) error {
			link := ""
			if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
				link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
			}
			return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
		},
	}

	if err = tomlSettings.NewDecoder(file).Decode(&cfg); err != nil {
		log.Error("[dot] Failed to decode configuration", "error", err)
		return err
	}

	return nil
}

// ExportConfig exports a dot configuration to a toml configuration file
func ExportConfig(cfg *Config, fp string) *os.File {
	var (
		newFile *os.File
		err     error
		raw     []byte
	)

	if raw, err = toml.Marshal(*cfg); err != nil {
		log.Error("[dot] Failed to marshal configuration", "error", err)
		os.Exit(1)
	}

	newFile, err = os.Create(filepath.Clean(fp))
	if err != nil {
		log.Error("[dot] Failed to create configuration file", "error", err)
		os.Exit(1)
	}

	_, err = newFile.Write(raw)
	if err != nil {
		log.Error("[dot] Failed to write to configuration file", "error", err)
		os.Exit(1)
	}

	if err := newFile.Close(); err != nil {
		log.Error("[dot] Failed to close configuration file", "error", err)
		os.Exit(1)
	}

	return newFile
}