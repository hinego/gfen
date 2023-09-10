package horm

var (
	HasOne    = "has_one"
	HasMany   = "has_many"
	BelongsTo = "belongs_to"
)

var (
	Model = Mixin{
		Column: []*Column{
			{
				Name:      "id",
				Type:      Int64,
				Increment: true,
				Primary:   true,
			},
			{
				Name: "created_at",
				Type: Int64,
			},
			{
				Name: "updated_at",
				Type: Int64,
			},
			{
				Name: "deleted_at",
				Type: DeletedAt,
			},
		},
	}
)
