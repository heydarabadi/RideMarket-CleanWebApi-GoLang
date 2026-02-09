package Models

// Role represents a user role in the RBAC system (e.g. admin, editor, support, finance)
type Role struct {
	BaseModel
	Name        string `gorm:"size:64;uniqueIndex;not null;comment:Unique short name of the role (e.g. admin, manager, viewer)"`
	DisplayName string `gorm:"size:100;not null;comment:Human-readable name shown in UI (e.g. System Administrator)"`
	Description string `gorm:"size:500;comment:Detailed explanation of what this role can do"`
	IsSystem    bool   `gorm:"default:false;index;comment:System-protected role (cannot be deleted or modified by normal users)"`
	IsDefault   bool   `gorm:"default:false;comment:New users automatically get this role (e.g. 'user' or 'guest')"`
	Priority    int    `gorm:"default:0;index;comment:Higher number = higher precedence (useful for role hierarchy / conflict resolution)"`
	Color       string `gorm:"size:7;comment:Hex color code for UI badges/tags (e.g. #FF5733)"`
	Icon        string `gorm:"size:100;comment:Icon name or SVG path for UI display (e.g. 'shield-lock')"`
	Scope       string `gorm:"size:50;index;comment:Optional scope like 'global', 'team', 'project' (for scoped RBAC)"`

	// Many-to-many relationship with Users
	Users []User `gorm:"many2many:user_roles;constraint:OnUpdate:No ACTION,OnDelete:NO ACTION;"`
}
