// Copyright 2020-present Sergio Andres Virviescas Santana, fasthttp
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// Package radix is a high performance HTTP routes storage.
package radix

import (
	"regexp"

	"github.com/valyala/fasthttp"
)

type nodeType uint8

type nodeHandler struct {
	tsr      bool
	handler  fasthttp.RequestHandler
	wildcard *nodeWildcard
}

type nodeWildcard struct {
	path     string
	paramKey string
	handler  fasthttp.RequestHandler
}

type node struct {
	nType nodeType

	path     string
	handlers map[string]*nodeHandler
	children []*node

	paramKeys  []string
	paramRegex *regexp.Regexp
}

type wildPath struct {
	path  string
	keys  []string
	start int
	end   int
	pType nodeType

	pattern string
	regex   *regexp.Regexp
}

// Tree is a routes storage
type Tree struct {
	root *node
}
