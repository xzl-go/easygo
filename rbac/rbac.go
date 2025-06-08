package rbac

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

type RBACManager struct {
	enforcer *casbin.Enforcer
}

func NewRBACManager(modelPath, policyPath string) (*RBACManager, error) {
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	adapter := fileadapter.NewAdapter(policyPath)
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return &RBACManager{
		enforcer: enforcer,
	}, nil
}

func (r *RBACManager) AddPolicy(sec, ptype string, rule []string) (bool, error) {
	return r.enforcer.AddPolicy(sec, ptype, rule)
}

func (r *RBACManager) RemovePolicy(sec, ptype string, rule []string) (bool, error) {
	return r.enforcer.RemovePolicy(sec, ptype, rule)
}

func (r *RBACManager) AddRoleForUser(user, role string) (bool, error) {
	return r.enforcer.AddRoleForUser(user, role)
}

func (r *RBACManager) DeleteRoleForUser(user, role string) (bool, error) {
	return r.enforcer.DeleteRoleForUser(user, role)
}

func (r *RBACManager) Enforce(sub, obj, act string) (bool, error) {
	return r.enforcer.Enforce(sub, obj, act)
}

func (r *RBACManager) GetRolesForUser(user string) ([]string, error) {
	return r.enforcer.GetRolesForUser(user)
}

func (r *RBACManager) GetPermissionsForUser(user string) ([][]string, error) {
	return r.enforcer.GetPermissionsForUser(user)
}

func (r *RBACManager) LoadPolicy() error {
	return r.enforcer.LoadPolicy()
}

func (r *RBACManager) SavePolicy() error {
	return r.enforcer.SavePolicy()
}
