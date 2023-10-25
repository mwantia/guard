package guard

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/miekg/dns"
)

type guard struct {
	Next   plugin.Handler
	Config *config
}

func (g guard) Name() string {
	return pluginName
}

func (g guard) ServeDNS(ctx context.Context, writer dns.ResponseWriter, response *dns.Msg) (int, error) {
	length := len(response.Question)

	for i := 0; i < length; i++ {
		question := response.Question[i]

		// Only resolve A or AAAA requests for now
		if question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA {

			fqdn := dns.Fqdn(question.Name)
			log.Debugf("Finding guard match for fqdn '%+v'", fqdn)

			for _, list := range g.Config.Lists {
				match, entry := list.IsMatch(fqdn)

				if match {
					answer := &dns.Msg{
						Answer: CreateGuardAnswers(question, entry.Address),
					}

					log.Debugf("Match found in entry '%+v'", entry.Content)
					metricsGuardRequestMatchCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

					answer.SetReply(response)
					_ = writer.WriteMsg(answer)

					return dns.RcodeSuccess, nil
				}
			}
		}
	}

	return plugin.NextOrFailure(g.Name(), g.Next, ctx, writer, response)
}

func CreateGuardAnswers(question dns.Question, address net.IP) []dns.RR {
	// Create a records header based on the initial question
	header := dns.RR_Header{
		Name:   question.Name,
		Class:  question.Qclass,
		Rrtype: question.Qtype,
	}

	if header.Rrtype == dns.TypeAAAA {
		return []dns.RR{
			&dns.AAAA{
				Hdr:  header,
				AAAA: net.IPv6zero,
			},
		}
	}

	return []dns.RR{
		&dns.A{
			Hdr: header,
			A:   address,
		},
	}
}
