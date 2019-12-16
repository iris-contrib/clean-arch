package http

import (
	"context"

	"github.com/kataras/iris/v12"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/iris-contrib/clean-arch/article"
	"github.com/iris-contrib/clean-arch/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type ArticleHandler struct {
	AUsecase article.Usecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewArticleHandler(r iris.Party, us article.Usecase) {
	handler := &ArticleHandler{
		AUsecase: us,
	}
	r.Get("/articles", handler.FetchArticle)
	r.Post("/articles", handler.Store)
	r.Get("/articles/{id:int64}", handler.GetByID)
	r.Delete("/articles/{id:int64}", handler.Delete)
}

// FetchArticle will fetch the article based on given params
func (a *ArticleHandler) FetchArticle(ctx iris.Context) {
	num := ctx.URLParamInt64Default("num", 1)
	cursor := ctx.URLParam("cursor")
	c := ctx.Request().Context()
	if c == nil {
		c = context.Background()
	}
	listAr, nextCursor, err := a.AUsecase.Fetch(c, cursor, num)

	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Header("X-Cursor", nextCursor)
	ctx.JSON(listAr)
}

// GetByID will get article by given id
func (a *ArticleHandler) GetByID(ctx iris.Context) {
	// always found, because of {id:int64}, otherwise router would already throw 404.
	id, _ := ctx.Params().GetInt64("id")

	c := ctx.Request().Context()
	if c == nil {
		c = context.Background()
	}

	art, err := a.AUsecase.GetByID(c, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(art)
}

func isRequestValid(m *models.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the article by given request body
func (a *ArticleHandler) Store(ctx iris.Context) {
	var article models.Article

	err := ctx.ReadJSON(&article)
	if err != nil {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(err.Error())
		return
	}

	if ok, err := isRequestValid(&article); !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(err.Error())
		return
	}

	c := ctx.Request().Context()
	if c == nil {
		c = context.Background()
	}

	err = a.AUsecase.Store(c, &article)

	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(article)
}

// Delete will delete article by given param
func (a *ArticleHandler) Delete(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")

	c := ctx.Request().Context()
	if c == nil {
		c = context.Background()
	}

	err := a.AUsecase.Delete(c, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}

func handleError(ctx iris.Context, err error) {
	if err == nil {
		return
	}

	statusCode := iris.StatusInternalServerError
	switch err {
	case models.ErrInternalServerError:
		statusCode = iris.StatusInternalServerError
	case models.ErrNotFound:
		statusCode = iris.StatusNotFound
	case models.ErrConflict:
		statusCode = iris.StatusConflict
	}

	ctx.Application().Logger().Error(err)

	ctx.StatusCode(statusCode)
	ctx.JSON(ResponseError{Message: err.Error()})
}
