package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("form is not valid when is should be")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required field missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid when all required fields exist")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	b := form.Has("a")
	if b {
		t.Error("form has field a when none was made")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	b = form.Has("a")
	if !b {
		t.Error("form does not have a field that was created")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "aa")

	r, _ := http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form := New(r.PostForm)

	b := form.MinLength("a", 1)
	if !b {
		t.Error("field doesn't meet min length when it should")
	}
	e := form.Errors.Get("a")
	if e != "" {
		t.Error("got error message when there should be none")
	}

	b = form.MinLength("a", 3)
	if b {
		t.Error("field meets min length when it should not")
	}
	e = form.Errors.Get("a")
	if e != "This field must be at least 3 characters long" {
		t.Error("got wrong error message")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "test@email.com")
	r, _ := http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form := New(r.PostForm)

	form.IsEmail("b")
	form.IsEmail("a")
	if form.Errors.Get("a") != "Invalid email address" {
		t.Error("did not get error message for invalid email address")
	}
	if form.Errors.Get("b") != "" {
		t.Error("got error message for valid email address")
	}
}
