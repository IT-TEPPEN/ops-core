package db

// IndexDefinition represents a database index definition.
type IndexDefinition struct {
	Name    string
	Table   string
	Columns []string
	Unique  bool
	Where   string // Optional WHERE clause for partial indexes
}

// GetIndexDefinitions returns all recommended database indexes for performance optimization.
func GetIndexDefinitions() []IndexDefinition {
	return []IndexDefinition{
		// Document indexes
		{
			Name:    "idx_documents_repository_id",
			Table:   "documents",
			Columns: []string{"repository_id"},
			Unique:  false,
		},
		{
			Name:    "idx_documents_owner",
			Table:   "documents",
			Columns: []string{"owner"},
			Unique:  false,
		},
		{
			Name:    "idx_documents_is_published",
			Table:   "documents",
			Columns: []string{"is_published"},
			Unique:  false,
		},
		{
			Name:    "idx_documents_published",
			Table:   "documents",
			Columns: []string{"is_published", "access_scope"},
			Unique:  false,
			Where:   "is_published = true",
		},
		{
			Name:    "idx_documents_created_at",
			Table:   "documents",
			Columns: []string{"created_at"},
			Unique:  false,
		},
		
		// Document version indexes
		{
			Name:    "idx_document_versions_document_id",
			Table:   "document_versions",
			Columns: []string{"document_id"},
			Unique:  false,
		},
		{
			Name:    "idx_document_versions_document_version",
			Table:   "document_versions",
			Columns: []string{"document_id", "version_number"},
			Unique:  true,
		},
		{
			Name:    "idx_document_versions_commit_hash",
			Table:   "document_versions",
			Columns: []string{"commit_hash"},
			Unique:  false,
		},
		
		// Execution record indexes
		{
			Name:    "idx_execution_records_document_id",
			Table:   "execution_records",
			Columns: []string{"document_id"},
			Unique:  false,
		},
		{
			Name:    "idx_execution_records_executor_id",
			Table:   "execution_records",
			Columns: []string{"executor_id"},
			Unique:  false,
		},
		{
			Name:    "idx_execution_records_status",
			Table:   "execution_records",
			Columns: []string{"status"},
			Unique:  false,
		},
		{
			Name:    "idx_execution_records_started_at",
			Table:   "execution_records",
			Columns: []string{"started_at"},
			Unique:  false,
		},
		{
			Name:    "idx_execution_records_completed_at",
			Table:   "execution_records",
			Columns: []string{"completed_at"},
			Unique:  false,
		},
		{
			Name:    "idx_execution_records_search",
			Table:   "execution_records",
			Columns: []string{"document_id", "executor_id", "status", "started_at"},
			Unique:  false,
		},
		
		// Execution record step indexes
		{
			Name:    "idx_execution_steps_record_id",
			Table:   "execution_steps",
			Columns: []string{"execution_record_id"},
			Unique:  false,
		},
		{
			Name:    "idx_execution_steps_record_step",
			Table:   "execution_steps",
			Columns: []string{"execution_record_id", "step_number"},
			Unique:  true,
		},
		
		// Attachment indexes
		{
			Name:    "idx_attachments_execution_record_id",
			Table:   "attachments",
			Columns: []string{"execution_record_id"},
			Unique:  false,
		},
		{
			Name:    "idx_attachments_step_id",
			Table:   "attachments",
			Columns: []string{"step_id"},
			Unique:  false,
		},
		
		// View history indexes
		{
			Name:    "idx_view_history_document_id",
			Table:   "view_history",
			Columns: []string{"document_id"},
			Unique:  false,
		},
		{
			Name:    "idx_view_history_user_id",
			Table:   "view_history",
			Columns: []string{"user_id"},
			Unique:  false,
		},
		{
			Name:    "idx_view_history_viewed_at",
			Table:   "view_history",
			Columns: []string{"viewed_at"},
			Unique:  false,
		},
		{
			Name:    "idx_view_history_aggregation",
			Table:   "view_history",
			Columns: []string{"document_id", "viewed_at"},
			Unique:  false,
		},
		
		// View statistics indexes
		{
			Name:    "idx_view_statistics_document_id",
			Table:   "view_statistics",
			Columns: []string{"document_id"},
			Unique:  true,
		},
		{
			Name:    "idx_view_statistics_total_views",
			Table:   "view_statistics",
			Columns: []string{"total_views"},
			Unique:  false,
		},
		{
			Name:    "idx_view_statistics_last_viewed_at",
			Table:   "view_statistics",
			Columns: []string{"last_viewed_at"},
			Unique:  false,
		},
		{
			Name:    "idx_view_statistics_popularity",
			Table:   "view_statistics",
			Columns: []string{"total_views", "last_viewed_at"},
			Unique:  false,
		},
		
		// User indexes
		{
			Name:    "idx_users_email",
			Table:   "users",
			Columns: []string{"email"},
			Unique:  true,
		},
		{
			Name:    "idx_users_username",
			Table:   "users",
			Columns: []string{"username"},
			Unique:  true,
		},
		
		// Group indexes
		{
			Name:    "idx_groups_name",
			Table:   "groups",
			Columns: []string{"name"},
			Unique:  false,
		},
		
		// Group member indexes
		{
			Name:    "idx_group_members_group_id",
			Table:   "group_members",
			Columns: []string{"group_id"},
			Unique:  false,
		},
		{
			Name:    "idx_group_members_user_id",
			Table:   "group_members",
			Columns: []string{"user_id"},
			Unique:  false,
		},
		{
			Name:    "idx_group_members_composite",
			Table:   "group_members",
			Columns: []string{"group_id", "user_id"},
			Unique:  true,
		},
	}
}

