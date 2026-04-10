package storage

import (
	"reflect"
	"testing"
)

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected []string
	}{
		{
			name:     "English stopwords filtered",
			title:    "The quick brown fox",
			expected: []string{"quick", "brown", "fox"},
		},
		{
			name:     "Spanish stopwords filtered",
			title:    "El gato con la pelota",
			expected: []string{"gato", "pelota"},
		},
		{
			name:     "Spanish article with meaningful keywords",
			title:    "El presidente habla sobre la economía",
			expected: []string{"presidente", "habla", "economía"},
		},
		{
			name:     "Mixed English and Spanish",
			title:    "Python y JavaScript son lenguajes",
			expected: []string{"python", "javascript", "lenguajes"},
		},
		{
			name:     "Short words filtered (less than 3 chars)",
			title:    "Go is a language",
			expected: []string{"language"},
		},
		{
			name:     "Common Spanish stopwords: que, una, con, de",
			title:    "Una noticia que habla de tecnología con detalles",
			expected: []string{"noticia", "habla", "tecnología", "detalles"},
		},
		{
			name:     "Lowercase normalization",
			title:    "BREAKING NEWS About TECHNOLOGY",
			expected: []string{"breaking", "news", "technology"},
		},
		{
			name:     "Spanish prepositions filtered",
			title:    "Artículo para usuarios desde España hasta México",
			expected: []string{"artículo", "usuarios", "españa", "méxico"},
		},
		{
			name:     "All stopwords, no keywords",
			title:    "El de la con y para",
			expected: []string{},
		},
		{
			name:     "Special characters stripped",
			title:    "Machine-learning & AI: the future!",
			expected: []string{"machine", "learning", "future"},
		},
		{
			name:     "Spanish verbs filtered (es, son, está, están)",
			title:    "OpenAI es líder tecnológico que está innovando",
			expected: []string{"openai", "líder", "tecnológico", "innovando"},
		},
		{
			name:     "Spanish pronouns filtered",
			title:    "Yo pienso que tú necesitas este libro",
			expected: []string{"pienso", "necesitas", "libro"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractKeywords(tt.title)
			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("extractKeywords(%q) = %v, expected %v", tt.title, result, tt.expected)
			}
		})
	}
}

func TestExtractKeywordsEmptyInput(t *testing.T) {
	result := extractKeywords("")
	if len(result) != 0 {
		t.Errorf("extractKeywords(\"\") should return empty slice, got %v", result)
	}
}

func TestExtractKeywordsNoValidWords(t *testing.T) {
	// Only stopwords and short words
	result := extractKeywords("a is in on by")
	if len(result) != 0 {
		t.Errorf("extractKeywords with only stopwords should return empty slice, got %v", result)
	}
}
