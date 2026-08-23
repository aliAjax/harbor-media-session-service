package application

import "harbor-sfu.local/harbor-sfu/internal/quality/domain"

func chooseBestLayer(layers []domain.LayerInput, bandwidth int64, loss uint64) (domain.LayerInput, bool) {
	var best *domain.LayerInput
	for i := range layers {
		if validLayer(layers[i]) && fitsLayerBudget(layers[i], bandwidth, loss) && betterLayer(best, layers[i]) {
			best = &layers[i]
		}
	}
	if best == nil {
		return domain.LayerInput{}, false
	}
	return *best, true
}
