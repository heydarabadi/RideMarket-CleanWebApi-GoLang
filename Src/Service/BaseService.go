package Service

import (
	"RideMarket-CleanWebApi-GoLang/Api/Dtos"
	"RideMarket-CleanWebApi-GoLang/Common"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Constants"
	"RideMarket-CleanWebApi-GoLang/Data/Database/DatabaseConfig"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"RideMarket-CleanWebApi-GoLang/pkg/ServiceErrors"
	"context"
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BaseService[T any, Tc any, Tu any, Tr any] struct {
	Database *gorm.DB
	Logger   Log.ILogger
	PreLoads []preload
}

type preload struct {
	string
}

func NewBaseService[T any, Tc any, Tu any, Tr any](cfg *Config.Config) *BaseService[T, Tc, Tu, Tr] {
	return &BaseService[T, Tc, Tu, Tr]{
		Database: DatabaseConfig.GetDb(),
		Logger:   Log.NewLogger(cfg),
	}
}

func (s *BaseService[T, Tc, Tu, Tr]) Create(ctx context.Context, req *Tc) (*Tr, error) {
	model, err := Common.TypeConverter[T](req)
	tx := s.Database.WithContext(ctx).Begin()
	err = tx.WithContext(ctx).Create(model).Error
	if err != nil {
		tx.WithContext(ctx).Rollback()
		s.Logger.Error(Log.Postgres, Log.Insert, err.Error(), nil)
	}
	err = tx.WithContext(ctx).Commit().Error
	if err != nil {
		tx.WithContext(ctx).Rollback()
		s.Logger.Error(Log.Postgres, Log.Commit, err.Error(), nil)
	}
	return Common.TypeConverter[Tr](model)

}

func (s *BaseService[T, Tc, Tu, Tr]) Update(ctx context.Context, id int, req *Tc) (*Tr, error) {
	updateMap, _ := Common.TypeConverter[map[string]interface{}](req)
	(*updateMap)["UpdatedBy"] = &sql.NullInt64{Int64: int64(ctx.Value(Constants.UserIdKey).(float64)), Valid: true}
	(*updateMap)["UpdatedAt"] = &sql.NullTime{Time: time.Now().UTC(), Valid: true}

	model := new(T)
	tx := s.Database.WithContext(ctx).Begin()

	if err := tx.Model(model).
		Where("id = ? and DeletedBy is null", id).
		Updates(*updateMap).Error; err != nil {
		tx.Rollback()
		s.Logger.Error(Log.Postgres, Log.Update, err.Error(), nil)
		return nil, err
	}

	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		tx.Rollback()
		s.Logger.Error(Log.Postgres, Log.Commit, err.Error(), nil)
		return nil, err
	}

	result, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *BaseService[T, Tc, Tu, Tr]) Delete(ctx context.Context, id int) error {
	tx := s.Database.WithContext(ctx).Begin()
	model := new(T)
	userId := tx.Statement.Context.Value(Constants.UserIdKey).(float64)

	deleteMap := map[string]interface{}{
		"DeletedBy": &sql.NullInt64{Int64: int64(userId), Valid: true},
		"DeletedAt": &sql.NullTime{time.Now().UTC(), true},
	}

	countAffected := tx.Model(model).
		Where("id = ? and DeletedBy is null", id).
		Updates(deleteMap).
		RowsAffected

	if countAffected == 0 {
		s.Logger.Error(Log.Postgres, Log.Update, ServiceErrors.RecordNotFound, nil)
		tx.Rollback()
		return &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.RecordNotFound}
	}
	tx.Commit()
	return nil
}

