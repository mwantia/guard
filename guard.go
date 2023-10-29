package guard

import (
	"context"
	"net"
	"strconv"

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
					address := entry.Address
					if address == nil {
						if question.Qtype == dns.TypeAAAA {

							d := g.Config.Defaults["default_ipv4_answer"]
							address = net.ParseIP(d)
						} else if question.Qtype == dns.TypeA {

							d := g.Config.Defaults["default_ipv6_answer"]
							address = net.ParseIP(d)
						}
					}

					answer := &dns.Msg{
						Answer: CreateGuardAnswers(question, address),
					}

					log.Debugf("Match found in entry '%+v'", entry.Content)
					metricsGuardRequestMatchCount.WithLabelValues(metrics.WithServer(ctx), list.Address, list.GuardType.ToString()).Inc()

					answer.SetReply(response)
					_ = writer.WriteMsg(answer)

					return dns.RcodeSuccess, nil
				}
			}
		}
	}

	next, err := strconv.ParseBool(g.Config.Defaults["next_or_failure"])
	if err != nil || next {
		return plugin.NextOrFailure(g.Name(), g.Next, ctx, writer, response)
	}

	return dns.RcodeNameError, nil
}

func CreateGuardAnswers(question dns.Question, address net.IP) []dns.RR {
	// Create a records header based on the initial question
	header := dns.RR_Header{
		Name:   question.Name,
		Class:  question.Qclass,
		Rrtype: question.Qtype,
		Ttl:    14400, // 4 hours
	}

	if header.Rrtype == dns.TypeAAAA {
		return []dns.RR{
			&dns.AAAA{
				Hdr:  header,
				AAAA: net.IPv6zero,
			},
		}
	} else if header.Rrtype == dns.TypeHTTPS {
		return []dns.RR{
			&dns.HTTPS{
				SVCB: dns.SVCB{
					Hdr:    header,
					Target: "0.0.0.0",
				},
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
