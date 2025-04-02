package constants

import "time"

const (
	DefaultTimeout = 2 * time.Second
)

type ModelValidationType = string
const (
	ValidationTypeCreate ModelValidationType = "create"
	ValidationTypeUpdate ModelValidationType = "update"
)
