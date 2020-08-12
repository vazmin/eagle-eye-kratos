package org

import (
	"context"
	orgpb "github.com/vazmin/eagle-eye-kratos/service/organization/api"
)

func (d *Dao) GetOrg(ctx context.Context, orgId string) (*orgpb.Organization, error) {
	return d.orgClient.GetOrganization(ctx, &orgpb.GetOrgReq{OrganizationId: orgId})
}
