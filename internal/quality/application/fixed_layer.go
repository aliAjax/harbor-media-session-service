package application

import "harbor-sfu.local/harbor-sfu/internal/quality/domain"

func validFixedLayer(_ string, _ []domain.LayerInput, _ int64, _ uint64) bool {
	return true
}
