package servant

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"net/http"
	"strings"
	"time"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
)

type TarsHttpMux struct {
	http.ServeMux
	servant   string
	beforefun tarsHTTPBeforeFunc
}

type tarsHttpMuxOption func(*TarsHttpMux)
type tarsHTTPBeforeFunc func(w http.ResponseWriter, r *http.Request) error

func NewTarsHttpMux(opts ...tarsHttpMuxOption) *TarsHttpMux {
	mux := &TarsHttpMux{}
	for _, opt := range opts {
		opt(mux)
	}
	return mux
}

func BeforeFunc(bf tarsHTTPBeforeFunc) tarsHttpMuxOption {
	return func(mux *TarsHttpMux) {
		mux.beforefun = bf
	}
}

func (mux *TarsHttpMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if mux.beforefun != nil {
		if err := mux.beforefun(w, r); err != nil {
			appzaplog.Error("before process failed", zap.Error(err))
			return
		}
	}
	h, pattern := mux.Handler(r)
	tw := &TarsResponseWriter{w, 0}
	startTime := time.Now().UnixNano() / 1000000
	h.ServeHTTP(tw, r)
	costTime := int64(time.Now().UnixNano()/1000000 - startTime)
	reqAddr := r.Header.Get("X-Forwarded-For-Pound")
	if reqAddr == "" {
		reqAddr = strings.SplitN(r.RemoteAddr, ":", 2)[0]
	}

	st := &httpStatInfo{
		reqAddr:    reqAddr,
		pattern:    pattern,
		statusCode: tw.statusCode,
		costTime:   costTime,
	}
	reportHttpStat(st)
}

type TarsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *TarsResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
