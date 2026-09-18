// Package core contains the universal domain records for the AI Core.
package core

import "time"

type ProjectStatus string

const ProjectActive ProjectStatus = "active"

type Project struct {
	ID        string
	Name      string
	Status    ProjectStatus
	CreatedAt time.Time
}
