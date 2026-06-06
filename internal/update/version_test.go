package update

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input      string
		wantMajor  int
		wantMinor  int
		wantPatch  int
		wantStable bool
		wantOK     bool
	}{
		{"1.0.0", 1, 0, 0, true, true},
		{"v1.0.0", 1, 0, 0, true, true},
		{"1.1.0", 1, 1, 0, true, true},
		{"2.0.0", 2, 0, 0, true, true},
		{"1.0.1", 1, 0, 1, true, true},
		{"1.9.9", 1, 9, 9, true, true},
		{"1.1.0-dev", 1, 1, 0, false, true},
		{"0.6.0-rc.2", 0, 6, 0, false, true},
		{"v0.6.0-rc.1", 0, 6, 0, false, true},
		{"1.1.0-SNAPSHOT-abc", 1, 1, 0, false, true},
		{"", 0, 0, 0, false, false},
		{"notaversion", 0, 0, 0, false, false},
		{"1.2", 0, 0, 0, false, false},
	}

	for _, tc := range tests {
		sv, stable, ok := parseVersion(tc.input)
		if ok != tc.wantOK {
			t.Errorf("parseVersion(%q): ok=%v, want %v", tc.input, ok, tc.wantOK)
			continue
		}
		if !ok {
			continue
		}
		if sv.major != tc.wantMajor || sv.minor != tc.wantMinor || sv.patch != tc.wantPatch {
			t.Errorf("parseVersion(%q): got %d.%d.%d, want %d.%d.%d",
				tc.input, sv.major, sv.minor, sv.patch,
				tc.wantMajor, tc.wantMinor, tc.wantPatch)
		}
		if stable != tc.wantStable {
			t.Errorf("parseVersion(%q): stable=%v, want %v", tc.input, stable, tc.wantStable)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"v1.0.0", "1.0.0", 0},
		{"1.1.0", "1.0.0", 1},
		{"1.0.0", "1.1.0", -1},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.9.9", "2.0.0", -1},
		{"1.0.0", "1.0.0", 0},
	}

	for _, tc := range tests {
		svA, _, _ := parseVersion(tc.a)
		svB, _, _ := parseVersion(tc.b)
		got := compareVersions(svA, svB)
		if got != tc.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestIsDevBuild(t *testing.T) {
	cases := []struct {
		v    string
		want bool
	}{
		{"1.1.0-dev", true},
		{"v1.1.0-dev", true},
		{"1.1.0-SNAPSHOT-abc", true},
		{"1.0.0", false},
		{"v1.0.0", false},
		{"0.6.0-rc.2", false},
		{"1.0.0-beta", false},
	}

	for _, tc := range cases {
		got := isDevBuild(tc.v)
		if got != tc.want {
			t.Errorf("isDevBuild(%q) = %v, want %v", tc.v, got, tc.want)
		}
	}
}

func TestNormalizeVersion(t *testing.T) {
	if got := normalizeVersion("v1.0.0"); got != "1.0.0" {
		t.Errorf("normalizeVersion(v1.0.0) = %q, want %q", got, "1.0.0")
	}
	if got := normalizeVersion("1.0.0"); got != "1.0.0" {
		t.Errorf("normalizeVersion(1.0.0) = %q, want %q", got, "1.0.0")
	}
}
