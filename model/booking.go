package model

import (
	"context"
	"fmt"
	"time"
	"workspace_booking/config"
	"workspace_booking/migration"
)

type Booking struct {
	Id                 int16                `json:"id"`
	CityId             int16                `json:"city_id"`
	BuildingId         int16                `json:"building_id"`
	FloorId            int16                `json:"floor_id"`
	Purpose            string               `json:"purpose"`
	UserId             int16                `json:"user_id"`
	UserIds            []int16              `json:"user_ids"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	FromDateTime       string               `json:"from_datetime"`
	ToDateTime         string               `json:"to_datetime"`
	SelectedWorkspaces []*SelectedWorkspace `json:"selected_workspaces"`
}

type BookingTiming struct {
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type SelectedWorkspace struct {
	Date  string  `json:"date"`
	Seats []int16 `json:"seats"`
}

type BookingDetail struct {
	Id                 int16     `json:"id"`
	CityId             int16     `json:"city_id"`
	BuildingId         int16     `json:"building_id"`
	FloorId            int16     `json:"floor_id"`
	UserId             int16     `json:"user_id"`
	CityName           string    `json:"city_name"`
	BuildingName       string    `json:"building_name"`
	FloorName          string    `json:"floor_name"`
	UserName           string    `json:"user_name"`
	FromDateTime       time.Time `json:"from_datetime"`
	ToDateTime         time.Time `json:"to_datetime"`
	Purpose            string    `json:"purpose"`
	BookingParticipant []*BookingParticipantDetail
	BookingWorkspace   []*BookingWorkspaceDetail
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
type Bookings struct {
	Bookings []*Booking
}

type AvailableWorkspaces struct {
	FloorDetails     Floor
	FromDateTime     string `json:"from_datetime"`
	ToDateTime       string `json:"to_datetime"`
	BookedWorkSpaces []*BookedWorkSpace
}

type BookedWorkSpace struct {
	BookedDate   time.Time `json:"date"`
	WorkspaceIds []int     `json:"seats"`
}

// InsertBooking will create the booking record in db
func (b *Booking) InsertBooking() error {

	dt := time.Now()
	query := "INSERT INTO bookings (city_id, building_id, floor_id, from_datetime, to_datetime, purpose, " +
		"user_id, created_at, updated_at)" +
		" VALUES ($1, $2, $3, $4,$5, $6, $7, $8, $9) RETURNING id, created_at, updated_at"
	d := migration.DbPool.QueryRow(
		context.Background(), query, b.CityId, b.BuildingId, b.FloorId, b.FromDateTime, b.ToDateTime,
		b.Purpose, b.UserId, dt, dt,
	)
	err := d.Scan(&b.Id, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func GetMyBookingDetails(isForPast bool, userId int) []*BookingDetail {
	currTime := time.Now()
	currentDate := config.SqlTimeFormat(currTime)
	query := "SELECT id, city_id, building_id, floor_id, user_id, (select name from cities where id = bookings.city_id) as city_name, (select name from buildings where id = bookings.building_id) as city_name, (select name from floors where id = bookings.floor_id) as floor_name, (select name from users where id = bookings.user_id) as user_name, from_datetime, to_datetime, purpose, created_at, updated_at FROM bookings WHERE user_id = $1"
	var condition string
	if isForPast {
		condition = " AND from_datetime >= $2"
	} else {
		condition = " AND from_datetime < $2"
	}
	finalQuery := query + condition
	// query all bookings data
	bookings, e := migration.DbPool.Query(context.Background(), finalQuery, userId, currentDate)
	if e != nil {
		return nil
	}
	defer bookings.Close()

	// declare BookingDetail array variable
	bookingsDetails := make([]*BookingDetail, 0)

	// iterate over bookings
	for bookings.Next() {
		booking := new(BookingDetail)
		e = bookings.Scan(&booking.Id, &booking.CityId, &booking.BuildingId, &booking.FloorId, &booking.UserId, &booking.CityName, &booking.BuildingName, &booking.FloorName, &booking.UserName, &booking.FromDateTime, &booking.ToDateTime, &booking.Purpose, &booking.CreatedAt, &booking.UpdatedAt)
		if e != nil {
			fmt.Println("Failed to get bookings_details record :", e)
			return []*BookingDetail{}
		}
		booking.BookingParticipant = GetBookingParticipantsDetailsByBookingId(booking.Id)
		booking.BookingWorkspace = GetBookingWorkspacesDetailsByBookingId(booking.Id)
		bookingsDetails = append(bookingsDetails, booking)
	}
	return bookingsDetails
}

func BookingTimestamp(t *BookingTiming) (string, string) {
	fromDateTime := ConvertDateTime(t.FromDate, t.StartTime)
	toDateTime := ConvertDateTime(t.ToDate, t.EndTime)
	return fromDateTime, toDateTime
}

func GetAvailableBookingSpace(floorId int, fromDate, toDate string) AvailableWorkspaces {
	floor := GetFloorByID(floorId)

	bookedWorkSpacesRecord := make([]*BookedWorkSpace, 0)

	rows, _ := migration.DbPool.Query(context.Background(), "SELECT from_datetime as date, array_agg(workspace_id) as seats from booking_workspaces where floor_id = $1 and from_datetime between $2 and $3 and to_datetime between $4 and $5 group by from_datetime", floorId, fromDate, toDate, fromDate, toDate)

	defer rows.Close()

	for rows.Next() {
		bookingRecord := new(BookedWorkSpace)
		rows.Scan(&bookingRecord.BookedDate, &bookingRecord.WorkspaceIds)
		bookedWorkSpacesRecord = append(bookedWorkSpacesRecord, bookingRecord)
	}
	availableWorkspaces := AvailableWorkspaces{FloorDetails: floor, FromDateTime: fromDate, ToDateTime: toDate, BookedWorkSpaces: bookedWorkSpacesRecord}

	return availableWorkspaces
}
