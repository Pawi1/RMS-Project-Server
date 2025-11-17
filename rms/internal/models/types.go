package models

type Role string

const (
    RoleAdmin    Role = "admin"
    RoleEditor   Role = "editor"
    RoleOperator Role = "operator"
    RoleMechanic Role = "mechanic"
    RoleViewer   Role = "viewer"
    RoleClient   Role = "client"
)