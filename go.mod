module g3-engine

go 1.21

replace (
	github.com/jfigge/guilib v0.0.3 => ../guilib
)

require (
	github.com/jfigge/guilib v0.0.3
	github.com/veandco/go-sdl2 v0.4.35
)
