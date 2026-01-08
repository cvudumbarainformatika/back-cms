package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// MenuController handles menu operations
type MenuController struct {
	db *sqlx.DB
}

// NewMenuController creates a new MenuController instance
func NewMenuController(db *sqlx.DB) *MenuController {
	return &MenuController{
		db: db,
	}
}

// FlexibleStringSlice can unmarshal from both JSON string and JSON array
type FlexibleStringSlice []string

// UnmarshalJSON implements custom unmarshaling for FlexibleStringSlice
func (f *FlexibleStringSlice) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*f = arr
		return nil
	}

	// If that fails, try to unmarshal as string (JSON-encoded array)
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Parse the string as JSON array
	if err := json.Unmarshal([]byte(str), &arr); err != nil {
		return err
	}

	*f = arr
	return nil
}

// MenuInput represents the input structure for menu operations
type MenuInput struct {
	ID        interface{}         `json:"id"` // Can be string like "menu-123456" for new items or number for existing
	Label     string              `json:"label"`
	Slug      string              `json:"slug"`
	To        string              `json:"to"`
	Icon      string              `json:"icon"`
	ParentID  interface{}         `json:"parentId"` // Can be string or number (not used in save, only for frontend reference)
	Position  string              `json:"position"`
	Order     int                 `json:"order"`
	IsActive  bool                `json:"isActive"`
	IsFixed   bool                `json:"isFixed"`
	IsDynamic bool                `json:"isDynamic"`
	Roles     FlexibleStringSlice `json:"roles"`
	Children  []MenuInput         `json:"children"`
}

// SaveMenusRequest represents the request body for saving menus
type SaveMenusRequest struct {
	Position string      `json:"position" binding:"required"`
	Menus    []MenuInput `json:"menus" binding:"required"`
}

// GetMenusByPosition returns menus by position (header, sidebar, footer)
// GET /api/v1/menus?position=header
func (mc *MenuController) GetMenusByPosition(c *gin.Context) {
	position := c.DefaultQuery("position", "header")

	// Validate position
	validPositions := map[string]bool{"header": true, "sidebar": true, "footer": true}
	if !validPositions[position] {
		utils.Error(c, http.StatusBadRequest, "invalid_position", "Invalid position. Must be: header, sidebar, or footer", nil)
		return
	}

	menus, err := models.GetMenusByPosition(mc.db, position)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch menus: "+err.Error(), nil)
		return
	}

	// Build hierarchical structure
	hierarchical := mc.buildHierarchy(menus)

	utils.Success(c, http.StatusOK, "Menus retrieved successfully", hierarchical)
}

// SaveMenus saves the entire menu structure for a position
// POST /api/v1/menus
func (mc *MenuController) SaveMenus(c *gin.Context) {
	var req SaveMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "validation_error", "Invalid request body: "+err.Error(), nil)
		return
	}

	// Validate position
	validPositions := map[string]bool{"header": true, "sidebar": true, "footer": true}
	if !validPositions[req.Position] {
		utils.Error(c, http.StatusBadRequest, "invalid_position", "Invalid position", nil)
		return
	}

	// Start transaction
	tx, err := mc.db.Beginx()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to start transaction", nil)
		return
	}
	defer tx.Rollback()

	// NOTE: We don't delete existing menus anymore to prevent duplication
	// Only new menus (with string ID like "menu-123") will be inserted
	// Existing menus (with numeric ID) will be skipped
	// To delete a menu, user must do it manually in the UI (future feature)

	// Create ID mapping for temporary IDs and existing IDs
	idMap := make(map[string]int64)

	// Save menus (insert new, skip existing)
	if err := mc.saveMenusRecursive(tx, req.Menus, nil, req.Position, idMap); err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to save menus: "+err.Error(), nil)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to commit transaction", nil)
		return
	}

	// Fetch updated menus
	menus, err := models.GetMenusByPosition(mc.db, req.Position)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch updated menus", nil)
		return
	}

	hierarchical := mc.buildHierarchy(menus)
	utils.Success(c, http.StatusOK, "Menus saved successfully", hierarchical)
}

