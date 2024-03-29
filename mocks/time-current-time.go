// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"
	timea "time"

	"github.com/bborbe/time"
)

type TimeCurrentTime struct {
	NowStub        func() timea.Time
	nowMutex       sync.RWMutex
	nowArgsForCall []struct {
	}
	nowReturns struct {
		result1 timea.Time
	}
	nowReturnsOnCall map[int]struct {
		result1 timea.Time
	}
	SetNowStub        func(timea.Time)
	setNowMutex       sync.RWMutex
	setNowArgsForCall []struct {
		arg1 timea.Time
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TimeCurrentTime) Now() timea.Time {
	fake.nowMutex.Lock()
	ret, specificReturn := fake.nowReturnsOnCall[len(fake.nowArgsForCall)]
	fake.nowArgsForCall = append(fake.nowArgsForCall, struct {
	}{})
	stub := fake.NowStub
	fakeReturns := fake.nowReturns
	fake.recordInvocation("Now", []interface{}{})
	fake.nowMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *TimeCurrentTime) NowCallCount() int {
	fake.nowMutex.RLock()
	defer fake.nowMutex.RUnlock()
	return len(fake.nowArgsForCall)
}

func (fake *TimeCurrentTime) NowCalls(stub func() timea.Time) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = stub
}

func (fake *TimeCurrentTime) NowReturns(result1 timea.Time) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = nil
	fake.nowReturns = struct {
		result1 timea.Time
	}{result1}
}

func (fake *TimeCurrentTime) NowReturnsOnCall(i int, result1 timea.Time) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = nil
	if fake.nowReturnsOnCall == nil {
		fake.nowReturnsOnCall = make(map[int]struct {
			result1 timea.Time
		})
	}
	fake.nowReturnsOnCall[i] = struct {
		result1 timea.Time
	}{result1}
}

func (fake *TimeCurrentTime) SetNow(arg1 timea.Time) {
	fake.setNowMutex.Lock()
	fake.setNowArgsForCall = append(fake.setNowArgsForCall, struct {
		arg1 timea.Time
	}{arg1})
	stub := fake.SetNowStub
	fake.recordInvocation("SetNow", []interface{}{arg1})
	fake.setNowMutex.Unlock()
	if stub != nil {
		fake.SetNowStub(arg1)
	}
}

func (fake *TimeCurrentTime) SetNowCallCount() int {
	fake.setNowMutex.RLock()
	defer fake.setNowMutex.RUnlock()
	return len(fake.setNowArgsForCall)
}

func (fake *TimeCurrentTime) SetNowCalls(stub func(timea.Time)) {
	fake.setNowMutex.Lock()
	defer fake.setNowMutex.Unlock()
	fake.SetNowStub = stub
}

func (fake *TimeCurrentTime) SetNowArgsForCall(i int) timea.Time {
	fake.setNowMutex.RLock()
	defer fake.setNowMutex.RUnlock()
	argsForCall := fake.setNowArgsForCall[i]
	return argsForCall.arg1
}

func (fake *TimeCurrentTime) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.nowMutex.RLock()
	defer fake.nowMutex.RUnlock()
	fake.setNowMutex.RLock()
	defer fake.setNowMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TimeCurrentTime) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ time.CurrentTime = new(TimeCurrentTime)
