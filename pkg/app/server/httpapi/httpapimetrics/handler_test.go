package httpapimetrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestLabels(t *testing.T) {

	tests := []struct {
		path      string
		reqMethod string
		status    int
		expected  prometheus.Labels
	}{
		{
			path:      "/api",
			reqMethod: "GET",
			status:    200,
			expected: prometheus.Labels{
				pathLabel:   "/api",
				codeLabel:   "200",
				methodLabel: "get",
			},
		},
		{
			path:      "/users",
			reqMethod: "POST",
			status:    404,
			expected: prometheus.Labels{
				pathLabel:   "/users",
				codeLabel:   "404",
				methodLabel: "post",
			},
		},
		{
			path:      "/tests",
			reqMethod: "PUT",
			status:    500,
			expected: prometheus.Labels{
				pathLabel:   "/tests",
				codeLabel:   "500",
				methodLabel: "put",
			},
		},
	}

	for _, test := range tests {
		result := labels(test.path, test.reqMethod, test.status)
		assert.Equal(t, test.expected, result)
	}
}

func TestSanitizeMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"GET", "get"},
		{"get", "get"},
		{"PUT", "put"},
		{"put", "put"},
		{"HEAD", "head"},
		{"head", "head"},
		{"POST", "post"},
		{"post", "post"},
		{"DELETE", "delete"},
		{"delete", "delete"},
		{"CONNECT", "connect"},
		{"connect", "connect"},
		{"OPTIONS", "options"},
		{"options", "options"},
		{"NOTIFY", "notify"},
		{"notify", "notify"},
		{"INVALID", "invalid"},
		{"invalid", "invalid"},
	}

	for _, test := range tests {
		result := sanitizeMethod(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestSanitizeCode(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{100, "100"},
		{101, "101"},
		{200, "200"},
		{0, "200"},
		{201, "201"},
		{202, "202"},
		{203, "203"},
		{204, "204"},
		{205, "205"},
		{206, "206"},
		{300, "300"},
		{301, "301"},
		{302, "302"},
		{304, "304"},
		{305, "305"},
		{307, "307"},
		{400, "400"},
		{401, "401"},
		{402, "402"},
		{403, "403"},
		{404, "404"},
		{405, "405"},
		{406, "406"},
		{407, "407"},
		{408, "408"},
		{409, "409"},
		{410, "410"},
		{411, "411"},
		{412, "412"},
		{413, "413"},
		{414, "414"},
		{415, "415"},
		{416, "416"},
		{417, "417"},
		{418, "418"},
		{500, "500"},
		{501, "501"},
		{502, "502"},
		{503, "503"},
		{504, "504"},
		{505, "505"},
		{428, "428"},
		{429, "429"},
		{431, "431"},
		{511, "511"},
		{600, "600"},
	}

	for _, test := range tests {
		result := sanitizeCode(test.input)
		assert.Equal(t, test.expected, result)
	}
}
