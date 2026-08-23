package quebec

import (
	q "harbor-sfu.local/harbor-sfu/internal/quality/application"
	"sync"
	"testing"
)

func overlap(w, r func()) {
	st := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(2)
	done.Add(2)
	go func() {
		defer done.Done()
		ready.Done()
		<-st
		for i := 0; i < 20000; i++ {
			w()
		}
	}()
	go func() {
		defer done.Done()
		ready.Done()
		<-st
		for i := 0; i < 20000; i++ {
			r()
		}
	}()
	ready.Wait()
	close(st)
	done.Wait()
}
func TestQuebecChangedStateRace(t *testing.T) {
	var s q.ChangedState
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 20000; i++ {
			s.Store(true)
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 20000; i++ {
			_ = s.Changed()
		}
	}()
	close(start)
	wg.Wait()
}
func TestQuebecCurrentStateRace(t *testing.T) {
	var s q.CurrentState
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 20000; i++ {
			s.Store("high")
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 20000; i++ {
			_ = s.Load()
		}
	}()
	close(start)
	wg.Wait()
}
func TestQuebecHysteresisStateRace(t *testing.T) {
	var s q.HysteresisState
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 20000; i++ {
			s.Store(2)
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 20000; i++ {
			_ = s.Level()
		}
	}()
	close(start)
	wg.Wait()
}
func TestQuebecMetricSnapshotRace(t *testing.T) {
	var s q.MetricState
	s.Set("loss", 0)
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 20000; i++ {
			s.Set("loss", 1)
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 20000; i++ {
			total := 0
			for _, value := range s.Snapshot() {
				total += value
			}
			_ = total
		}
	}()
	close(start)
	wg.Wait()
}
