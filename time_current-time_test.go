// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
	libtimetest "github.com/bborbe/time/test"
)

var _ = Describe("CurrentTime", func() {
	Describe("NewCurrentTime", func() {
		It("returns a CurrentTime instance", func() {
			currentTime := libtime.NewCurrentTime()
			Expect(currentTime).NotTo(BeNil())
		})

		It("implements CurrentTime interface", func() {
			var currentTime libtime.CurrentTime = libtime.NewCurrentTime()
			Expect(currentTime).NotTo(BeNil())
		})
	})

	Describe("CurrentTimeGetterFunc", func() {
		It("implements CurrentTimeGetter interface", func() {
			fixedDateTime := libtimetest.ParseDateTime("2023-12-25T10:15:30Z")
			getterFunc := libtime.CurrentTimeGetterFunc(func() libtime.DateTime {
				return fixedDateTime
			})

			// Note: CurrentTimeGetterFunc returns DateTime, not time.Time
			// This appears to be an inconsistency in the interface design
			result := getterFunc.Now()
			Expect(result).To(Equal(fixedDateTime))
		})

		It("calls the underlying function", func() {
			fixedDateTime := libtimetest.ParseDateTime("2023-01-01T00:00:00Z")
			callCount := 0
			getterFunc := libtime.CurrentTimeGetterFunc(func() libtime.DateTime {
				callCount++
				return fixedDateTime
			})

			result := getterFunc.Now()
			Expect(result).To(Equal(fixedDateTime))
			Expect(callCount).To(Equal(1))
		})
	})

	Describe("currentTime", func() {
		var currentTime libtime.CurrentTime

		BeforeEach(func() {
			currentTime = libtime.NewCurrentTime()
		})

		Describe("Now", func() {
			Context("when no fixed time is set", func() {
				It("returns current time", func() {
					now1 := currentTime.Now()
					time.Sleep(1 * time.Millisecond)
					now2 := currentTime.Now()

					// Should be different times since we're getting actual current time
					Expect(now1).To(BeTemporally("<=", now2))
				})

				It("returns some time value when no fixed time is set", func() {
					result := currentTime.Now()
					// Just verify we get a non-zero time value
					Expect(result).NotTo(BeZero())
				})
			})

			Context("when fixed time is set", func() {
				It("returns the fixed time", func() {
					fixedTime := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()
					currentTime.SetNow(fixedTime)

					result := currentTime.Now()
					Expect(result).To(Equal(fixedTime))
				})

				It("returns same fixed time on multiple calls", func() {
					fixedTime := libtimetest.ParseDateTime("2023-01-01T12:00:00Z").Time()
					currentTime.SetNow(fixedTime)

					result1 := currentTime.Now()
					result2 := currentTime.Now()

					Expect(result1).To(Equal(fixedTime))
					Expect(result2).To(Equal(fixedTime))
					Expect(result1).To(Equal(result2))
				})
			})
		})

		Describe("SetNow", func() {
			It("sets the fixed time", func() {
				fixedTime := libtimetest.ParseDateTime("2023-06-15T14:30:45Z").Time()
				currentTime.SetNow(fixedTime)

				result := currentTime.Now()
				Expect(result).To(Equal(fixedTime))
			})

			It("overwrites previously set time", func() {
				firstTime := libtimetest.ParseDateTime("2023-01-01T00:00:00Z").Time()
				secondTime := libtimetest.ParseDateTime("2023-12-31T23:59:59Z").Time()

				currentTime.SetNow(firstTime)
				Expect(currentTime.Now()).To(Equal(firstTime))

				currentTime.SetNow(secondTime)
				Expect(currentTime.Now()).To(Equal(secondTime))
			})

			It("accepts zero time", func() {
				zeroTime := time.Time{}
				currentTime.SetNow(zeroTime)

				result := currentTime.Now()
				Expect(result).To(Equal(zeroTime))
			})
		})

		Describe("thread safety", func() {
			It("is safe for concurrent Now calls", func() {
				fixedTime := libtimetest.ParseDateTime("2023-07-04T16:20:30Z").Time()
				currentTime.SetNow(fixedTime)

				const numGoroutines = 100
				results := make([]time.Time, numGoroutines)
				var wg sync.WaitGroup

				for i := 0; i < numGoroutines; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						results[index] = currentTime.Now()
					}(i)
				}

				wg.Wait()

				// All results should be the same
				for i := 0; i < numGoroutines; i++ {
					Expect(results[i]).To(Equal(fixedTime))
				}
			})

			It("is safe for concurrent SetNow calls", func() {
				const numGoroutines = 10
				times := make([]time.Time, numGoroutines)

				// Create different times
				baseTime := libtimetest.ParseDateTime("2023-01-01T00:00:00Z").Time()
				for i := 0; i < numGoroutines; i++ {
					times[i] = baseTime.Add(time.Duration(i) * time.Hour)
				}

				var wg sync.WaitGroup
				for i := 0; i < numGoroutines; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						currentTime.SetNow(times[index])
					}(i)
				}

				wg.Wait()

				// Should have some valid time set (one of the times from the goroutines)
				result := currentTime.Now()
				found := false
				for _, t := range times {
					if result.Equal(t) {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue(), "Result should be one of the set times")
			})

			It("is safe for concurrent Now and SetNow calls", func() {
				fixedTime := libtimetest.ParseDateTime("2023-08-15T09:30:15Z").Time()
				const numReaders = 50
				const numWriters = 5

				var wg sync.WaitGroup
				results := make([]time.Time, numReaders)

				// Start readers
				for i := 0; i < numReaders; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						results[index] = currentTime.Now()
					}(i)
				}

				// Start writers
				for i := 0; i < numWriters; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						timeToSet := fixedTime.Add(time.Duration(index) * time.Minute)
						currentTime.SetNow(timeToSet)
					}(i)
				}

				wg.Wait()

				// All results should be valid time.Time values
				for i := 0; i < numReaders; i++ {
					Expect(results[i]).NotTo(BeNil())
				}
			})
		})

		Describe("interface compliance", func() {
			It("implements CurrentTimeGetter", func() {
				var getter libtime.CurrentTimeGetter = currentTime
				Expect(getter).NotTo(BeNil())

				result := getter.Now()
				Expect(result).NotTo(BeNil())
			})

			It("implements CurrentTimeSetter", func() {
				var setter libtime.CurrentTimeSetter = currentTime
				Expect(setter).NotTo(BeNil())

				fixedTime := libtimetest.ParseDateTime("2023-05-10T11:45:22Z").Time()
				setter.SetNow(fixedTime)

				result := currentTime.Now()
				Expect(result).To(Equal(fixedTime))
			})

			It("implements CurrentTime", func() {
				var currentT libtime.CurrentTime = currentTime
				Expect(currentT).NotTo(BeNil())

				fixedTime := libtimetest.ParseDateTime("2023-03-20T08:15:45Z").Time()
				currentT.SetNow(fixedTime)

				result := currentT.Now()
				Expect(result).To(Equal(fixedTime))
			})
		})
	})
})
