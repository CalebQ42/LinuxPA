package ui

import (
	"image"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// A rectangle with borders.
// Default only has black borders.
type Rect struct {
	x      int32
	y      int32
	h      int32
	w      int32
	radius int32

	borderWidth int32
	borderColor color.RGBA
	fillColor   color.RGBA

	tex *rl.Texture2D
}

func NewRect(x, y, height, width int32) *Rect {
	return &Rect{
		x:           x,
		y:           y,
		h:           height,
		w:           width,
		radius:      0,
		borderWidth: 2,
		borderColor: rl.Black,
		fillColor:   rl.Blank,
	}
}

func (r *Rect) SetPosition(x, y int32) {
	r.x, r.y = x, y
}

func (r *Rect) SetSize(h, w int32) {
	r.h, r.w = h, w
	r.Unload()
}

func (r *Rect) SetBorderRadius(radius int32) {
	r.radius = radius
	r.Unload()
}

func (r *Rect) SetBorderWidth(w int32) {
	r.borderWidth = w
	r.Unload()
}

func (r *Rect) SetBorderColor(c color.RGBA) {
	r.borderColor = c
	r.Unload()
}

func (r *Rect) SetFillColor(c color.RGBA) {
	r.fillColor = c
	r.Unload()
}

func (r *Rect) Unload() {
	if r.tex != nil {
		rl.UnloadTexture(*r.tex)
		r.tex = nil
	}
}

func (r *Rect) Draw() {
	if r.tex != nil {
		rl.DrawTextureEx(*r.tex, rl.NewVector2(float32(r.x), float32(r.y)), 0, 1, rl.White)
		return
	}
	// Create image at 16x resolution then resize down so it looks better
	img := rl.NewImageFromImage(image.NewAlpha(image.Rect(0, 0, int(r.w*4), int(r.h*4))))

	// Filled areas
	if r.fillColor != rl.Blank {
		rl.ImageDrawRectangle(img, 4*r.radius, 0, 4*(r.w-(2*r.radius)), 4*r.h, r.fillColor)
		if r.radius > 0 {
			rl.ImageDrawRectangle(img, 4*(r.w-r.radius), 4*r.radius, 4*r.radius, 4*(r.h-(2*r.radius)), r.fillColor)
			rl.ImageDrawRectangle(img, 0, 4*r.radius, 4*r.radius, 4*(r.h-(2*r.radius)), r.fillColor)
		}
	}

	if r.borderColor != rl.Blank {
		// Horizontal lines
		rl.ImageDrawRectangle(img, 4*r.radius, 0, 4*(r.w-(2*r.radius)), r.borderWidth*4, r.borderColor)
		rl.ImageDrawRectangle(img, 4*r.radius, (4*r.h)-(r.borderWidth*4), 4*(r.w-(2*r.radius)), r.borderWidth*4, r.borderColor)
		// Vertical lines
		rl.ImageDrawRectangle(img, 0, 4*r.radius, r.borderWidth*4, 4*(r.h-(2*r.radius)), r.borderColor)
		rl.ImageDrawRectangle(img, (4*r.w)-(4*r.borderWidth), 4*r.radius, r.borderWidth*4, 4*(r.h-(2*r.radius)), r.borderColor)
	}
	if r.radius > 0 {
		r.placeCorners(r.radius, img)
	}

	rl.ImageResize(img, r.w, r.h)
	tmp := rl.LoadTextureFromImage(img)
	r.tex = &tmp
	rl.DrawTextureEx(*r.tex, rl.NewVector2(float32(r.x), float32(r.y)), 0, 1, r.borderColor)
}

func (r Rect) placeCorners(radius float32, img *rl.Image) {
	// Setup corner
	corner := rl.NewImageFromImage(image.NewAlpha(image.Rect(0, 0, int(radius*4), int(radius*4))))
	if r.borderWidth > 0 {
		rl.ImageDrawCircleV(corner, 4*radius, 4*radius, 4*radius, r.borderColor)
	}
	rl.ImageDrawCircle(corner, 4*radius, 4*radius, (4*radius)-(4*r.borderWidth), r.fillColor)
	// Add corners to actual image
	rl.ImageDraw(img, corner, rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.White)
	rl.ImageRotateCW(corner)
	rl.ImageDraw(img, corner, rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.NewRectangle(float32(4*(r.w-radius)), 0, float32(corner.Width), float32(corner.Height)), rl.White)
	rl.ImageRotateCW(corner)
	rl.ImageDraw(img, corner, rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.NewRectangle(float32(4*(r.w-radius)), float32(4*(r.h-radius)), float32(corner.Width), float32(corner.Height)), rl.White)
	rl.ImageRotateCW(corner)
	rl.ImageDraw(img, corner, rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.NewRectangle(0, float32(4*(r.h-radius)), float32(corner.Width), float32(corner.Height)), rl.White)
}
