package application

import "harbor-sfu.local/harbor-sfu/internal/quality/domain"

func betterLayer(current *domain.LayerInput, candidate domain.LayerInput) bool {
	return current == nil
}
