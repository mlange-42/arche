module github.com/mlange-42/arche/benchmark

go 1.21

toolchain go1.21.3

require (
	github.com/marioolofo/go-gameengine-ecs v0.9.0
	github.com/mlange-42/arche v0.10.1
	github.com/pkg/profile v1.7.0
	github.com/unitoftime/ecs v0.0.2-0.20240109122000-af4227c75194
	github.com/wfranczyk/ento v0.1.0
	github.com/yohamta/donburi v1.3.4
)

replace github.com/mlange-42/arche v0.10.1 => ../

require (
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/google/pprof v0.0.0-20230406165453-00490a63f317 // indirect
	github.com/unitoftime/cod v0.0.0-20230616173404-085cf4fe3918 // indirect
)
