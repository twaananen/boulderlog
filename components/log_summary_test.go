package components

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestLogSummaryEscapesResponse(t *testing.T) {
	var output bytes.Buffer
	if err := LogSummary("<script>alert('test')</script>", 5, nil, nil).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "<script>alert") || !strings.Contains(output.String(), "&lt;script&gt;") {
		t.Fatal("response must render as escaped text")
	}
}
