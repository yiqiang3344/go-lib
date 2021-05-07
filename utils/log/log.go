package cLog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2/server"
	"github.com/opentracing/opentracing-go"
	cTime "github.com/yiqiang3344/go-lib/utils/time"
	"github.com/yiqiang3344/go-lib/utils/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SrvRequest struct {
	Time   time.Time
	Url    string
	Header interface{}
	Body   interface{}
}

type SrvResponse struct {
	Time       time.Time
	StatusCode int
	Data       string
}

func (u *SrvRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("time", cTime.MilliDatetimeByTime(u.Time))
	enc.AddString("url", u.Url)
	headJson, _ := json.Marshal(u.Header)
	enc.AddString("header", string(headJson))
	bodyJson, _ := json.Marshal(u.Body)
	enc.AddString("body", string(bodyJson))
	return nil
}

func (u *SrvResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("time", cTime.MilliDatetimeByTime(u.Time))
	enc.AddInt("status_code", u.StatusCode)
	enc.AddString("data", u.Data)
	return nil
}

var Logger *zap.Logger
var RequestID string

func initRequestId(ctx context.Context) {
	parent := opentracing.SpanFromContext(ctx)
	if parent == nil {
		RequestID = uuid.GenUuId()
	} else {
		RequestID = strings.Split(fmt.Sprint(parent.Context()), ":")[0]
	}
}

func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(cTime.MilliDatetimeByTime(t))
}

func InitLogger(project string) {
	hook := lumberjack.Logger{
		Filename:   "log/app.log", // 日志文件路径
		MaxSize:    512,           // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 5,             // 日志文件最多保存多少个备份
		MaxAge:     30,            // 文件最多保存多少天
		Compress:   false,         // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "at",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     encodeTime,                     // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	//caller := zap.AddCaller()
	// 开启文件及行号
	//development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("project", project))
	// 构造日志
	//Logger = zap.New(core, caller, development, filed)
	Logger = zap.New(core, filed)
}

func debugTrace() zap.Field {
	_, file, line, _ := runtime.Caller(2)
	return zap.String("at", file+":"+strconv.Itoa(line))
}

func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		initRequestId(ctx)

		startTime := time.Now()
		request := &SrvRequest{
			Time:   startTime,
			Url:    req.Endpoint(),
			Header: req.Header(),
			Body:   req.Body(),
		}

		ret := fn(ctx, req, rsp)
		endTime := time.Now()

		rspStr, _ := json.Marshal(rsp)
		reponse := &SrvResponse{
			Time:       endTime,
			StatusCode: 200,
			Data:       string(rspStr),
		}

		AccessLog(
			zap.Object("request", request),
			zap.Object("response", reponse),
			zap.Float64("response_time", float64(endTime.Sub(startTime).Microseconds())/1e6),
		)
		return ret
	}
}

func AccessLog(args ...zap.Field) {
	args = append(args,
		debugTrace(),
		zap.String("category", "access"),
		zap.String("request_id", RequestID),
	)
	Logger.Info("", args...)
}

func WebClientLog(args ...zap.Field) {
	args = append(args,
		debugTrace(),
		zap.String("category", "web_client"),
		zap.String("request_id", RequestID),
	)
	Logger.Info("", args...)
}

func DebugLog(message, messageTag string, args ...zap.Field) {
	args = append(args,
		debugTrace(),
		zap.String("message_tag", messageTag),
		zap.String("category", "debug"),
		zap.String("request_id", RequestID),
	)
	Logger.Debug(message, args...)
}

func BizLog(message, messageTag string, args ...zap.Field) {
	args = append(args,
		debugTrace(),
		zap.String("message_tag", messageTag),
		zap.String("category", "biz"),
		zap.String("request_id", RequestID),
	)
	Logger.Info(message, args...)
}

func ErrorLog(message, messageTag string, args ...zap.Field) {
	args = append(args,
		debugTrace(),
		zap.String("message_tag", messageTag),
		zap.String("category", "error"),
		zap.String("request_id", RequestID),
	)
	Logger.Error(message, args...)
}

func FatalLog(message, messageTag string, args ...zap.Field) {
	args = append(args,
		debugTrace(),
		zap.String("message_tag", messageTag),
		zap.String("category", "fatal"),
		zap.String("request_id", RequestID),
	)
	Logger.Fatal(message, args...)
}
