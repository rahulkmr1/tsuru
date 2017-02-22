// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/router"
	check "gopkg.in/check.v1"
)

func (s *S) TestRoutersListNoContent(c *check.C) {
	err := config.Unset("routers")
	c.Assert(err, check.IsNil)
	defer config.Set("routers:fake:type", "fake")
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/plans/routers", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Authorization", "bearer "+s.token.GetValue())
	m := RunServer(true)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusNoContent)
}

func (s *S) TestRoutersList(c *check.C) {
	config.Set("routers:router1:type", "foo")
	config.Set("routers:router2:type", "bar")
	defer config.Unset("routers:router1:type")
	defer config.Unset("routers:router2:type")
	recorder := httptest.NewRecorder()
	expected := []router.PlanRouter{
		{Name: "fake", Type: "fake", Default: true},
		{Name: "fake-tls", Type: "fake-tls"},
		{Name: "router1", Type: "foo"},
		{Name: "router2", Type: "bar"},
	}
	request, err := http.NewRequest("GET", "/routers", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Authorization", "bearer "+s.token.GetValue())
	m := RunServer(true)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.Header().Get("Content-Type"), check.Equals, "application/json")
	var routers []router.PlanRouter
	err = json.Unmarshal(recorder.Body.Bytes(), &routers)
	c.Assert(err, check.IsNil)
	c.Assert(routers, check.DeepEquals, expected)
}

func (s *S) TestRoutersListNoAppCreatePermission(c *check.C) {
	config.Set("routers:router1:type", "foo")
	config.Set("routers:router2:type", "bar")
	defer config.Unset("routers:router1:type")
	defer config.Unset("routers:router2:type")
	token := userWithPermission(c)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/routers", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Authorization", "bearer "+token.GetValue())
	m := RunServer(true)
	m.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusForbidden)
}