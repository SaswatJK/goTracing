package main

import "fmt"
import "image"
import "image/color"
import "image/jpeg"
import "os"

const VIEWPORT_WIDTH int = 500
const VIEWPORT_HEIGHT int = 500
const IMAGE_WIDTH int = 2000
const IMAGE_HEIGHT int = 2000

type vec3 struct {
	r float32
	g float32
	b float32
}

type Camera struct {
	position  vec3
	direction vec3
	right     vec3
}

type Scene struct {
	sceneImage *image.RGBA
	camera     *Camera
}

/*
   We need to go through each 'pixel' of the image, and then for each pixel, we will be looping through different 'shapes' that are in the 'scene'. For that, I will need to add a pointer to an array of shapes. These shapes should have the same 'structure' (in code).
   So the quesiton is, how do we loop through the pixels? To loop through the pixels, I need to shoot a 'ray' to the pixel. The problem is, where is the ray? For that, I need to calculate the closest and the farthest planes of the scene, as well as the farthest left & right AND farthest up and down, which shoudld just be the VIEWPORT_WIDTH, and VIEWPORT_HEIGHT.
  Okay, so what actually matters is the 'size' of the 'scene'. Meaning, the VIEWPORT_WIDTH x VIEWPORT_HEIGHT scene, what does it actaully represent in the 'world'? What coordinates?
*/

func initializeScene() *Scene {
	return &Scene{
		sceneImage: image.NewRGBA(image.Rect(0, 0, VIEWPORT_WIDTH, VIEWPORT_HEIGHT)),
		camera: &Camera{
			position:  vec3{0.0, 0.0, 0.0},
			direction: vec3{0.0, 0.0, -1.0},
			right:     vec3{1.0, 0.0, 0.0},
		},
	}
}

func colorScene(currScene *Scene) { //Color it white for now.
	var pixelColor color.Color
	pixelColor = color.RGBA{125, 100, 230, 255}
	for i := 0; i < VIEWPORT_HEIGHT; i++ {
		for j := 0; j < VIEWPORT_WIDTH; j++ {
			var x, y float32 //Basically they are the normalized coordinates??
			//Firstly we convert these things to 0-1 normalization. 0.5 is for pointing to the fake 'center' of the pixel.
			x = (float32(j) + 0.5) / float32(VIEWPORT_WIDTH)
			y = (float32(i) + 0.5) / float32(VIEWPORT_HEIGHT)
			//Convert the 0-1 to NDC (-1 to +1).
			x = (2.0*x - 1.0)
			y = (2.0*y - 1.0)
			//Converting this to the real world coordinates now.
			x = (float32(IMAGE_WIDTH) * 0.5)
			y = (float32(IMAGE_HEIGHT) * 0.5)
			//Now we need to convert this to a ray. We will need an 'FOV', the camera direction, which we already have, and the camera position.
			currScene.sceneImage.Set(i, i, pixelColor)
		}
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
