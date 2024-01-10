package services

// If a type is mergeable, the Servicer instance passed to the Merge method is merged into the original type instance
type MergeableServicer interface {
	Merge(A Servicer) (merged bool)
}
