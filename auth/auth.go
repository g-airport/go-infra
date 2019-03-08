package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	gCtx "github.com/onlyLeoll/go-infra/context"
	iEnv "github.com/onlyLeoll/go-infra/env"
	iErr "github.com/onlyLeoll/go-infra/errors"
	iLog "github.com/onlyLeoll/go-infra/log"

	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

var (
	accessLog *iLog.Logger
)

func init() {
	path, _ := filepath.Abs(fmt.Sprintf("%s%saccess.log",
		iEnv.LogDir, string(filepath.Separator)))
	accessLog, _ = iLog.NewLogger(path, "info")
}

type User struct {
	ID     int32
	Rights uint64
}

type Auth struct {
	Rights                  uint64
	LogRequest, LogResponse bool
}

func UserAuth(ctx context.Context) (*User, error) {
	token := GetToken(ctx)
	if token == "" {
		return nil, ErrNoToken
	}
	return GetUserByToken(token)
}

func GetToken(ctx context.Context) string {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return ""
	}
	return md["Token"]
}

func GetUserByToken(token string) (*User, error) {
	return nil, nil
}

// user - server
func Check(uRights, sRights uint64) error {
	if uRights&sRights == RightsNone {
		return ErrNoPermission
	}
	return nil
}

func ServerAuth(ctx context.Context, rights uint64) (*User, error) {
	if rights == RightsNone {
		return nil, nil
	}

	u, err := UserAuth(ctx)
	if err != nil {
		return nil, err
	}

	if err := Check(u.Rights, rights); err != nil {
		return u, err
	}

	return u, err
}

func ServerAuthWrapper(rights map[string]Auth) server.HandlerWrapper {
	serviceName := filepath.Base(os.Args[0])
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {

			method := req.Method()
			r, ok := rights[method]
			if !ok {
				accessLog.Warnw("no method found", "method", method)
				return ErrNoMethod
			}

			u, err := ServerAuth(ctx, r.Rights)
			if err != nil {
				accessLog.Warn("try to call auth failed",
					"user", u,
					"method", method,
					"err", err,
				)
				return err
			}

			ctx = gCtx.SetUser(ctx, *u)

			//todo
			addrText := ""

			start := time.Now()
			err = fn(ctx, req, rsp)
			end := time.Now()

			code := 200
			if er := iErr.ParseRPCError(err); er != nil {
				code = er.Code
			}

			methodName := method
			if sp := strings.Split(methodName, "."); len(sp) > 1 {
				methodName = sp[1]
			}

			kv := []interface{}{
				"user", u,
				"method", method,
				"cost", end.Sub(start).String(),
				"remote_addr", addrText,
				"time", time.Now(),
				"status", code,
				"elapsed_time", end.Sub(start).Seconds() * 1000,
				"source_srv", serviceName,
				"source_ip", addrText,
				"interface_name", methodName,
			}

			if r.LogRequest {
				reqRaw := req.Body()
				kv = append(kv, "request_content", &reqRaw)
			}

			if r.LogResponse {
				kv = append(kv, "response_content", &rsp)
			}

			kv = append(kv, "err", err)
			accessLog.Infow("called", kv...)

			return err
		}
	}
}
