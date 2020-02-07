package news

import "testing"

func TestGetData(t *testing.T) {
	err := GetData()
	if err != nil {
		t.Error(err)
	}
}

func TestCheck(t *testing.T) {
	err := Check()
	if err != nil {
		t.Error(err)
	}
}
func TestCheckOtherThings(t *testing.T) {
	err := CheckOtherThings()
	if err != nil {
		t.Error(err)
	}
}

func TestSendMessage(t *testing.T) {
	err := SendMessage("HELLLLLLLooooooooooo")
	if err != nil {
		t.Error(err)
	}
}
