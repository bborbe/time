// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
	libtimetest "github.com/bborbe/time/test"
)

var _ = Describe("WaiterUntil", func() {
	var currentDateTime libtime.CurrentDateTime

	BeforeEach(func() {
		currentDateTime = libtime.NewCurrentDateTime()
	})

	Describe("NewWaiterUntil", func() {
		It("returns a WaiterUntil instance", func() {
			waiter := libtime.NewWaiterUntil(currentDateTime)
			Expect(waiter).NotTo(BeNil())
		})

		It("implements WaiterUntil interface", func() {
			var waiter libtime.WaiterUntil = libtime.NewWaiterUntil(currentDateTime)
			Expect(waiter).NotTo(BeNil())
		})

		It("accepts CurrentDateTimeGetter interface", func() {
			getter := libtime.CurrentDateTimeGetterFunc(func() libtime.DateTime {
				return libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
			})
			waiter := libtime.NewWaiterUntil(getter)
			Expect(waiter).NotTo(BeNil())
		})
	})

	Describe("WaiterUntilFunc", func() {
		It("implements WaiterUntil interface", func() {
			callCount := 0
			waiterFunc := libtime.WaiterUntilFunc(func(ctx context.Context, until libtime.DateTime) error {
				callCount++
				return nil
			})

			var waiter libtime.WaiterUntil = waiterFunc
			Expect(waiter).NotTo(BeNil())

			until := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
			err := waiter.WaitUntil(context.Background(), until)
			Expect(err).To(BeNil())
			Expect(callCount).To(Equal(1))
		})

		It("calls the underlying function with correct parameters", func() {
			var receivedCtx context.Context
			var receivedUntil libtime.DateTime

			waiterFunc := libtime.WaiterUntilFunc(func(ctx context.Context, until libtime.DateTime) error {
				receivedCtx = ctx
				receivedUntil = until
				return nil
			})

			expectedCtx := context.Background()
			expectedUntil := libtimetest.ParseDateTime("2023-12-25T15:30:45Z")

			err := waiterFunc.WaitUntil(expectedCtx, expectedUntil)
			Expect(err).To(BeNil())
			Expect(receivedCtx).To(Equal(expectedCtx))
			Expect(receivedUntil).To(Equal(expectedUntil))
		})
	})

	Describe("WaitUntil", func() {
		var waiter libtime.WaiterUntil

		BeforeEach(func() {
			waiter = libtime.NewWaiterUntil(currentDateTime)
		})

		Describe("with future time", func() {
			It("waits until the specified time", func(ctx context.Context) {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				// Wait until 50ms in the future
				until := now.Add(libtime.Duration(50 * time.Millisecond))

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), until)
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				// Should wait for the specified duration
				expectedDuration := 50 * time.Millisecond
				Expect(elapsed).To(BeNumerically(">=", expectedDuration))
				Expect(elapsed).To(BeNumerically("<", expectedDuration+200*time.Millisecond))
			}, SpecTimeout(2*time.Second))

			It("waits for exact duration without buffer", func(ctx context.Context) {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				// Wait until 1ms in the future
				until := now.Add(libtime.Duration(1 * time.Millisecond))

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), until)
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				// Should wait for exactly 1ms
				expectedDuration := 1 * time.Millisecond
				Expect(elapsed).To(BeNumerically(">=", expectedDuration))
				Expect(elapsed).To(BeNumerically("<", expectedDuration+200*time.Millisecond))
			}, SpecTimeout(2*time.Second))

			It("returns nil when wait completes successfully", func(ctx context.Context) {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				until := now.Add(libtime.Duration(1 * time.Millisecond))
				err := waiter.WaitUntil(context.Background(), until)
				Expect(err).To(BeNil())
			}, SpecTimeout(2*time.Second))
		})

		Describe("with past time", func() {
			It("returns immediately when until time is in the past", func() {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				// Until time is 1 hour in the past
				until := now.Add(libtime.Duration(-1 * time.Hour))

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), until)
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				Expect(elapsed).To(BeNumerically("<", 10*time.Millisecond)) // Should be almost instant
			})

			It("returns immediately when until time equals current time", func() {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				until := now // Same time

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), until)
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				Expect(elapsed).To(BeNumerically("<", 10*time.Millisecond)) // Should be almost instant
			})

			It("returns immediately when until time is significantly in the past", func() {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				until := libtimetest.ParseDateTime("2020-01-01T00:00:00Z") // Years in the past

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), until)
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				Expect(elapsed).To(BeNumerically("<", 10*time.Millisecond))
			})
		})

		Describe("with context cancellation", func() {
			It("returns context error when context is cancelled during wait", func() {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				until := now.Add(libtime.Duration(1 * time.Hour)) // Far future
				ctx, cancel := context.WithCancel(context.Background())

				// Start the wait in a goroutine
				errChan := make(chan error, 1)
				go func() {
					errChan <- waiter.WaitUntil(ctx, until)
				}()

				// Cancel after a short delay
				time.Sleep(50 * time.Millisecond)
				cancel()

				// Should get cancelled error
				var err error
				Eventually(errChan).Should(Receive(&err))
				Expect(err).To(Equal(context.Canceled))
			})

			It("returns context deadline exceeded when context times out", func() {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				until := now.Add(libtime.Duration(1 * time.Hour)) // Far future
				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
				defer cancel()

				err := waiter.WaitUntil(ctx, until)
				Expect(err).To(Equal(context.DeadlineExceeded))
			})

			It("returns context error immediately for past times even with cancelled context", func() {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				until := now.Add(libtime.Duration(-1 * time.Hour)) // Past time
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately

				err := waiter.WaitUntil(ctx, until)
				Expect(err).To(BeNil()) // Past times skip waiting and return nil
			})
		})

		Describe("with dynamic current time", func() {
			It("calculates wait duration based on current time at call time", func(ctx context.Context) {
				// Set initial time
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				// Set until time 50ms in the future from initial time
				until := now.Add(libtime.Duration(50 * time.Millisecond))

				// Change current time to be closer to until time (30ms in the future)
				currentDateTime.SetNow(now.Add(libtime.Duration(20 * time.Millisecond)))

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), until)
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				// Should wait for remaining 30ms
				expectedDuration := 30 * time.Millisecond
				Expect(elapsed).To(BeNumerically(">=", expectedDuration))
				Expect(elapsed).To(BeNumerically("<", expectedDuration+200*time.Millisecond))
			}, SpecTimeout(2*time.Second))

			It("uses current time from getter at call time", func() {
				callCount := 0
				getter := libtime.CurrentDateTimeGetterFunc(func() libtime.DateTime {
					callCount++
					return libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				})

				waiter := libtime.NewWaiterUntil(getter)
				until := libtimetest.ParseDateTime("2023-12-25T09:00:00Z") // Past time

				err := waiter.WaitUntil(context.Background(), until)
				Expect(err).To(BeNil())
				Expect(callCount).To(Equal(1)) // Should call getter once
			})
		})

		Describe("edge cases", func() {
			It("handles very short wait durations", func(ctx context.Context) {
				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)

				until := now.Add(libtime.Duration(1 * time.Nanosecond))

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), until)
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				// Should wait for 1ns
				expectedDuration := 1 * time.Nanosecond
				Expect(elapsed).To(BeNumerically(">=", expectedDuration))
			}, SpecTimeout(2*time.Second))

			It("handles zero time values", func() {
				currentDateTime.SetNow(libtime.DateTime{}) // Zero time
				until := libtime.DateTime{}                // Zero time

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), until)
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				// Zero times are equal, so this should return immediately
				Expect(elapsed).To(BeNumerically("<", 10*time.Millisecond))
			})

			It("handles different timezones correctly", func(ctx context.Context) {
				utc := time.UTC
				est := time.FixedZone("EST", -5*3600)

				nowUTC := time.Date(2023, 12, 25, 15, 0, 0, 0, utc)
				untilEST := time.Date(2023, 12, 25, 10, 0, 1, 0, est) // 1 second later in EST (same moment + 1s)

				currentDateTime.SetNow(libtime.DateTime(nowUTC))

				start := time.Now()
				err := waiter.WaitUntil(context.Background(), libtime.DateTime(untilEST))
				elapsed := time.Since(start)

				Expect(err).To(BeNil())
				// Should wait for 1 second
				expectedDuration := 1 * time.Second
				Expect(elapsed).To(BeNumerically(">=", expectedDuration))
				Expect(elapsed).To(BeNumerically("<", expectedDuration+200*time.Millisecond))
			}, SpecTimeout(3*time.Second))
		})

		Describe("interface compliance", func() {
			It("implements WaiterUntil interface", func() {
				var w libtime.WaiterUntil = waiter
				Expect(w).NotTo(BeNil())

				now := libtimetest.ParseDateTime("2023-12-25T10:00:00Z")
				currentDateTime.SetNow(now)
				until := now.Add(libtime.Duration(-1 * time.Hour)) // Past time for quick test

				err := w.WaitUntil(context.Background(), until)
				Expect(err).To(BeNil())
			})
		})
	})
})
