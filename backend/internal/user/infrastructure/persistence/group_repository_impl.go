package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/repository"
	"opscore/backend/internal/user/domain/value_object"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GroupRepositoryImpl is a PostgreSQL implementation of the GroupRepository interface
type GroupRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewGroupRepositoryImpl creates a new GroupRepositoryImpl
func NewGroupRepositoryImpl(db *pgxpool.Pool) repository.GroupRepository {
	return &GroupRepositoryImpl{db: db}
}

// Save persists a new group or updates an existing one
func (r *GroupRepositoryImpl) Save(ctx context.Context, group entity.Group) error {
	query := `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := r.db.Exec(ctx, query,
		group.ID().String(),
		group.Name(),
		group.Description(),
		group.CreatedAt(),
		group.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save group: %w", err)
	}

	// Save group-member associations
	err = r.saveGroupMembers(ctx, group.ID().String(), group.MemberIDs())
	if err != nil {
		return fmt.Errorf("failed to save group members: %w", err)
	}

	return nil
}

// saveGroupMembers saves the group-member associations
func (r *GroupRepositoryImpl) saveGroupMembers(ctx context.Context, groupID string, memberIDs []value_object.UserID) error {
	// Delete existing associations
	deleteQuery := `DELETE FROM user_groups WHERE group_id = $1;`
	_, err := r.db.Exec(ctx, deleteQuery, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete existing group members: %w", err)
	}

	// Insert new associations
	if len(memberIDs) > 0 {
		insertQuery := `INSERT INTO user_groups (user_id, group_id) VALUES ($1, $2);`
		batch := &pgx.Batch{}
		for _, memberID := range memberIDs {
			batch.Queue(insertQuery, memberID.String(), groupID)
		}

		results := r.db.SendBatch(ctx, batch)
		defer results.Close()

		for range memberIDs {
			_, err := results.Exec()
			if err != nil {
				return fmt.Errorf("failed to insert group member: %w", err)
			}
		}
	}

	return nil
}

// FindByID retrieves a group by its ID
func (r *GroupRepositoryImpl) FindByID(ctx context.Context, id value_object.GroupID) (entity.Group, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM groups
		WHERE id = $1;
	`

	var groupID, name, description string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, id.String()).Scan(
		&groupID,
		&name,
		&description,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find group by ID: %w", err)
	}

	// Get member IDs
	memberIDs, err := r.getGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}

	return r.toDomainEntity(groupID, name, description, memberIDs, createdAt, updatedAt)
}

// FindByMemberID retrieves all groups that contain the specified user as a member
func (r *GroupRepositoryImpl) FindByMemberID(ctx context.Context, userID value_object.UserID) ([]entity.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.created_at, g.updated_at
		FROM groups g
		INNER JOIN user_groups ug ON g.id = ug.group_id
		WHERE ug.user_id = $1
		ORDER BY g.created_at DESC;
	`

	rows, err := r.db.Query(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query groups by member ID: %w", err)
	}
	defer rows.Close()

	var groups []entity.Group
	for rows.Next() {
		var groupID, name, description string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&groupID, &name, &description, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group row: %w", err)
		}

		// Get member IDs
		memberIDs, err := r.getGroupMemberIDs(ctx, groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to get group members: %w", err)
		}

		group, err := r.toDomainEntity(groupID, name, description, memberIDs, createdAt, updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to domain entity: %w", err)
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over group rows: %w", err)
	}

	return groups, nil
}

// FindAll retrieves all groups
func (r *GroupRepositoryImpl) FindAll(ctx context.Context) ([]entity.Group, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM groups
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var groups []entity.Group
	for rows.Next() {
		var groupID, name, description string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&groupID, &name, &description, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group row: %w", err)
		}

		// Get member IDs
		memberIDs, err := r.getGroupMemberIDs(ctx, groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to get group members: %w", err)
		}

		group, err := r.toDomainEntity(groupID, name, description, memberIDs, createdAt, updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to domain entity: %w", err)
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over group rows: %w", err)
	}

	return groups, nil
}

// Update updates an existing group
func (r *GroupRepositoryImpl) Update(ctx context.Context, group entity.Group) error {
	query := `
		UPDATE groups
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4;
	`

	result, err := r.db.Exec(ctx, query,
		group.Name(),
		group.Description(),
		group.UpdatedAt(),
		group.ID().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("group with ID %s not found", group.ID().String())
	}

	// Update group-member associations
	err = r.saveGroupMembers(ctx, group.ID().String(), group.MemberIDs())
	if err != nil {
		return fmt.Errorf("failed to save group members: %w", err)
	}

	return nil
}

// Delete removes a group by its ID
func (r *GroupRepositoryImpl) Delete(ctx context.Context, id value_object.GroupID) error {
	// First delete group-member associations
	deleteMembersQuery := `DELETE FROM user_groups WHERE group_id = $1;`
	_, err := r.db.Exec(ctx, deleteMembersQuery, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete group members: %w", err)
	}

	// Then delete the group
	deleteQuery := `DELETE FROM groups WHERE id = $1;`
	result, err := r.db.Exec(ctx, deleteQuery, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("group with ID %s not found", id.String())
	}

	return nil
}

// getGroupMemberIDs retrieves the member IDs for a group
func (r *GroupRepositoryImpl) getGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	query := `
		SELECT user_id
		FROM user_groups
		WHERE group_id = $1
		ORDER BY user_id;
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to query group members: %w", err)
	}
	defer rows.Close()

	var memberIDs []string
	for rows.Next() {
		var memberID string
		if err := rows.Scan(&memberID); err != nil {
			return nil, fmt.Errorf("failed to scan member ID: %w", err)
		}
		memberIDs = append(memberIDs, memberID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over member ID rows: %w", err)
	}

	return memberIDs, nil
}

// toDomainEntity converts database data to a domain entity
func (r *GroupRepositoryImpl) toDomainEntity(
	groupID, name, description string,
	memberIDStrs []string,
	createdAt, updatedAt time.Time,
) (entity.Group, error) {
	id, err := value_object.NewGroupID(groupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %w", err)
	}

	memberIDs := make([]value_object.UserID, 0, len(memberIDStrs))
	for _, uidStr := range memberIDStrs {
		uid, err := value_object.NewUserID(strings.TrimSpace(uidStr))
		if err != nil {
			return nil, fmt.Errorf("invalid user ID: %w", err)
		}
		memberIDs = append(memberIDs, uid)
	}

	return entity.ReconstructGroup(id, name, description, memberIDs, createdAt, updatedAt), nil
}
