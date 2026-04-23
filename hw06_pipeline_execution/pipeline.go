package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	firstIn := make(Bi)

	go func() {
		defer close(firstIn)
		for {
			select {
			case <-done:
				drain(in)
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case firstIn <- val:
				case <-done:
					drain(in)
					return
				}
			}
		}
	}()

	currentIn := In(firstIn)
	for _, stage := range stages {
		currentIn = runStage(currentIn, done, stage)
	}

	return currentIn
}

func runStage(in In, done In, stage Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		stageOut := stage(in)

		for {
			select {
			case <-done:
				drain(stageOut)
				return
			case val, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					drain(stageOut)
					return
				}
			}
		}
	}()

	return out
}

func drain(ch In) {
	for range ch {
	}
}
