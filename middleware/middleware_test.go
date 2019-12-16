package middleware_test

import (
	test "net/http/httptest"
	"testing"

	"github.com/iris-contrib/clean-arch/middleware"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
)

func TestIrisCORSMiddleware(t *testing.T) {
	app := iris.New()
	req := test.NewRequest(iris.MethodGet, "/", nil)
	rec := test.NewRecorder()
	ctx := app.ContextPool.Acquire(rec, req)

	h := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // allows everything, use that to change the hosts.
	})

	h(ctx)
	app.ContextPool.Release(ctx)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware(t *testing.T) {
	app := iris.New()
	req := test.NewRequest(iris.MethodGet, "/", nil)
	rec := test.NewRecorder()
	ctx := app.ContextPool.Acquire(rec, req)

	m := middleware.InitMiddleware()
	ctx.Do([]iris.Handler{m.CORS, func(ctx iris.Context) {
		ctx.Header("X-Header", "OK")
	}})
	app.ContextPool.Release(ctx)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "OK", rec.Header().Get("X-Header"))
}
