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

var _ = Describe("Layout", func() {
	var layout libtime.Layout
	Context("Format", func() {
		var input time.Time
		var output string
		JustBeforeEach(func() {
			output = layout.Format(input)
		})
		BeforeEach(func() {
			input = time.Unix(1654248772, 123456789)
		})
		Context("RFC3339Nano", func() {
			BeforeEach(func() {
				layout = libtime.RFC3339Nano
			})
			It("has correct output", func() {
				Expect(output).To(Equal("2022-06-03T09:32:52.123456789Z"))
			})
		})
		Context("RFC3339", func() {
			BeforeEach(func() {
				layout = libtime.RFC3339
			})
			It("has correct output", func() {
				Expect(output).To(Equal("2022-06-03T09:32:52Z"))
			})
		})
		Context("DateLayout", func() {
			BeforeEach(func() {
				layout = libtime.DateLayout
			})
			It("has correct output", func() {
				Expect(output).To(Equal("2022-06-03"))
			})
		})
		Context("SecondLayout", func() {
			BeforeEach(func() {
				layout = libtime.SecondLayout
			})
			It("has correct output", func() {
				Expect(output).To(Equal("1654248772"))
			})
		})
		Context("MilliLayout", func() {
			BeforeEach(func() {
				layout = libtime.MilliLayout
			})
			It("has correct output", func() {
				Expect(output).To(Equal("1654248772123"))
			})
		})
		Context("MicroLayout", func() {
			BeforeEach(func() {
				layout = libtime.MicroLayout
			})
			It("has correct output", func() {
				Expect(output).To(Equal("1654248772123456"))
			})
		})
		Context("NanoLayout", func() {
			BeforeEach(func() {
				layout = libtime.NanoLayout
			})
			It("has correct output", func() {
				Expect(output).To(Equal("1654248772123456789"))
			})
		})
	})
	Context("Parse", func() {
		var ctx context.Context
		var input interface{}
		var output *time.Time
		var err error
		JustBeforeEach(func() {
			output, err = layout.Parse(ctx, input)
		})
		BeforeEach(func() {
			ctx = context.Background()
		})
		Context("RFC3339Nano", func() {
			BeforeEach(func() {
				layout = libtime.RFC3339Nano
				input = "2022-06-03T09:32:52.123456789Z"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.Unix()).To(Equal(int64(1654248772)))
				Expect(output.Nanosecond()).To(Equal(123456789))
			})
		})
		Context("RFC3339", func() {
			BeforeEach(func() {
				layout = libtime.RFC3339
				input = "2022-06-03T09:32:52Z"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.Unix()).To(Equal(int64(1654248772)))
			})
		})
		Context("DateLayout", func() {
			BeforeEach(func() {
				layout = libtime.DateLayout
				input = "2022-06-03"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.Year()).To(Equal(2022))
				Expect(output.Month()).To(Equal(time.June))
				Expect(output.Day()).To(Equal(3))
			})
		})
		Context("SecondLayout", func() {
			BeforeEach(func() {
				layout = libtime.SecondLayout
			})
			Context("with valid int64 string", func() {
				BeforeEach(func() {
					input = "1654248772"
				})
				It("parses successfully", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(output).ToNot(BeNil())
					Expect(output.Unix()).To(Equal(int64(1654248772)))
				})
			})
			Context("with int64", func() {
				BeforeEach(func() {
					input = int64(1654248772)
				})
				It("parses successfully", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(output).ToNot(BeNil())
					Expect(output.Unix()).To(Equal(int64(1654248772)))
				})
			})
			Context("with invalid string", func() {
				BeforeEach(func() {
					input = "invalid"
				})
				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(output).To(BeNil())
				})
			})
		})
		Context("MilliLayout", func() {
			BeforeEach(func() {
				layout = libtime.MilliLayout
			})
			Context("with valid int64 string", func() {
				BeforeEach(func() {
					input = "1654248772123"
				})
				It("parses successfully", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(output).ToNot(BeNil())
					Expect(output.Unix()).To(Equal(int64(1654248772)))
					Expect(output.UnixMilli()).To(Equal(int64(1654248772123)))
				})
			})
			Context("with int64", func() {
				BeforeEach(func() {
					input = int64(1654248772123)
				})
				It("parses successfully", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(output).ToNot(BeNil())
					Expect(output.UnixMilli()).To(Equal(int64(1654248772123)))
				})
			})
			Context("with invalid string", func() {
				BeforeEach(func() {
					input = "invalid"
				})
				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(output).To(BeNil())
				})
			})
		})
		Context("MicroLayout", func() {
			BeforeEach(func() {
				layout = libtime.MicroLayout
			})
			Context("with valid int64 string", func() {
				BeforeEach(func() {
					input = "1654248772123456"
				})
				It("parses successfully", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(output).ToNot(BeNil())
					Expect(output.Unix()).To(Equal(int64(1654248772)))
					Expect(output.UnixMicro()).To(Equal(int64(1654248772123456)))
				})
			})
			Context("with int64", func() {
				BeforeEach(func() {
					input = int64(1654248772123456)
				})
				It("parses successfully", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(output).ToNot(BeNil())
					Expect(output.UnixMicro()).To(Equal(int64(1654248772123456)))
				})
			})
			Context("with invalid string", func() {
				BeforeEach(func() {
					input = "invalid"
				})
				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(output).To(BeNil())
				})
			})
		})
		Context("NanoLayout", func() {
			BeforeEach(func() {
				layout = libtime.NanoLayout
			})
			Context("with valid int64 string", func() {
				BeforeEach(func() {
					input = "1654248772123456789"
				})
				It("parses successfully", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(output).ToNot(BeNil())
					Expect(output.Unix()).To(Equal(int64(1654248772)))
					Expect(output.UnixNano()).To(Equal(int64(1654248772123456789)))
				})
			})
			Context("with int64", func() {
				BeforeEach(func() {
					input = int64(1654248772123456789)
				})
				It("parses successfully", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(output).ToNot(BeNil())
					Expect(output.UnixNano()).To(Equal(int64(1654248772123456789)))
				})
			})
			Context("with invalid string", func() {
				BeforeEach(func() {
					input = "invalid"
				})
				It("returns error", func() {
					Expect(err).To(HaveOccurred())
					Expect(output).To(BeNil())
				})
			})
		})
		Context("with unsupported type", func() {
			BeforeEach(func() {
				layout = libtime.RFC3339
				input = 123 // int instead of string
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
			})
		})
		Context("with invalid string format", func() {
			BeforeEach(func() {
				layout = libtime.RFC3339
				input = "invalid-libtime-format"
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
			})
		})
	})
	Context("String", func() {
		It("returns string representation", func() {
			layout = libtime.RFC3339Nano
			Expect(layout.String()).To(Equal("2006-01-02T15:04:05.999999999Z07:00"))
		})
	})
})

