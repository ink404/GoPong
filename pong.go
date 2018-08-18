package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

type color struct {
	r, g, b byte
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}

}

type pos struct {
	x, y int
}

type ball struct {
	pos
	radius int
	// x,y velocity
	xv, yv int
	color  color
}

type paddle struct {
	pos
	// height, width
	w, h int
	color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, color{255, 255, 255}, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle) {
	ball.x += ball.xv
	ball.y += ball.yv

	// checking if the ball is within the window
	if ball.y-ball.radius < 0 || ball.y+ball.radius > winHeight {
		ball.yv *= -1
	}
	if ball.x-ball.radius < 0 || ball.x+ball.radius > winWidth {
		ball.x, ball.y = 300, 300
	}

	// checking if the ball and paddle intersect
	if ball.x < leftPaddle.x+leftPaddle.w/2 {
		if ball.y > leftPaddle.y-leftPaddle.h/2 && ball.y < leftPaddle.y+leftPaddle.h/2 {
			ball.xv *= -1
		}
	}
	if ball.x > rightPaddle.x-rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 && ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.xv *= -1
		}
	}

}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2
	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, color{255, 255, 255}, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= 10
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += 10
	}
}

func (paddle *paddle) aiUpdate(ball *ball) {
	paddle.y = ball.y
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func main() {

	// Added after EP06 to address macosx issues
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("PONG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{50, 100}, 20, 100, color{255, 255, 255}}
	player2 := paddle{pos{winWidth - 50, 100}, 20, 100, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 10, 10, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	// startT := time.Now()
	// Changd after EP 06 to address MacOSX
	// OSX requires that you consume events for windows to open and work properly
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)
		
		// draws
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		// updates
		ball.update(&player1, &player2)
		player1.update(keyState)
		player2.aiUpdate(&ball)
		// render updates
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		// elapsedT := float32(time.Since(startT).Seconds()) / 1000
		sdl.Delay(16)
	}

}
