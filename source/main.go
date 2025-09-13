package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

const VIEWPORT_WIDTH int = 500
const VIEWPORT_HEIGHT int = 500
const IMAGE_WIDTH int = 2000
const IMAGE_HEIGHT int = 2000

type vec3 struct {
	r float32
	g float32
	b float32
}

func vecMagnitude(v vec3) float32 {
	mag := math.Sqrt((float64(v.r) * float64(v.r)) + (float64(v.g) * float64(v.g)) + (float64(v.b) * float64(v.b)))
	return float32(mag)
}

func vecNormalize(v vec3) vec3 { //Why inline it? Cause I am doing it millions of times, cause it's a game?
	var temp vec3
	mag := math.Sqrt((float64(v.r) * float64(v.r)) + (float64(v.g) * float64(v.g)) + (float64(v.b) * float64(v.b)))
	temp.r = v.r / float32(mag)
	temp.g = v.g / float32(mag)
	temp.b = v.b / float32(mag)
	return temp
}

func vecAdd(v1, v2 vec3) vec3 {
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

func vecSub(v1, v2 vec3) vec3 { //v1 - v2
	var nv2 vec3
	nv2 = v2
	nv2 = negate(nv2)
	temp := vecAdd(v1, nv2)
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

type Sphere struct {
	center          vec3
	color           vec3
	radius          float32
	transparency    float32
	refractiveIndex float32
}

type Ray struct {
	origin       vec3
	direction    vec3
	color        vec3
	prevRI       float32
	transparency float32
}

type Scene struct {
	sceneImage *image.RGBA
	camera     *Camera
	objects    *[]Sphere
}

/*
   We need to go through each 'pixel' of the image, and then for each pixel, we will be looping through different 'shapes' that are in the 'scene'. For that, I will need to add a pointer to an array of shapes. These shapes should have the same 'structure' (in code).
   So the quesiton is, how do we loop through the pixels? To loop through the pixels, I need to shoot a 'ray' to the pixel. The problem is, where is the ray? For that, I need to calculate the closest and the farthest planes of the scene, as well as the farthest left & right AND farthest up and down, which shoudld just be the VIEWPORT_WIDTH, and VIEWPORT_HEIGHT.
  Okay, so what actually matters is the 'size' of the 'scene'. Meaning, the VIEWPORT_WIDTH x VIEWPORT_HEIGHT scene, what does it actaully represent in the 'world'? What coordinates?
*/

func initializeScene() *Scene {
	return &Scene{
		sceneImage: image.NewRGBA(image.Rect(0, 0, VIEWPORT_WIDTH, VIEWPORT_HEIGHT)),
		camera: &Camera{ //We will normalize all the directions.
			position: vec3{0.0, 0.0, 0.0},
			view:     vecNormalize(vec3{0.0, 0.0, -100.0}),
			right:    vecNormalize(vec3{1.0, 0.0, 0.0}),
			up:       vecNormalize(vec3{0.0, 1.0, 0.0}),
		},
		objects: &[]Sphere{
			{
				center:          vec3{0.0, 00.0, -55.0},
				color:           vec3{0.0, 0.0, 0.3},
				radius:          30.0,
				transparency:    1.0,
				refractiveIndex: 1.33,
			},
			{
				center:          vec3{0.0, 00.0, -100.0},
				color:           vec3{1.0, 0.0, 0.0},
				radius:          30.0,
				transparency:    0.0,
				refractiveIndex: 0.0,
			},
		},
	}
}

func extractColor(vColor vec3) color.RGBA {
	var temp color.RGBA
	temp = color.RGBA{uint8(vColor.r * 255), uint8(vColor.g * 255), uint8(vColor.b * 255), 255}
	return temp
}

func hitSphere(r *Ray, object *Sphere) {
	var OC vec3 = vecSub(object.center, r.origin)
	a := dot(r.direction, r.direction)
	b := -2.0 * dot(r.direction, OC)
	c := dot(OC, OC) - (object.radius * object.radius)
	discriminant := b*b - 4*a*c
	if discriminant > 0.0 { // > means 2 real solutions, (the ray goes through the sphere and goes out of it), = means 1 solution (tangent to the sphere).
		//Solving the equatin gives us 't' from which we can find the point of intersection and find the normal by subtracting from that point, the center.
		r.transparency = object.transparency
		t := (-b - float32(math.Sqrt(float64(discriminant)))) / (2 * a)
		//I realized that this refraction may not work because of the fact that the ray origin is still coming from the old origin.
		var intersectionPoint vec3 = mulWithScalar(r.origin, t)
		var normal vec3 = vecNormalize(vecSub(intersectionPoint, object.center))
		immediateColor := mulWithScalar(object.color, 1.0)
		//immediateColor := r.direction
		r.color = vecAdd(immediateColor, r.color)
		if r.transparency == 0.0 {
			return
		}
		incidentAngleRadians := math.Acos(float64(dot(r.direction, normal) / (vecMagnitude(r.direction) * vecMagnitude(normal))))
		//Snell's law: n1sintheta1 = n2sintheta2 || ni = refractive indices of i || thetai = angles with normal of i.
		refractiveAngleRadiansSin := (r.prevRI / object.refractiveIndex) * float32(math.Sin(incidentAngleRadians))
		refractiveAngleRadians := math.Asin(float64(refractiveAngleRadiansSin))
		//normal = negate(normal)
		tangent := vecNormalize(vecSub(r.direction, mulWithScalar(normal, dot(r.direction, normal)))) //Projecting the incident vector to the tangent plane by removing the 'normal' component. (v.n) -> length of v projected to n (since n is normalized) || (v.n)n = vector pinting in directon n with the length of the projection. Subtract that from v, and we get the remaining (the tangent).
		refractedNormalComponent := mulWithScalar(normal, float32(math.Cos(refractiveAngleRadians)))
		refractedTangentComponent := mulWithScalar(tangent, float32(math.Sin(refractiveAngleRadians)))
		r.prevRI = object.refractiveIndex
		r.origin = intersectionPoint
		r.direction = vecAdd(refractedNormalComponent, refractedTangentComponent)
		OC = vecSub(object.center, r.origin)
		a = dot(r.direction, r.direction)
		b = -2.0 * dot(r.direction, OC)
		c = dot(OC, OC) - (object.radius * object.radius)
		t = (-b + float32(math.Sqrt(float64(discriminant)))) / (2 * a)
		intersectionPoint = mulWithScalar(r.origin, t)
		normal = vecNormalize(vecSub(object.center, intersectionPoint))
		incidentAngleRadians = math.Acos(float64(dot(r.direction, normal) / (vecMagnitude(r.direction) * vecMagnitude(normal))))
		//Snell's law: n1sintheta1 = n2sintheta2 || ni = refractive indices of i || thetai = angles with normal of i.
		refractiveAngleRadiansSin = (r.prevRI / 1.0) * float32(math.Sin(incidentAngleRadians))
		refractiveAngleRadians = math.Asin(float64(refractiveAngleRadiansSin))
		//normal = negate(normal)
		tangent = vecNormalize(vecSub(r.direction, mulWithScalar(normal, dot(r.direction, normal)))) //Projecting the incident vector to the tangent plane by removing the 'normal' component. (v.n) -> length of v projected to n (since n is normalized) || (v.n)n = vector pinting in directon n with the length of the projection. Subtract that from v, and we get the remaining (the tangent).
		refractedNormalComponent = mulWithScalar(normal, float32(math.Cos(refractiveAngleRadians)))
		refractedTangentComponent = mulWithScalar(tangent, float32(math.Sin(refractiveAngleRadians)))
		r.origin = intersectionPoint
		r.direction = vecAdd(refractedNormalComponent, refractedTangentComponent)
		//r.direction = vecAdd(mulWithScalar(vecAdd(refractedNormalComponent, refractedTangentComponent), 0), r.direction)
		//randSource := rand.NewSource(1171)
		//newRand := rand.New(randSource)
		//newRand.float32 gives a 0-1 float32 value in the interval [0.0, 1.0) (Half Open)
		//var dirOffset vec3 = vec3{newRand.Float32() * 0.2, newRand.Float32() * 0.2, newRand.Float32() * 0.2}
	}
	if discriminant == 0.0 { //tangent
		immediateColor := r.direction
		r.transparency = object.transparency
		if r.transparency == 0 {
			r.color = vecAdd(immediateColor, r.color)
			return
		}
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
			rayDir = vecSub(mulWithScalar(currScene.camera.right, x), mulWithScalar(currScene.camera.up, y))
			rayDir = vecAdd(rayDir, currScene.camera.view)
			var primaryRay Ray
			primaryRay.direction = rayDir
			primaryRay.origin = currScene.camera.position
			primaryRay.color = vec3{0.0, 0.0, 0.0}
			primaryRay.prevRI = 1.0
			primaryRay.transparency = 1.0
			for i := range *currScene.objects {
				hitSphere(&primaryRay, &(*currScene.objects)[i])
			}
			pixelColor = extractColor(primaryRay.color)
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
