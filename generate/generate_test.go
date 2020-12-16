package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeName(t *testing.T) {
	results := map[string]string{
		"-test-":  "  test - , ,",
		"unknown": "",
		"ministère-de-l-agriculture-et-de-l-alimentation.xml":                 "Ministère-de-l'agriculture-et-de-l'alimentation.xml",
		"commission-nationale-de-l-informatique-et-des-libertés.xml":          "Commission-nationale-de-l'informatique-et-des-libertés.xml",
		"ministère-de-l-économie-des-finances-et-de-la-relance-industrie.xml": "Ministère-de-l'économie-des-finances-et-de-la-relance-Industrie.xml",
		"ministère-de-l-intérieur.xml":                                        "Ministère-de-l'intérieur.xml",
		"autorité-nationale-des-jeux.xml":                                     "Autorité-nationale-des-jeux.xml",
		"èéôû-ç-éàè":                                                          "èéôû'ç ,ÉÀÈ",
	}

	for expectedOutput, input := range results {
		helper(t, expectedOutput, input)
	}

}

func helper(t *testing.T, expectedOutput string, input string) {
	assert.Equal(t, expectedOutput, sanitizeName(input))
}
