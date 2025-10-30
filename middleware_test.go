package headerctx

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func EchoDummy(c echo.Context) error {
	return c.String(200, "OK")
}

func EchoContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func Test_NotRequire(t *testing.T) {
	injectHeaders := InjectHeaders(false, "header-a", "header-b")(EchoDummy)

	c := EchoContext()
	c.Request().Header.Set("header-a", "value-a")

	if err := injectHeaders(c); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// check echo.Context
	if val := c.Get("header-a"); val != "value-a" {
		t.Errorf("expected value-a, got %s", val)
	}

	if val := c.Get("header-b"); val != nil {
		t.Errorf("expected nil, got %s", val)
	}

	// check context.Context
	ctx := c.Request().Context()

	if val := ctx.Value(HeaderCtxKey("header-a")); val != "value-a" {
		t.Errorf("expected value-a, got %s", val)
	}

	if val := ctx.Value(HeaderCtxKey("header-b")); val != nil {
		t.Errorf("expected nil, got %s", val)
	}
}

func Test_Require(t *testing.T) {
	injectHeaders := InjectHeaders(true, "header-a", "header-b")(EchoDummy)

	c := EchoContext()
	c.Request().Header.Set("header-a", "value-a")

	// missing header (fail)
	if err := injectHeaders(c); err == nil {
		t.Fatalf("expected err, got nil")
	}

	c.Request().Header.Set("header-b", "value-b")

	// both headers (pass)
	if err := injectHeaders(c); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// check echo.Context
	if val := c.Get("header-a"); val != "value-a" {
		t.Errorf("expected value-a, got %s", val)
	}

	if val := c.Get("header-b"); val != "value-b" {
		t.Errorf("expected value-b, got %s", val)
	}

	// check context.Context
	ctx := c.Request().Context()

	if val := ctx.Value(HeaderCtxKey("header-a")); val != "value-a" {
		t.Errorf("expected value-a, got %s", val)
	}

	if val := ctx.Value(HeaderCtxKey("header-b")); val != "value-b" {
		t.Errorf("expected value-b, got %s", val)
	}
}
