package storage

import (
	"errors"
	"fmt"
)

// Code is a stable machine-readable storage failure category.
type Code string

const (
	CodePathUnavailable           Code = "storage_path_unavailable"
	CodePermissionDenied          Code = "storage_permission_denied"
	CodeOpenFailed                Code = "storage_open_failed"
	CodeBusy                      Code = "storage_busy"
	CodeCorrupt                   Code = "storage_corrupt"
	CodePragmaInvalid             Code = "storage_pragma_invalid"
	CodeMigrationInvalid          Code = "storage_migration_invalid"
	CodeMigrationFailed           Code = "storage_migration_failed"
	CodeMigrationChecksumMismatch Code = "storage_migration_checksum_mismatch"
	CodeSchemaTooNew              Code = "storage_schema_too_new"
	CodeTransactionFailed         Code = "storage_transaction_failed"
	CodeCommitFailed              Code = "storage_commit_failed"
	CodeClosed                    Code = "storage_closed"
)

// Error retains an internal cause without exposing it through user-facing projections.
type Error struct {
	Code      Code
	Operation string
	Retryable bool
	Cause     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Operation)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// Wrap classifies a storage operation while preserving the internal cause.
func Wrap(code Code, operation string, cause error, retryable bool) error {
	return &Error{Code: code, Operation: operation, Retryable: retryable, Cause: cause}
}

// IsCode reports whether err contains a classified storage error with code.
func IsCode(err error, code Code) bool {
	var storageErr *Error
	return errors.As(err, &storageErr) && storageErr.Code == code
}

// SafeError is the only storage error representation allowed across UI boundaries.
type SafeError struct {
	Code        string `json:"code"`
	SafeMessage string `json:"safeMessage"`
	Retryable   bool   `json:"retryable"`
}

// Project removes implementation details and returns fixed user-safe text.
func Project(err error) SafeError {
	var storageErr *Error
	if !errors.As(err, &storageErr) {
		return SafeError{
			Code:        string(CodeOpenFailed),
			SafeMessage: safeMessage(CodeOpenFailed),
			Retryable:   false,
		}
	}
	return SafeError{
		Code:        string(storageErr.Code),
		SafeMessage: safeMessage(storageErr.Code),
		Retryable:   storageErr.Retryable,
	}
}

func safeMessage(code Code) string {
	switch code {
	case CodePathUnavailable:
		return "无法确定本地数据存储位置。"
	case CodePermissionDenied:
		return "无法访问本地数据目录，请检查目录权限。"
	case CodeBusy:
		return "本地数据存储正被占用，请稍后重试。"
	case CodeCorrupt:
		return "本地数据存储无法读取，应用不会自动修改现有数据。"
	case CodePragmaInvalid:
		return "本地数据存储配置验证失败。"
	case CodeMigrationInvalid:
		return "本地数据结构记录无效。"
	case CodeMigrationFailed:
		return "本地数据结构升级失败。"
	case CodeMigrationChecksumMismatch:
		return "本地数据结构历史校验失败。"
	case CodeSchemaTooNew:
		return "本地数据由更新版本的 ItemAll 创建。"
	case CodeTransactionFailed:
		return "本地数据操作未能完成。"
	case CodeCommitFailed:
		return "本地数据提交失败。"
	case CodeClosed:
		return "本地数据存储已关闭。"
	case CodeOpenFailed:
		fallthrough
	default:
		return "无法打开本地数据存储。"
	}
}
