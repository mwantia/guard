package guard

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestGuard(tst *testing.T) {
	controller := caddy.NewTestController("dns", `
   	guard {
			url "https://raw.githubusercontent.com/ph00lt0/blocklists/master/blocklist.txt" adguard 5
			file "./examples/tests/blocklist_hosts" hosts 5
   	}
  `)
	controller.Next()

	config, err := CreateConfig(controller.Dispenser)
	if err != nil {
		tst.Error(err)
	}

	g := guard{
		Next:   test.ErrorHandler(),
		Config: config,
	}

	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	ctx := context.TODO()
	response := new(dns.Msg)

	response.SetQuestion("google.de.", dns.TypeA)

	recorder := dnstest.NewRecorder(&test.ResponseWriter{})
	g.ServeDNS(ctx, recorder, response)

	answers := recorder.Msg.Answer
	if len(answers) == 0 {
		tst.Errorf("No answers received for question 'google.de.'")
	}

	fmt.Println("Amount of answers recived:", len(answers))

	aAnswer := recorder.Msg.Answer[0].(*dns.A).A
	if !net.IPv4zero.Equal(aAnswer) {
		tst.Errorf("Answer doesn't match with the defined ip")
	}

	fmt.Println("The answer for question 'google.de.' received was:", aAnswer)
}
