package xray

import (
	participant "harbor-sfu.local/harbor-sfu/internal/participant/domain"
	rtcpapp "harbor-sfu.local/harbor-sfu/internal/rtcp/application"
	rtcpdomain "harbor-sfu.local/harbor-sfu/internal/rtcp/domain"
	rtpapp "harbor-sfu.local/harbor-sfu/internal/rtp/application"
	rtpdomain "harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"sync"
	"testing"
	"time"
)

func TestXrayRemapperConcurrentReset(t *testing.T) {
	r := rtpapp.NewRemapper()
	p := rtpdomain.Packet{SSRC: 41}
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(2)
	done.Add(2)
	go func() {
		defer done.Done()
		ready.Done()
		<-start
		for i := 0; i < 20000; i++ {
			r.Map(p)
		}
	}()
	go func() {
		defer done.Done()
		ready.Done()
		<-start
		for i := 0; i < 20000; i++ {
			r.Reset()
		}
	}()
	ready.Wait()
	close(start)
	done.Wait()
}

func TestXrayJitterConcurrentAccess(t *testing.T) {
	j := rtpdomain.NewJitterBuffer(32)
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(2)
	done.Add(2)
	go func() {
		defer done.Done()
		ready.Done()
		<-start
		for i := 0; i < 20000; i++ {
			j.Push(rtpdomain.Packet{Sequence: uint16(i)})
		}
	}()
	go func() {
		defer done.Done()
		ready.Done()
		<-start
		for i := 0; i < 20000; i++ {
			_ = j.Len()
		}
	}()
	ready.Wait()
	close(start)
	done.Wait()
}

func TestXrayRTCPConcurrentSnapshot(t *testing.T) {
	c := rtcpapp.NewController()
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(2)
	done.Add(2)
	go func() {
		defer done.Done()
		ready.Done()
		<-start
		for i := 0; i < 20000; i++ {
			c.Observe(rtcpdomain.Report{SSRC: uint32(i % 64)})
		}
	}()
	go func() {
		defer done.Done()
		ready.Done()
		<-start
		for i := 0; i < 20000; i++ {
			_ = c.Snapshot()
		}
	}()
	ready.Wait()
	close(start)
	done.Wait()
}

func TestXrayParticipantSessionConcurrentClose(t *testing.T) {
	s := &participant.Session{LastSeen: time.Now()}
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(2)
	done.Add(2)
	go func() {
		defer done.Done()
		ready.Done()
		<-start
		for i := 0; i < 20000; i++ {
			s.Close()
		}
	}()
	go func() {
		defer done.Done()
		ready.Done()
		<-start
		for i := 0; i < 20000; i++ {
			_ = s.Touch()
			_ = s.Expired(time.Now(), time.Minute)
		}
	}()
	ready.Wait()
	close(start)
	done.Wait()
}
