package helper

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic ASCII
		{"Hello World", "hello-world"},
		{"This is a test!", "this-is-a-test"},
		{"  Leading and trailing  ", "leading-and-trailing"},
		{"Multiple   spaces", "multiple-spaces"},
		{"Special@#*&Characters", "special-characters"},
		{"MiXeD CaSe", "mixed-case"},
		{"--Already-Slug--", "already-slug"},
		{"underscores_are_not-handled", "underscores-are-not-handled"},
		{"123 Numbers 456", "123-numbers-456"},
		{"", ""},
		{"---", ""},

		// Unicode and international scripts
		{"Café del Mar", "café-del-mar"},
		{"日本語のテキスト", "日本語のテキスト"},
		{"Привет мир", "привет-мир"},
		{"مرحبا بالعالم", "مرحبا-بالعالم"},
		{"你好，世界", "你好-世界"},
		{"😀 Smile Emoji 😀", "smile-emoji"},
		{"Grüße, München!", "grüße-münchen"},
		{"naïve approach", "naïve-approach"},
		{"São Paulo", "são-paulo"},
		{"L'année dernière", "l-année-dernière"},
		{"123✨abc", "123-abc"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := Slugify(tt.input)
			if actual != tt.expected {
				t.Errorf("Slugify(%q) = %q; expected %q", tt.input, actual, tt.expected)
			}
		})
	}
}
