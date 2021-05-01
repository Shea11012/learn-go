package main

import "sync"

var lock = &sync.Mutex{}

type single struct {}

var singleInstance *single

func getInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &single{}
		}
	}

	return singleInstance
}