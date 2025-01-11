package ui

import (
	"fmt"
	"image"
	"image/color"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// A rectangle with borders.
// Default only has black borders.
type Rect struct {
	x         int32
	y         int32
	h         int32
	w         int32
	radius    int32
	reload    bool
	shown     bool
	rendering bool

	borderWidth int32
	borderColor color.RGBA
	fillColor   color.RGBA

	curTweens map[string]*Tween
	tex       rl.RenderTexture2D
}

func NewRect(x, y, height, width int32) *Rect {
	r := &Rect{
		x:           x,
		y:           y,
		h:           height,
		w:           width,
		radius:      0,
		borderWidth: 2,
		borderColor: rl.Black,
		fillColor:   rl.Blank,
		reload:      true,
		curTweens:   make(map[string]*Tween),
	}
	r.tex = rl.LoadRenderTexture(width, height)
	rl.SetTextureFilter(r.tex.Texture, rl.FilterAnisotropic8x|rl.FilterBilinear)
	return r
}

func (r *Rect) Copy() *Rect {
	return &Rect{
		x:           r.x,
		y:           r.y,
		h:           r.h,
		w:           r.w,
		radius:      r.radius,
		borderWidth: r.borderWidth,
		borderColor: r.borderColor,
		fillColor:   r.fillColor,
	}
}

func (r *Rect) addTween(key string, old, new int32) {
	r.curTweens[key] = NewTween(old, new, 1000*time.Millisecond)
}

func (r *Rect) SetPosition(x, y int32) {
	if !r.shown {
		r.x, r.y = x, y
		return
	}
	if r.x != x {
		r.addTween("x", r.x, x)
	}
	if r.y != y {
		r.addTween("y", r.y, y)
	}
}

func (r *Rect) SetSize(h, w int32) {
	if !r.shown {
		r.h, r.w = h, w
		return
	}
	if r.h != h {
		r.addTween("h", r.h, h)
		r.reload = true
	}
	if r.w != w {
		r.addTween("w", r.w, w)
		r.reload = true
	}
}

func (r *Rect) SetBorderRadius(radius int32) {
	if !r.shown {
		r.radius = radius
		return
	}
	if r.radius != radius {
		r.addTween("radius", r.radius, radius)
		r.reload = true
	}
}

func (r *Rect) SetBorderWidth(w int32) {
	if !r.shown {
		r.borderWidth = w
		return
	}
	if r.borderWidth != w {
		r.addTween("borderWidth", r.borderWidth, w)
		r.reload = true
	}
}

func (r *Rect) SetBorderColor(c color.RGBA) {
	if r.borderColor != c {
		r.reload = true
	}
	r.borderColor = c
}

func (r *Rect) SetFillColor(c color.RGBA) {
	if r.fillColor != c {
		r.reload = true
	}
	r.fillColor = c
}

func (r *Rect) Draw() {
	r.shown = true
	if !r.reload && len(r.curTweens) == 0 {
		rl.DrawTextureEx(r.tex.Texture, rl.NewVector2(float32(r.x), float32(r.y)), 0, 1, rl.White)
		return
	}
	var curRect *Rect
	updateTex := false
	if len(r.curTweens) != 0 {
		curRect = r.Copy()
		tween, ok := r.curTweens["x"]
		if ok {
			curRect.x = tween.CurVal()
			if tween.Ended() {
				r.x = curRect.x
				delete(r.curTweens, "x")
			}
		}
		tween, ok = r.curTweens["y"]
		if ok {
			curRect.y = tween.CurVal()
			if tween.Ended() {
				r.y = curRect.y
				delete(r.curTweens, "y")
			}
		}
		tween, ok = r.curTweens["h"]
		if ok {
			curRect.h = tween.CurVal()
			if curRect.h != r.h {
				updateTex = true
			}
			fmt.Println("h", curRect.h, r.h)
			if tween.Ended() {
				r.h = curRect.h
				delete(r.curTweens, "h")
			}
		}
		tween, ok = r.curTweens["w"]
		if ok {
			curRect.w = tween.CurVal()
			if curRect.w != r.w {
				updateTex = true
			}
			fmt.Println("w", curRect.w, r.w)
			if tween.Ended() {
				r.w = curRect.w
				delete(r.curTweens, "w")
			}
		}
		tween, ok = r.curTweens["radius"]
		if ok {
			curRect.radius = tween.CurVal()
			if curRect.radius != r.radius {
				updateTex = true
			}
			if tween.Ended() {
				r.radius = curRect.radius
				delete(r.curTweens, "radius")
			}
		}
		tween, ok = r.curTweens["borderWidth"]
		if ok {
			curRect.borderWidth = tween.CurVal()
			if curRect.borderWidth != r.borderWidth {
				updateTex = true
			}
			if tween.Ended() {
				r.borderWidth = curRect.borderWidth
				delete(r.curTweens, "borderWidth")
			}
		}
	} else {
		curRect = r
	}
	if r.reload || updateTex {
		// rl.UnloadTexture(r.tex)
		// r.tex = rl.LoadTextureFromImage(curRect.buildImage())
		// rl.UpdateTexture(r.tex, rl.LoadImageColors(curRect.buildImage()))
		rl.BeginTextureMode(r.tex)
		rl.UpdateTextureRec(r.tex.Texture, rl.NewRectangle(0, 0, float32(curRect.w), float32(curRect.h)), rl.LoadImageColors(curRect.buildImage()))
		rl.EndTextureMode()
	}
	rl.DrawTextureEx(r.tex.Texture, rl.NewVector2(float32(curRect.x), float32(curRect.y)), 0, 1, rl.White)
	r.reload = false
}

func (r Rect) buildImage() *rl.Image {
	// Create image at 16x resolution then resize down so it looks better
	img := rl.NewImageFromImage(image.NewAlpha(image.Rect(0, 0, int(r.w), int(r.h))))
	// Filled areas
	if r.fillColor != rl.Blank {
		rl.ImageDrawRectangle(img, r.radius, 0, (r.w - (2 * r.radius)), r.h, r.fillColor)
		if r.radius > 0 {
			rl.ImageDrawRectangle(img, (r.w - r.radius), r.radius, r.radius, (r.h - (2 * r.radius)), r.fillColor)
			rl.ImageDrawRectangle(img, 0, r.radius, r.radius, (r.h - (2 * r.radius)), r.fillColor)
		}
	}

	if r.borderColor != rl.Blank {
		// Horizontal lines
		rl.ImageDrawRectangle(img, r.radius, 0, (r.w - (2 * r.radius)), r.borderWidth, r.borderColor)
		rl.ImageDrawRectangle(img, r.radius, (r.h)-(r.borderWidth), (r.w - (2 * r.radius)), r.borderWidth, r.borderColor)
		// Vertical lines
		rl.ImageDrawRectangle(img, 0, r.radius, r.borderWidth, (r.h - (2 * r.radius)), r.borderColor)
		rl.ImageDrawRectangle(img, (r.w)-(r.borderWidth), r.radius, r.borderWidth, (r.h - (2 * r.radius)), r.borderColor)
	}
	if r.radius > 0 {
		r.placeCorners(r.radius, img)
	}
	return img
}

func (r Rect) placeCorners(radius int32, img *rl.Image) {
	// Setup corner
	corner := rl.NewImageFromImage(image.NewAlpha(image.Rect(0, 0, int(radius), int(radius))))
	if r.borderWidth > 0 {
		rl.ImageDrawCircle(corner, radius, radius, radius, r.borderColor)
	}
	rl.ImageDrawCircle(corner, radius, radius, (radius)-(r.borderWidth), r.fillColor)
	// Add corners to actual image
	rl.ImageDraw(img, corner, rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.White)
	rl.ImageRotateCW(corner)
	rl.ImageDraw(img, corner, rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.NewRectangle(float32((r.w-radius)), 0, float32(corner.Width), float32(corner.Height)), rl.White)
	rl.ImageRotateCW(corner)
	rl.ImageDraw(img, corner, rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.NewRectangle(float32((r.w-radius)), float32((r.h-radius)), float32(corner.Width), float32(corner.Height)), rl.White)
	rl.ImageRotateCW(corner)
	rl.ImageDraw(img, corner, rl.NewRectangle(0, 0, float32(corner.Width), float32(corner.Height)), rl.NewRectangle(0, float32((r.h-radius)), float32(corner.Width), float32(corner.Height)), rl.White)
}
