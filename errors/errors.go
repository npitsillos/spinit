package errors

import (
	"errors"
)

var (
	ErrNoDockerFile                = errors.New("no Dockerfile found in directory")
	ErrAccessingHomeDir            = errors.New("can't access home directory")
	ErrInvalidBuildCtxLocalMount   = errors.New("invalid buildCtx local mount dir")
	ErrInvalidDockerfileLocalMount = errors.New("invalid dockerfile local mount dir")
)
