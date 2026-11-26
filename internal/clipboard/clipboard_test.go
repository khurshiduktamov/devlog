package clipboard

import "testing"

func TestCandidates(t *testing.T) {
	cases := []struct {
		goos    string
		wayland bool
		first   string
	}{
		{"darwin", false, "pbcopy"},
		{"windows", false, "powershell.exe"},
		{"linux", false, "xclip"},
		{"linux", true, "wl-copy"},
	}
	for _, c := range cases {
		got := candidates(c.goos, c.wayland)
		if len(got) == 0 || got[0].name != c.first {
			t.Errorf("%s wayland=%v: first tool = %v, want %s", c.goos, c.wayland, got, c.first)
		}
	}
}