func (s *BaseService[T, Tc, Tu, Tr]) GetById(ctx context.Context, id int) (*Tr, error) {
	model := new(T)
	err := s.Database.
		Where("id = ? and DeletedBy is null", id).
		First(model).Error

	if err != nil {
		s.Logger.Error(Log.Postgres, Log.Select, err.Error(), nil)
		return nil, err
	}
	res, err := Common.TypeConverter[Tr](model)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func getQuery[T any](filter *Dtos.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	query := make([]string, 0)
	query = append(query, "DeletedBy is null")

	if filter.Filter != nil {
		for name, filter := range filter.Filter {
			_, ok := typeT.FieldByName(name)
			if ok {
				switch filter.Type {
				case "contains":
					query = append(query, fmt.Sprintf("%s ILIKE '%%%s%%'", name, filter.From))

				case "notcontains":
					query = append(query, fmt.Sprintf("%s NOT ILIKE '%%%s%%'", name, filter.From))

				case "startswith", "startwith", "startsWith":
					query = append(query, fmt.Sprintf("%s ILIKE '%s%%'", name, filter.From))

				case "endswith", "endwith", "endsWith":
					query = append(query, fmt.Sprintf("%s ILIKE '%%%s'", name, filter.From))

				case "eq", "equals", "e":
					query = append(query, fmt.Sprintf("%s = '%s'", name, filter.From))

				case "ne", "neq", "notequals", "notEqual":
					query = append(query, fmt.Sprintf("%s <> '%s'", name, filter.From))

				case "gt", "greaterthan", "greaterThan":
					query = append(query, fmt.Sprintf("%s > '%s'", name, filter.From))

				case "gte", "ge", "gteq", "greaterthanequal", "greaterThanEqual":
					query = append(query, fmt.Sprintf("%s >= '%s'", name, filter.From))

				case "lt", "lessthan", "lessThan":
					query = append(query, fmt.Sprintf("%s < '%s'", name, filter.From))

				case "lte", "le", "lteq", "lessthanequal", "lessThanEqual":
					query = append(query, fmt.Sprintf("%s <= '%s'", name, filter.From))

				case "inrange", "between", "range":
					if filter.To != "" {
						query = append(query, fmt.Sprintf("%s BETWEEN '%s' AND '%s'", name, filter.From, filter.To))
					} else {
						query = append(query, fmt.Sprintf("%s >= '%s'", name, filter.From))
					}

				case "isempty":
					query = append(query, fmt.Sprintf("(%s IS NULL OR %s = '')", name, name))

				case "isnotempty", "notempty":
					query = append(query, fmt.Sprintf("(%s IS NOT NULL AND %s <> '')", name, name))

				default:
					// unknown → skip silently
					break
				}
			}
		}
	}

	return strings.Join(query, " AND ")
}

func getSort[T any](filter *Dtos.DynamicFilter) string {
	t := new(T)
	typeT := reflect.TypeOf(*t)
	sort := make([]string, 0)

	if filter.Sort != nil {
		for _, tp := range *filter.Sort {
			fld, ok := typeT.FieldByName(tp.ColId)
			if ok && (tp.Sort == "asc" || tp.Sort == "desc") {
				sort = append(sort, fmt.Sprintf("%s %s", fld.Name, tp.Sort))
			}
		}
	}

	return strings.Join(sort, ", ")
}

func PreLoad(db *gorm.DB, preloads []preload) *gorm.DB {
	for _, item := range preloads {
		db = db.Preload(item.string)
	}
	return db
}

func Paginate[T any, Tr any](pagination *Dtos.PaginationInputWithFilter, preloads []preload, db *gorm.DB) (*Dtos.PagedList[Tr], error) {
	model := new(T)
	var items *[]T
	var rItems *[]Tr
	db = PreLoad(db, preloads)
	query := getQuery[T](&pagination.DynamicFilter)
	sort := getSort[T](&pagination.DynamicFilter)

	var totalRows int64 = 0
	db.Model(model).Where(query).Count(&totalRows)

	err := db.Where(query).Offset(int(pagination.GetOffset())).Limit(int(pagination.GetPageSize())).Order(sort).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	rItems, err = Common.TypeConverter[[]Tr](items)
	if err != nil {
		return nil, err
	}

	return NewPagedList(rItems, totalRows, pagination.PageNumber, pagination.PageSize), err

}

func NewPagedList[T any](items *[]T, count int64, pageNumber int64, pageSize int64) *Dtos.PagedList[T] {
	pl := &Dtos.PagedList[T]{
		PageNumber: pageNumber,
		TotalRows:  count,
		Items:      items,
	}
	pl.TotalPages = int64(math.Ceil(float64(count) / float64(pageSize)))
	pl.HasNextPage = pl.PageNumber < pl.TotalPages
	pl.HasPreviousPage = pl.PageNumber > 1
	return pl
}

func (s *BaseService[T, Tc, Tu, Tr]) GetByFilter(ctx context.Context, req *Dtos.PaginationInputWithFilter) (*Dtos.PagedList[Tr], error) {
	return Paginate[T, Tr](req, s.PreLoads, s.Database)
}
