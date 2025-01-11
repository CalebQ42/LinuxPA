package ui

type num interface {
	int32 | float32
}

type Tween[T num] struct {
	start    T
	nextMult float32
	end      T
}
