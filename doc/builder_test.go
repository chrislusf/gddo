// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package doc

import (
	"go/ast"
	"testing"
)

var badSynopsis = []string{
	"+build !release",
	"COPYRIGHT Jimmy Bob",
	"### Markdown heading",
	"-*- indent-tabs-mode: nil -*-",
	"vim:set ts=2 sw=2 et ai ft=go:",
}

func TestBadSynopsis(t *testing.T) {
	for _, s := range badSynopsis {
		if synopsis(s) != "" {
			t.Errorf(`synopsis(%q) did not return ""`, s)
		}
	}
}

const readme = `
    $ go get github.com/user/repo/pkg1
    [foo](http://gopkgdoc.appspot.com/pkg/github.com/user/repo/pkg2)
    [foo](http://go.pkgdoc.org/github.com/user/repo/pkg3)
    [foo](http://godoc.org/github.com/user/repo/pkg4)
    <http://go.pkgdoc.org/github.com/user/repo/pkg5>
    [foo](http://godoc.org/github.com/user/repo/pkg6#Export)
    http://gowalker.org/github.com/user/repo/pkg7
    Build Status: [![Build Status](https://drone.io/github.com/user/repo1/status.png)](https://drone.io/github.com/user/repo1/latest)
    'go get example.org/package1' will install package1.
    (http://go.pkgdoc.org/example.org/package2 "Package2's documentation on GoPkgDoc").
    import "example.org/package3"
`

var expectedReferences = []string{
	"github.com/user/repo/pkg1",
	"github.com/user/repo/pkg2",
	"github.com/user/repo/pkg3",
	"github.com/user/repo/pkg4",
	"github.com/user/repo/pkg5",
	"github.com/user/repo/pkg6",
	"github.com/user/repo/pkg7",
	"github.com/user/repo1",
	"example.org/package1",
	"example.org/package2",
	"example.org/package3",
}

func TestReferences(t *testing.T) {
	references := make(map[string]bool)
	addReferences(references, []byte(readme))
	for _, r := range expectedReferences {
		if !references[r] {
			t.Errorf("missing %s", r)
		}
		delete(references, r)
	}
	for r := range references {
		t.Errorf("extra %s", r)
	}
}

var simpleImporterTests = []string{
	"code.google.com/p/biogo.foobar",
	"code.google.com/p/google-api-go-client/foobar/v3",
	"git.gitorious.org/go-pkg/foobar.git",
	"github.com/quux/go-foobar",
	"github.com/quux/go.foobar",
	"github.com/quux/foobar.go",
	"github.com/quux/foobar-go",
	"github.com/quux/foobar",
	"foobar",
	"quux/foobar",
}

func TestSimpleImporter(t *testing.T) {
	for _, path := range simpleImporterTests {
		m := make(map[string]*ast.Object)
		obj, _ := simpleImporter(m, path)
		if obj.Name != "foobar" {
			t.Errorf("simpleImporter(%q) = %q, want %q", path, obj.Name, "foobar")
		}
	}
}
