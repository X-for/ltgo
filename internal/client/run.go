package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/X-for/ltgo/internal/models"
)

type InterpretResponse struct {
	InterpretID string `json:"interpret_id"`
}

type CheckResponse struct {
	State            string `json:"state"`
	StatusCode       int    `json:"status_code"` // 10:CE, 11:WA, 12:RE, 13:TLE, 14:MLE, 15:OK?
	StatusMsg        string `json:"status_msg"`
	CompileError     string `json:"compile_error"`
	FullCompileError string `json:"full_compile_error"`
	RuntimeError     string `json:"runtime_error"`

	CodeAnswer     []string `json:"code_answer"`          // 你的返回值
	ExpectedOutput []string `json:"expected_code_answer"` // 预期返回值
	InputFormatted []string `json:"input_formatted"`      // 输入参数 (非常重要!)
	StdOutput      []string `json:"std_output_list"`      // 你的 fmt.Println 输出

	CorrectAnswer bool `json:"correct_answer"`

	ElapsedTime    int `json:"elapsed_time"`
	TotalCorrect   int `json:"total_correct"`
	TotalTestcases int `json:"total_testcases"`
}

// RunCode 提交运行任务
func (c *Client) RunCode(q *models.QuestionDetail, code string, lang string) (string, error) {
	// 1. 构造请求 Payload
	payload := map[string]interface{}{
		"lang":        lang,
		"question_id": q.QuestionFrontendID, // 注意：有些时候这里需要 QuestionID (后端ID)，而非 FrontendID
		"typed_code":  code,
		"data_input":  q.SampleTestCase,
	}

	body, _ := json.Marshal(payload)
	path := fmt.Sprintf("/problems/%s/interpret_solution/", q.TitleSlug)

	respBody, err := c.Post(path, body)
	if err != nil {
		return "", err
	}

	//fmt.Println("DEBUG CHECK:", string(respBody))

	var resp InterpretResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", err
	}

	return resp.InterpretID, nil
}

// CheckResult 轮询结果
func (c *Client) CheckResult(interpretID string) (*CheckResponse, error) {
	path := fmt.Sprintf("/submissions/detail/%s/check/", interpretID)

	// 轮询几次，每次间隔 1-2 秒
	for i := 0; i < 20; i++ {
		respBody, err := c.Get(path)
		if err != nil {
			return nil, err
		}
		//fmt.Println("DEBUG CHECK:", string(respBody))

		var resp CheckResponse
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, err
		}

		if resp.State == "SUCCESS" {
			return &resp, nil
		}

		// 还没跑完，等一下
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("timeout waiting for result")
}
