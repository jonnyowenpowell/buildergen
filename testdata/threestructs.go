package threestructs

import "time"

type X struct {
	a int           `btest:""`
	B int           `btest:""`
	c int           `btest:"opt"`
	d int           `btest:"ignore"`
	e time.Duration `btest:""`
	f []struct {
		l, n, m []int `json:"l"`
	} `btest:""`
}

type y struct {
	a string `btest:""`
	b string
}

type Z struct {
	a string
}
