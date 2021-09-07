package api

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danielgatis/go-between/internal/raft"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
)

func checkRaftLeader(raftNode *raft.Node) echo.MiddlewareFunc {
	forwardToLeader := func(c echo.Context) error {
		r := c.Request()
		r.URL.Host = string(raftNode.Leader())
		r.URL.Scheme = "http"
		r.RequestURI = ""
		client := &http.Client{}

		r.Header.Set("X-Forwarded-For", strings.Join([]string{r.Form.Get("X-Forwarded-For"), r.RemoteAddr}, ","))

		resp, err := client.Do(r.WithContext(context.Background()))
		if err != nil {
			return eris.Wrap(err, "client.Do")
		}

		w := c.Response()
		h := w.Header()

		for k, vv := range resp.Header {
			for _, v := range vv {
				h.Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)

		if _, err := io.Copy(w, resp.Body); err != nil {
			return eris.Wrap(err, "io.Copy")
		}

		resp.Body.Close()
		return nil
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if raftNode.Leader() == "" {
				return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "cluster not bootstrapped"})
			}

			if !raftNode.IsLeader() {
				if err := forwardToLeader(c); err != nil {
					return eris.Wrap(err, "forwardToLeader")
				}

				return nil
			}

			return next(c)
		}
	}
}

func requestLogger(log logrus.FieldLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now().UTC()

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now().UTC()

			p := req.URL.Path

			bytesIn := req.Header.Get(echo.HeaderContentLength)

			fields := (map[string]interface{}{
				"time_rfc3339":  time.Now().UTC().Format(time.RFC3339),
				"remote_ip":     c.RealIP(),
				"host":          req.Host,
				"uri":           req.RequestURI,
				"method":        req.Method,
				"path":          p,
				"referer":       req.Referer(),
				"user_agent":    req.UserAgent(),
				"status":        res.Status,
				"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
				"latency_human": stop.Sub(start).String(),
				"bytes_in":      bytesIn,
				"bytes_out":     strconv.FormatInt(res.Size, 10),
			})

			if err != nil {
				fields["error"] = err.Error()
			}

			log.WithFields(fields).Info("Handled request")

			return nil
		}
	}
}
