package cmd

import "testing"

func TestEncoding(t *testing.T) {
	expectedEncoding := "Tm90aWZ5UFBNOkFnaWxlRGEkaGJvYXJkMQ=="
	u := "NotifyPPM"
	p := "AgileDa$hboard1"

	result := EncodeCreds(u, p)
	if result != expectedEncoding {
		t.Errorf("Expected encoding %s but found %s\n", expectedEncoding, result)
	}

}
