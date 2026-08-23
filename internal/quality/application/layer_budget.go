package application

import "harbor-sfu.local/harbor-sfu/internal/quality/domain"

func fitsLayerBudget(layer domain.LayerInput, bandwidth int64, loss uint64) bool {
	return layer.Bitrate <= bandwidth && loss <= 100
}
