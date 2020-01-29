package coindesk

import "testing"

func TestGetData(t *testing.T) {
	err := GetData()
	if err != nil {
		t.Error(err)
	}
}
