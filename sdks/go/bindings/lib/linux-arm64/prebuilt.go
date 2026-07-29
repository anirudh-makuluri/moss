//go:build linux && arm64

package mosscore_platform

/*
#cgo CFLAGS: -I${SRCDIR}/../../include
#cgo LDFLAGS: -L${SRCDIR} -lmoss -lstdc++ -ldl -lm -lpthread
*/
import "C"
