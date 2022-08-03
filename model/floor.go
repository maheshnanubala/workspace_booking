package model

import (
	"context"
	"fmt"
	"time"
	"workspace_booking/migration"
)

// Building struct
type Floor struct {
	Id              int              `json:"id"`
	Name            string           `json:"name"`
	TotalWorkSpace  int              `json:"total_workspace"`
	TotalConference int              `json:"total_conference"`
	BuildingId      int64            `json:"building_id"`
	BuildingName    string           `json:"building_name"`
	FloorWorkSpaces []FloorWorkSpace `json:"workspaces"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// Floors struct
type Floors struct {
	Floors []*Floor
}

func GetAllFloors() []*Floor {
	rows, _ := migration.DbPool.Query(context.Background(), "select floors.id, floors.name, floors.total_workspace, floors.total_conference, floors.building_id, buildings.name as building_name, floors.created_at, floors.updated_at from floors LEFT JOIN buildings ON floors.building_id = buildings.id")
	defer rows.Close()

	// declare empty post variable
	floors := make([]*Floor, 0)
	// iterate over rows
	for rows.Next() {
		floor := new(Floor)
		err := rows.Scan(&floor.Id, &floor.Name, &floor.TotalWorkSpace, &floor.TotalConference, &floor.BuildingId, &floor.BuildingName, &floor.CreatedAt, &floor.UpdatedAt)
		allFloorWorkspaces := GetAllworkspaces(floor.Id)
		floor.FloorWorkSpaces = allFloorWorkspaces
		if err != nil {
			fmt.Println("Failed to get floors record :", err)
			return []*Floor{}
		}
		floors = append(floors, floor)
	}
	return floors
}

func (floor *Floor) CreateFloor() error {
	dt := time.Now()
	query := "INSERT INTO floors (name, total_workspace, total_conference, building_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at"
	building := migration.DbPool.QueryRow(context.Background(), "select name from buildings where id = $1", &floor.BuildingId)
	err := building.Scan(&floor.BuildingName)
	if err != nil {
		return err
	}
	d := migration.DbPool.QueryRow(context.Background(), query, &floor.Name, &floor.TotalWorkSpace, &floor.TotalConference, &floor.BuildingId, dt, dt)
	err = d.Scan(&floor.Id, &floor.CreatedAt, &floor.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetFloorByID(floorId int) Floor {
	floor := Floor{}
	// rows := migration.DbPool.QueryRow(context.Background(), "select floors.id, floors.name, floors.total_workspace, floors.building_id, buildings.name as building_name, floors.created_at, floors.updated_at from floors LEFT JOIN buildings ON floors.building_id = buildings.id where floors.id = $1", floorId)
	rows := migration.DbPool.QueryRow(context.Background(), "select * from floors where id = $1", floorId)
	err := rows.Scan(&floor.Id, &floor.Name, &floor.TotalWorkSpace, &floor.TotalConference, &floor.BuildingId, &floor.CreatedAt, &floor.UpdatedAt)
	allFloorWorkspaces := GetAllworkspaces(floor.Id)
	floor.FloorWorkSpaces = allFloorWorkspaces
	if err != nil {
		fmt.Println("Failed to get floors record :", err)
		return Floor{}
	}
	return floor
}
