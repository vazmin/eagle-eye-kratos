// Code generated by protoc-gen-bm v0.1, DO NOT EDIT.
// source: api.proto

/*
Package api is a generated blademaster stub package.
This code was generated with kratos/tool/protobuf/protoc-gen-bm v0.1.

package 命名使用 {appid}.{version} 的方式, version 形如 v1, v2 ..

It is generated from these files:
	api.proto
*/
package api

import (
	"context"

	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
)
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

// to suppressed 'imported but not used warning'
var _ *bm.Context
var _ context.Context
var _ binding.StructValidator

var PathOrganizationSvcPing = "/eagle.organization.v1.OrganizationSvc/Ping"
var PathOrganizationSvcGetOrganization = "/organization"
var PathOrganizationSvcAddOrganization = "/organization"
var PathOrganizationSvcUpdateOrganization = "/organization"
var PathOrganizationSvcDeleteOrganization = "/organization"

// OrganizationSvcBMServer is the server API for OrganizationSvc service.
type OrganizationSvcBMServer interface {
	Ping(ctx context.Context, req *google_protobuf1.Empty) (resp *google_protobuf1.Empty, err error)

	GetOrganization(ctx context.Context, req *GetOrgReq) (resp *Organization, err error)

	AddOrganization(ctx context.Context, req *Organization) (resp *google_protobuf1.Empty, err error)

	UpdateOrganization(ctx context.Context, req *Organization) (resp *google_protobuf1.Empty, err error)

	DeleteOrganization(ctx context.Context, req *Organization) (resp *google_protobuf1.Empty, err error)
}

var OrganizationSvcSvc OrganizationSvcBMServer

func organizationSvcPing(c *bm.Context) {
	p := new(google_protobuf1.Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := OrganizationSvcSvc.Ping(c, p)
	c.JSON(resp, err)
}

func organizationSvcGetOrganization(c *bm.Context) {
	p := new(GetOrgReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := OrganizationSvcSvc.GetOrganization(c, p)
	c.JSON(resp, err)
}

func organizationSvcAddOrganization(c *bm.Context) {
	p := new(Organization)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := OrganizationSvcSvc.AddOrganization(c, p)
	c.JSON(resp, err)
}

func organizationSvcUpdateOrganization(c *bm.Context) {
	p := new(Organization)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := OrganizationSvcSvc.UpdateOrganization(c, p)
	c.JSON(resp, err)
}

func organizationSvcDeleteOrganization(c *bm.Context) {
	p := new(Organization)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := OrganizationSvcSvc.DeleteOrganization(c, p)
	c.JSON(resp, err)
}

// RegisterOrganizationSvcBMServer Register the blademaster route
func RegisterOrganizationSvcBMServer(e *bm.Engine, server OrganizationSvcBMServer) {
	OrganizationSvcSvc = server
	e.GET("/eagle.organization.v1.OrganizationSvc/Ping", organizationSvcPing)
	e.GET("/organization", organizationSvcGetOrganization)
	e.POST("/organization", organizationSvcAddOrganization)
	e.PUT("/organization", organizationSvcUpdateOrganization)
	e.DELETE("/organization", organizationSvcDeleteOrganization)
}
