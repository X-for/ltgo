package client

import (
	"errors"
	"strings"

	"github.com/X-for/ltgo/internal/models"
)

func (c *Client) GetQuestions(limit, skip int) (*models.QuestionListResponse, error) {
	// 专门针对 CN 的 V2 Query
	// 删除了 filters 参数定义和传参
	query := `
    query problemsetQuestionListV2($categorySlug: String, $limit: Int, $skip: Int) {
        problemsetQuestionListV2(
            categorySlug: $categorySlug
            limit: $limit
            skip: $skip
        ) {
            questions {
                questionFrontendId
                title
                translatedTitle
                titleSlug
                difficulty
                status
                paidOnly
            }
        }
    }`

	vars := map[string]interface{}{
		"categorySlug": "",
		"skip":         skip,
		"limit":        limit,
		// "filters":      map[string]interface{}{},
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
			questionId
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

// GetQuestionSlugByID 根据题目 ID (FrontendID) 查找 Slug
// 注意：这是一个比较重的操作，因为可能需要拉取很多题目
func (c *Client) GetQuestionSlugByID(id string) (string, error) {
	// 假设题目在前 3000 个里 (绝大多数情况够用了)
	resp, err := c.GetQuestions(3000, 0)
	if err != nil {
		return "", err
	}

	// 处理 V2 和 V1 兼容性
	questions := resp.Data.ProblemsetQuestionListV2.Questions
	if len(questions) == 0 {
		questions = resp.Data.ProblemsetQuestionList.Questions
	}

	for _, q := range questions {
		if q.QuestionFrontendID == id {
			return q.TitleSlug, nil
		}
	}

	return "", errors.New("question ID not found in the first 3000 questions")
}

// SearchQuestions 搜索题目 (支持 ID 或 标题关键字)
func (c *Client) SearchQuestions(keyword string) ([]models.Question, error) {
	// 构造 V2 Query，注意增加了 $searchKeyword 参数
	query := `
    query problemsetQuestionListV2($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionFilterInput) {
        problemsetQuestionListV2(
            categorySlug: $categorySlug
            limit: $limit
            skip: $skip
            filters: $filters
        ) {
            questions {
                questionFrontendId
                title
                translatedTitle
                titleSlug
                difficulty
                status
                paidOnly
            }
        }
    }`

	// 构造 Filters
	// 经验证，CN V2 的搜索其实是放在 filters.searchKeywords 里的
	// 或者 filters.search
	// 让我们先试一种最通用的：直接在 filters 里传 search

	vars := map[string]interface{}{
		"categorySlug": "",
		"limit":        20, // 限制返回 20 个搜索结果
		"skip":         0,
		"filters": map[string]interface{}{
			"searchKeywords": keyword, // [新增] 专门针对 CN 站点的搜索参数
			"search":         keyword, // [保留] 兼容 Global 站点
		},
	}

	var resp models.QuestionListResponse
	if err := c.GraphQL(query, vars, &resp); err != nil {
		return nil, err
	}

	questions := resp.Data.ProblemsetQuestionListV2.Questions
	if len(questions) == 0 {
		questions = resp.Data.ProblemsetQuestionList.Questions
	}

	return questions, nil
}

// SearchQuestionsByKeyword 在本地过滤题目
func (c *Client) SearchQuestionsByKeyword(keyword string) ([]models.Question, error) {
	// 1. 获取所有题目 (或者前 3000 个)
	// 实际上大多数人不需要这么全，我们可以先取 2000
	all, err := c.GetQuestions(2000, 0)
	if err != nil {
		return nil, err
	}

	questions := all.Data.ProblemsetQuestionListV2.Questions
	if len(questions) == 0 {
		questions = all.Data.ProblemsetQuestionList.Questions
	}

	// 2. 内存过滤
	var matched []models.Question
	keyword = strings.ToLower(keyword)

	for _, q := range questions {
		// 匹配 ID (精确匹配)
		if q.QuestionFrontendID == keyword {
			// 如果 ID 完全匹配，直接返回这一个
			return []models.Question{q}, nil
		}

		// 匹配 Title 或 Slug (模糊匹配)
		if strings.Contains(strings.ToLower(q.Title), keyword) ||
			strings.Contains(strings.ToLower(q.TitleSlug), keyword) {
			matched = append(matched, q)
		}
	}

	return matched, nil
}
