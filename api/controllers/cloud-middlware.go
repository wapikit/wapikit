//go:build managed_cloud
// +build managed_cloud

package controller

import (
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/interfaces"
)

func mountServices(app *interfaces.App, org *model.Organization) {
	// no-op for cloud edition
}
