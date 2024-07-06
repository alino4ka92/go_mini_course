package main

type Rectangle struct {
	width, height int
}

func (r Rectangle) Area() int {
	return r.width * r.height
}
