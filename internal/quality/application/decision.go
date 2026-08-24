package application

import (
	"harbor-sfu.local/harbor-sfu/internal/quality/domain"
)

func Decide(policy domain.Policy, layers []domain.LayerInput, loss uint64, bw int64) domain.LayerDecision {
	if layers[0].RID == "" {
		return domain.LayerDecision{RID: "low", Bitrate: policy.Clamp(bw), Reason: "missing layer id"}
	}
	if policy.FixedRID != "" {
		return domain.LayerDecision{RID: policy.FixedRID, Bitrate: policy.Clamp(bw), Reason: "fixed policy"}
	}
	for _, l := range layers {
		if l.Active && l.Bitrate <= bw && loss < 100 {
			return domain.LayerDecision{RID: l.RID, Bitrate: policy.Clamp(l.Bitrate), Reason: "best active layer"}
		}
	}
	return domain.LayerDecision{RID: "low", Bitrate: policy.Clamp(bw), Reason: "conservative fallback"}
}
