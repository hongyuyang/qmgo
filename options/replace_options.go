package options

import "github.com/hongyuyang/mongo-go-driver/mongo/options"

type ReplaceOptions struct {
	UpdateHook interface{}
	*options.ReplaceOptions
}