var _ = Describe("Layouts", func() {
	var layouts libtime.Layouts
	var ctx context.Context
	var input interface{}
	var output *time.Time
	var err error
	JustBeforeEach(func() {
		output, err = layouts.Parse(ctx, input)
	})
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("with multiple layouts", func() {
		BeforeEach(func() {
			layouts = libtime.Layouts{libtime.RFC3339, libtime.SecondLayout}
		})
		Context("when first layout matches", func() {
			BeforeEach(func() {
				input = "2022-06-03T09:32:52Z"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.Unix()).To(Equal(int64(1654248772)))
			})
		})
		Context("when second layout matches", func() {
			BeforeEach(func() {
				input = "1654248772"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.Unix()).To(Equal(int64(1654248772)))
			})
		})
		Context("when no layout matches", func() {
			BeforeEach(func() {
				input = "invalid-format"
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("parse 'invalid-format' with any layouts failed"))
			})
		})
	})
	Context("with empty layouts", func() {
		BeforeEach(func() {
			layouts = libtime.Layouts{}
			input = "2022-06-03T09:32:52Z"
		})
		It("returns error", func() {
			Expect(err).To(HaveOccurred())
			Expect(output).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("parse '2022-06-03T09:32:52Z' with any layouts failed"))
		})
	})
})

