package horm

var (
	HasOne    = "has_one"
	HasMany   = "has_many"
	BelongsTo = "belongs_to"
)
var (
	CacheHard int = 2 // 程序启动时需要先加载缓存，再进行下一步
	CacheSoft int = 1 // 程序启动时使用协程加载缓存，不影响下一步
	CacheOff  int = 0
)
var (
	Model = Mixin{
		Column: []*Column{
			{
				Name:      "id",
				Type:      Int64,
				Increment: true,
				Primary:   true,
				Desc:      "ID",
			},
			{
				Name:      "created_at",
				Type:      Int64,
				Desc:      "创建时间",
				HideTable: true,
			},
			{
				Name:      "updated_at",
				Type:      Int64,
				Desc:      "更新时间",
				HideTable: true,
			},
			{
				Name:      "deleted_at",
				Type:      DeletedAt,
				Desc:      "删除时间",
				HideTable: true,
			},
		},
	}
)
