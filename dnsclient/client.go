package dnsclient

import (
	"fmt"
	"github.com/miekg/dns"
	"math/rand"
)

func assembleMsg(typ uint16, qs string) *dns.Msg {
	msg := &dns.Msg{}
	msg.Id = uint16(rand.Int())
	msg.MsgHdr.Response = false
	msg.MsgHdr.Opcode = dns.OpcodeQuery
	msg.MsgHdr.RecursionDesired = true

	var qarr []dns.Question

	q := dns.Question{}
	q.Name = qs
	q.Qtype = typ
	q.Qclass = dns.ClassINET

	qarr = append(qarr, q)

	msg.Question = qarr

	return msg

}

func SendAndRcv(rhost string, qs string, typ uint16) *dns.Msg {
	//fmt.Println(rhost,qs,typ)
	m := assembleMsg(typ, qs)

	//fmt.Println(m.String())

	msg, err := dns.Exchange(m, rhost)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return msg

}
