package main

import (
	"bufio"
	"errors"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/tarm/serial"
	"log"
	"os"
	"strconv"
	"time"
)

// Serial configuration
const port0 = "/dev/ttyACM0"
const port1 = "/dev/ttyACM1"
const baud = 115200
const timeout = 1
const threads = 2

// Graphics configuration
const screenWidth = 1920
const screenHeight = 1080
const lineWidth = 16
const lineThickness = 12

// Plot configuration
const yAverage = 10
const yDistAverage = 4
const yMinDefault = yAverage - yDistAverage
const yMaxDefault = yAverage + yDistAverage
const yMinSmoothed = yAverage - float32(yDistAverage)/10
const yMaxSmoothed = yAverage + float32(yDistAverage)/10
const sampleCount = 1
const serialBufSize = 1024
const chanBufSize = 16

// Don't change these
const xHeight = screenHeight - 100
const yWidth = 50
const xLen = screenWidth / lineWidth

// Global variables
var port = port0
var ser *serial.Port
var yMin = yMinDefault
var yMax = yMaxDefault
var ys = make([]float32, xLen)
var values = make(chan float32, chanBufSize)

// Copy of Arduino map function
func remap(x float32, inMin float32, inMax float32, outMin float32, outMax float32) float32 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// Open serial port
func serialOpen() {
	conf := &serial.Config{Name: port, Baud: baud, ReadTimeout: timeout * time.Millisecond}
	var err error
	ser, err = serial.OpenPort(conf)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if port == port0 {
				port = port1
				serialOpen()
			} else {
				log.Fatalln("serial.OpenPort() failed:", err)
			}
		}
	}
}

// Serial reading goroutine
func serialRead(ser *serial.Port) {
	reader := bufio.NewReaderSize(ser, serialBufSize)
	data := make(chan string, threads*sampleCount*2)
	for range [threads]int{} {
		go parseInput(data, values)
	}
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println("serial.Read() failed:", err)
			continue
		}
		if len(line) >= 4 {
			if len(data) < cap(data) {
				data <- string(line)
			}
		}
	}
}

func parseInput(data chan string, values chan float32) {
	avg, count := float32(0.0), 0
	for {
		line := <-data
		v, err := strconv.ParseFloat(line, 32)
		val := float32(v)
		if err != nil {
			log.Println("Float conversion failed:", err)
			continue
		}
		if val >= yMinSmoothed && val <= yMaxSmoothed {
			val = yAverage
		}
		val = rl.Clamp(val, float32(yMin), float32(yMax))
		if count == sampleCount {
			avg /= sampleCount
			avg = remap(avg, float32(yMin), float32(yMax), 50, xHeight)
			if len(values) < cap(values) {
				values <- avg
			}
			count = 0
			avg = 0
		} else {
			count++
			avg += val
		}
	}
}

func main() {
	// Configure window
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetConfigFlags(rl.FlagInterlacedHint)
	rl.InitWindow(screenWidth, screenHeight, "Real-time Plotter")
	rl.ToggleFullscreen()
	defer rl.CloseWindow()

	// Initialize serial
	serialOpen()
	go serialRead(ser)

	// Main loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw X and Y axes
		rl.DrawLine(yWidth, xHeight, screenWidth-50, xHeight, rl.Black)
		rl.DrawLine(yWidth, 50, yWidth, screenHeight-50, rl.Black)

		// Update ys
		ys = append(ys[1:], <-values)

		// Draw curve
		var lastX float32 = yWidth
		for i := 0; i < xLen-1; i++ {
			start := rl.Vector2{X: lastX, Y: ys[i]}
			lastX += lineWidth
			end := rl.Vector2{X: lastX, Y: ys[i+1]}
			rl.DrawLineBezier(start, end, lineThickness, rl.Blue)
		}

		// Draw axis labels
		rl.DrawText("X", screenWidth-60, xHeight+20, 20, rl.Black)
		rl.DrawText("Y", 20, yWidth+10, 20, rl.Black)
		rl.DrawText("Y: m/s^2", yWidth+20, xHeight+20, 20, rl.Black)

		//Draw reference points
		rl.DrawText(strconv.Itoa(yAverage), 20, xHeight/2, 20, rl.Black)
		rl.DrawText(strconv.Itoa((yMin+3*yMax)/4), 20, xHeight/4, 20, rl.Black)
		rl.DrawText(strconv.Itoa((3*yMin+yMax)/4), 20, xHeight*3/4, 20, rl.Black)
		rl.DrawText(strconv.Itoa(yMin), 20, xHeight-10, 20, rl.Black)

		// Adjust Y range
		if rl.IsKeyPressed(rl.KeyUp) && yMin > 1 {
			yMin--
			yMax++
		} else if rl.IsKeyPressed(rl.KeyDown) && yMin < 9 {
			yMin++
			yMax--
		}

		rl.DrawFPS(screenWidth-80, 0)
		rl.EndDrawing()
	}
}
