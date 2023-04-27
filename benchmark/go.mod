module github.com/mlange-42/arche/benchmark

go 1.20

require (
	github.com/marioolofo/go-gameengine-ecs v0.9.0
	github.com/mlange-42/arche v0.7.0
	github.com/pkg/profile v1.7.0
	github.com/wfranczyk/ento v0.1.0
	github.com/yohamta/donburi v1.3.4
)

replace github.com/mlange-42/arche v0.7.0 => ../

require (
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/google/pprof v0.0.0-20230406165453-00490a63f317 // indirect
)
