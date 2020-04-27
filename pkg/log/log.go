package log

import (
	"encoding/json"
	"fmt"
	"github.com/Abramovic/logrus_influxdb"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"go-shop-v2/app/email"
	"math"
	"os"
	"time"
)

const format = "2006-01-02 15:04:05"

type Option struct {
	AppName        string
	Path           string
	MaxAge         time.Duration
	RotationTime   time.Duration
	Email          email.Mailer
	To             string
	InfluxDBConfig *logrus_influxdb.Config
}

func isProd() bool {
	return os.Getenv(gin.EnvGinMode) == "release"
}

func Setup(opt Option) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: format,
	})
	logrus.SetReportCaller(true)
	hook := logHook{}
	logrus.AddHook(hook)
	if isProd() {
		writer, _ := rotatelogs.New(
			opt.Path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(opt.Path),
			rotatelogs.WithMaxAge(opt.MaxAge),
			rotatelogs.WithRotationTime(opt.RotationTime),
		)
		logrus.SetOutput(writer)
	}

	if opt.Email != nil {
		logrus.AddHook(&emailHook{opt})
	}

	if opt.InfluxDBConfig != nil {
		hook, err := logrus_influxdb.NewInfluxDB(opt.InfluxDBConfig)
		if err == nil {
			logrus.AddHook(hook)
		} else {
			logrus.Error(err)
		}
	}
}

type emailHook struct {
	opt Option
}
type logHook struct {
}

func (e logHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

func (e logHook) Fire(entry *logrus.Entry) error {
	var fileInfo []byte
	if entry.Caller != nil {
		fileInfo, _ = json.MarshalIndent(entry.Caller, "", "\t")
		entry.WithFields(logrus.Fields{
			"fileInfo": string(fileInfo),
		})
	}
	return nil
}

func (e emailHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

func (e emailHook) Fire(entry *logrus.Entry) error {
	if e.opt.Email != nil && e.opt.To != "" {
		subject := e.opt.AppName + " - " + entry.Level.String()
		fields, _ := json.MarshalIndent(entry.Data, "", "\t")
		var fileInfo []byte
		if entry.Caller != nil {
			fileInfo, _ = json.MarshalIndent(entry.Caller, "", "\t")
		}
		contents := fmt.Sprintf(`<ul>
				<li>错误信息：%s</li>
				<li>附加字段：%s</li>
				<li>文件信息：%s</li>
			</ul>`, entry.Message, fields, fileInfo)

		return e.opt.Email.Send(e.opt.To, subject, contents)
	}
	return nil
}

func Logger() gin.HandlerFunc {
	if !isProd() {
		return gin.Logger()
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "N/A"
	}
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logrus.WithFields(logrus.Fields{
			"hostname":    hostname,
			"statusCode":  statusCode,
			"latency":     latency, // time to process
			"clientIP":    clientIP,
			"method":      c.Request.Method,
			"path":        path,
			"referer":     referer,
			"dataLength":  dataLength,
			"userAgent":   clientUserAgent,
			"measurement": "access_log",
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format(format), c.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, latency)
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
