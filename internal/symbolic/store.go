package symbolic

// Store is the persistence boundary for canonical symbolic records.
// GetEntity is intentionally narrow: Service needs it to validate that an
// entity reference exists within the requested project.
type Store interface {
	SaveEntity(Entity) error
	GetEntity(projectID, id string) (Entity, error)
	ListEntities(projectID string) ([]Entity, error)

	SaveFact(Fact) error
	ListFacts(projectID string) ([]Fact, error)

	SaveRelation(Relation) error
	ListRelations(projectID string) ([]Relation, error)

	SaveConstraint(Constraint) error
	ListConstraints(projectID string) ([]Constraint, error)
}
