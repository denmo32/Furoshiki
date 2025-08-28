package component

// DirtyLevelはウィジェットのダーティ状態のレベルを示します。
// これにより、再描画のみが必要か、レイアウトの再計算まで必要かを効率的に管理できます。
type DirtyLevel int

const (
	// LevelCleanはウィジェットがダーティでないことを示します。
	LevelClean DirtyLevel = iota
	// LevelRedrawDirtyはウィジェットの再描画のみが必要なことを示します。
	LevelRedrawDirty
	// LevelRelayoutDirtyはウィジェットのレイアウト再計算と再描画が必要なことを示します。
	LevelRelayoutDirty
)

// Dirtyはウィジェットのダーティ状態を管理します。
// 状態の伝播ロジックは含まず、状態管理に専念します。
type Dirty struct {
	level DirtyLevel
}

// NewDirtyは、新しいDirtyコンポーネントを生成します。
func NewDirty() *Dirty {
	return &Dirty{level: LevelClean}
}

// MarkDirtyはウィジェットのダーティレベルを設定します。
// より高いレベルのダーティ状態で上書きされることはありません。
func (d *Dirty) MarkDirty(relayout bool) {
	level := LevelRedrawDirty
	if relayout {
		level = LevelRelayoutDirty
	}
	if d.level < level {
		d.level = level
	}
}

// IsDirtyはウィジェットが再描画または再レイアウトを必要とするかを返します。
func (d *Dirty) IsDirty() bool {
	return d.level > LevelClean
}

// NeedsRelayoutはウィジェットがレイアウトの再計算を必要とするかを返します。
func (d *Dirty) NeedsRelayout() bool {
	return d.level == LevelRelayoutDirty
}

// ClearDirtyはダーティ状態をリセットします。
func (d *Dirty) ClearDirty() {
	d.level = LevelClean
}

// Levelは現在のダーティレベルを返します。
func (d *Dirty) Level() DirtyLevel {
	return d.level
}
