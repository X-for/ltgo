package client

import (
	"errors"

	"github.com/X-for/ltgo/internal/models"
)

func (c *Client) GetQuestions(limit, skip int) (*models.QuestionListResponse, error) {
	query := `
    query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
        problemsetQuestionList: questionList(
            categorySlug: $categorySlug
            limit: $limit
            skip: $skip
            filters: $filters
        ) {
            total: totalNum
            questions: data {
                questionFrontendId
                title
                titleSlug
                difficulty
                status
                isPaidOnly
            }
        }
    }`

	vars := map[string]interface{}{
		"categorySlug": "",
		"skip":         skip,
		"limit":        limit,
		"filters":      map[string]interface{}{},
	}

	var resp models.QuestionListResponse
	if err := c.GraphQL(query, vars, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetQuestionDetail 获取单题详情 (描述 + 代码模板)
func (c *Client) GetQuestionDetail(titleSlug string) (*models.QuestionDetail, error) {
	query := `
    query questionData($titleSlug: String!) {
        question(titleSlug: $titleSlug) {
            questionFrontendId
            title
            titleSlug
            content
            translatedContent
            difficulty
            sampleTestCase
            codeSnippets {
                lang
                langSlug
                code
            }
        }
    }`

	vars := map[string]interface{}{
		"titleSlug": titleSlug,
	}

	var resp models.QuestionDetailResponse
	if err := c.GraphQL(query, vars, &resp); err != nil {
		return nil, err
	}

	if resp.Data.Question.Title == "" {
		return nil, errors.New("question not found")
	}

	return &resp.Data.Question, nil
}
