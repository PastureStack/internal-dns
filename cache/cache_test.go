package cache

import (
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestInsertHitAndExpire(t *testing.T) {
	cache := New(2, 60)
	question := dns.Question{Name: "api.pasture.internal.", Qtype: dns.TypeA, Qclass: dns.ClassINET}
	message := new(dns.Msg)
	message.SetQuestion(question.Name, question.Qtype)
	message.Answer = []dns.RR{&dns.A{
		Hdr: dns.RR_Header{Name: question.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 30},
		A:   []byte{10, 0, 0, 8},
	}}

	cache.InsertMessage(Key(question, false, false), message, time.Minute)
	if got, _ := cache.Hit(question, false, false, 42); got == nil || got.Id != 42 || len(got.Answer) != 1 {
		t.Fatalf("unexpected cache hit: %#v", got)
	}

	cache.Remove(Key(question, false, false))
	cache.InsertMessage(Key(question, false, false), message, -time.Second)
	if got, _ := cache.Hit(question, false, false, 43); got != nil {
		t.Fatalf("expired hit = %#v", got)
	}
}

func TestCapacityEvictsOnlyExcessEntries(t *testing.T) {
	cache := New(2, 60)
	for _, name := range []string{"a.", "b.", "c."} {
		question := dns.Question{Name: name, Qtype: dns.TypeA, Qclass: dns.ClassINET}
		message := new(dns.Msg)
		message.SetQuestion(name, dns.TypeA)
		cache.InsertMessage(Key(question, false, false), message, time.Minute)
	}
	if got := len(cache.m); got != 2 {
		t.Fatalf("cache size = %d, want 2", got)
	}
}

func TestKeyRRsetIncludesSRVPort(t *testing.T) {
	header := dns.RR_Header{Name: "service.", Rrtype: dns.TypeSRV, Class: dns.ClassINET, Ttl: 60}
	first := &dns.SRV{Hdr: header, Priority: 1, Weight: 2, Port: 80, Target: "a."}
	second := &dns.SRV{Hdr: header, Priority: 1, Weight: 2, Port: 81, Target: "a."}
	if KeyRRset([]dns.RR{first}) == KeyRRset([]dns.RR{second}) {
		t.Fatal("SRV keys must differ by port")
	}
}
