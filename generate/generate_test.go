package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeName(t *testing.T) {
	results := map[string]string{
		"-test-":  "  test - , ,",
		"unknown": "",
		"Ministere-de-l-agriculture-et-de-l'alimentation.xml":                 "Ministère-de-l'agriculture-et-de-l'alimentation.xml",
		"Commission-nationale-de-l'informatique-et-des-libertes.xml":          "Commission-nationale-de-l'informatique-et-des-libertés.xml",
		"Ministere-de-l'economie-des-finances-et-de-la-relance-Industrie.xml": "Ministère-de-l'conomie-des-finances-et-de-la-relance-Industrie.xml",
		"Ministere-de-l'interieur.xml":                                        "Ministère-de-l'intérieur.xml",
		"Autorite-nationale-des-jeux.xml":                                     "Autorité-nationale-des-jeux.xml",
		"eeou-c-":                                                             "èéôû'ç ,",
	}

	for expectedOutput, input := range results {
		helper(t, expectedOutput, input)
	}

}

func helper(t *testing.T, expectedOutput string, input string) {
	assert.Equal(t, expectedOutput, sanitizeName(input))
}
