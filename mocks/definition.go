//go:build mocks

package mocks

import (
	// ensures this package is in vendors folder and fixes a bug in go:generate that appears because of use of reflection in mocks generation
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=./patron/log.go	-package=mocks	-mock_names=Logger=Logger	github.com/beatlabs/patron/log	Logger
