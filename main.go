package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

func draw(color termbox.Attribute) {
	x, y := termbox.Size()
	r := []rune(" ")[0]

	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			termbox.SetCell(i, j, r, color, color)
		}
	}
	termbox.Flush()
}

type signal struct{}

func killEvent() <-chan signal {
	ch := make(chan signal)

	go func(chan<- signal) {
		for {
			if ev := termbox.PollEvent(); ev.Type == termbox.EventKey {
				ch <- signal{}
			}
		}
	}(ch)

	return ch
}

var colors = func() func() termbox.Attribute {
	red := termbox.ColorRed
	blue := termbox.ColorBlue

	current := red
	return func() termbox.Attribute {
		res := current
		if current == red {
			current = blue
		} else {
			current = red
		}
		return res
	}
}()

func loop() {
	ticker := time.Tick(1 * time.Second)
	kill := killEvent()

	draw(colors())
	for {
		select {
		case <-ticker:
			draw(colors())
		case <-kill:
			return
		}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	loop()
}
