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

func (s IntervalSlice) Normalize() IntervalSlice {
	var n IntervalSlice
	if len(s) == 0 {
		return n
	}

	s.Sort()

	n = append(n, s[0])

	for i := 1; i < len(s); i++ {
		last := n[len(n)-1]
		if s[i].Start > last.Stop {
			n = append(n, s[i])
			continue
		}

		stop := s[i].Stop
		if last.Stop > stop {
			stop = last.Stop
		}
		n[len(n)-1] = Interval{last.Start, stop}
	}

	return n
}

func (s IntervalSlice) Contain(sub IntervalSlice) bool {
	j := 0
	for i := 0; i < len(s); i++ {
		for ; j < len(s); j++ {
			if sub[j].Start > s[i].Stop {
				continue
			}
			break
		}
		if j == len(s) {
			return false
		}
		if sub[j].Start < s[i].Start || sub[j].Stop > s[i].Stop {
			return false
		}
	}
	return true
}

func (s IntervalSlice) Equal(o IntervalSlice) bool {
	if len(s) != len(o) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i].Start != o[i].Start || s[i].Stop != o[i].Stop {
			return false
		}
	}
	return true
}

func (s IntervalSlice) Compare(o IntervalSlice) int {
	if s.Equal(o) {
		return 0
	}
	if s.Contain(o) {
		return 1
	}
	return -1
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (s *IntervalSlice) InsertInterval(interval Interval) {
	var (
		count int
		i     int
	)

	*s = append(*s, interval)
	total := len(*s)
	for i = total - 1; i > 0; i-- {
		if (*s)[i].Stop < (*s)[i-1].Start {
			(*s)[i], (*s)[i-1] = (*s)[i-1], (*s)[i]
		} else if (*s)[i].Start > (*s)[i-1].Stop {
			break
		} else {
			(*s)[i-1].Start = min((*s)[i-1].Start, (*s)[i].Start)
			(*s)[i-1].Stop = max((*s)[i-1].Stop, (*s)[i].Stop)
			count++
		}
	}
	if count > 0 {
		i++
		if i+count < total {
			copy((*s)[i:], (*s)[i+count:])
		}
		*s = (*s)[:total-count]
	}
}
