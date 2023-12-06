package login

import (
	"context"
	"io"
	"main/internal/views/dto"
	"testing"

	"github.com/PuerkitoBio/goquery"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestViewLoginForm_Get(t *testing.T) {
	// Pipe the rendered template into goquery.
	r, w := io.Pipe()

	go func() {
		_ = LoginForm(LoginFormData{Defaults: dto.UserLoginDTO{}}).Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	assert.Equal(t, 0, doc.Find(".text-error").Length(), "Found .text-error when we're not supposed to")

	// get the form inputs  and make sure the values are empty!
	emailField := doc.Find("input[type='email']").First()

	assert.NotNil(t, emailField, "Could not find email input field")
	val, exists := emailField.Attr("value")
	assert.True(t, exists, "Missing input field value item")

	assert.Equal(t, "", val, "Value is not empty")

}

func TestViewLoginForm_WithErrors(t *testing.T) {
	// Pipe the rendered template into goquery.
	r, w := io.Pipe()

	prefilled := "test"

	go func() {

		_ = LoginForm(LoginFormData{Defaults: dto.UserLoginDTO{Email: prefilled}, Errors: map[string]error{
			"email": validation.NewError("", "Email or password is invalid "),
		}}).Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	// Should have our error message
	assert.Equal(t, 1, doc.Find(".text-error").Length(), "Missing .text-error")

	// get the form inputs  and make sure the values are empty!
	emailField := doc.Find("input[type='email']").First()
	assert.NotNil(t, emailField, "Could not find email input field")
	val, exists := emailField.Attr("value")
	assert.True(t, exists, "Missing input field value item")

	assert.Equal(t, prefilled, val, "Value does not match previous entry")
}
