package mis

import (
	"math"
)

//Mov is path ganrator
//mov get the reltive place on the line and return the cordinates
//z is in the range [0,255]
//when 255 is sign to end the movment
type Mov func(z int) (x, y int)

func MovL(x, y, fx, fy int) Mov {
	return func(z int) (int, int) {
		return mid(x, fx, z), mid(y, fy, z)
	}
}

func MovArc(x, y, fx, fy int) Mov {
	return func(z int) (int, int) {
		return arc1(x, fx, z), arc2(y, fy, z)
	}
}

func mid(x, fx, z int) int {
	return (fx*z + x*(255-z)) / 255
}

func arc1(x, fx, z int) int {
	y := float64(x)
	theta := float64(z) * 3.141592 / 2 / 255
	y += math.Sin(theta) * float64(fx-x)
	return int(y)
}

func arc2(x, fx, z int) int {
	y := float64(x)
	theta := float64(z) * 3.141592 / 2 / 255
	y += (1 - math.Cos(theta)) * float64(fx-x)
	return int(y)
}

//sMov is speed configorions
//given a frame i it return x place on Mov
//see Mov type for informion about the path
//x need to be in the range [0,255]
//when 255 is sign to end the movment
type Spd func(i int) (x int)

//sig gnral sigmoid speed funcion
func Sig(long float64) Spd {
	return func(i int) int {
		return int(_sig(float64(i) / long))
	}
}

func _sig(x float64) float64 {
	return x * x / (x*x - x + 0.5) * 128
}

func Acl(long float64) func(i int) (x int) {
	return func(i int) int {
		return int(_acl(float64(i) / long))
	}
}

func _acl(x float64) float64 {
	return x * x * 256
}

func DeAcl(long float64) func(i int) (x int) {
	return func(i int) int {
		return int(_dacl(float64(i) / long))
	}
}

func _dacl(x float64) float64 {
	return (1 - (1-x)*(1-x)) * 256
}

//Apply run the Mov on
//right use of apply will be somthing like
func Apply(s Spd, m Mov, x, y *int) func() bool {
	i := 0
	return func() bool {
		*x, *y = m(s(i))
		i++
		if s(i) < 255 {
			return true
		}
		return false
	}
}
