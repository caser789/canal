package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Like MySQL GTID Interval struct, [start, stop), left closed and right open
// See MySQL rpl_gtid.h
type Interval struct {
	// first GTID of this interval
	Start int64
	// first GTID of next interval
	Stop int64
}

func (i Interval) String() string {
	if i.Stop == i.Start+1 {
		return fmt.Sprintf("%d", i.Start)
	}

	return fmt.Sprintf("%d-%d", i.Start, i.Stop-1)
}

// Interval is [start, stop), but the GTID string's format is [n] or [n1-n2], closed interval
func parseInterval(s string) (i Interval, err error) {
	nums := strings.Split(s, "-")
	switch len(nums) {
	case 1:
		i.Start, err = strconv.ParseInt(nums[0], 10, 64)
		i.Stop = i.Start + 1
	case 2:
		i.Start, err = strconv.ParseInt(nums[0], 10, 64)
		if err == nil {
			i.Stop, err = strconv.ParseInt(nums[1], 10, 64)
			i.Stop++
		}
	default:
		err = fmt.Errorf("invalid interval format, must be n[-n]")
	}

	if err != nil {
		return
	}

	if i.Stop <= i.Start {
		err = fmt.Errorf("invalid interval format, must be n[-n] and end must >= start")
	}

	return
}

type IntervalSlice []Interval

func (s IntervalSlice) Len() int {
	return len(s)
}

func (s IntervalSlice) Less(i, j int) bool {
	if s[i].Start < s[j].Start {
		return true
	}

	if s[i].Start > s[j].Start {
		return false
	}

	return s[i].Stop < s[j].Stop
}

func (s IntervalSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s IntervalSlice) Sort() {
	sort.Sort(s)
}

func (s IntervalSlice) Normalize() (o IntervalSlice) {
	if len(s) == 0 {
		return
	}

	s.Sort()

	o = append(o, s[0])
	for i := 1; i < len(s); i++ {
		last := o[len(o)-1]
		if s[i].Start > last.Stop {
			o = append(o, s[i])
			continue
		}
		stop := s[i].Stop
		if last.Stop > stop {
			stop = last.Stop
		}
		o[len(o)-1] = Interval{Start: last.Start, Stop: stop}
	}

	return o
}
