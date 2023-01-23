package pipeline

import "github.com/1ort/checkr/proxy"

type PipeLine interface {
	Run(in <-chan proxy.Proxy, done <-chan any) <-chan proxy.Proxy
	Plug(stage PipeLine) PipeLine
}

type BasePipeLine struct {
}

func (b BasePipeLine)

type sequentialPipeline struct {
	stages []PipeLine
}

// func BuildPipeline(stages ...Pipeline) Pipeline {

// }

func (p *sequentialPipeline) Run(in <-chan proxy.Proxy, done <-chan any) <-chan proxy.Proxy {
	for _, child := range p.childrenPipelines {
		in = child.Run(in, done)
	}
	return in
}

// FetchCounter := NewProxyCounter("found")
// pipeline := BuildPipeline(
// 	DefaultProxyProvider(),
// 	FetchCounter,

// )
// pipeline := NewPipeLine().Plug().Plug().Plug()
