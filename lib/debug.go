package lib

type counter struct {
	count int
	inc   int
	max   int
}

func (c *counter) Next() int {
	if c.max == 0 {
		return 0
	}
	//
	cur := c.count
	if (cur == 0 && c.inc < 0) || (cur == c.max-1 && c.inc > 0) {
		c.inc *= -1
	}
	c.count += c.inc
	return cur
}

var states = []string{
	"[+----]",
	"[-+---]",
	"[--+--]",
	"[---+-]",
	"[----+]",
}
var loadingState = counter{0, 1, len(states)}

func InteractiveLoading() string {
	return states[loadingState.Next()]
}
