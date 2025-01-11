package ui

import (
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
		curRect.drawTexture()
		rl.EndTextureMode()
	}
	rl.DrawTexturePro(r.tex.Texture, rl.NewRectangle(0, 0, float32(curRect.w), float32(curRect.h)), rl.NewRectangle(float32(curRect.x), float32(curRect.y), float32(curRect.w), float32(curRect.h)), rl.NewVector2(0, 0), 0, rl.White)
	// rl.DrawTextureEx(r.tex.Texture, rl.NewVector2(float32(curRect.x), float32(curRect.y)), 0, 1, rl.White)
	r.reload = false
}

func (r Rect) drawTexture() {
	rl.ClearBackground(rl.Blank)
	// Filled areas
	if r.fillColor != rl.Blank {
		rl.DrawRectangle(r.radius, 0, (r.w - (2 * r.radius)), r.h, r.fillColor)
		if r.radius > 0 {
			rl.DrawRectangle((r.w - r.radius), r.radius, r.radius, (r.h - (2 * r.radius)), r.fillColor)
			rl.DrawRectangle(0, r.radius, r.radius, (r.h - (2 * r.radius)), r.fillColor)
		}
	}

	if r.borderColor != rl.Blank {
		// Horizontal lines
		rl.DrawRectangle(r.radius, 0, (r.w - (2 * r.radius)), r.borderWidth, r.borderColor)
		rl.DrawRectangle(r.radius, (r.h)-(r.borderWidth), (r.w - (2 * r.radius)), r.borderWidth, r.borderColor)
		// Vertical lines
		rl.DrawRectangle(0, r.radius, r.borderWidth, (r.h - (2 * r.radius)), r.borderColor)
		rl.DrawRectangle((r.w)-(r.borderWidth), r.radius, r.borderWidth, (r.h - (2 * r.radius)), r.borderColor)
	}
	if r.radius > 0 {
		r.placeCorners()
	}
}

func (r Rect) placeCorners() {
	//bl
	if r.borderWidth > 0 {
		rl.DrawRing(rl.NewVector2(float32(r.radius), float32(r.radius)), float32(r.radius-r.borderWidth), float32(r.radius), 180, 270, 20, r.borderColor)
	}
	rl.DrawCircleSector(rl.NewVector2(float32(r.radius), float32(r.radius)), float32(r.radius-r.borderWidth), 180, 270, 20, r.fillColor)
	//tl
	if r.borderWidth > 0 {
		rl.DrawRing(rl.NewVector2(float32(r.radius), float32(r.h-r.radius)), float32(r.radius-r.borderWidth), float32(r.radius), 90, 180, 20, r.borderColor)
	}
	rl.DrawCircleSector(rl.NewVector2(float32(r.radius), float32(r.h-r.radius)), float32(r.radius-r.borderWidth), 90, 180, 20, r.fillColor)
	//tr
	if r.borderWidth > 0 {
		rl.DrawRing(rl.NewVector2(float32(r.w-r.radius), float32(r.h-r.radius)), float32(r.radius-r.borderWidth), float32(r.radius), 0, 90, 20, r.borderColor)
	}
	rl.DrawCircleSector(rl.NewVector2(float32(r.w-r.radius), float32(r.h-r.radius)), float32(r.radius-r.borderWidth), 0, 90, 20, r.fillColor)
	//br
	if r.borderWidth > 0 {
		rl.DrawRing(rl.NewVector2(float32(r.w-r.radius), float32(r.radius)), float32(r.radius-r.borderWidth), float32(r.radius), 270, 360, 20, r.borderColor)
	}
	rl.DrawCircleSector(rl.NewVector2(float32(r.w-r.radius), float32(r.radius)), float32(r.radius-r.borderWidth), 270, 360, 20, r.fillColor)
}
