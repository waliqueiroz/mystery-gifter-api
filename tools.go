//go:build tools
// +build tools

package main

import (
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "go.uber.org/mock/mockgen"
)
