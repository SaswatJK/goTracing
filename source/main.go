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

func add(v1, v2 vec3) vec3 {
	var temp vec3
	temp.r = v1.r + v2.r
	temp.g = v1.g + v2.g
	temp.b = v1.b + v2.b
	return temp
}

func negate(v vec3) vec3 {
	var temp vec3
	temp.r = -v.r
	temp.g = -v.g
	temp.b = -v.b
	return temp
}

func sub(v1, v2 vec3) vec3 {
	var nv2 vec3
	nv2 = v2
	nv2 = negate(nv2)
	temp := add(v1, nv2)
	return temp
}

func mulWithScalar(v vec3, s float32) vec3 {
	var temp vec3
	temp.r = v.r * s
	temp.g = v.g * s
	temp.b = v.b * s
	return temp
}

func dot(v1, v2 vec3) float32 {
	return (v1.r*v2.r + v1.b*v2.b + v1.g*v2.g)
}

func clampMax(value float32, clampValue float32) float32 {
	if value < clampValue {
		return clampValue
	}
	return value
}

type Camera struct {
	position vec3
	view     vec3
	right    vec3
	up       vec3
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
			position: vec3{0.0, 0.0, 0.0},
			view:     vec3{0.0, 0.0, -1.0},
			right:    vec3{1.0, 0.0, 0.0},
			up:       vec3{0.0, 1.0, 0.0},
		},
	}
}

func colorScene(currScene *Scene) { //Color it white for now.
	var pixelColor color.Color
	pixelColor = color.RGBA{125, 100, 230, 255}
	for i := 0; i < VIEWPORT_HEIGHT; i++ {
		for j := 0; j < VIEWPORT_WIDTH; j++ {
			var x, y float32 //Basically they are the normalized coordinates??
			//q(t) = o + td
			x = (2.0 * (float32(j) + 0.5) / float32(VIEWPORT_WIDTH)) - 1.0
			y = (2.0 * (float32(i) + 0.5) / float32(VIEWPORT_HEIGHT)) - 1.0
			//s(x,y) = a*f*x*r - f*y*u + v || u = up || v = view || a = aspect ratio = 1 || r = right || f = focal length = tan(phi/2) ||
			//d = normalized (s)
			var rayDir vec3
			rayDir = sub(mulWithScalar(currScene.camera.right, x), mulWithScalar(currScene.camera.up, y))
			rayDir = add(rayDir, currScene.camera.view)
			var posRayDir vec3
			posRayDir.r = clampMax(rayDir.r, 0.0)
			posRayDir.g = clampMax(rayDir.g, 0.0)
			posRayDir.b = clampMax(rayDir.b, 0.0)
			pixelColor = color.RGBA{uint8(posRayDir.r * 255), uint8(posRayDir.g * 255), uint8(posRayDir.b * 255), 255}
			currScene.sceneImage.Set(j, i, pixelColor)
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
