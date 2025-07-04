package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"

	"nebulai-cpu/logger"
	"nebulai-cpu/matrix"

	"encoding/base64"
	"strings"

	"golang.org/x/net/proxy"
)

type UserInfo struct {
	Email       string  `json:"email"`
	FinishPoint float64 `json:"finish_point"`
	UpdatedAt   string  `json:"UpdatedAt"`
}

type SubmitTaskResponse struct {
	CalcStatus bool   `json:"calc_status"`
	Seed1      int64  `json:"seed1"`
	Seed2      int64  `json:"seed2"`
	MatrixSize int    `json:"matrix_size"`
	TaskID     string `json:"task_id"`
}

// 创建支持http/https/socks5代理的http.Client
func NewHttpClient(proxyStr string) *http.Client {
	if proxyStr == "" {
		return &http.Client{}
	}
	u, err := url.Parse(proxyStr)
	if err != nil {
		return &http.Client{}
	}
	if u.Scheme == "socks5" {
		dialer, err := proxy.SOCKS5("tcp", u.Host, nil, proxy.Direct)
		if err != nil {
			logger.LogError("[代理] socks5创建失败: %v", err)
			return &http.Client{}
		}
		tr := &http.Transport{
			Dial: dialer.Dial,
		}
		return &http.Client{Transport: tr}
	} else {
		tr := &http.Transport{
			Proxy: http.ProxyURL(u),
		}
		return &http.Client{Transport: tr}
	}
}

