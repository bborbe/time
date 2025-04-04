// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/bborbe/time"
)

type CurrentDateTime struct {
	NowStub        func() time.DateTime
	nowMutex       sync.RWMutex
	nowArgsForCall []struct {
	}
	nowReturns struct {
		result1 time.DateTime
	}
	nowReturnsOnCall map[int]struct {
		result1 time.DateTime
	}
	SetNowStub        func(time.DateTime)
	setNowMutex       sync.RWMutex
	setNowArgsForCall []struct {
		arg1 time.DateTime
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CurrentDateTime) Now() time.DateTime {
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

func (fake *CurrentDateTime) NowCallCount() int {
	fake.nowMutex.RLock()
	defer fake.nowMutex.RUnlock()
	return len(fake.nowArgsForCall)
}

func (fake *CurrentDateTime) NowCalls(stub func() time.DateTime) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = stub
}

func (fake *CurrentDateTime) NowReturns(result1 time.DateTime) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = nil
	fake.nowReturns = struct {
		result1 time.DateTime
	}{result1}
}

func (fake *CurrentDateTime) NowReturnsOnCall(i int, result1 time.DateTime) {
	fake.nowMutex.Lock()
	defer fake.nowMutex.Unlock()
	fake.NowStub = nil
	if fake.nowReturnsOnCall == nil {
		fake.nowReturnsOnCall = make(map[int]struct {
			result1 time.DateTime
		})
	}
	fake.nowReturnsOnCall[i] = struct {
		result1 time.DateTime
	}{result1}
}

func (fake *CurrentDateTime) SetNow(arg1 time.DateTime) {
	fake.setNowMutex.Lock()
	fake.setNowArgsForCall = append(fake.setNowArgsForCall, struct {
		arg1 time.DateTime
	}{arg1})
	stub := fake.SetNowStub
	fake.recordInvocation("SetNow", []interface{}{arg1})
	fake.setNowMutex.Unlock()
	if stub != nil {
		fake.SetNowStub(arg1)
	}
}

func (fake *CurrentDateTime) SetNowCallCount() int {
	fake.setNowMutex.RLock()
	defer fake.setNowMutex.RUnlock()
	return len(fake.setNowArgsForCall)
}

func (fake *CurrentDateTime) SetNowCalls(stub func(time.DateTime)) {
	fake.setNowMutex.Lock()
	defer fake.setNowMutex.Unlock()
	fake.SetNowStub = stub
}

func (fake *CurrentDateTime) SetNowArgsForCall(i int) time.DateTime {
	fake.setNowMutex.RLock()
	defer fake.setNowMutex.RUnlock()
	argsForCall := fake.setNowArgsForCall[i]
	return argsForCall.arg1
}

func (fake *CurrentDateTime) Invocations() map[string][][]interface{} {
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

func (fake *CurrentDateTime) recordInvocation(key string, args []interface{}) {
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

var _ time.CurrentDateTime = new(CurrentDateTime)
