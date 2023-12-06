package cmd


import "testing"

func TestWelcomeMessage(t *testing.T) {
	got := welcomeMessage() 
	want := "Welcome To the Interviewer"
	
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
