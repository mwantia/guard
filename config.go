package guard

import (
	"strings"
	"time"

	"github.com/coredns/caddy/caddyfile"
)

type config struct {
	Lists    []GuardList
	Defaults map[string]string
}

func CreateConfig(caddy caddyfile.Dispenser) (*config, error) {
	config := config{}
	config.Defaults = map[string]string{
		"default_refresh_frequency": "0s",
		"default_ipv4_answer":       "0.0.0.0",
		"default_ipv6_answer":       "::",
		"next_or_failure":           "true",
	}

	var err error

	for caddy.NextBlock() {
		val := caddy.Val()
		_, found := config.Defaults[val]

		if found {

			caddy.NextArg()
			new := caddy.Val()
			config.Defaults[val] = new
		} else {

			var list GuardList

			if strings.EqualFold(val, "directory") {
				list, err = ParseConfigList(caddy, Directory, &config)

			} else if strings.EqualFold(val, "file") {
				list, err = ParseConfigList(caddy, File, &config)

			} else if strings.EqualFold(val, "url") {
				list, err = ParseConfigList(caddy, Url, &config)

			}

			if err != nil {
				return nil, err
			}

			if len(list.ListType) > 0 && len(list.Address) > 0 {

				list.Setup()
				config.Lists = append(config.Lists, list)
			}
		}
	}

	return &config, nil
}

func ParseConfigList(caddy caddyfile.Dispenser, listType ListType, config *config) (GuardList, error) {
	// Create empty entry for returns
	empty := GuardList{}

	defaultFrequency, err := time.ParseDuration(config.Defaults["default_refresh_frequency"])
	if err != nil {
		return empty, err
	}

	if caddy.NextArg() {
		address := caddy.Val()

		if caddy.NextArg() {
			guardType, result := ParseGuardType(caddy.Val())
			// Default type 'Adguard'
			if !result {
				guardType = AdGuard
			}

			if caddy.NextArg() {
				customFrequency, err := time.ParseDuration(caddy.Val())
				if err != nil {
					return empty, err
				}

				if customFrequency >= 0 {
					return GuardList{
						ListType:  listType,
						Address:   address,
						GuardType: guardType,
						Frequency: customFrequency,
					}, nil
				}
			}

			return GuardList{
				ListType:  listType,
				Address:   address,
				GuardType: guardType,
				Frequency: defaultFrequency,
			}, nil
		}

		return GuardList{
			ListType:  listType,
			Address:   address,
			GuardType: AdGuard,
			Frequency: defaultFrequency,
		}, nil
	}

	return empty, nil
}
