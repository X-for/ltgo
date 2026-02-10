package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/X-for/ltgo/internal/models"
)

type AutoCompleteResponse struct {
	State string `json:"state"` // "success"
	Data  []struct {
		QuestionID         string `json:"question_id"`
		QuestionFrontendID string `json:"question_frontend_id"`
		Title              string `json:"title"`
		TitleSlug          string `json:"title_slug"`
		Difficulty         int    `json:"difficulty"` // 注意：这个接口返回的 difficulty 是 int (1,2,3)
		PaidOnly           bool   `json:"paid_only"`
		IsFavor            bool   `json:"is_favor"`
		Status             string `json:"status"` // "ac", null, etc.
	} `json:"data"`
}

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

// SearchQuestions 搜索流程：先拿 ID 列表，再拿详情
func (c *Client) SearchQuestions(keyword string) ([]models.Question, error) {
	// 1. 调用 REST API 获取符合条件的 Question ID 列表 (后端 ID，不是 FrontendID)
	path := fmt.Sprintf("/problems/api/filter-questions/all/?search_keywords=%s", keyword)
	respBody, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	// 解析 ID 列表 (例如: [1, 15, 203...])
	var questionIDs []int // 注意：这是后端 ID (QuestionID)，通常是数字
	if err := json.Unmarshal(respBody, &questionIDs); err != nil {
		return nil, fmt.Errorf("failed to parse search result IDs: %v", err)
	}

	if len(questionIDs) == 0 {
		return []models.Question{}, nil
	}

	// 限制一下数量，别一次查太多，取前 20 个
	if len(questionIDs) > 20 {
		questionIDs = questionIDs[:20]
	}

	// 2. 构造 GraphQL 请求，批量获取这些 ID 对应的题目详情
	// 我们使用 filters.questionIds 来精确查询
	query := `
    query problemsetQuestionListV2($filters: QuestionFilterInput, $limit: Int) {
        problemsetQuestionListV2(
            filters: $filters
            limit: $limit
        ) {
            questions {
                questionId
                questionFrontendId
                title
                titleSlug
                difficulty
                status
                paidOnly
            }
        }
    }`

	// 将 []int 转换为 []string，因为 GraphQL 参数通常是字符串数组
	var qidStrings []string
	for _, id := range questionIDs {
		qidStrings = append(qidStrings, fmt.Sprintf("%d", id))
	}

	vars := map[string]interface{}{
		"limit": 20,
		"filters": map[string]interface{}{
			// 这里是关键：用 questionIds 过滤器
			"questionIds": qidStrings,
		},
	}

	var resp models.QuestionListResponse
	if err := c.GraphQL(query, vars, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch question details: %v", err)
	}

	// 3. 排序优化 (可选)
	// GraphQL 返回的顺序可能和我们传进去的 ID 顺序不一样 (即搜索结果的相关性顺序)
	// 为了保持搜索的最佳匹配度，我们最好按 questionIDs 的顺序重新排一下
	// 但为了简单，先直接返回即可。

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
