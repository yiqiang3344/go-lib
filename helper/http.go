package helper

import (
	"bytes"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
func PostJson(ctx context.Context, url string, data interface{}) (bool, int, string) {
	resp, err := httpRequest(
		ctx,
		"POST",
		url,
		data,
		"application/json",
		5*time.Second,
	)
	var statusCode int
	result := ""
	if resp != nil {
		statusCode = resp.StatusCode
		_result, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		result = string(_result)
	} else {
		statusCode = 500
		if err != nil {
			result = err.Error()
		}
	}
	return true, statusCode, result
}

func httpRequest(ctx context.Context, method string, url string, data interface{}, contentType string, timeout time.Duration) (*http.Response, error) {
	span := NewInnerSpan(RunFuncName(), ctx)
	if span != nil {
		defer span.Finish()
	}

	c := &http.Client{Timeout: timeout}
	var res *http.Response
	var req *http.Request
	var err error
	if method == "POST" {
		jsonStr, _ := json.Marshal(data)
		req, err = http.NewRequest(
			method,
			url,
			bytes.NewBuffer(jsonStr),
		)
		if err != nil {
			AddSpanError(span, err)
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
	} else {
		req, err = http.NewRequest(
			method,
			url,
			nil,
		)
		if err != nil {
			AddSpanError(span, err)
			return nil, err
		}
	}
	if span != nil {
		err = span.Tracer().Inject(span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
		if err != nil {
			FatalLog("Could not inject span context into header:"+err.Error(), "")
		}
	}
	startTime := time.Now()
	res, err = c.Do(req)
	if err != nil {
		AddSpanError(span, err)
		addErrorLog(startTime, time.Now(), req, data, err)
		return nil, err
	}

	addSuccessLog(startTime, time.Now(), res, data)
	return res, nil
}

func addLog(startTime time.Time, endTime time.Time, url string, header interface{}, data interface{}, statusCode int, result string) {
	request := &SrvRequest{
		Time:   startTime,
		Url:    url,
		Header: header,
		Body:   data,
	}
	reponse := &SrvResponse{
		Time:       endTime,
		StatusCode: statusCode,
		Data:       string(result),
	}
	WebClientLog(
		zap.Object("request", request),
		zap.Object("response", reponse),
		zap.Float64("response_time", float64(endTime.Sub(startTime).Microseconds())/1e6),
	)
}

func addErrorLog(startTime time.Time, endTime time.Time, req *http.Request, data interface{}, err error) {
	var statusCode int
	result := ""
	statusCode = 500
	result = err.Error()
	addLog(startTime, endTime, req.URL.String(), req.Header, data, statusCode, result)
}

func addSuccessLog(startTime time.Time, endTime time.Time, rep *http.Response, data interface{}) {
	var statusCode int
	result := ""
	statusCode = rep.StatusCode
	_result, _ := ioutil.ReadAll(rep.Body)
	defer rep.Body.Close()
	result = string(_result)
	addLog(startTime, endTime, rep.Request.URL.String(), rep.Request.Header, data, statusCode, result)
}
