package client

import (
	"errors"
	"strings"

	"github.com/X-for/ltgo/internal/models"
)

type SearchOptions struct {
	Keyword    string
	Difficulty string // "EASY", "MEDIUM", "HARD"
	Status     string // "TO_DO", "SOLVED", "ATTEMPTED"
	Tag        string // e.g. "array", "dynamic-programming"
	FrontendID string // id of problem
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

// SearchQuestions 严格复刻抓包请求
func (c *Client) SearchQuestions(opts SearchOptions) ([]models.Question, error) {
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
	// 构造 Filters
	filters := map[string]interface{}{
		"filterCombineType":   "ALL",
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
		// 下面是动态填充的部分
		"difficultyFilter": map[string]interface{}{
			"difficulties": []string{},
			"operator":     "IS",
		},
		"statusFilter": map[string]interface{}{
			"questionStatuses": []string{},
			"operator":         "IS",
		},
		"languageFilter": map[string]interface{}{
			"languageSlugs": []string{},
			"operator":      "IS",
		},
		"topicFilter": map[string]interface{}{
			"topicSlugs": []string{},
			"operator":   "IS",
		},
	}

	// 填充动态过滤条件
	if opts.Difficulty != "" {
		filters["difficultyFilter"].(map[string]interface{})["difficulties"] = []string{strings.ToUpper(opts.Difficulty)}
	}
	if opts.Status != "" {
		filters["statusFilter"].(map[string]interface{})["questionStatuses"] = []string{strings.ToUpper(opts.Status)}
	}
	if opts.Tag != "" {
		filters["topicFilter"].(map[string]interface{})["topicSlugs"] = []string{opts.Tag}
	}
	if opts.Keyword == "" {
		opts.Keyword = opts.FrontendID
	}

	vars := map[string]interface{}{
		"skip":          0,
		"limit":         20,
		"categorySlug":  category,
		"searchKeyword": opts.Keyword,
		"sortBy": map[string]interface{}{
			"sortField": "CUSTOM",
			"sortOrder": "ASCENDING",
		},
		"filters": filters,
	}

	var resp models.QuestionListResponse
	if err := c.GraphQL(query, vars, &resp); err != nil {
		return nil, err
	}

	questions := resp.Data.ProblemsetQuestionListV2.Questions
	if len(questions) == 0 {
		questions = resp.Data.ProblemsetQuestionList.Questions
	}
	// [新增] 客户端精确过滤 ID
	if opts.FrontendID != "" {
		var exactMatch []models.Question
		for _, q := range questions {
			if q.QuestionFrontendID == opts.FrontendID {
				exactMatch = append(exactMatch, q)
				break // 找到一个就够了，ID 是唯一的
			}
		}
		// 如果找到了，就只返回这一条
		// 如果没找到（可能是 filters 已经过滤太狠了，或者是翻页问题），那就返回空，或者返回原始列表（取决于策略）
		// 这里我们选择：如果找到了就精确返回；没找到就返回空（因为用户明确要求了 ID）
		return exactMatch, nil
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
