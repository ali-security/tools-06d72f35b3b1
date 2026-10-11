// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"net"
	"testing"
)

// TestListenImplicitHostIsLoopback checks that an address with no host, as
// produced by "gopls serve -port=N" or passed as "-listen=:N", is bound to
// the loopback interface instead of all network interfaces.
func TestListenImplicitHostIsLoopback(t *testing.T) {
	ln, err := listen(":0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	ip := ln.Addr().(*net.TCPAddr).IP
	if !ip.IsLoopback() {
		t.Errorf("listen(\":0\") bound %v, want a loopback address", ln.Addr())
	}
}

// TestListenExplicitHost checks that an explicit host is honored, including
// 0.0.0.0 for callers that deliberately want all network interfaces.
func TestListenExplicitHost(t *testing.T) {
	for _, tc := range []struct {
		addr string
		want func(net.IP) bool
	}{
		{"127.0.0.1:0", net.IP.IsLoopback},
		{"0.0.0.0:0", net.IP.IsUnspecified},
	} {
		ln, err := listen(tc.addr)
		if err != nil {
			t.Fatal(err)
		}
		ip := ln.Addr().(*net.TCPAddr).IP
		ln.Close()
		if !tc.want(ip) {
			t.Errorf("listen(%q) bound %v", tc.addr, ip)
		}
	}
}
