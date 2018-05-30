// Copyright 2018 The oksvg Authors. All rights reserved.
//
// created: 2018 by S.R.Wiley
//_
// Implements SVG style matrix transformations.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/transform
package oksvg

import (
	"math"

	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

type Matrix2D struct {
	A, B, C, D, E, F float64
}

// matrix3 is a full 3x3 float64 matrix
// used for inverting
type matrix3 [9]float64

func otherPair(i int) (a, b int) {
	switch i {
	case 0:
		a, b = 1, 2
	case 1:
		a, b = 0, 2
	case 2:
		a, b = 0, 1
	}
	return
}

func (m *matrix3) coFact(i, j int) float64 {
	ai, bi := otherPair(i)
	aj, bj := otherPair(j)
	a, b, c, d := m[ai+aj*3], m[bi+bj*3], m[ai+bj*3], m[bi+aj*3]
	return a*b - c*d
}

func (m *matrix3) Invert() *matrix3 {
	var cofact matrix3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sign := float64(1 - (i+j%2)%2*2) // "checkerboard of minuses" grid
			cofact[i+j*3] = m.coFact(i, j) * sign
		}
	}
	deteriminate := m[0]*cofact[0] + m[1]*cofact[1] + m[2]*cofact[2]
	// transpose cofact
	for i := 0; i < 2; i++ {
		for j := i + 1; j < 3; j++ {
			cofact[i+j*3], cofact[j+i*3] = cofact[j+i*3], cofact[i+j*3]
		}
	}
	for i := 0; i < 9; i++ {
		cofact[i] /= deteriminate
	}
	return &cofact
}

func (a Matrix2D) Invert() Matrix2D {
	n := &matrix3{a.A, a.C, a.E, a.B, a.D, a.F, 0, 0, 1}
	n = n.Invert()
	return Matrix2D{A: n[0], C: n[1], E: n[2], B: n[3], D: n[4], F: n[5]}
}

func (a Matrix2D) Mult(b Matrix2D) Matrix2D {
	return Matrix2D{
		A: a.A*b.A + a.C*b.B,
		B: a.B*b.A + a.D*b.B,
		C: a.A*b.C + a.C*b.D,
		D: a.B*b.C + a.D*b.D,
		E: a.A*b.E + a.C*b.F + a.E,
		F: a.B*b.E + a.D*b.F + a.F}
}

var Identity = Matrix2D{1, 0, 0, 1, 0, 0}

// TFixed transforms a fixed.Point26_6 by the matrix
func (m Matrix2D) TFixed(a fixed.Point26_6) (b fixed.Point26_6) {
	b.X = fixed.Int26_6((float64(a.X)*m.A + float64(a.Y)*m.C) + m.E*64)
	b.Y = fixed.Int26_6((float64(a.X)*m.B + float64(a.Y)*m.D) + m.F*64)
	return
}

func (m Matrix2D) Transform(x1, y1 float64) (x2, y2 float64) {
	x2 = x1*m.A + y1*m.C + m.E
	y2 = x1*m.B + y1*m.D + m.F
	return
}

func (a Matrix2D) Scale(x, y float64) Matrix2D {
	return a.Mult(Matrix2D{
		A: x,
		B: 0,
		C: 0,
		D: y,
		E: 0,
		F: 0})
}

func (a Matrix2D) SkewY(theta float64) Matrix2D {
	return a.Mult(Matrix2D{
		A: 1,
		B: math.Tan(theta),
		C: 0,
		D: 1,
		E: 0,
		F: 0})
}

func (a Matrix2D) SkewX(theta float64) Matrix2D {
	return a.Mult(Matrix2D{
		A: 1,
		B: 0,
		C: math.Tan(theta),
		D: 1,
		E: 0,
		F: 0})
}

func (a Matrix2D) Translate(x, y float64) Matrix2D {
	return a.Mult(Matrix2D{
		A: 1,
		B: 0,
		C: 0,
		D: 1,
		E: x,
		F: y})
}

func (a Matrix2D) Rotate(theta float64) Matrix2D {
	return a.Mult(Matrix2D{
		A: math.Cos(theta),
		B: math.Sin(theta),
		C: -math.Sin(theta),
		D: math.Cos(theta),
		E: 0,
		F: 0})
}

type MatrixAdder struct {
	rasterx.Adder
	M Matrix2D
}

func (t *MatrixAdder) Reset() {
	t.M = Identity
}

func (t *MatrixAdder) Start(a fixed.Point26_6) {
	t.Adder.Start(t.M.TFixed(a))
}

// Line adds a linear segment to the current curve.
func (t *MatrixAdder) Line(b fixed.Point26_6) {
	t.Adder.Line(t.M.TFixed(b))
}

// QuadBezier adds a quadratic segment to the current curve.
func (t *MatrixAdder) QuadBezier(b, c fixed.Point26_6) {
	t.Adder.QuadBezier(t.M.TFixed(b), t.M.TFixed(c))
}

// CubeBezier adds a cubic segment to the current curve.
func (t *MatrixAdder) CubeBezier(b, c, d fixed.Point26_6) {
	t.Adder.CubeBezier(t.M.TFixed(b), t.M.TFixed(c), t.M.TFixed(d))
}
