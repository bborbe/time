// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
)

var _ = Describe("WaiterDuration", func() {
	Describe("NewWaiterDuration", func() {
		It("returns a WaiterDuration instance", func() {
			waiter := libtime.NewWaiterDuration()
			Expect(waiter).NotTo(BeNil())
		})

		It("implements WaiterDuration interface", func() {
			var waiter libtime.WaiterDuration = libtime.NewWaiterDuration()
			Expect(waiter).NotTo(BeNil())
		})
	})

	Describe("WaiterDurationFunc", func() {
		It("implements WaiterDuration interface", func() {
			callCount := 0
			waiterFunc := libtime.WaiterDurationFunc(func(ctx context.Context, duration libtime.Duration) error {
				callCount++
				return nil
			})

			var waiter libtime.WaiterDuration = waiterFunc
			Expect(waiter).NotTo(BeNil())

			err := waiter.Wait(context.Background(), libtime.Duration(1*time.Millisecond))
			Expect(err).To(BeNil())
			Expect(callCount).To(Equal(1))
		})

		It("calls the underlying function with correct parameters", func() {
			var receivedCtx context.Context
			var receivedDuration libtime.Duration

			waiterFunc := libtime.WaiterDurationFunc(func(ctx context.Context, duration libtime.Duration) error {
				receivedCtx = ctx
				receivedDuration = duration
				return nil
			})

			expectedCtx := context.Background()
			expectedDuration := libtime.Duration(5 * time.Second)

			err := waiterFunc.Wait(expectedCtx, expectedDuration)
			Expect(err).To(BeNil())
			Expect(receivedCtx).To(Equal(expectedCtx))
			Expect(receivedDuration).To(Equal(expectedDuration))
		})
	})

	Describe("Wait", func() {
		var waiter libtime.WaiterDuration

		BeforeEach(func() {
			waiter = libtime.NewWaiterDuration()
		})

		Describe("with positive duration", func() {
			It("waits for the specified duration", func() {
				start := time.Now()
				duration := libtime.Duration(50 * time.Millisecond)

				err := waiter.Wait(context.Background(), duration)

				elapsed := time.Since(start)
				Expect(err).To(BeNil())
				Expect(elapsed).To(BeNumerically(">=", 50*time.Millisecond))
				Expect(elapsed).To(BeNumerically("<", 200*time.Millisecond)) // Allow some tolerance
			})

			It("returns nil when wait completes successfully", func() {
				duration := libtime.Duration(1 * time.Millisecond)

				err := waiter.Wait(context.Background(), duration)
				Expect(err).To(BeNil())
			})

			It("waits for longer durations", func() {
				start := time.Now()
				duration := libtime.Duration(100 * time.Millisecond)

				err := waiter.Wait(context.Background(), duration)

				elapsed := time.Since(start)
				Expect(err).To(BeNil())
				Expect(elapsed).To(BeNumerically(">=", 100*time.Millisecond))
				Expect(elapsed).To(BeNumerically("<", 300*time.Millisecond))
			})
		})

		Describe("with zero or negative duration", func() {
			It("returns immediately with zero duration", func() {
				start := time.Now()
				duration := libtime.Duration(0)

				err := waiter.Wait(context.Background(), duration)

				elapsed := time.Since(start)
				Expect(err).To(BeNil())
				Expect(elapsed).To(BeNumerically("<", 10*time.Millisecond)) // Should be almost instant
			})

			It("returns immediately with negative duration", func() {
				start := time.Now()
				duration := libtime.Duration(-5 * time.Second)

				err := waiter.Wait(context.Background(), duration)

				elapsed := time.Since(start)
				Expect(err).To(BeNil())
				Expect(elapsed).To(BeNumerically("<", 10*time.Millisecond)) // Should be almost instant
			})
		})

		Describe("with context cancellation", func() {
			It("returns context error when context is cancelled before wait", func() {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately

				duration := libtime.Duration(1 * time.Second)
				err := waiter.Wait(ctx, duration)

				Expect(err).To(Equal(context.Canceled))
			})

			It("returns context error when context is cancelled during wait", func() {
				ctx, cancel := context.WithCancel(context.Background())

				// Start the wait in a goroutine
				errChan := make(chan error, 1)
				go func() {
					duration := libtime.Duration(500 * time.Millisecond)
					errChan <- waiter.Wait(ctx, duration)
				}()

				// Cancel the context after a short delay
				time.Sleep(50 * time.Millisecond)
				cancel()

				// Should get cancelled error
				var err error
				Eventually(errChan).Should(Receive(&err))
				Expect(err).To(Equal(context.Canceled))
			})

			It("returns context deadline exceeded when context times out", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
				defer cancel()

				duration := libtime.Duration(200 * time.Millisecond) // Longer than context timeout
				err := waiter.Wait(ctx, duration)

				Expect(err).To(Equal(context.DeadlineExceeded))
			})

			It("completes successfully when context timeout is longer than wait duration", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
				defer cancel()

				duration := libtime.Duration(50 * time.Millisecond) // Shorter than context timeout
				err := waiter.Wait(ctx, duration)

				Expect(err).To(BeNil())
			})
		})

		Describe("edge cases", func() {
			It("handles very short durations", func() {
				start := time.Now()
				duration := libtime.Duration(1 * time.Nanosecond)

				err := waiter.Wait(context.Background(), duration)

				elapsed := time.Since(start)
				Expect(err).To(BeNil())
				// Even 1 nanosecond should complete quickly
				Expect(elapsed).To(BeNumerically("<", 50*time.Millisecond))
			})

			It("handles microsecond durations", func() {
				start := time.Now()
				duration := libtime.Duration(100 * time.Microsecond)

				err := waiter.Wait(context.Background(), duration)

				elapsed := time.Since(start)
				Expect(err).To(BeNil())
				Expect(elapsed).To(BeNumerically("<", 50*time.Millisecond))
			})

			It("handles context that is already done", func() {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				// With positive duration, should check context and return cancelled
				duration := libtime.Duration(1 * time.Second)
				err := waiter.Wait(ctx, duration)

				Expect(err).To(Equal(context.Canceled))
			})

			It("handles multiple concurrent waits", func() {
				const numWaiters = 5
				duration := libtime.Duration(50 * time.Millisecond)
				errChan := make(chan error, numWaiters)

				start := time.Now()
				for i := 0; i < numWaiters; i++ {
					go func() {
						errChan <- waiter.Wait(context.Background(), duration)
					}()
				}

				// Collect all results
				for i := 0; i < numWaiters; i++ {
					var err error
					Eventually(errChan).Should(Receive(&err))
					Expect(err).To(BeNil())
				}

				elapsed := time.Since(start)
				// All should complete around the same time
				Expect(elapsed).To(BeNumerically(">=", 50*time.Millisecond))
				Expect(elapsed).To(BeNumerically("<", 200*time.Millisecond))
			})
		})

		Describe("interface compliance", func() {
			It("implements WaiterDuration interface", func() {
				var w libtime.WaiterDuration = waiter
				Expect(w).NotTo(BeNil())

				err := w.Wait(context.Background(), libtime.Duration(1*time.Millisecond))
				Expect(err).To(BeNil())
			})
		})
	})
})
