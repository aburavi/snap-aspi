package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	jwtoken "github.com/aburavi/snaputils/jwtoken"
	"github.com/aburavi/snaputils/utils"
	"github.com/spf13/viper"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting      = errors.New("inconsistent mapping between route and handler (programmer error)")
	ErrUnknown         = errors.New("unknown user error")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("grpc: item not found")
	ErrConflict        = errors.New("grpc: item conflict")
	ErrUnAuthorized    = errors.New("Unauthorized, Authentication required")
	ErrForbidden       = errors.New("Forbidden User")
)

func checkAuthorization(token string, logger log.Logger) (jwt.MapClaims, error) {
	//pubKey := os.Getenv("KEYCLOAK_PUBLICKEY")
	publicKey := viper.GetString("PUBLICKEY")
	level.Info(logger).Log("token", fmt.Sprintf("%s", token))
	jwtToken := jwtoken.NewJWT([]byte(publicKey), logger)

	content, err := jwtToken.Validate(token)
	if err != nil {
		level.Info(logger).Log("parse token failed", fmt.Sprintf("%s", err.Error()))
		return nil, err
	}
	//djwt := content.(map[string]interface{})
	return content.(jwt.MapClaims), nil
}

func decodeNoTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	req := map[string]interface{}{}
	header := map[string]string{}

	vals := r.Body
	if vals == nil {
		return nil, errors.New("Unknown/Empty Body,...")
	}
	defer r.Body.Close()
	body, err := io.ReadAll(vals)
	if err != nil {
		return nil, err
	}

	for name, values := range r.Header {
		for _, value := range values {
			header[name] = value
		}
	}

	req["header"] = header
	req["body"] = body
	req["path"] = r.URL.Path

	level.Info(logger).Log("decode payload", fmt.Sprintf("%v", req["body"].(byte)))
	//level.Info(logger).Log("decode params", fmt.Sprintf("%v", req["params"].(map[string][]byte)))
	return req, nil
}

func errorEncoder(_ context.Context, failed error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	err, ok := failed.(*utils.Error)
	if ok {
		switch err.Code {
		case utils.Code(utils.NotExist.String()):
			statusCode = http.StatusNotFound
		default:
		}
	} else {
		err = &utils.Error{Code: utils.Code(utils.Internal.String()), Message: failed.Error()}
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorWrapper{err})
}

type errorWrapper struct {
	Error error `json:"error"`
}

func ProcessMapString(key string, data map[string]interface{}) string {
	result, ok := data[key].(string)
	if ok {
		return result
	}
	return ""
}
