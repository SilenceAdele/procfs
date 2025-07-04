// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package procfs

import (
	"math"
	"os"
	"testing"
)

func TestProcStat(t *testing.T) {
	p, err := getProcFixtures(t).Proc(26231)
	if err != nil {
		t.Fatal(err)
	}

	s, err := p.Stat()
	if err != nil {
		t.Fatal(err)
	}

	// pid stat int fields
	for _, test := range []struct {
		name string
		want int
		have int
	}{
		{name: "pid", want: 26231, have: s.PID},
		{name: "user time", want: 1677, have: int(s.UTime)},
		{name: "system time", want: 44, have: int(s.STime)},
		{name: "start time", want: 82375, have: int(s.Starttime)},
		{name: "virtual memory size", want: 56274944, have: int(s.VSize)},
		{name: "resident set size", want: 1981, have: s.RSS},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}

	// pid stat uint64 fields
	for _, test := range []struct {
		name string
		want uint64
		have uint64
	}{
		{name: "RSS Limit", want: 18446744073709551615, have: s.RSSLimit},
		{name: "Start Code", want: 4194304, have: s.StartCode},
		{name: "End Code", want: 6294284, have: s.EndCode},
		{name: "Start Stack", want: 140736914091744, have: s.StartStack},
		{name: "delayacct_blkio_ticks", want: 31, have: s.DelayAcctBlkIOTicks},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}

	// pid stat uint fields
	for _, test := range []struct {
		name string
		want uint
		have uint
	}{
		{name: "processor", want: 0, have: s.Processor},
		{name: "rt_priority", want: 0, have: s.RTPriority},
		{name: "policy", want: 0, have: s.Policy},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}
}

func TestProcStatLimits(t *testing.T) {
	p, err := getProcFixtures(t).Proc(26232)
	if err != nil {
		t.Fatal(err)
	}

	s, err := p.Stat()
	if err != nil {
		t.Errorf("want not error, have %s", err)
	}

	// max values of stat int fields
	for _, test := range []struct {
		name string
		want int
		have int
	}{
		{name: "waited for children user time", want: math.MinInt64, have: s.CUTime},
		{name: "waited for children system time", want: math.MaxInt64, have: s.CSTime},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}
}

func TestProcStatComm(t *testing.T) {
	s1, err := testProcStat(26231)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := "vim", s1.Comm; want != have {
		t.Errorf("want comm %s, have %s", want, have)
	}

	s2, err := testProcStat(584)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := "(a b ) ( c d) ", s2.Comm; want != have {
		t.Errorf("want comm %s, have %s", want, have)
	}
}

func TestProcStatVirtualMemory(t *testing.T) {
	s, err := testProcStat(26231)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := 56274944, int(s.VirtualMemory()); want != have {
		t.Errorf("want virtual memory %d, have %d", want, have)
	}
}

func TestProcStatResidentMemory(t *testing.T) {
	s, err := testProcStat(26231)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := 1981*os.Getpagesize(), s.ResidentMemory(); want != have {
		t.Errorf("want resident memory %d, have %d", want, have)
	}
}

func TestProcStatStartTime(t *testing.T) {
	s, err := testProcStat(26231)
	if err != nil {
		t.Fatal(err)
	}

	time, err := s.StartTime()
	if err != nil {
		t.Fatal(err)
	}
	if want, have := 1418184099.75, time; want != have {
		t.Errorf("want start time %f, have %f", want, have)
	}
}

func TestProcStatCPUTime(t *testing.T) {
	s, err := testProcStat(26231)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := 17.21, s.CPUTime(); want != have {
		t.Errorf("want cpu time %f, have %f", want, have)
	}
}

func testProcStat(pid int) (ProcStat, error) {
	fs, err := NewFS(procTestFixtures)
	if err != nil {
		return ProcStat{}, err
	}
	p, err := fs.Proc(pid)
	if err != nil {
		return ProcStat{}, err
	}

	return p.Stat()
}