// saveMenusRecursive recursively saves menus and their children
// Only inserts NEW menus (string ID or null), skips EXISTING menus (numeric ID)
func (mc *MenuController) saveMenusRecursive(tx *sqlx.Tx, menus []MenuInput, parentID *int64, position string, idMap map[string]int64) error {
	for _, menuInput := range menus {
		// Skip fixed menus completely (don't insert or update them)
		if menuInput.IsFixed {
			// But we still need to map their ID for children reference
			if menuInput.ID != nil {
				var idStr string
				switch v := menuInput.ID.(type) {
				case string:
					idStr = v
				case float64:
					idStr = fmt.Sprintf("%.0f", v)
				case int, int64:
					idStr = fmt.Sprintf("%v", v)
				}

				// For fixed menus, map their actual ID
				if idStr != "" {
					if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
						idMap[idStr] = id
					}
				}
			}

			// Process children of fixed menu (they might not be fixed)
			if len(menuInput.Children) > 0 {
				var fixedParentID *int64
				if menuInput.ID != nil {
					var idStr string
					switch v := menuInput.ID.(type) {
					case string:
						idStr = v
					case float64:
						idStr = fmt.Sprintf("%.0f", v)
					case int, int64:
						idStr = fmt.Sprintf("%v", v)
					}

					if mappedID, exists := idMap[idStr]; exists {
						fixedParentID = &mappedID
					}
				}

				if err := mc.saveMenusRecursive(tx, menuInput.Children, fixedParentID, position, idMap); err != nil {
					return err
				}
			}
			continue
		}

		// Check if this is an existing menu (numeric ID) or new menu (string ID / null)
		var isExistingMenu bool
		var existingMenuID int64

		if menuInput.ID != nil {
			var idStr string
			switch v := menuInput.ID.(type) {
			case string:
				idStr = v
			case float64:
				idStr = fmt.Sprintf("%.0f", v)
			case int, int64:
				idStr = fmt.Sprintf("%v", v)
			}

			// Try to parse as int64 - if successful, it's an existing menu
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				isExistingMenu = true
				existingMenuID = id
				idMap[idStr] = id
			}
		}

		if isExistingMenu {
			// This menu already exists in database - SKIP insert, just process children
			if len(menuInput.Children) > 0 {
				if err := mc.saveMenusRecursive(tx, menuInput.Children, &existingMenuID, position, idMap); err != nil {
					return err
				}
			}
			continue
		}

		// This is a NEW menu - proceed with insert
		rolesJSON, err := json.Marshal(menuInput.Roles)
		if err != nil {
			return err
		}

		menu := &models.Menu{
			Label:    menuInput.Label,
			Slug:     menuInput.Slug,
			To:       menuInput.To,
			Icon:     menuInput.Icon,
			ParentID: parentID, // Use parentID from recursion parameter
			Position: position,
			Order:    menuInput.Order,
			IsActive: menuInput.IsActive,
			IsFixed:  menuInput.IsFixed,
			Roles:    string(rolesJSON),
		}

		// Insert menu
		insertQuery := "INSERT INTO menus (label, slug, `to`, icon, parent_id, position, `order`, is_active, is_fixed, roles, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())"
		result, err := tx.Exec(insertQuery, menu.Label, menu.Slug, menu.To, menu.Icon, menu.ParentID, menu.Position, menu.Order, menu.IsActive, menu.IsFixed, menu.Roles)
		if err != nil {
			return err
		}

		// Get inserted ID
		insertedID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Map the new ID for children reference
		if menuInput.ID != nil {
			var idStr string
			switch v := menuInput.ID.(type) {
			case string:
				idStr = v
			case float64:
				idStr = fmt.Sprintf("%.0f", v)
			case int, int64:
				idStr = fmt.Sprintf("%v", v)
			}

			if idStr != "" {
				idMap[idStr] = insertedID
			}
		}

		// Recursively save children with this menu's ID as parent
		if len(menuInput.Children) > 0 {
			if err := mc.saveMenusRecursive(tx, menuInput.Children, &insertedID, position, idMap); err != nil {
				return err
			}
		}
	}

	return nil
}

// buildHierarchy builds a hierarchical menu structure from flat list
func (mc *MenuController) buildHierarchy(menus []models.Menu) []models.Menu {
	if len(menus) == 0 {
		return []models.Menu{}
	}

	// Create a map of menu ID to menu pointer
	menuMap := make(map[int64]*models.Menu)

	// First pass: populate map and initialize children
	for i := range menus {
		menus[i].Children = []models.Menu{}
		menuMap[menus[i].ID] = &menus[i]
	}

	// Second pass: build parent-child relationships
	var roots []models.Menu
	for i := range menus {
		if menus[i].ParentID == nil {
			// This is a root menu - we'll add it later after children are populated
			continue
		} else {
			// This is a child - add to parent's children
			if parent, exists := menuMap[*menus[i].ParentID]; exists {
				parent.Children = append(parent.Children, menus[i])
			}
		}
	}

	// Third pass: collect all root menus (now with their children populated)
	for i := range menus {
		if menus[i].ParentID == nil {
			roots = append(roots, menus[i])
		}
	}

	return roots
}

// GetMenuByID returns a single menu by ID
// GET /api/v1/menus/:id
func (mc *MenuController) GetMenuByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid menu ID", nil)
		return
	}

	menu, err := models.FindMenuByID(mc.db, id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch menu", nil)
		return
	}

	if menu == nil {
		utils.Error(c, http.StatusNotFound, "not_found", "Menu not found", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Menu retrieved successfully", menu)
}

// DeleteMenu deletes a menu by ID
// DELETE /api/v1/menus/:id
func (mc *MenuController) DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid menu ID", nil)
		return
	}

	// Check if menu exists
	menu, err := models.FindMenuByID(mc.db, id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch menu", nil)
		return
	}

	if menu == nil {
		utils.Error(c, http.StatusNotFound, "not_found", "Menu not found", nil)
		return
	}

	// Check if menu is fixed
	if menu.IsFixed {
		utils.Error(c, http.StatusForbidden, "forbidden", "Cannot delete fixed menu", nil)
		return
	}

	// Delete menu (will cascade delete children if database has foreign key constraint)
	// Or we can manually delete children first
	if err := mc.deleteMenuAndChildren(mc.db, id); err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to delete menu: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Menu deleted successfully", nil)
}

// deleteMenuAndChildren recursively deletes a menu and all its children
func (mc *MenuController) deleteMenuAndChildren(db *sqlx.DB, menuID int64) error {
	// First, find all children
	var children []models.Menu
	query := "SELECT id FROM menus WHERE parent_id = ?"
	if err := db.Select(&children, query, menuID); err != nil {
		return err
	}

	// Recursively delete children
	for _, child := range children {
		if err := mc.deleteMenuAndChildren(db, child.ID); err != nil {
			return err
		}
	}

	// Delete the menu itself
	deleteQuery := "DELETE FROM menus WHERE id = ?"
	_, err := db.Exec(deleteQuery, menuID)
	return err
}
