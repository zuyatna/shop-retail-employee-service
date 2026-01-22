package usecase

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/dto/employee"
	"golang.org/x/crypto/bcrypt"
)

type IDGenerator interface {
	NewID() (string, error)
}

type EmployeeUsecase struct {
	repo        EmployeeRepository
	storageRepo StorageRepository
	idGen       IDGenerator
	ctxTimeout  time.Duration
}

func NewEmployeeUsecase(repo EmployeeRepository, storageRepo StorageRepository, idGen IDGenerator, timeout time.Duration) *EmployeeUsecase {
	return &EmployeeUsecase{
		repo:        repo,
		storageRepo: storageRepo,
		idGen:       idGen,
		ctxTimeout:  timeout,
	}
}

func (uc *EmployeeUsecase) Register(ctx context.Context, req employee.CreateEmployeeRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	existing, err := uc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() != EmployeeNotFoundError {
			return "", fmt.Errorf("failed to check existing email: %w", err)
		}
	}
	if existing != nil {
		return "", fmt.Errorf("email already exists")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	id, err := uc.idGen.NewID()
	if err != nil {
		return "", fmt.Errorf("failed to generate ID: %w", err)
	}

	birthDate, err := parseBirthDate(req.BirthDate)
	if err != nil {
		return "", fmt.Errorf("invalid birth date format: %w", err)
	}

	// Create domain newEmployee
	newEmployee, err := domain.NewEmployee(domain.NewEmployeeParams{
		ID:             domain.EmployeeID(id),
		Name:           req.Name,
		Email:          domain.Email(req.Email),
		HashedPassword: string(hashedBytes),
		Role:           domain.Role(req.Role),
		Position:       req.Position,
		Salary:         req.Salary,
		BirthDate:      birthDate,
		Address:        req.Address,
		City:           req.City,
		Province:       req.Province,
		PhoneNumber:    req.PhoneNumber,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create newEmployee domain: %w", err)
	}

	if err := uc.repo.Save(ctx, newEmployee); err != nil {
		return "", fmt.Errorf("failed to save newEmployee: %w", err)
	}
	log.Printf("New employee registered: ID=%s, Email=%s \n", id, req.Email)

	return id, nil
}

func (uc *EmployeeUsecase) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	if id == "" {
		return nil, fmt.Errorf("findByID ID cannot be empty")
	}

	findByID, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get findByID by ID: %w", err)
	}

	if findByID == nil {
		return nil, fmt.Errorf(EmployeeNotFoundError)
	}

	return findByID, nil
}

func (uc *EmployeeUsecase) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	findByEmail, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get findByEmail by email: %w", err)
	}

	if findByEmail == nil {
		return nil, fmt.Errorf(EmployeeNotFoundError)
	}

	return findByEmail, nil
}

func (uc *EmployeeUsecase) GetAll(ctx context.Context) ([]*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	employees, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all employees: %w", err)
	}

	return employees, nil
}

func (uc *EmployeeUsecase) UpdateProfile(ctx context.Context, id string, req employee.UpdateEmployeeRequest) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	findByID, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find findByID: %w", err)
	}

	updateIfPresent(req.Name, findByID.SetName)
	updateIfPresent(req.Position, findByID.SetPosition)
	updateIfPresent(req.Salary, findByID.SetSalary)
	updateIfPresent(req.Address, findByID.SetAddress)
	updateIfPresent(req.City, findByID.SetCity)
	updateIfPresent(req.Province, findByID.SetProvince)
	updateIfPresent(req.PhoneNumber, findByID.SetPhoneNumber)
	updateIfPresent(req.Photo, findByID.SetPhoto)

	if err := uc.repo.Update(ctx, findByID); err != nil {
		return fmt.Errorf("failed to update findByID in repo: %w", err)
	}
	emailLog := ""
	if req.Email != nil {
		emailLog = *req.Email
	}
	log.Printf("Update employee: ID=%s, Email=%s \n", id, emailLog)

	return nil
}

func (uc *EmployeeUsecase) UploadPhoto(ctx context.Context, employeeID string, file io.ReadSeeker, headerSize int64, contentType string, fileName string) error {
	{
		ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
		defer cancel()

		const MaxSize = 5 * 1024 * 1024 // 5 MB
		if headerSize > MaxSize {
			return fmt.Errorf("photo size exceeds the maximum limit of 5 MB")
		}

		ext := strings.ToLower(filepath.Ext(fileName))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return fmt.Errorf("unsupported file format: %s", ext)
		}

		if contentType != "image/jpeg" && contentType != "image/png" {
			return fmt.Errorf("unsupported content type: %s", contentType)
		}

		cfg, _, err := image.DecodeConfig(file)
		if err != nil {
			return fmt.Errorf("invalid image or corrupted file: %v", err)
		}

		if cfg.Width > 800 || cfg.Height > 800 {
			return fmt.Errorf("image dimensions exceed the maximum limit of 800x800 pixels")
		}

		// Reset file pointer after reading
		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}

		uniqueID, _ := uc.idGen.NewID()
		storageFileName := fmt.Sprintf("%s%s", uniqueID, ext)

		fileURL, err := uc.storageRepo.UploadFile(ctx, storageFileName, contentType, file, headerSize)
		if err != nil {
			return fmt.Errorf("failed to upload photo to storage: %w", err)
		}

		existingEmployee, err := uc.repo.FindByID(ctx, employeeID)
		if err != nil {
			return fmt.Errorf("failed to find employee by ID: %w", err)
		}
		if existingEmployee == nil {
			return fmt.Errorf(EmployeeNotFoundError)
		}

		existingEmployee.SetPhoto(fileURL)

		if err := uc.repo.Update(ctx, existingEmployee); err != nil {
			return fmt.Errorf("failed to update employee photo URL: %w", err)
		}

		log.Printf("Uploaded photo for employee ID=%s, URL=%s \n", employeeID, fileURL)

		return nil
	}
}

func (uc *EmployeeUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	findByID, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find findByID: %w", err)
	}

	if findByID == nil {
		return fmt.Errorf(EmployeeNotFoundError)
	}

	findByID.Delete()

	if err := uc.repo.Update(ctx, findByID); err != nil {
		return fmt.Errorf("failed to soft delete findByID: %w", err)
	}
	log.Printf("Deleted employee: ID=%s, Email=%s \n", id, findByID.Email())

	return nil
}

func parseBirthDate(dateStr string) (*time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateStr)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func updateIfPresent[T any](val *T, setter func(T)) {
	if val != nil {
		setter(*val)
	}
}
