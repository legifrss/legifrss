package generate

import (
	"testing"

	"github.com/gorilla/feeds"
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
		"ministère-de-l-éducation-nationale-de-la-jeunesse-et-des-sports.xml": "ministère-de-l-éducation-nationale-de-la-jeunesse-et-des-sports-sports.xml",
	}

	for expectedOutput, input := range results {
		helper(t, expectedOutput, input)
	}

}

func helper(t *testing.T, expectedOutput string, input string) {
	assert.Equal(t, expectedOutput, sanitizeName(input))
}

func TestMergeFeeds(t *testing.T) {
	feed1 := feeds.AtomFeed{
		Entries: []*feeds.AtomEntry{
			{Id: "1"}, {Id: "2"}, {Id: "3"}},
	}
	feed2 := feeds.AtomFeed{
		Entries: []*feeds.AtomEntry{
			{Id: "4", Content: &feeds.AtomContent{Content: "test"}}},
	}
	f := mergeFeeds(feed2, feed1, 2)
	assert.Equal(t, "4", f.Entries[0].Id)
	assert.Equal(t, "test", f.Entries[0].Content.Content)
	assert.Equal(t, 2, len(f.Entries))
}
