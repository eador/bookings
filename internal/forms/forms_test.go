package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	isValid := form.Valid()
	if !isValid {
		t.Error("form is not valid when is should be")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required field missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")
	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid when all required fields exist")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	b := form.Has("a")
	if b {
		t.Error("form has field a when none was made")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	b = form.Has("a")
	if !b {
		t.Error("form does not have a field that was created")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "aa")
	form := New(postedData)

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
	form := New(postedData)

	form.IsEmail("b")
	form.IsEmail("a")
	if form.Errors.Get("a") != "Invalid email address" {
		t.Error("did not get error message for invalid email address")
	}
	if form.Errors.Get("b") != "" {
		t.Error("got error message for valid email address")
	}
}