func SubmitTask(result1, result2, taskId, token string) (*SubmitTaskResponse, error) {
	url := "https://nebulai.network/open_compute/finish/task"
	headers := map[string]string{
		"Referer":      "https://nebulai.network/_next/static/chunks/5573.48499b61fe0d9727.js",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Content-Type": "application/json",
		"token":        token,
	}
	data := map[string]string{
		"result_1": result1,
		"result_2": result2,
		"task_id":  taskId,
	}
	jsonData, _ := json.Marshal(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("failed to call api")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Data SubmitTaskResponse `json:"data"`
	}
	logger.LogInfo("[SubmitTask] HTTP Body: %s", string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func QueryUserInfo(token, jwtToken string) (*UserInfo, error) {
	url := "https://nebulai.network/open_compute/get/user_info"
	headers := map[string]string{
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"token":              jwtToken,
		"authorization":      fmt.Sprintf("Bearer %s", token),
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"priority":           "u=1, i",
		"referer":            "https://nebulai.network/opencompute",
		"sec-ch-ua":          `"Chromium";v="136", "Microsoft Edge";v="136", "Not.A/Brand";v="99"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("failed to call api")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Data UserInfo `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func GetComputeToken(token string) (string, error) {
	url := "https://nebulai.network/open_compute/login/token"
	headers := map[string]string{
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"authorization":      fmt.Sprintf("Bearer %s", token),
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"cache-control":      "no-cache",
		"content-type":       "application/json",
		"dnt":                "1",
		"origin":             "https://nebulai.network",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://nebulai.network/opencompute?invite_by=ghHaK5",
		"sec-ch-ua":          `"Not.A/Brand";v="99", "Chromium";v="136"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("failed to call api")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Data struct {
			Jwt string `json:"jwt"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result.Data.Jwt, nil
}

func StartTask(token, jwtToken string) error {
	url := "https://nebulai.network/open_compute/start/task"
	headers := map[string]string{
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"authorization":      fmt.Sprintf("Bearer %s", token),
		"token":              jwtToken,
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"cache-control":      "no-cache",
		"content-type":       "application/json",
		"dnt":                "1",
		"origin":             "https://nebulai.network",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://nebulai.network/opencompute",
		"sec-ch-ua":          `"Not.A/Brand";v="99", "Chromium";v="136"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.LogInfo("[StartTask] HTTP Body: %s", string(body))
	if resp.StatusCode != 200 {
		return errors.New("failed to call api")
	}
	return nil
}

// 多账号任务主循环
func RunAccountTask(token, jwtToken, proxy string, idx int, gpuEnabled bool) {
	client := NewHttpClient(proxy)
	defer func() {
		if r := recover(); r != nil {
			logger.LogError("[账号%d] panic: %v", idx, r)
		}
	}()
	StartTaskWithClient(token, jwtToken, client)
	var result1, result2, taskId string
	count := 0
	for {
		if IsTokenExpired(jwtToken) {
			var err error
			jwtToken, err = GetComputeTokenWithClient(token, client)
			if err != nil {
				logger.LogError("[账号%d] 获取JWT Token失败: %v", idx, err)
				return
			}
		}
		userInfo, err := QueryUserInfoWithClient(token, jwtToken, client)
		if err != nil {
			logger.LogError("[账号%d] 查询用户信息失败: %v", idx, err)
			return
		}
		if IsExpiredOver24Hours(userInfo.UpdatedAt) {
			StartTaskWithClient(token, jwtToken, client)
		}
		logger.LogInfo("[账号%d] 提交任务：%s", idx, taskId)
		data, err := SubmitTaskWithClient(result1, result2, taskId, jwtToken, client)
		if err != nil {
			logger.LogError("[账号%d] 提交任务失败: %v", idx, err)
			return
		}
		logger.LogInfo("[账号%d] 提交任务返回：%s 状态：%v", idx, taskId, data.CalcStatus)
		if data.CalcStatus {
			logger.LogInfo("[账号%d] 任务计算提交成功", idx)
			logger.LogInfo("[账号%d] 开始下一轮计算", idx)
			if count%10 == 0 {
				userInfo, _ := QueryUserInfoWithClient(token, jwtToken, client)
				logger.LogInfo("[账号%d] 账号%s 现在已经挖到了 %6.f NEB", idx, userInfo.Email, userInfo.FinishPoint)
			}
			count++
		}
		res1, res2, err := matrix.AutoCalculateResult(data.Seed1, data.Seed2, data.MatrixSize, gpuEnabled)
		if err != nil {
			logger.LogWarning("[账号%d] 计算失败: %v，已降级为CPU", idx, err)
			res1, res2 = matrix.CalculateResult(data.Seed1, data.Seed2, data.MatrixSize)
		}
		result1 = fmt.Sprintf("%f", res1)
		result2 = fmt.Sprintf("%f", res2)
		taskId = data.TaskID
	}
}

// 下面所有API函数增加WithClient版本，client参数替换原有http.Client创建

func SubmitTaskWithClient(result1, result2, taskId, token string, client *http.Client) (*SubmitTaskResponse, error) {
	url := "https://nebulai.network/open_compute/finish/task"
	headers := map[string]string{
		"Referer":      "https://nebulai.network/_next/static/chunks/5573.48499b61fe0d9727.js",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Content-Type": "application/json",
		"token":        token,
	}
	data := map[string]string{
		"result_1": result1,
		"result_2": result2,
		"task_id":  taskId,
	}
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("failed to call api")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Data SubmitTaskResponse `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func QueryUserInfoWithClient(token, jwtToken string, client *http.Client) (*UserInfo, error) {
	url := "https://nebulai.network/open_compute/get/user_info"
	headers := map[string]string{
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"token":              jwtToken,
		"authorization":      fmt.Sprintf("Bearer %s", token),
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"priority":           "u=1, i",
		"referer":            "https://nebulai.network/opencompute",
		"sec-ch-ua":          `"Chromium";v="136", "Microsoft Edge";v="136", "Not.A/Brand";v="99"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("failed to call api")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Data UserInfo `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func GetComputeTokenWithClient(token string, client *http.Client) (string, error) {
	url := "https://nebulai.network/open_compute/login/token"
	headers := map[string]string{
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"authorization":      fmt.Sprintf("Bearer %s", token),
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"cache-control":      "no-cache",
		"content-type":       "application/json",
		"dnt":                "1",
		"origin":             "https://nebulai.network",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://nebulai.network/opencompute?invite_by=ghHaK5",
		"sec-ch-ua":          `"Not.A/Brand";v="99", "Chromium";v="136"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("failed to call api")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Data struct {
			Jwt string `json:"jwt"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result.Data.Jwt, nil
}

func StartTaskWithClient(token, jwtToken string, client *http.Client) error {
	url := "https://nebulai.network/open_compute/start/task"
	headers := map[string]string{
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"authorization":      fmt.Sprintf("Bearer %s", token),
		"token":              jwtToken,
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"cache-control":      "no-cache",
		"content-type":       "application/json",
		"dnt":                "1",
		"origin":             "https://nebulai.network",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://nebulai.network/opencompute",
		"sec-ch-ua":          `"Not.A/Brand";v="99", "Chromium";v="136"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.LogInfo("[StartTask] HTTP Body: %s", string(body))
	if resp.StatusCode != 200 {
		return errors.New("failed to call api")
	}
	return nil
}

// 判断token是否过期
func IsTokenExpired(token string) bool {
	if token == "" {
		return true
	}
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return true
	}
	payload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return true
	}
	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return true
	}
	if exp, ok := claims["exp"].(float64); ok {
		return int64(exp)*1000 < time.Now().UnixMilli()
	}
	return false
}

// 判断时间是否超过24小时
func IsExpiredOver24Hours(isoDateStr string) bool {
	t, err := time.Parse(time.RFC3339, isoDateStr)
	if err != nil {
		return false
	}
	diff := time.Since(t)
	return diff.Hours() > 24
}

// 工具函数：四舍五入到n位小数
func round(f float64, n int) float64 {
	pow := math.Pow(10, float64(n))
	return math.Round(f*pow) / pow
}

// 保证文件以正确的Go语法结尾
