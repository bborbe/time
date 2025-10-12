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

var _ = Describe("CurrentDateTime", func() {
	Describe("NewCurrentDateTime", func() {
		It("returns a CurrentDateTime instance", func() {
			currentDateTime := libtime.NewCurrentDateTime()
			Expect(currentDateTime).NotTo(BeNil())
		})

		It("implements CurrentDateTime interface", func() {
			var currentDateTime libtime.CurrentDateTime = libtime.NewCurrentDateTime()
			Expect(currentDateTime).NotTo(BeNil())
		})
	})

	Describe("CurrentDateTimeGetterFunc", func() {
		It("implements CurrentDateTimeGetter interface", func() {
			fixedTime := libtimetest.ParseDateTime("2023-12-25T10:15:30Z")
			getterFunc := libtime.CurrentDateTimeGetterFunc(func() libtime.DateTime {
				return fixedTime
			})

			var getter libtime.CurrentDateTimeGetter = getterFunc
			Expect(getter.Now()).To(Equal(fixedTime))
		})

		It("calls the underlying function", func() {
			fixedTime := libtimetest.ParseDateTime("2023-01-01T00:00:00Z")
			callCount := 0
			getterFunc := libtime.CurrentDateTimeGetterFunc(func() libtime.DateTime {
				callCount++
				return fixedTime
			})

			result := getterFunc.Now()
			Expect(result).To(Equal(fixedTime))
			Expect(callCount).To(Equal(1))
		})
	})

	Describe("currentDateTime", func() {
		var currentDateTime libtime.CurrentDateTime

		BeforeEach(func() {
			currentDateTime = libtime.NewCurrentDateTime()
		})

		Describe("Now", func() {
			Context("when no fixed time is set", func() {
				It("returns current time", func() {
					now1 := currentDateTime.Now()
					time.Sleep(1 * time.Millisecond)
					now2 := currentDateTime.Now()

					// Should be different times since we're getting actual current time
					Expect(now1.Time()).To(BeTemporally("<=", now2.Time()))
				})

				It("returns some time value when no fixed time is set", func() {
					result := currentDateTime.Now()
					// Just verify we get a non-zero time value
					Expect(result).NotTo(BeZero())
				})
			})

			Context("when fixed time is set", func() {
				It("returns the fixed time", func() {
					fixedTime := libtimetest.ParseDateTime("2023-12-25T10:15:30Z")
					currentDateTime.SetNow(fixedTime)

					result := currentDateTime.Now()
					Expect(result).To(Equal(fixedTime))
				})

				It("returns same fixed time on multiple calls", func() {
					fixedTime := libtimetest.ParseDateTime("2023-01-01T12:00:00Z")
					currentDateTime.SetNow(fixedTime)

					result1 := currentDateTime.Now()
					result2 := currentDateTime.Now()

					Expect(result1).To(Equal(fixedTime))
					Expect(result2).To(Equal(fixedTime))
					Expect(result1).To(Equal(result2))
				})
			})
		})

		Describe("SetNow", func() {
			It("sets the fixed time", func() {
				fixedTime := libtimetest.ParseDateTime("2023-06-15T14:30:45Z")
				currentDateTime.SetNow(fixedTime)

				result := currentDateTime.Now()
				Expect(result).To(Equal(fixedTime))
			})

			It("overwrites previously set time", func() {
				firstTime := libtimetest.ParseDateTime("2023-01-01T00:00:00Z")
				secondTime := libtimetest.ParseDateTime("2023-12-31T23:59:59Z")

				currentDateTime.SetNow(firstTime)
				Expect(currentDateTime.Now()).To(Equal(firstTime))

				currentDateTime.SetNow(secondTime)
				Expect(currentDateTime.Now()).To(Equal(secondTime))
			})

			It("accepts zero time", func() {
				zeroTime := libtime.DateTime{}
				currentDateTime.SetNow(zeroTime)

				result := currentDateTime.Now()
				Expect(result).To(Equal(zeroTime))
			})
		})

		Describe("thread safety", func() {
			It("is safe for concurrent Now calls", func() {
				fixedTime := libtimetest.ParseDateTime("2023-07-04T16:20:30Z")
				currentDateTime.SetNow(fixedTime)

				const numGoroutines = 100
				results := make([]libtime.DateTime, numGoroutines)
				var wg sync.WaitGroup

				for i := 0; i < numGoroutines; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						results[index] = currentDateTime.Now()
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
				times := make([]libtime.DateTime, numGoroutines)

				// Create different times
				for i := 0; i < numGoroutines; i++ {
					times[i] = libtimetest.ParseDateTime("2023-01-01T00:00:00Z").
						Add(libtime.Duration(time.Duration(i) * time.Hour))
				}

				var wg sync.WaitGroup
				for i := 0; i < numGoroutines; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						currentDateTime.SetNow(times[index])
					}(i)
				}

				wg.Wait()

				// Should have some valid time set (one of the times from the goroutines)
				result := currentDateTime.Now()
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
				fixedTime := libtimetest.ParseDateTime("2023-08-15T09:30:15Z")
				const numReaders = 50
				const numWriters = 5

				var wg sync.WaitGroup
				results := make([]libtime.DateTime, numReaders)

				// Start readers
				for i := 0; i < numReaders; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						results[index] = currentDateTime.Now()
					}(i)
				}

				// Start writers
				for i := 0; i < numWriters; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						timeToSet := fixedTime.Add(
							libtime.Duration(time.Duration(index) * time.Minute),
						)
						currentDateTime.SetNow(timeToSet)
					}(i)
				}

				wg.Wait()

				// All results should be valid DateTime values
				for i := 0; i < numReaders; i++ {
					Expect(results[i]).NotTo(BeNil())
				}
			})
		})

		Describe("interface compliance", func() {
			It("implements CurrentDateTimeGetter", func() {
				var getter libtime.CurrentDateTimeGetter = currentDateTime
				Expect(getter).NotTo(BeNil())

				result := getter.Now()
				Expect(result).NotTo(BeNil())
			})

			It("implements CurrentDateTimeSetter", func() {
				var setter libtime.CurrentDateTimeSetter = currentDateTime
				Expect(setter).NotTo(BeNil())

				fixedTime := libtimetest.ParseDateTime("2023-05-10T11:45:22Z")
				setter.SetNow(fixedTime)

				result := currentDateTime.Now()
				Expect(result).To(Equal(fixedTime))
			})

			It("implements CurrentDateTime", func() {
				var currentDT libtime.CurrentDateTime = currentDateTime
				Expect(currentDT).NotTo(BeNil())

				fixedTime := libtimetest.ParseDateTime("2023-03-20T08:15:45Z")
				currentDT.SetNow(fixedTime)

				result := currentDT.Now()
				Expect(result).To(Equal(fixedTime))
			})
		})
	})
})
