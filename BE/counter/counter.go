package counter

type Counter struct {
	n int
}

func (c *Counter) Add()       { c.n = c.n + 1 }
func (c *Counter) Value() int { return c.n }
