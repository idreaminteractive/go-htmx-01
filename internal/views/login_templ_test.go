package views

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestLoginForm_Get(t *testing.T) {
	// Pipe the rendered template into goquery.
	r, w := io.Pipe()

	go func() {
		_ = LoginForm(UserLoginDTO{}, UserLoginFormErrors{}).Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}
	// expect that the value of our input fields are empty
	if doc.Find(".text-error").Length() > 0 {
		t.Error("Found an error rendered when it was empty!")
	}

	// get the form inputs  and make sure the values are empty!
	emailField := doc.Find("input[type='email']").First()

	val, exists := emailField.Attr("value")
	if !exists {
		t.Error("Could not find value")
	}
	if val != "" {
		t.Error("Value is not empty")
	}

}

func TestLoginForm_WithErrors(t *testing.T) {
	// Pipe the rendered template into goquery.
	r, w := io.Pipe()

	prefilled := "test"
	go func() {
		_ = LoginForm(UserLoginDTO{Email: prefilled}, UserLoginFormErrors{Message: "Error in post"}).Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}
	// expect that the value of our input fields are empty
	if doc.Find(".text-error").Length() != 1 {
		t.Error("Missing our error message")
	}

	// get the form inputs  and make sure the values are empty!
	emailField := doc.Find("input[type='email']").First()

	val, exists := emailField.Attr("value")
	if !exists {
		t.Error("Could not find field")
	}

	if val != prefilled {
		t.Error("Value is not prefilled")
	}

}
