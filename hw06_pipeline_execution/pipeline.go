package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrap(stage Stage, in In, done In) Out {
	merge := make(Bi)

	go func() {
		defer close(merge)

		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}

				merge <- val
			}
		}
	}()

	return stage(merge)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = wrap(stage, out, done)
	}

	return out
}
