package application

import "harbor-sfu.local/harbor-sfu/internal/quality/domain"

func validLayer(layer domain.LayerInput) bool { return layer.Active }
