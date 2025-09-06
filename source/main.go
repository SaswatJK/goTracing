package main

import "fmt"
import "image"
import "image/color"
import "image/jpeg"
import "os"

const IMG_WIDTH = 500
const IMG_HEIGHT = 500

type vec3 struct {
	r float32
	g float32
	b float32
}

type Camera struct {
	position  vec3
	direction vec3
}

type Scene struct {
	sceneImage *image.RGBA
	camera     *Camera
}

func initializeScene() *Scene {
	return &Scene{
		sceneImage: image.NewRGBA(image.Rect(0, 0, IMG_WIDTH, IMG_HEIGHT)),
		camera: &Camera{
			position:  vec3{0.0, 0.0, 0.0},
			direction: vec3{0.0, 0.0, -1.0},
		},
	}
}

func colorScene(currScene *Scene) { //Color it white for now.
	var pixelColor color.Color
	pixelColor = color.RGBA{125, 100, 230, 255}
	for i := 0; i < 25; i++ {
		currScene.sceneImage.Set(i, i, pixelColor)
	}
}

func main() {

	currentScene := initializeScene()
	colorScene(currentScene)
	file, err := os.Create("output.jpg") //Since this will already be in the build directory, don't need to do relative path.
	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = jpeg.Encode(file, currentScene.sceneImage, &jpeg.Options{Quality: 100})
	fmt.Println("Hello world")
}
