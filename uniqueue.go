package graphpipe

type uniqueue struct {
	queue  []int
	queued []bool
}

func newUniqueue(size int) *uniqueue {
	return &uniqueue{queue: make([]int, 0, size), queued: make([]bool, size)}
}

func (u *uniqueue) Push(id int) bool {
	if u.queued[id] {
		return false
	}
	u.queue = append(u.queue, id)
	u.queued[id] = true
	return true
}

func (u *uniqueue) Pop() int {
	i := u.queue[0]
	u.queue = u.queue[1:]
	u.queued[i] = false
	return i
}

func (u *uniqueue) Len() int {
	return len(u.queue)
}