var _ = Describe("Layout Edge Cases", func() {
	var layout libtime.Layout
	var ctx context.Context
	var input interface{}
	var output *time.Time
	var err error
	JustBeforeEach(func() {
		output, err = layout.Parse(ctx, input)
	})
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("Boundary conditions", func() {
		Context("SecondLayout with zero", func() {
			BeforeEach(func() {
				layout = libtime.SecondLayout
				input = "0"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.Unix()).To(Equal(int64(0)))
			})
		})
		Context("SecondLayout with negative value", func() {
			BeforeEach(func() {
				layout = libtime.SecondLayout
				input = "-1654248772"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.Unix()).To(Equal(int64(-1654248772)))
			})
		})
		Context("MilliLayout with zero", func() {
			BeforeEach(func() {
				layout = libtime.MilliLayout
				input = "0"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.UnixMilli()).To(Equal(int64(0)))
			})
		})
		Context("MilliLayout with negative value", func() {
			BeforeEach(func() {
				layout = libtime.MilliLayout
				input = "-1654248772123"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.UnixMilli()).To(Equal(int64(-1654248772123)))
			})
		})
		Context("MicroLayout with zero", func() {
			BeforeEach(func() {
				layout = libtime.MicroLayout
				input = "0"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.UnixMicro()).To(Equal(int64(0)))
			})
		})
		Context("MicroLayout with negative value", func() {
			BeforeEach(func() {
				layout = libtime.MicroLayout
				input = "-1654248772123456"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.UnixMicro()).To(Equal(int64(-1654248772123456)))
			})
		})
		Context("NanoLayout with zero", func() {
			BeforeEach(func() {
				layout = libtime.NanoLayout
				input = "0"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.UnixNano()).To(Equal(int64(0)))
			})
		})
		Context("NanoLayout with negative value", func() {
			BeforeEach(func() {
				layout = libtime.NanoLayout
				input = "-1654248772123456789"
			})
			It("parses successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.UnixNano()).To(Equal(int64(-1654248772123456789)))
			})
		})
	})
	Context("Error handling edge cases", func() {
		Context("SecondLayout with overflow string", func() {
			BeforeEach(func() {
				layout = libtime.SecondLayout
				input = "92233720368547758070" // Greater than max int64
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
			})
		})
		Context("MilliLayout with overflow string", func() {
			BeforeEach(func() {
				layout = libtime.MilliLayout
				input = "92233720368547758070" // Greater than max int64
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
			})
		})
		Context("MicroLayout with overflow string", func() {
			BeforeEach(func() {
				layout = libtime.MicroLayout
				input = "92233720368547758070" // Greater than max int64
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
			})
		})
		Context("NanoLayout with overflow string", func() {
			BeforeEach(func() {
				layout = libtime.NanoLayout
				input = "92233720368547758070" // Greater than max int64
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
			})
		})
		Context("SecondLayout with empty string", func() {
			BeforeEach(func() {
				layout = libtime.SecondLayout
				input = ""
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
			})
		})
		Context("RFC3339 with empty string", func() {
			BeforeEach(func() {
				layout = libtime.RFC3339
				input = ""
			})
			It("returns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(output).To(BeNil())
			})
		})
		Context("SecondLayout with float type", func() {
			BeforeEach(func() {
				layout = libtime.SecondLayout
				input = 1654248772.5
			})
			It("rounds to int64", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(output).ToNot(BeNil())
				Expect(output.Unix()).To(Equal(int64(1654248773)))
			})
		})
	})
})
