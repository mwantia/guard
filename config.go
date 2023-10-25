package guard

import (
	"strconv"
	"strings"

	"github.com/coredns/caddy/caddyfile"
)

type config struct {
	Lists []GuardList
}

func CreateConfig(caddy caddyfile.Dispenser) (*config, error) {
	config := config{}
	var err error

	for caddy.NextBlock() {
		val := caddy.Val()
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

	return &config, nil
}

func ParseConfigList(caddy caddyfile.Dispenser, listType ListType, config *config) (GuardList, error) {
	// Create empty entry for returns
	empty := GuardList{}

	if caddy.NextArg() {
		address := caddy.Val()

		if caddy.NextArg() {
			guardType := GetGuardType(caddy.Val())

			if caddy.NextArg() {
				frequency, err := strconv.Atoi(caddy.Val())
				if err != nil {
					return empty, err
				}

				if frequency > 0 {
					return GuardList{
						ListType:  listType,
						Address:   address,
						GuardType: guardType,
						Frequency: frequency,
					}, nil
				}
			}

			return GuardList{
				ListType:  listType,
				Address:   address,
				GuardType: guardType,
				Frequency: 0,
			}, nil
		}

		return GuardList{
			ListType:  listType,
			Address:   address,
			GuardType: AdGuard,
			Frequency: 0,
		}, nil
	}

	return empty, nil
}

func GetGuardType(guard string) GuardType {
	if len(guard) > 0 {
		switch guard {
		case "adguard":
			return AdGuard

		case "hosts":
			return Hosts

		case "regex":
			return Regex
		}
	}

	return 0
}
