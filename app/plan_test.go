// Copyright 2014 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"sort"

	"github.com/tsuru/config"

	"launchpad.net/gocheck"
)

func (s *S) TestPlanAdd(c *gocheck.C) {
	p := Plan{
		Name:     "plan1",
		Memory:   9223372036854775807,
		Swap:     1024,
		CpuShare: 100,
	}
	err := p.Save()
	c.Assert(err, gocheck.IsNil)
	defer s.conn.Plans().RemoveId(p.Name)
	var plan Plan
	err = s.conn.Plans().FindId(p.Name).One(&plan)
	c.Assert(err, gocheck.IsNil)
	c.Assert(plan, gocheck.DeepEquals, p)
}

func (s *S) TestPlanAddInvalid(c *gocheck.C) {
	invalidPlans := []Plan{
		{
			Memory:   9223372036854775807,
			Swap:     1024,
			CpuShare: 100,
		},
		{
			Name:   "plan1",
			Memory: 9223372036854775807,
			Swap:   1024,
		},
	}
	for _, p := range invalidPlans {
		err := p.Save()
		c.Assert(err, gocheck.FitsTypeOf, PlanValidationError{})
	}
}

func (s *S) TestPlanAddDupp(c *gocheck.C) {
	p := Plan{
		Name:     "plan1",
		Memory:   9223372036854775807,
		Swap:     1024,
		CpuShare: 100,
	}
	defer s.conn.Plans().RemoveId(p.Name)
	err := p.Save()
	c.Assert(err, gocheck.IsNil)
	err = p.Save()
	c.Assert(err, gocheck.Equals, ErrPlanAlreadyExists)
}

type planList []Plan

func (l planList) Len() int           { return len(l) }
func (l planList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l planList) Less(i, j int) bool { return l[i].Name < l[j].Name }

func (s *S) TestPlansList(c *gocheck.C) {
	expected := []Plan{
		s.defaultPlan,
		{Name: "plan1", Memory: 1, Swap: 2, CpuShare: 3},
		{Name: "plan2", Memory: 3, Swap: 4, CpuShare: 5},
	}
	err := s.conn.Plans().Insert(expected[1])
	c.Assert(err, gocheck.IsNil)
	err = s.conn.Plans().Insert(expected[2])
	c.Assert(err, gocheck.IsNil)
	defer s.conn.Plans().RemoveId(expected[1].Name)
	defer s.conn.Plans().RemoveId(expected[2].Name)
	plans, err := PlansList()
	c.Assert(err, gocheck.IsNil)
	sort.Sort(planList(plans))
	c.Assert(plans, gocheck.DeepEquals, expected)
}

func (s *S) TestPlanRemove(c *gocheck.C) {
	plans := []Plan{
		{Name: "plan1", Memory: 1, Swap: 2, CpuShare: 3},
		{Name: "plan2", Memory: 3, Swap: 4, CpuShare: 5},
	}
	err := s.conn.Plans().Insert(plans[0])
	c.Assert(err, gocheck.IsNil)
	err = s.conn.Plans().Insert(plans[1])
	c.Assert(err, gocheck.IsNil)
	defer s.conn.Plans().RemoveId(plans[0].Name)
	defer s.conn.Plans().RemoveId(plans[1].Name)
	err = PlanRemove(plans[0].Name)
	c.Assert(err, gocheck.IsNil)
	var dbPlans []Plan
	err = s.conn.Plans().Find(nil).All(&dbPlans)
	c.Assert(err, gocheck.IsNil)
	sort.Sort(planList(dbPlans))
	c.Assert(dbPlans, gocheck.DeepEquals, []Plan{
		s.defaultPlan,
		{Name: "plan2", Memory: 3, Swap: 4, CpuShare: 5},
	})
}

func (s *S) TestPlanRemoveInvalid(c *gocheck.C) {
	err := PlanRemove("xxxx")
	c.Assert(err, gocheck.Equals, ErrPlanNotFound)
}

func (s *S) TestDefaultPlan(c *gocheck.C) {
	p, err := defaultPlan()
	c.Assert(err, gocheck.IsNil)
	c.Assert(*p, gocheck.DeepEquals, s.defaultPlan)
}

func (s *S) TestDefaultPlanWithoutDefault(c *gocheck.C) {
	s.conn.Plans().RemoveAll(nil)
	defer s.conn.Plans().Insert(s.defaultPlan)
	config.Set("docker:memory", 12)
	config.Set("docker:swap", 32)
	defer config.Unset("docker:memory")
	defer config.Unset("docker:swap")
	p, err := defaultPlan()
	c.Assert(err, gocheck.IsNil)
	expected := Plan{
		Name:     "autogenerated",
		Memory:   12,
		Swap:     20,
		CpuShare: 100,
	}
	c.Assert(*p, gocheck.DeepEquals, expected)
}

func (s *S) TestFindPlanByName(c *gocheck.C) {
	p := Plan{
		Name:     "plan1",
		Memory:   9223372036854775807,
		Swap:     1024,
		CpuShare: 100,
	}
	err := p.Save()
	c.Assert(err, gocheck.IsNil)
	defer s.conn.Plans().RemoveId(p.Name)
	dbPlan, err := findPlanByName(p.Name)
	c.Assert(err, gocheck.IsNil)
	c.Assert(*dbPlan, gocheck.DeepEquals, p)
}
