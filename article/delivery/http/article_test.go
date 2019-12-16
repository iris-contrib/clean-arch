package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	articleHttp "github.com/iris-contrib/clean-arch/article/delivery/http"
	"github.com/iris-contrib/clean-arch/article/mocks"
	"github.com/iris-contrib/clean-arch/models"
)

func TestFetch(t *testing.T) {
	var mockArticle models.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	mockListArticle := make([]*models.Article, 0)
	mockListArticle = append(mockListArticle, &mockArticle)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListArticle, "10", nil)

	app := iris.New()
	req, err := http.NewRequest(iris.MethodGet, "/article?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx := app.ContextPool.Acquire(rec, req)
	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	handler.FetchArticle(ctx)
	app.ContextPool.Release(ctx)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, iris.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.Usecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", models.ErrInternalServerError)

	app := iris.New()
	req, err := http.NewRequest(iris.MethodGet, "/article?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx := app.ContextPool.Acquire(rec, req)
	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	handler.FetchArticle(ctx)
	app.ContextPool.Release(ctx)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, iris.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockArticle models.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockArticle.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(&mockArticle, nil)

	app := iris.New()
	req, err := http.NewRequest(iris.MethodGet, "/article/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx := app.ContextPool.Acquire(rec, req)
	ctx.Params().Set("id", strconv.Itoa(num))
	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	handler.GetByID(ctx)
	app.ContextPool.Release(ctx)

	assert.Equal(t, iris.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockArticle := models.Article{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockArticle := mockArticle
	tempMockArticle.ID = 0
	mockUCase := new(mocks.Usecase)

	j, err := json.Marshal(tempMockArticle)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil)

	app := iris.New()
	req, err := http.NewRequest(iris.MethodPost, "/article", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	rec := httptest.NewRecorder()
	req.URL.Path = "/article"
	ctx := app.ContextPool.Acquire(rec, req)

	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	handler.Store(ctx)
	app.ContextPool.Release(ctx)

	assert.Equal(t, iris.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockArticle models.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockArticle.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	app := iris.New()
	req, err := http.NewRequest(iris.MethodDelete, "/article/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	ctx := app.ContextPool.Acquire(rec, req)
	ctx.Params().Set("id", strconv.Itoa(num))
	handler := articleHttp.ArticleHandler{
		AUsecase: mockUCase,
	}
	handler.Delete(ctx)
	app.ContextPool.Release(ctx)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)
}
