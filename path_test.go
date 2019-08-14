// Copyright 2013 Julien Schmidt. All rights reserved.
// Based on the path package, Copyright 2009 The Go Authors.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package router

import (
	"runtime"
	"testing"

	"github.com/valyala/fasthttp"
)

var cleanTests = []struct {
	path, result string
}{
	// Already clean
	{"/", "/"},
	{"/abc", "/abc"},
	{"/a/b/c", "/a/b/c"},
	{"/abc/", "/abc/"},
	{"/a/b/c/", "/a/b/c/"},

	// missing root
	{"", "/"},
	{"abc", "/abc"},
	{"abc/def", "/abc/def"},
	{"a/b/c", "/a/b/c"},

	// Remove doubled slash
	{"//", "/"},
	{"/abc//", "/abc/"},
	{"/abc/def//", "/abc/def/"},
	{"/a/b/c//", "/a/b/c/"},
	{"/abc//def//ghi", "/abc/def/ghi"},
	{"//abc", "/abc"},
	{"///abc", "/abc"},
	{"//abc//", "/abc/"},

	// Remove . elements
	{".", "/"},
	{"./", "/"},
	{"/abc/./def", "/abc/def"},
	{"/./abc/def", "/abc/def"},
	{"/abc/.", "/abc/"},

	// Remove .. elements
	{"..", "/"},
	{"../", "/"},
	{"../../", "/"},
	{"../..", "/"},
	{"../../abc", "/abc"},
	{"/abc/def/ghi/../jkl", "/abc/def/jkl"},
	{"/abc/def/../ghi/../jkl", "/abc/jkl"},
	{"/abc/def/..", "/abc"},
	{"/abc/def/../..", "/"},
	{"/abc/def/../../..", "/"},
	{"/abc/def/../../..", "/"},
	{"/abc/def/../../../ghi/jkl/../../../mno", "/mno"},

	// Combinations
	{"abc/./../def", "/def"},
	{"abc//./../def", "/def"},
	{"abc/../../././../def", "/def"},
}

func TestPathClean(t *testing.T) {
	for _, test := range cleanTests {
		if s := CleanPath(test.path); s != test.result {
			t.Errorf("CleanPath(%s) = %s, want %s", test.path, s, test.result)
		}
		if s := CleanPath(test.result); s != test.result {
			t.Errorf("CleanPath(%s) = %s, want %s", test.result, s, test.result)
		}
	}
}

func TestPathCleanMallocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping malloc count in short mode")
	}
	if runtime.GOMAXPROCS(0) > 1 {
		t.Log("skipping AllocsPerRun checks; GOMAXPROCS>1")
		return
	}

	for _, test := range cleanTests {
		allocs := testing.AllocsPerRun(100, func() { CleanPath(test.result) })
		if allocs > 0 {
			t.Errorf("CleanPath(%q): %v allocs, want zero", test.result, allocs)
		}
	}
}

func TestGetOptionalPath(t *testing.T) {
	handler := func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	expectedPaths := []string{
		"/show/:name/",
		"/show/:name/:surname/",
		"/show/:name/:surname/at/",
		"/show/:name/:surname/at/:address/",
		"/show/:name/:surname/at/:address/:id/",
		"/show/:name/:surname/at/:address/:id/:phone/",
	}
	r := New()
	r.GET("/show/:name/:surname?/at/:address?/:id/:phone?", handler)

	for _, path := range expectedPaths {
		ctx := new(fasthttp.RequestCtx)

		h, _ := r.Lookup("GET", path, ctx)

		if h == nil {
			t.Errorf("Expected optional path '%s' is not registered", path)
		}
	}
}

func BenchmarkCleanPathWithBuffer(b *testing.B) {
	path := "/../bench/"
	cpb := acquireCleanPathBuffer()

	for i := 0; i < b.N; i++ {
		cleanPathWithBuffer(cpb, path)
		cpb.reset()
	}
}

func BenchmarkCleanPath(b *testing.B) {
	path := "/../bench/"

	for i := 0; i < b.N; i++ {
		CleanPath(path)
	}
}