// GenerateCreateIndexSQL generates SQL for creating an index.
func GenerateCreateIndexSQL(idx IndexDefinition) string {
	sql := "CREATE"
	if idx.Unique {
		sql += " UNIQUE"
	}
	sql += " INDEX IF NOT EXISTS " + idx.Name + " ON " + idx.Table + " ("
	
	for i, col := range idx.Columns {
		if i > 0 {
			sql += ", "
		}
		sql += col
	}
	sql += ")"
	
	if idx.Where != "" {
		sql += " WHERE " + idx.Where
	}
	
	sql += ";"
	return sql
}

// GenerateDropIndexSQL generates SQL for dropping an index.
func GenerateDropIndexSQL(indexName string) string {
	return "DROP INDEX IF EXISTS " + indexName + ";"
}

// IndexRationale provides documentation for why certain indexes are recommended.
var IndexRationale = map[string]string{
	"idx_documents_repository_id": "Speeds up queries filtering documents by repository",
	"idx_documents_owner": "Optimizes queries filtering documents by owner",
	"idx_documents_is_published": "Improves performance when filtering published/unpublished documents",
	"idx_documents_published": "Partial index for published documents, optimizes common queries",
	"idx_documents_created_at": "Enables efficient sorting and filtering by creation date",
	
	"idx_document_versions_document_id": "Speeds up retrieving all versions of a document",
	"idx_document_versions_document_version": "Ensures uniqueness and fast lookup by version number",
	"idx_document_versions_commit_hash": "Enables efficient lookup by Git commit hash",
	
	"idx_execution_records_document_id": "Optimizes queries for execution records by document",
	"idx_execution_records_executor_id": "Speeds up filtering by executor (user)",
	"idx_execution_records_status": "Improves queries filtering by execution status",
	"idx_execution_records_started_at": "Enables efficient date-based queries and sorting",
	"idx_execution_records_completed_at": "Optimizes queries for completed executions",
	"idx_execution_records_search": "Composite index for complex search queries",
	
	"idx_execution_steps_record_id": "Speeds up retrieving steps for an execution record",
	"idx_execution_steps_record_step": "Ensures uniqueness and fast lookup of specific steps",
	
	"idx_attachments_execution_record_id": "Optimizes queries for attachments by execution record",
	"idx_attachments_step_id": "Speeds up retrieving attachments for a specific step",
	
	"idx_view_history_document_id": "Optimizes view history queries by document",
	"idx_view_history_user_id": "Speeds up user view history retrieval",
	"idx_view_history_viewed_at": "Enables efficient time-based queries",
	"idx_view_history_aggregation": "Optimizes aggregation queries for statistics",
	
	"idx_view_statistics_document_id": "Unique index ensures one statistics record per document",
	"idx_view_statistics_total_views": "Enables efficient sorting by popularity",
	"idx_view_statistics_last_viewed_at": "Optimizes queries for recently viewed documents",
	"idx_view_statistics_popularity": "Composite index for popular documents queries",
	
	"idx_users_email": "Ensures email uniqueness and fast user lookup by email",
	"idx_users_username": "Ensures username uniqueness and fast user lookup",
	
	"idx_groups_name": "Speeds up group lookup by name",
	
	"idx_group_members_group_id": "Optimizes queries for group members",
	"idx_group_members_user_id": "Speeds up finding groups for a user",
	"idx_group_members_composite": "Ensures uniqueness and optimizes membership checks",
}
