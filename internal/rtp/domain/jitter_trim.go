package domain

func (j *JitterBuffer) trimNewest() {
	if j.max < 1 {
		j.packets = nil
		return
	}
	j.packets = j.packets[:j.max]
}
