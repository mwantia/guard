package guard

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type ListType string

const (
	Directory ListType = "directory"
	File      ListType = "file"
	Url       ListType = "url"
)

type GuardEntry struct {
	Content   string
	Modifiers string
	Address   net.IP
	Regex     bool
	Exact     bool
}

type GuardList struct {
	ListType  ListType
	Address   string
	GuardType GuardType
	Frequency time.Duration
	Entries   []GuardEntry
}

func (list *GuardList) Setup() bool {

	if list.Refresh() {
		if list.Frequency > 0 {

			log.Info("Enabling refresh for list '", list.Address, "' with frequency '", list.Frequency, "'")
			ticker := time.NewTicker(list.Frequency)

			go func() {
				for range ticker.C {

					log.Info("Refreshing list: '", list.Address, "'")
					list.Refresh()
				}
			}()
		}

		return true
	}

	return false
}

func (list *GuardList) Refresh() bool {
	// Handle empty lists
	if list == nil {
		return false
	}

	list.Entries = list.Entries[:0]

	switch list.ListType {
	case Directory:
		o, err := os.Open(list.Address)
		if err != nil {
			return false
		}

		files, err := o.ReadDir(0)
		if err != nil {
			return false
		}

		for _, f := range files {
			file, err := os.Open(list.Address + "/" + f.Name())
			if err != nil {
				continue
			}

			entries := LoadEntriesFromFile(file, list.GuardType)
			list.Entries = append(list.Entries, entries...)

			log.Info("Read entries from '", list.Address+"/"+f.Name(), "': '", len(entries), "'")
		}

	case File:
		file, err := os.Open(list.Address)
		if err != nil {
			return false
		}

		entries := LoadEntriesFromFile(file, list.GuardType)
		list.Entries = append(list.Entries, entries...)

		log.Info("Read entries from '", list.Address, "': '", len(entries), "'")

	case Url:
		response, err := http.Get(list.Address)
		if err != nil {
			return false
		}

		defer response.Body.Close()
		entries := LoadEntriesFromFile(response.Body, list.GuardType)
		list.Entries = append(list.Entries, entries...)

		log.Info("Read entries from '", list.Address, "': '", len(entries), "'")
	}

	return true
}

func LoadEntriesFromFile(reader io.Reader, guardType GuardType) []GuardEntry {
	var entries []GuardEntry

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {

		line := scanner.Text()
		if len(line) > 0 {

			switch guardType {
			case AdGuard:
				// Comments in adguard files start with !
				if strings.HasPrefix(line, "!") || strings.HasPrefix(line, "$") {
					continue
				}
				// Ignore allowlisted entries since we can't separate it
				if strings.HasPrefix(line, "@@") {
					continue
				}

				line = strings.TrimSpace(strings.Split(line, "#")[0])
				data := strings.Split(line, "^")
				content := CleanupAdguardContent(data[0])

				modifiers := ""
				if len(data) > 1 {
					modifiers = data[1]
				}

				if len(content) == 0 {
					content = modifiers
					modifiers = ""
				}

				regex := false
				if strings.Contains(content, "*") {
					regex = true
				}

				address := net.IPv4zero
				if strings.HasPrefix(data[0], "||") {

					entries = append(entries, GuardEntry{
						Address:   address,
						Modifiers: modifiers,
						Content:   content,
						Exact:     false,
						Regex:     regex,
					})
				} else if strings.HasPrefix(data[0], "|") {

					entries = append(entries, GuardEntry{
						Address:   address,
						Modifiers: modifiers,
						Content:   content,
						Exact:     true,
						Regex:     regex,
					})
				} else {

					entries = append(entries, GuardEntry{
						Address:   address,
						Modifiers: modifiers,
						Content:   content,
						Exact:     true,
						Regex:     regex,
					})
				}
			case Hosts:
				// Comments in hosts files start with !
				if !strings.HasPrefix(line, "#") {
					data := strings.Fields(line)

					var address net.IP
					var content string
					// Allow single lines without any ip defined
					if len(data) == 1 {
						address = net.IPv4zero
						content = data[0]
					} else {
						address = net.ParseIP(data[0])
						content = data[1]
					}

					entries = append(entries, GuardEntry{
						Address: address,
						Content: content,
						Exact:   true,
					})
				}
			case Regex:
				// Comments in hosts files start with !
				if !strings.HasPrefix(line, "#") {

					entries = append(entries, GuardEntry{
						Address: net.IPv4zero,
						Content: line,
						Exact:   false,
						Regex:   true,
					})
				}
			}
		}
	}

	return entries
}

func CleanupAdguardContent(content string) string {
	cleanup := strings.TrimLeft(content, "|")

	cleanup = strings.TrimPrefix(cleanup, "http:")
	cleanup = strings.TrimPrefix(cleanup, "https:")

	cleanup = strings.Trim(cleanup, "/")

	return cleanup
}

func (list *GuardList) IsMatch(fqdn string) (bool, GuardEntry) {
	// Create empty entry for returns
	empty := GuardEntry{}
	// Ignore empty requests
	if len(fqdn) == 0 {
		return false, empty
	}
	// Handle empty lists
	if list == nil || len(list.Entries) == 0 {
		return false, empty
	}

	for _, entry := range list.Entries {

		c := dns.Fqdn(entry.Content)
		if entry.Exact {

			if entry.Regex {
				// Not yet implemented
				return false, empty
			} else {

				if fqdn == c {

					return true, entry
				}
			}
		} else {

			if strings.HasSuffix(fqdn, c) {

				return true, entry
			}
		}
	}

	return false, empty
}
