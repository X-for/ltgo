package client

import (
	"errors"

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

// SearchQuestions ...
// SearchQuestions 严格复刻抓包请求
func (c *Client) SearchQuestions(keyword string) ([]models.Question, error) {
	query := `
    query problemsetQuestionListV2($filters: QuestionFilterInput, $limit: Int, $searchKeyword: String, $skip: Int, $sortBy: QuestionSortByInput, $categorySlug: String) {
      problemsetQuestionListV2(
        filters: $filters
        limit: $limit
        searchKeyword: $searchKeyword
        skip: $skip
        sortBy: $sortBy
        categorySlug: $categorySlug
      ) {
        questions {
          titleSlug
          title
          translatedTitle
          questionFrontendId
          paidOnly
          difficulty
          status
        }
      }
    }`

	// 这里的 categorySlug 对于 CN 站必须是 "all-code-essentials"
	// 对于 COM 站，如果报错，可能需要改成 "" 或者不传，目前先优先保证 CN
	category := "all-code-essentials"
	if c.cfg.Site != "cn" {
		category = ""
	}

	vars := map[string]interface{}{
		"skip":          0,
		"limit":         20, // 抓包是100，我们这里改小点也可以，或者你改成100也行
		"categorySlug":  category,
		"searchKeyword": keyword,
		"sortBy": map[string]interface{}{
			"sortField": "CUSTOM",
			"sortOrder": "ASCENDING",
		},
		"filters": map[string]interface{}{
			"filterCombineType": "ALL",
			"statusFilter": map[string]interface{}{
				"questionStatuses": []string{},
				"operator":         "IS",
			},
			"difficultyFilter": map[string]interface{}{
				"difficulties": []string{},
				"operator":     "IS",
			},
			"languageFilter": map[string]interface{}{
				"languageSlugs": []string{},
				"operator":      "IS",
			},
			"topicFilter": map[string]interface{}{
				"topicSlugs": []string{},
				"operator":   "IS",
			},
			"acceptanceFilter":    map[string]interface{}{},
			"frequencyFilter":     map[string]interface{}{},
			"frontendIdFilter":    map[string]interface{}{},
			"lastSubmittedFilter": map[string]interface{}{},
			"publishedFilter":     map[string]interface{}{},
			"companyFilter": map[string]interface{}{
				"companySlugs": []string{},
				"operator":     "IS",
			},
			"positionFilter": map[string]interface{}{
				"positionSlugs": []string{},
				"operator":      "IS",
			},
			"positionLevelFilter": map[string]interface{}{
				"positionLevelSlugs": []string{},
				"operator":           "IS",
			},
			"contestPointFilter": map[string]interface{}{
				"contestPoints": []string{},
				"operator":      "IS",
			},
			"premiumFilter": map[string]interface{}{
				"premiumStatus": []string{},
				"operator":      "IS",
			},
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
//func (c *Client) SearchQuestionsByKeyword(keyword string) ([]models.Question, error) {
//	// 1. 获取所有题目 (或者前 3000 个)
//	// 实际上大多数人不需要这么全，我们可以先取 2000
//	all, err := c.GetQuestions(2000, 0)
//	if err != nil {
//		return nil, err
//	}
//
//	questions := all.Data.ProblemsetQuestionListV2.Questions
//	if len(questions) == 0 {
//		questions = all.Data.ProblemsetQuestionList.Questions
//	}
//
//	// 2. 内存过滤
//	var matched []models.Question
//	keyword = strings.ToLower(keyword)
//
//	for _, q := range questions {
//		// 匹配 ID (精确匹配)
//		if q.QuestionFrontendID == keyword {
//			// 如果 ID 完全匹配，直接返回这一个
//			return []models.Question{q}, nil
//		}
//
//		// 匹配 Title 或 Slug (模糊匹配)
//		if strings.Contains(strings.ToLower(q.Title), keyword) ||
//			strings.Contains(strings.ToLower(q.TitleSlug), keyword) {
//			matched = append(matched, q)
//		}
//	}
//
//	return matched, nil
//}
