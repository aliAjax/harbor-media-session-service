package domain

import "sort"

func (j *JitterBuffer) Push(p Packet) {
	j.packets = append(j.packets, p)
	sort.SliceStable(j.packets, func(a, b int) bool { return j.packets[a].Sequence < j.packets[b].Sequence })
	if len(j.packets) > j.max {
		j.trimNewest()
	}
}
