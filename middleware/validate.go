package middleware

import (
	"fmt"
	"reflect"
	"strings"
	"net/http"

	"encoding/json"

	"go.opentelemetry.io/otel"

	"github.com/go-playground/validator/v10"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	locale string
	trans ut.Translator
)

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
	    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	    if name == "-" {
	        return ""
	    }

	    return name
	})


	locale = "en"

	enTranslator := en.New()
	uni := ut.New(enTranslator, enTranslator)

	trans, _ = uni.GetTranslator(locale)
	en_translations.RegisterDefaultTranslations(validate, trans)
}

/*
Codes we care about

200 OK,
400 BadRequest
401 Unauthorized - Authentication denied,
403 Forbidden - Authorization denined,
404 Not Found
500 Internal Server Error,
503 Service Unavailable
*/

type FieldError struct {
	Path string `json:"path"`
	Err string `json:"err"`
}

type HttpClientErrorResponse struct {
	StatusCode int `json:"status_code"`
	Method string `json:"method"`
	Url string `json:"url"`
	Errors []FieldError `json:"errors"`
}

func WithRequestValidation(request interface{}) MiddlewareHandler{
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			tr := otel.Tracer("request")
			ctx, span := tr.Start(ctx, "request-validation")
			defer span.End()

	    err := validate.Struct(request)
			if err == nil {
				// No validation error, continue chain
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			response := HttpClientErrorResponse{
				StatusCode: http.StatusBadRequest,
				Method: r.Method,
				Url: fmt.Sprintf("%s%s", r.Host, r.URL.RequestURI()),
				// Body: ...
			}

			for _, verr := range err.(validator.ValidationErrors) {
				e := FieldError{
					Path: verr.Field(),
					Err: verr.Translate(trans),
				}

				response.Errors = append(response.Errors, e)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.StatusCode)

			d, err := json.Marshal(response)
			if err != nil {
				panic(err) // TODO FIXME
			}
			w.Write(d)

		})
	}
}

func WithResponseValidation(response interface{}) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO

		/*

			ctx := r.Context()
			tr := otel.Tracer("request")
			ctx, span := tr.Start(ctx, "response-validation")
			defer span.End()

	    err := validate.Struct(response)
			if err == nil {
				// No validation error, continue chain
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			response := HttpClientErrorResponse{
				StatusCode: http.StatusBadRequest,
				Method: r.Method,
				Url: fmt.Sprintf("%s%s", r.Host, r.URL.RequestURI()),
				// Body: ...
			}

			for _, verr := range err.(validator.ValidationErrors) {
				e := FieldError{
					Path: verr.Field(),
					Err: verr.Translate(trans),
				}

				response.Errors = append(response.Errors, e)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.StatusCode)

			d, err := json.Marshal(response)
			if err != nil {
				panic(err) // TODO FIXME
			}
			w.Write(d)
		*/
		})
	}
}





// TODO move this to better place
func WithRequestParser(request interface{}) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tr := otel.Tracer("request")
			ctx, span := tr.Start(ctx, "middleware.request-parser")
			defer span.End()

			next.ServeHTTP(w, r)
		})
	}
}
