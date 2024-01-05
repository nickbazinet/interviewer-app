package cmd


import (
	"testing"
)

func TestGetCategory(t *testing.T) {
	got, err := getCategory("local")
	want := []string{"ansible", "terraform", "aws"}

	var topic []string
	for _, category := range got {
		topic = append(topic, category.Name)	
	}

	if err != nil {
		t.Errorf("Error was not null %s", err)
	}


	if len(got) != len(topic) {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetCategory_nonExistant(t *testing.T) {
	_, err := getCategory("anyrandomstringthatdoesntexist")
	
	if err == nil {
		t.Error("Got a nil error. Should have existing error for non-supported value")
	}
}

func TestGetCategory_chatgpt(t *testing.T) {
	_, err := getCategory("chatgpt")

	if err != nil {
		t.Errorf("Error was not null %s", err)
	}
}

