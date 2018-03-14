package param_test

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/outo/metronome/param"
	"time"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("bpm unit test", func() {

	table.DescribeTable("Interval calculates the time span "+
		"between beats' beginnings, based on this bpm value",

		func(givenBpm param.Bpm, expectedInterval time.Duration) {
			gomega.Expect(givenBpm.Interval()).To(gomega.Equal(expectedInterval))
		},
		table.Entry("60 bpm", param.Bpm(60), time.Second),
		table.Entry("120 bpm", param.Bpm(120), 500*time.Millisecond),
		table.Entry("12000 bpm", param.Bpm(12000), 5*time.Millisecond),
		table.Entry("10 bpm", param.Bpm(10), 6*time.Second),
	)

})
