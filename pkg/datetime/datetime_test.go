package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatToAgo(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name           string
		input          time.Time
		expectedOutput string
	}{
		{
			name:           "10 minutes ago",
			input:          now.Add(-10 * time.Minute),
			expectedOutput: "10 minutes ago",
		},
		{
			name:           "2 hours ago",
			input:          now.Add(-2 * time.Hour),
			expectedOutput: "2 hours ago",
		},
		{
			name:           "1 day ago",
			input:          now.Add(-24 * time.Hour),
			expectedOutput: "1 day ago",
		},
		{
			name:           "1 week ago",
			input:          now.Add(-7 * 24 * time.Hour),
			expectedOutput: "1 week ago",
		},
		{
			name:           "4 minutes from now",
			input:          now.Add(5 * time.Minute),
			expectedOutput: "4 minutes from now",
		},
		{
			name:           "6 days from now",
			input:          now.Add(7 * 24 * time.Hour),
			expectedOutput: "6 days from now",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualOutput := FormatToAgo(testCase.input)

			assert.Equal(t, testCase.expectedOutput, actualOutput)
		})
	}
}

func TestFormatToDateTime(t *testing.T) {
	input := time.Date(2024, 12, 12, 12, 30, 0, 0, time.UTC)
	actualOutput := FormatToDateTime(input)

	assert.Equal(t, "2024-12-12 12:30", actualOutput)
}
