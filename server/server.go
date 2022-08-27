package server

import (
	"sync"
)

func Start() (err error) {

	dependencies, err := NewDependencies()
	if err != nil {
		return
	}

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go func() {
		startHTTPServer(dependencies)
		wg.Done()
	}()

	wg.Wait()
	return
}
