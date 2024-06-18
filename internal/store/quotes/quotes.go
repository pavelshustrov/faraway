package quotes

import "math/rand"

type QuoteStore struct {
	quotes []string
}

func NewStore() *QuoteStore {
	return &QuoteStore{
		quotes: []string{
			"quote_1",
			"quote_2",
			"quote_3",
			"quote_4",
			"quote_5",
		},
	}
}

func (qs QuoteStore) GetRandomQuote() string {
	return qs.quotes[rand.Intn(len(qs.quotes))]
}
