package queue

type Queue struct {
	Len  int
	Data []string
}

func (q Queue) Empty() bool {
	return q.Len == 0
}

func (q *Queue) Push(val string) {
	q.Data = append(q.Data, val)
	q.Len++
}

func (q *Queue) Front() string{
	return q.Data[0]
}

func (q *Queue) Pop() (ele string) {
	ele = q.Data[0]
	q.Data = q.Data[1:]
	q.Len--
	return
}

