package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Registry struct {
	Requests atomic.Uint64
	Messages atomic.Uint64
	Packets  atomic.Uint64
	Labels   map[string]string
}

func (r *Registry) Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "harbor_requests_total %d\nharbor_signaling_messages_total %d\nharbor_rtp_packets_total %d\n", r.Requests.Load(), r.Messages.Load(), r.Packets.Load())
}
