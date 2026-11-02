package validate

import "testing"

func TestEmail(t *testing.T) {
	if !Email("") || !Email("a@b.co") {
		t.Fatal("valid email rejected")
	}
	if Email("nope") || Email("a@b") || Email("a b@c.com") || Email("@b.co") {
		t.Fatal("bad email accepted")
	}
}

func TestCheckFile(t *testing.T) {
	if err := CheckFile("note.pdf", "application/pdf", 4); err != nil {
		t.Fatal(err)
	}
	if err := CheckFile("note.pdf", "application/octet-stream", 4); err != nil {
		t.Fatal(err)
	}
	if err := CheckFile("note.pdf", "image/png", 4); err == nil {
		t.Fatal("mismatched type was accepted")
	}
	if err := CheckFile("note.txt", "text/plain", 4); err == nil {
		t.Fatal("txt was accepted")
	}
	if err := CheckFile("run.exe", "application/octet-stream", 4); err == nil {
		t.Fatal("exe was accepted")
	}
	if err := CheckFile("big.pdf", "application/pdf", MaxBytes+1); err == nil {
		t.Fatal("oversize file was accepted")
	}
}
