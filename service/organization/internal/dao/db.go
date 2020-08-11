package dao

import (
	"context"
	"github.com/didi/gendry/builder"
	pb "github.com/vazmin/eagle-eye-kratos/service/organization/api"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
)

const (
	_tableName = "organizations"
	_selectByOrgIdSQL = `SELECT organization_id, name, contact_name, contact_email, contact_phone FROM organizations WHERE organization_id = ?`
	_insertOrgSQL = `insert into organizations(organization_id, name, contact_name, contact_email, contact_phone) values(?,?,?,?,?)`
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

func (d *dao) RawOrg(ctx context.Context, orgId string) (org *pb.Organization, err error) {
	query := d.db.QueryRow(ctx, _selectByOrgIdSQL, orgId)
	org = &pb.Organization{}
	err = query.Scan(&org.Id, &org.Name, &org.ContactName, &org.ContactEmail, &org.ContactPhone)
	if err == sql.ErrNoRows {
		err = nil
		org = nil
		return
	}
	return
}

func (d *dao) RawInsertOrg(ctx context.Context, org *pb.Organization) (err error) {
	_, err = d.db.Exec(ctx, _insertOrgSQL, org.Id, org.Name, org.ContactName, org.ContactEmail, org.ContactPhone)
	return
}

func (d *dao) RawUpdateOrg(ctx context.Context, org *pb.Organization) (err error) {
	set := map[string]interface{}{
		"name": org.Name,
		"contact_name": org.ContactName,
		"contact_email": org.ContactEmail,
		"contact_phone": org.ContactPhone,
	}
	updateSet := builder.OmitEmpty(set, []string{"name", "contact_name", "contact_email", "contact_phone"})
	update, vals, err := builder.BuildUpdate(_tableName, map[string]interface{}{"organization_id": org.Id}, updateSet)
	if err != nil {return err}
	_, err = d.db.Exec(ctx, update, vals...)
	return err
}

func (d *dao) RawDeleteOrg(ctx context.Context, org *pb.Organization) (err error) {
	buildDelete, vals, err := builder.BuildDelete(_tableName, map[string]interface{}{"organization_id": org.Id})
	if err != nil {return err}
	_, err = d.db.Exec(ctx, buildDelete, vals...)
	return err
}
