// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IOs interface {
		Run(ctx context.Context, arg string, workSpace string) error
		KillAll() error
	}
)

var (
	localOs IOs
)

func Os() IOs {
	if localOs == nil {
		panic("implement not found for interface IOs, forgot register?")
	}
	return localOs
}

func RegisterOs(i IOs) {
	localOs = i
}
