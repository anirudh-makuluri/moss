//go:build darwin && arm64

package mosscore_platform

/*
#cgo CFLAGS: -I${SRCDIR}/../../include
#cgo LDFLAGS: -L${SRCDIR} -lmoss -lc++ -framework Security -framework SystemConfiguration
*/
import "C"
