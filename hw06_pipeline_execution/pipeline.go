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
				return
			case val, ok := <-in:
				if !ok {
					return
				}

				select {
				case firstIn <- val:
				case <-done:
					return
				}
			}
		}
	}()

	currentIn := In(firstIn)
	for i := 0; i < len(stages); i++ {
		currentIn = stages[i](currentIn)
	}

	return currentIn
}
