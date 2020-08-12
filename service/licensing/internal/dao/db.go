package dao

import (
	"context"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	pb "github.com/vazmin/eagle-eye-kratos/service/licensing/api"
)

const (
	_tableName = "licenses"
)

var (
	selectFields = []string{"license_id", "organization_id", "license_type",
		"product_name", "license_max", "license_allocated", "comment"}
)

func NewDB() (db *sql.DB, cf func(), err error) {
	var (
		cfg sql.Config
		ct paladin.TOML
	)
	if err = paladin.Get("db.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = sql.NewMySQL(&cfg)
	cf = func() {db.Close()}
	return
}

func (d *dao) RawLicensesByOrg(ctx context.Context, orgId string) (licenses []*pb.License, err error) {
	cond, values, err := builder.BuildSelect(_tableName, map[string]interface{}{"organization_id": orgId}, selectFields)
	if err != nil {
		return nil, err
	}
	query, err := d.db.Query(ctx, cond, values...)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	err = scanner.Scan(query, &licenses)
	return
}

func (d *dao) RawLicense(ctx context.Context, orgId string, licenseId string) (license *pb.License, err error) {
	where := map[string]interface{}{"organization_id": orgId, "license_id": licenseId}
	cond, values, err := builder.BuildSelect(_tableName, where, selectFields)
	if err != nil {
		return nil, err
	}
	query, err := d.db.Query(ctx, cond, values...)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	err = scanner.Scan(query, &license)
	return
}

func toMapList(license *pb.License) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"license_id": license.LicenseId,
			"organization_id": license.OrganizationId,
			"license_type": license.LicenseType,
			"product_name": license.ProduceName,
			"license_max": license.LicenseMax,
			"license_allocated": license.LicenseAllocated,
			"comment": license.Comment,
		},
	}
}
func toMap(license *pb.License) map[string]interface{} {
	return map[string]interface{}{
		"license_id": license.LicenseId,
		"organization_id": license.OrganizationId,
		"license_type": license.LicenseType,
		"product_name": license.ProduceName,
		"license_max": license.LicenseMax,
		"license_allocated": license.LicenseAllocated,
		"comment": license.Comment,
	}
}

func (d *dao) RawInsertLicense(ctx context.Context, license *pb.License) error {
	insert, values, err := builder.BuildInsert(_tableName, toMapList(license))
	if err != nil {
		return err
	}
	_, err = d.db.Exec(ctx, insert, values...)
	return err
}

func (d *dao) RawUpdateLicense(ctx context.Context, license *pb.License) (err error) {
	set := toMap(license)
	updateSet := builder.OmitEmpty(set, []string{"license_type", "product_name", "license_max", "license_allocated", "comment"})
	update, vals, err := builder.BuildUpdate(_tableName,
		map[string]interface{}{"license_id": license.LicenseId, "organization_id": license.OrganizationId}, updateSet)
	if err != nil {return err}
	_, err = d.db.Exec(ctx, update, vals...)
	return err
}

func (d *dao) RawDeleteLicense(ctx context.Context, license *pb.License) (err error) {
	buildDelete, vals, err := builder.BuildDelete(_tableName,
		map[string]interface{}{"license_id": license.LicenseId, "organization_id": license.OrganizationId})
	if err != nil {return err}
	_, err = d.db.Exec(ctx, buildDelete, vals...)
	return err
}