package application

import "errors"

func finishManifestCommit(primary, cleanup error) error { return primary }
func FinishManifestCommit(primary, cleanup error) error { return finishManifestCommit(primary, cleanup) }

var _ = errors.Is
