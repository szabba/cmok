package configuration

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/szabba/cmok/auth"
	"github.com/szabba/cmok/userlist"
)

type Config struct {
	userlist.AuthConfig
	userlist.AccessConfig
}

func Parse(r io.Reader) (Config, error) {
	wrap := func(err error) error {
		return fmt.Errorf("problem parsing config: %s", err)
	}
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	rc := new(rawConfig)
	err := dec.Decode(rc)
	if err != nil {
		return Config{}, wrap(err)
	}
	if rc.AuthConfig == nil {
		return Config{}, wrap(fmt.Errorf("missing users config"))
	}
	if rc.AccessConfig == nil {
		return Config{}, wrap(fmt.Errorf("missing permissions config"))
	}
	return Config{
		AuthConfig: userlist.AuthConfig{
			Users: map[auth.User]userlist.UserConfig{
				"ci":  {Password: "pass"},
				"dev": {Password: "pass"},
			},
		},

		AccessConfig: userlist.AccessConfig{
			Permissions: map[auth.User]userlist.Permissions{
				"ci":  userlist.All(),
				"dev": userlist.Read(),
			},
		},
	}, nil
}

type rawConfig struct {
	*userlist.AuthConfig
	*userlist.AccessConfig
}
