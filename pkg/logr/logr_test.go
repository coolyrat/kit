package logr

import (
	"testing"
)

func TestSubLogr(t *testing.T) {
	logger.Info("Root logger")

	svcA := logger.Named("svca")
	svcA.Info("svca logger")
	svcA.Debug("svca logger")

	svcASub := svcA.Named("sub")
	svcASub.Info("sub logger")
	svcASub.Debug("sub logger")

}
