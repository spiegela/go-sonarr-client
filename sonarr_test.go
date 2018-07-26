package sonarr

import "testing"

func TestAppendURL(t *testing.T) {
	endpoint := "api/series"
	expectedURL := "http://192.168.1.25:8989/" + endpoint
	expectedURLWithBase := "http://192.168.1.25:8989/sonarr/" + endpoint

	testInput := []string{
		"http://192.168.1.25:8989",
		"http://192.168.1.25:8989/",
		"http://192.168.1.25:8989/sonarr",
	}

	for _, input := range testInput {
		if fullURL := appendEndpoint(input, endpoint); fullURL != expectedURL && fullURL != expectedURLWithBase {
			t.Errorf("expected either\n\t- '%s'\n\t- '%s'\ngot\n\t- '%s'", expectedURL, expectedURLWithBase, fullURL)
		}
	}
}
