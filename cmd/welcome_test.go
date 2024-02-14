package cmd


import "testing"

func TestWelcomeMessage(t *testing.T) {
	got := welcomeMessage() 
	want := "Welcome To the Interviewer\n"
	
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
