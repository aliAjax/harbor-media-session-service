package application

import (
	"harbor-sfu.local/harbor-sfu/internal/quality/domain"
)

func Decide(policy domain.Policy, layers []domain.LayerInput, loss uint64, bw int64) domain.LayerDecision {
	if policy.FixedRID != "" && validFixedLayer(policy.FixedRID, layers, bw, loss) {
		return domain.LayerDecision{RID: policy.FixedRID, Bitrate: policy.Clamp(bw), Reason: "fixed policy"}
	}
	if best, ok := chooseBestLayer(layers, bw, loss); ok {
		return domain.LayerDecision{RID: best.RID, Bitrate: policy.Clamp(best.Bitrate), Reason: "best active layer"}
	}
	return domain.LayerDecision{RID: "low", Bitrate: policy.Clamp(bw), Reason: "conservative fallback"}
}
