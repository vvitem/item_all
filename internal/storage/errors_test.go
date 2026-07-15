package storage

import (
	"errors"
	"strings"
	"testing"
)

func TestProjectReturnsSafeStorageError(t *testing.T) {
	err := Wrap(CodePermissionDenied, "open", errors.New(`C:\Users\private\itemall.db`), true)
	got := Project(err)
	if got.Code != string(CodePermissionDenied) {
		t.Fatalf("Code = %q", got.Code)
	}
	if got.SafeMessage != "无法访问本地数据目录，请检查目录权限。" {
		t.Fatalf("message = %q", got.SafeMessage)
	}
	if !got.Retryable {
		t.Fatal("Retryable = false")
	}
	if strings.Contains(got.SafeMessage, "Users") || strings.Contains(got.SafeMessage, "itemall.db") {
		t.Fatalf("leak: %q", got.SafeMessage)
	}
}

func TestProjectUnknownErrorIsGeneric(t *testing.T) {
	got := Project(errors.New("raw driver detail"))
	if got.Code != string(CodeOpenFailed) || got.SafeMessage != "无法打开本地数据存储。" || got.Retryable {
		t.Fatalf("projection = %+v", got)
	}
}

func TestIsCodeTraversesWrappedError(t *testing.T) {
	if !IsCode(Wrap(CodeBusy, "write", errors.New("locked"), true), CodeBusy) {
		t.Fatal("missing code")
	}
}
