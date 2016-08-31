// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-cli/ssh"
)

type FakeSession struct {
	StartStub        func() ([]string, error)
	startMutex       sync.RWMutex
	startArgsForCall []struct{}
	startReturns     struct {
		result1 []string
		result2 error
	}
	FinishStub        func() error
	finishMutex       sync.RWMutex
	finishArgsForCall []struct{}
	finishReturns     struct {
		result1 error
	}
}

func (fake *FakeSession) Start() ([]string, error) {
	fake.startMutex.Lock()
	fake.startArgsForCall = append(fake.startArgsForCall, struct{}{})
	fake.startMutex.Unlock()
	if fake.StartStub != nil {
		return fake.StartStub()
	} else {
		return fake.startReturns.result1, fake.startReturns.result2
	}
}

func (fake *FakeSession) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *FakeSession) StartReturns(result1 []string, result2 error) {
	fake.StartStub = nil
	fake.startReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeSession) Finish() error {
	fake.finishMutex.Lock()
	fake.finishArgsForCall = append(fake.finishArgsForCall, struct{}{})
	fake.finishMutex.Unlock()
	if fake.FinishStub != nil {
		return fake.FinishStub()
	} else {
		return fake.finishReturns.result1
	}
}

func (fake *FakeSession) FinishCallCount() int {
	fake.finishMutex.RLock()
	defer fake.finishMutex.RUnlock()
	return len(fake.finishArgsForCall)
}

func (fake *FakeSession) FinishReturns(result1 error) {
	fake.FinishStub = nil
	fake.finishReturns = struct {
		result1 error
	}{result1}
}

var _ ssh.Session = new(FakeSession)
