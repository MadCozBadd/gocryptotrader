package news

import "testing"

func TestCheckOtherThings(t *testing.T) {
	err := CheckOtherThings()
	if err != nil {
		t.Error(err)
	}
}

func TestSendMessage(t *testing.T) {
	err := SendMessage("HELLoooooo")
	if err != nil {
		t.Error(err)
	}
}

func TestWriteFile(t *testing.T) {
	err := WriteFile([]string{"bitcoin", "litecoin", "ethereum"})
	if err != nil {
		t.Error(err)
	}
}

func TestReadFile(t *testing.T) {
	_, err := ReadFile()
	if err != nil {
		t.Error(err)
	}
}
