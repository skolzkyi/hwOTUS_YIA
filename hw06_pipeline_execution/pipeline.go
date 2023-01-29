package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	chainStep := in
	for i := range stages {
		chainStep = func(in In) (out Out) {
			workChannel := make(Bi)
			go func() {
				defer close(workChannel)
				for {
					select {
					case <-done:
						return
					case item, ok := <-in:
						if !ok {
							return
						}
						workChannel <- item
					}
				}
			}()
			return stages[i](workChannel)
		}(chainStep)
	}
	return chainStep
}
