package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/X-for/ltgo/internal/models"
)

// SubmitResponse 提交任务的响应 (和 InterpretResponse 几乎一样)
type SubmitResponse struct {
	SubmissionID int64 `json:"submission_id"`
}

// SubmitCheckResponse 提交结果的详情
type SubmitCheckResponse struct {
	State            string `json:"state"`       // SUCCESS
	StatusMsg        string `json:"status_msg"`  // Accepted, Wrong Answer
	StatusCode       int    `json:"status_code"` // 10=Accepted, 11=WA
	CompileError     string `json:"compile_error"`
	FullCompileError string `json:"full_compile_error"`
	RuntimeError     string `json:"runtime_error"`

	TotalCorrect   int `json:"total_correct"`
	TotalTestcases int `json:"total_testcases"`

	// 性能数据
	StatusRuntime     string  `json:"status_runtime"`     // "4 ms"
	RuntimePercentile float64 `json:"runtime_percentile"` // 击败比例
	StatusMemory      string  `json:"status_memory"`      // "5 MB"
	MemoryPercentile  float64 `json:"memory_percentile"`

	// 错误时的详细信息
	InputFormatted string `json:"input_formatted"` // 出错的那个 case 的输入
	CodeOutput     string `json:"code_output"`     // 你的输出
	ExpectedOutput string `json:"expected_output"` // 预期输出
	StdOutput      string `json:"std_output"`      // 你的打印
}

// SubmitCode 提交代码进行判题
func (c *Client) SubmitCode(q *models.QuestionDetail, code string, lang string) (int64, error) {
	payload := map[string]interface{}{
		"lang":        lang,
		"question_id": q.QuestionID, // ⚠️ 注意：提交通常需要 QuestionID (后端ID)，不是 FrontendID
		"typed_code":  code,
	}

	body, _ := json.Marshal(payload)
	path := fmt.Sprintf("/problems/%s/submit/", q.TitleSlug)

	respBody, err := c.Post(path, body)
	if err != nil {
		return 0, err
	}

	var resp SubmitResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return 0, err
	}

	return resp.SubmissionID, nil
}

// CheckSubmission 轮询提交结果
func (c *Client) CheckSubmission(submissionID int64) (*SubmitCheckResponse, error) {
	path := fmt.Sprintf("/submissions/detail/%d/check/", submissionID)

	for i := 0; i < 20; i++ {
		respBody, err := c.Get(path)
		if err != nil {
			return nil, err
		}

		var resp SubmitCheckResponse
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, err
		}

		if resp.State == "SUCCESS" {
			return &resp, nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("timeout waiting for submission result")
}
