//go:build cgo && windows && amd64

package mosscore

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: ${SRCDIR}/lib/windows-amd64/moss.lib
*/
import "C"
