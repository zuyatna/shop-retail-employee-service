package repo

import (
	"context"
	"errors"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAttendanceRepo struct {
	collection *mongo.Collection
}

func NewMongoAttendanceRepo(db *mongo.Database) *MongoAttendanceRepo {
	return &MongoAttendanceRepo{
		collection: db.Collection("employee_attendances"),
	}
}

type attendanceModel struct {
	ID           string     `bson:"_id,omitempty" json:"id"`
	EmployeeID   string     `bson:"employee_id" json:"employee_id"`
	EmployeeName string     `bson:"employee_name" json:"employee_name"`
	Location     string     `bson:"location" json:"location"`
	CheckIn      time.Time  `bson:"check_in" json:"check_in"`
	CheckOut     *time.Time `bson:"check_out,omitempty" json:"check_out,omitempty"`
	IsLate       bool       `bson:"is_late" json:"is_late"`
	Date         time.Time  `bson:"date" json:"date"`
	UpdatedAt    time.Time  `bson:"updated_at" json:"updated_at"`
}

func (r *MongoAttendanceRepo) Save(ctx context.Context, attendance *domain.Attendance) error {
	model := attendanceModel{
		ID:           attendance.ID,
		EmployeeID:   attendance.EmployeeID,
		EmployeeName: attendance.EmployeeName,
		Location:     attendance.Location,
		CheckIn:      attendance.CheckIn,
		CheckOut:     attendance.CheckOut,
		IsLate:       attendance.IsLate,
		Date:         attendance.Date,
		UpdatedAt:    time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, model)
	return err
}

func (r *MongoAttendanceRepo) Update(ctx context.Context, attendance *domain.Attendance) error {
	filter := bson.M{"_id": attendance.ID}
	update := bson.M{
		"$set": bson.M{
			"check_out":  attendance.CheckOut,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoAttendanceRepo) FindByEmployeeIDAndDate(ctx context.Context, employeeID string, date time.Time) (*domain.Attendance, error) {
	filter := bson.M{
		"employee_id": employeeID,
		"date":        date,
	}

	var model attendanceModel
	err := r.collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, err
	}

	// Convert to domain.Attendance
	return &domain.Attendance{
		ID:           model.ID,
		EmployeeID:   model.EmployeeID,
		EmployeeName: model.EmployeeName,
		Location:     model.Location,
		CheckIn:      model.CheckIn,
		CheckOut:     model.CheckOut,
		IsLate:       model.IsLate,
		Date:         model.Date,
	}, nil
}
