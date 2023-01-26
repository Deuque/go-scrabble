package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Deuque/slack-scrabble/util"
)

type Scrabbler struct {
	FetchWord func() (*string, error)
}

func NewMockScrabbler() Scrabbler {
	return Scrabbler{
		func() (*string, error) {
			testWord := []string{"flavor", "kettle", "Rainbow"}
			ri := util.RandomInt(len(testWord) - 1)
			return &testWord[ri], nil
		},
	}
}

func NewHttpScrabbler() Scrabbler {
	return Scrabbler{
		func() (*string, error) {
			req, err := http.NewRequest(http.MethodGet, "https://random-word-api.herokuapp.com/word", nil)
			if err != nil {
				return nil, err
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}

			var data []string
			dec := json.NewDecoder(res.Body)
			dec.Decode(&data)

			return &data[0], nil
		},
	}
}

func (s *Scrabbler) ScrabbleWord(word string) string {
	wordSplit := strings.Split(word, "")
	resultSplit := wordSplit
	wordLength := len(wordSplit)

	for i := 0; i < wordLength; i++ {
		position := i

		for position == i {
			position = util.RandomInt(wordLength - 1)
		}

		temp := resultSplit[position]
		resultSplit[position] = wordSplit[i]
		resultSplit[i] = temp
	}

	return strings.Join(resultSplit, "")
}

func (s *Scrabbler) CheckAnswer(word, answer string) bool {
	return strings.EqualFold(word, answer)
}
