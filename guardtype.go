package guard

import "strings"

type GuardType int

const (
	AdGuard GuardType = 1
	Hosts   GuardType = 2
	Regex   GuardType = 4
)

var (
	guardTypeMap = map[string]GuardType{
		"adguard": AdGuard,
		"hosts":   Hosts,
		"regex":   Regex,
	}
)

func ParseGuardType(s string) (GuardType, bool) {
	guard, result := guardTypeMap[strings.ToLower(s)]
	return guard, result
}

func (guard GuardType) ToString() string {
	switch guard {
	case 1:
		return "adguard"
	case 2:
		return "hosts"
	case 4:
		return "regex"
	}

	return ""
}
