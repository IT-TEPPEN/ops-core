package repository

import "context"

// Context keys for filters. Repository implementations can use these keys
// to extract filter criteria from the context passed by Apply.
const (
	filterKeyRepositoryID         = "document_filter_repository_id"
	filterKeyProviderRepositoryID = "document_filter_provider_repository_id"
	filterKeyOwner                = "document_filter_owner"
	filterKeyRepository           = "document_filter_repository"
)

// repositoryIDFilter filters by internal repository ID.
type repositoryIDFilter struct {
	repoID string
}

// NewRepositoryIDFilter creates a filter for repository ID.
func NewRepositoryIDFilter(repoID string) Filter {
	return repositoryIDFilter{repoID: repoID}
}

func (f repositoryIDFilter) Apply(ctx context.Context) context.Context {
	return context.WithValue(ctx, filterKeyRepositoryID, f.repoID)
}

// providerRepositoryIDFilter filters by provider repository ID.
type providerRepositoryIDFilter struct {
	providerRepoID string
}

// NewProviderRepositoryIDFilter creates a filter for provider repository ID.
func NewProviderRepositoryIDFilter(providerRepoID string) Filter {
	return providerRepositoryIDFilter{providerRepoID: providerRepoID}
}

func (f providerRepositoryIDFilter) Apply(ctx context.Context) context.Context {
	return context.WithValue(ctx, filterKeyProviderRepositoryID, f.providerRepoID)
}

// ownerFilter filters by repository owner.
type ownerFilter struct {
	owner string
}

// NewOwnerFilter creates a filter for repository owner.
func NewOwnerFilter(owner string) Filter {
	return ownerFilter{owner: owner}
}

func (f ownerFilter) Apply(ctx context.Context) context.Context {
	return context.WithValue(ctx, filterKeyOwner, f.owner)
}

// repositoryNameFilter filters by repository name.
type repositoryNameFilter struct {
	repository string
}

// NewRepositoryNameFilter creates a filter for repository name.
func NewRepositoryNameFilter(repository string) Filter {
	return repositoryNameFilter{repository: repository}
}

func (f repositoryNameFilter) Apply(ctx context.Context) context.Context {
	return context.WithValue(ctx, filterKeyRepository, f.repository)
}

// ParsedFilters holds extracted filter values from context.
type ParsedFilters struct {
	RepositoryID         string
	ProviderRepositoryID string
	Owner                string
	Repository           string
}

// ParseFilters extracts filter values from context populated by Filter.Apply.
func ParseFilters(ctx context.Context) ParsedFilters {
	get := func(key string) string {
		if v := ctx.Value(key); v != nil {
			if s, ok := v.(string); ok {
				return s
			}
		}
		return ""
	}

	return ParsedFilters{
		RepositoryID:         get(filterKeyRepositoryID),
		ProviderRepositoryID: get(filterKeyProviderRepositoryID),
		Owner:                get(filterKeyOwner),
		Repository:           get(filterKeyRepository),
	}
}

// ApplyFilters applies all provided filters to the context in order, returning the final context.
func ApplyFilters(ctx context.Context, filters ...Filter) context.Context {
	cur := ctx
	for _, f := range filters {
		cur = f.Apply(cur)
	}
	return cur
}
