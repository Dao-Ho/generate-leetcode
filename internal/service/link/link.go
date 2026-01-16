package link

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type LinkService struct {
	client *http.Client
}

func NewLinkService() *LinkService {
	return &LinkService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *LinkService) GetLink(command string, difficulty string) (string, error) {
	switch command {
	case "random":
		return s.getRandomLeetcode(difficulty)
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
}

type graphQLRequest struct {
	Query string `json:"query"`
}

type problemsetResponse struct {
	Data struct {
		ProblemsetQuestionList struct {
			Questions []struct {
				Title      string `json:"title"`
				TitleSlug  string `json:"titleSlug"`
				Difficulty string `json:"difficulty"`
				IsPaidOnly bool   `json:"isPaidOnly"`
			} `json:"questions"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}

func (s *LinkService) getRandomLeetcode(difficulty string) (string, error) {
	query := graphQLRequest{
		Query: `query problemsetQuestionList {
			problemsetQuestionList: questionList(
				categorySlug: ""
				limit: 3000
				skip: 0
				filters: {}
			) {
				questions: data {
					title
					titleSlug
					difficulty
					isPaidOnly
				}
			}
		}`,
	}

	body, err := json.Marshal(query)
	if err != nil {
		return "", fmt.Errorf("failed to marshal query: %w", err)
	}

	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch problems: %w", err)
	}
	defer resp.Body.Close()

	var data problemsetResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Filter questions
	var filtered []struct {
		Title      string
		TitleSlug  string
		Difficulty string
	}

	for _, q := range data.Data.ProblemsetQuestionList.Questions {
		if q.IsPaidOnly {
			continue
		}
		// Filter by difficulty if specified
		if difficulty != "" && q.Difficulty != difficulty {
			continue
		}
		filtered = append(filtered, struct {
			Title      string
			TitleSlug  string
			Difficulty string
		}{q.Title, q.TitleSlug, q.Difficulty})
	}

	if len(filtered) == 0 {
		return "", fmt.Errorf("no questions found")
	}

	question := filtered[rand.Intn(len(filtered))]

	return fmt.Sprintf("*%s* (%s)\nhttps://leetcode.com/problems/%s/",
		question.Title,
		question.Difficulty,
		question.TitleSlug,
	), nil
}
