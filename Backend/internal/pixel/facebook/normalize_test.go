package facebook

import "testing"

func TestNormalizePhoneIndonesia(t *testing.T) {
	cases := map[string]string{
		"081234567890":        "6281234567890",
		"+62 812-3456-7890":   "6281234567890",
		"6281234567890":       "6281234567890",
	}
	for in, want := range cases {
		got := NormalizePhone(in, "62")
		if got != want {
			t.Fatalf("%q => %q, want %q", in, got, want)
		}
	}
}

func TestHashEmailDeterministic(t *testing.T) {
	a := HashEmail("Test@Example.com")
	b := HashEmail("  test@example.com  ")
	if a == "" || a != b {
		t.Fatalf("hash mismatch: %q vs %q", a, b)
	}
}

func TestFormatFBC(t *testing.T) {
	got := FormatFBC("IwAR2xxx", 1700000000)
	want := "fb.1.1700000000.IwAR2xxx"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestBuildServerEventPurchaseRequiresValue(t *testing.T) {
	props := map[string]any{
		"url": "https://example.com/thanks", "user_agent": "Mozilla/5.0",
		"client_ip": "1.2.3.4", "fbp": "fb.1.1.1",
		"em": []string{HashEmail("a@b.com")},
	}
	_, err := BuildServerEvent("Purchase", "e1", 1, props)
	if err == nil {
		t.Fatal("expected error without value/currency")
	}
}

func TestBuildServerEventPageViewTierC(t *testing.T) {
	props := map[string]any{
		"url": "https://example.com/", "user_agent": "Mozilla/5.0",
		"client_ip": "1.2.3.4", "fbp": "fb.1.1.1",
	}
	ev, err := BuildServerEvent("PageView", "e1", 1, props)
	if err != nil {
		t.Fatal(err)
	}
	if ev.UserData.FBP == "" || ev.EventSourceURL == "" {
		t.Fatal("missing required website fields")
	}
	if ev.UserData.QualityTier(true) != "C" {
		t.Fatalf("tier %s want C", ev.UserData.QualityTier(true))
	}
}
