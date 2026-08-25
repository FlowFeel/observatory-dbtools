package asset

import (
	"testing"
)

func TestExtractFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Observatory direct URL",
			input:    "https://observatory.wiki/w/images/9/9f/Observatory-logo-favicon.png",
			expected: "Observatory-logo-favicon.png",
		},
		{
			name:     "Observatory thumb URL",
			input:    "https://observatory.wiki/w/images/thumb/a/ab/Banner.jpg/800px-Banner.jpg",
			expected: "Banner.jpg",
		},
		{
			name:     "Localhost dev URL",
			input:    "http://localhost:8080/w/images/1/12/Author_Photo.png",
			expected: "Author_Photo.png",
		},
		{
			name:     "Wikimedia Commons URL",
			input:    "https://upload.wikimedia.org/wikipedia/commons/4/47/PNG_transparency_demonstration_1.png",
			expected: "PNG_transparency_demonstration_1.png",
		},
		{
			name:     "Already canonical bare filename",
			input:    "Johndoe.png",
			expected: "Johndoe.png",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ExtractFilename(tc.input)
			if actual != tc.expected {
				t.Errorf("ExtractFilename(%q) = %q, expected %q", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestRemediateWikitext(t *testing.T) {
	input := `{{Article
|Cover image=https://observatory.wiki/w/images/9/9f/Observatory-logo-favicon.png
|Author avatar=http://localhost:8085/w/images/a/ab/Jane_Doe.jpg
|Text=Some content with an inline image [[https://upload.wikimedia.org/wikipedia/commons/4/47/Sample.png|thumb]]
}}`

	expected := `{{Article
|Cover image=Observatory-logo-favicon.png
|Author avatar=Jane_Doe.jpg
|Text=Some content with an inline image [[Sample.png|thumb]]
}}`

	actual := RemediateWikitext(input)
	if actual != expected {
		t.Errorf("RemediateWikitext mismatch:\nGOT:\n%s\nEXPECTED:\n%s", actual, expected)
	}
}

func TestDetectHardcodedURLs(t *testing.T) {
	content := "Visit https://observatory.wiki/w/images/a/b/Test.jpg and https://upload.wikimedia.org/wikipedia/commons/c/d/Sample.png"
	urls := DetectHardcodedURLs(content)
	if len(urls) != 2 {
		t.Fatalf("expected 2 URLs, got %d: %v", len(urls), urls)
	}
}
