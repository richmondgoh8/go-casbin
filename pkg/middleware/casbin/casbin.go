package adapter

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jmoiron/sqlx"
)

type IAdapter interface {
	// LoadPolicy loads all policy rules from the storage.
	LoadPolicy(model model.Model) error
	// SavePolicy saves all policy rules to the storage.
	SavePolicy(model model.Model) error

	// AddPolicy adds a policy rule to the storage.
	// This is part of the Auto-Save feature.
	AddPolicy(sec string, ptype string, rule []string) error
	// RemovePolicy removes a policy rule from the storage.
	// This is part of the Auto-Save feature.
	RemovePolicy(sec string, ptype string, rule []string) error
	// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
	// This is part of the Auto-Save feature.
	RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error
}

type PGCasbinAdapter struct {
	db          *sqlx.DB
	schema      string
	controllers []string
}

type CasbinRule struct {
	ID    string
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

func NewCasbinAdapter(db *sqlx.DB, schema string) *PGCasbinAdapter {
	return &PGCasbinAdapter{
		db:     db,
		schema: schema,
	}
}

func (a *PGCasbinAdapter) LoadPolicy(model model.Model) error {
	query := fmt.Sprintf(`SELECT p.policy_type, lower(r.role_name), lower(p.controller), lower(p.action_type) FROM %s.%s p inner join test.roles r on p.role_id = r.id`, a.schema, "permissions")
	rows, err := a.db.Query(query)
	if err != nil {
		return err
	}

	var count int = 0

	for rows.Next() {
		count++
		var rule CasbinRule
		if rowErr := rows.Scan(
			&rule.PType,
			&rule.V0,
			&rule.V1,
			&rule.V2,
		); rowErr != nil {
			return err
		}
		loadPolicyLine(rule, model)
	}
	fmt.Printf("Loaded %d Policies\n", count)

	return nil
}

func (a *PGCasbinAdapter) SavePolicy(model model.Model) error {
	return errors.New("not implemented")
}

func (a *PGCasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

func (a *PGCasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

func (a *PGCasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}

func loadPolicyLine(rule CasbinRule, model model.Model) {
	lineText := rule.PType
	if rule.V0 != "" {
		lineText += ", " + rule.V0
	}
	if rule.V1 != "" {
		lineText += ", " + rule.V1
	}
	if rule.V2 != "" {
		lineText += ", " + rule.V2
	}
	if rule.V3 != "" {
		lineText += ", " + rule.V3
	}
	if rule.V4 != "" {
		lineText += ", " + rule.V4
	}
	if rule.V5 != "" {
		lineText += ", " + rule.V5
	}

	persist.LoadPolicyLine(lineText, model)
}
