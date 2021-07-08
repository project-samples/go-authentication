package cassandra

import (
	"context"
	"errors"

	// "database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/core-go/signup"
	"github.com/gocql/gocql"
)

type SignUpRepository struct {
	Cassandra          *gocql.ClusterConfig
	UserTable          string
	PasswordTable      string
	Status             signup.UserStatusConf
	MaxPasswordAge     int32
	MaxPasswordAgeName string

	UserIdName       string
	UserName         string
	ContactName      string
	StatusName       string
	PasswordName     string
	SignedUpTimeName string

	UpdatedTimeName string
	UpdatedByName   string
	VersionName     string

	GenderMapper signup.GenderMapper
	Schema       *signup.SignUpSchemaConfig
}

func NewSignUpRepositoryByConfig(db *gocql.ClusterConfig, userTable, passwordTable string, statusConfig signup.UserStatusConf, maxPasswordAge int32, c *signup.SignUpSchemaConfig, options ...signup.GenderMapper) *SignUpRepository {
	if len(c.UserName) == 0 {
		c.UserName = "username"
	}
	if len(c.Contact) == 0 {
		c.Contact = "email"
	}
	if len(c.Password) == 0 {
		c.Password = "password"
	}
	if len(c.Status) == 0 {
		c.Status = "status"
	}
	var genderMapper signup.GenderMapper
	if len(options) > 0 {
		genderMapper = options[0]
	}
	c.UserId = strings.ToLower(c.UserId)
	c.UserName = strings.ToLower(c.UserName)
	c.Contact = strings.ToLower(c.Contact)
	c.Password = strings.ToLower(c.Password)
	c.Status = strings.ToLower(c.Status)
	c.SignedUpTime = strings.ToLower(c.SignedUpTime)
	c.Language = strings.ToLower(c.Language)
	c.MaxPasswordAge = strings.ToLower(c.MaxPasswordAge)
	c.DateOfBirth = strings.ToLower(c.DateOfBirth)
	c.GivenName = strings.ToLower(c.GivenName)
	c.MiddleName = strings.ToLower(c.MiddleName)
	c.FamilyName = strings.ToLower(c.FamilyName)
	c.Gender = strings.ToLower(c.Gender)
	c.CreatedTime = strings.ToLower(c.CreatedTime)
	c.CreatedBy = strings.ToLower(c.CreatedBy)
	c.UpdatedTime = strings.ToLower(c.UpdatedTime)
	c.UpdatedBy = strings.ToLower(c.UpdatedBy)
	c.Version = strings.ToLower(c.Version)

	userName := c.UserName
	contact := c.Contact
	password := c.Password
	status := c.Status
	r := &SignUpRepository{
		Cassandra:          db,
		UserTable:          userTable,
		PasswordTable:      passwordTable,
		Status:             statusConfig,
		MaxPasswordAge:     maxPasswordAge,
		GenderMapper:       genderMapper,
		Schema:             c,
		MaxPasswordAgeName: c.MaxPasswordAge,
		UserIdName:         c.UserId,
		UserName:           userName,
		ContactName:        contact,
		PasswordName:       password,
		StatusName:         status,
		SignedUpTimeName:   c.SignedUpTime,

		UpdatedTimeName: c.UpdatedBy,
		UpdatedByName:   c.UpdatedBy,
		VersionName:     c.Version,
	}
	return r
}

func NewSignUpRepository(db *gocql.ClusterConfig, userTable, passwordTable string, statusConfig signup.UserStatusConf, maxPasswordAge int32, maxPasswordAgeName string, userId string, options ...string) *SignUpRepository {
	var contactName string
	if len(options) > 0 && len(options[0]) > 0 {
		contactName = options[0]
	}
	if len(contactName) == 0 {
		contactName = "email"
	}
	return &SignUpRepository{
		Cassandra:          db,
		UserTable:          userTable,
		PasswordTable:      passwordTable,
		Status:             statusConfig,
		MaxPasswordAge:     maxPasswordAge,
		MaxPasswordAgeName: maxPasswordAgeName,
		UserIdName:         userId,
		UserName:           "username",
		ContactName:        contactName,
		PasswordName:       "password",
		StatusName:         "status",
	}
}

func (s *SignUpRepository) Activate(ctx context.Context, id string) (bool, error) {
	version := 3
	if s.Status.Registered == s.Status.Verifying {
		version = 2
	}
	return s.updateStatus(ctx, id, s.Status.Verifying, s.Status.Activated, version, "")
}

func (s *SignUpRepository) SentVerifiedCode(ctx context.Context, id string) (bool, error) {
	if s.Status.Registered == s.Status.Verifying {
		return true, nil
	}
	return s.updateStatus(ctx, id, s.Status.Registered, s.Status.Verifying, 2, "")
}
func (s *SignUpRepository) CheckUserName(ctx context.Context, userName string) (bool, error) {
	session, er0 := s.Cassandra.CreateSession()
	if er0 != nil {
		return false, er0
	}
	query := fmt.Sprintf("Select %s from %s where %s = ?", s.UserName, s.UserTable, s.UserName)
	var username string
	err := session.Query(query, userName).Scan(&username)
	if err != nil {
		return false, err
	}
	return true, err
}

func (s *SignUpRepository) CheckUserNameAndContact(ctx context.Context, userName string, contact string) (bool, bool, error) {
	return s.existUserNameAndField(ctx, userName, s.ContactName, contact)
}

func (s *SignUpRepository) existUserNameAndField(ctx context.Context, userName string, fieldName string, fieldValue string) (bool, bool, error) {
	session, er0 := s.Cassandra.CreateSession()
	if er0 != nil {
		return false, false, er0
	}
	queryUsername := fmt.Sprintf("select %s from %s where %s = ? ALLOW FILTERING", s.UserName, s.UserTable, s.UserName)
	queryEmail := fmt.Sprintf("select %s from %s where %s = ? ALLOW FILTERING", fieldName, s.UserTable, fieldName)
	var userNameResult string
	var emailResult string
	nameExist := false
	emailExist := false
	rowsUsername := session.Query(queryUsername, userName).Iter()
	rowsEmail := session.Query(queryEmail, fieldValue).Iter()
	rowsUsername.Scan(&userNameResult)
	rowsEmail.Scan(&emailResult)
	if userNameResult == userName || emailResult == fieldValue {
		nameExist = true
		emailExist = true
		myErr := errors.New("Username or email already exists!")
		return false, false, myErr
	}
	defer session.Close()
	return nameExist, emailExist, nil
}

func (s *SignUpRepository) Save(ctx context.Context, userId string, info signup.SignUpInfo) (bool, error) {
	session, er0 := s.Cassandra.CreateSession()
	if er0 != nil {
		return false, er0
	}
	// tạo đối tượng user với các thuộc tính trong đối tượng.
	user := make(map[string]interface{})
	user[s.UserIdName] = userId
	user[s.UserName] = info.Username
	user[s.ContactName] = info.Contact
	user[s.StatusName] = s.Status.Registered
	if s.MaxPasswordAge > 0 && len(s.MaxPasswordAgeName) > 0 {
		user[s.MaxPasswordAgeName] = s.MaxPasswordAge
	}
	// t
	if s.Schema != nil {
		user = signup.BuildMap(ctx, user, userId, info, *s.Schema, s.GenderMapper)
	}

	if s.UserTable != s.PasswordTable && len(info.Password) > 0 {
		// kiểm tra password và tạo đối tượng pass
		pass := make(map[string]interface{})
		pass[s.UserIdName] = userId
		pass[s.PasswordName] = info.Password
		// tạo query lưu dữ liệu cho table user
		query, value := BuildInsert(s.UserTable, user)
		// tạo query lưu dữ liệu password cho table authentication
		passQuery, passValue := BuildInsert(s.PasswordTable, pass)
		// lưu user
		er1 := session.Query(query, value...).Exec()
		if er1 != nil {
			return false, er1
		}
		// lưu password
		er2 := session.Query(passQuery, passValue...).Exec()
		if er2 != nil {
			return false, er2
		}
		return false, nil
	}
	if len(info.Password) > 0 {
		user[s.PasswordName] = info.Password
	}
	query, value := BuildInsert(s.UserTable, user)
	er3 := session.Query(query, value...).Exec()
	if er3 != nil {
		return false, er3
	}
	defer session.Close()
	return true, er3
}

func (s *SignUpRepository) SavePasswordAndActivate(ctx context.Context, userId, password string) (bool, error) {
	session, er0 := s.Cassandra.CreateSession()
	if er0 != nil {
		return false, er0
	}
	user := make(map[string]interface{})
	user[s.UserIdName] = userId
	user[s.PasswordName] = password
	query, value := BuildInsert(s.PasswordTable, user)
	er1 := session.Query(query, value...).Exec()
	if er1 != nil {
		return false, er1
	}
	return s.Activate(ctx, userId)
}

func (s *SignUpRepository) updateStatus(ctx context.Context, id string, from, to string, version int, password string) (bool, error) {
	session, er0 := s.Cassandra.CreateSession()
	if er0 != nil {
		return false, er0
	}
	user := make(map[string]interface{})
	user[s.StatusName] = to
	if len(s.UpdatedTimeName) > 0 {
		user[s.UpdatedTimeName] = time.Now()
	}
	if len(s.UpdatedByName) > 0 {
		user[s.UpdatedByName] = id
	}
	if len(s.VersionName) > 0 && version > 0 {
		user[s.VersionName] = version
	}
	if s.UserTable == s.PasswordTable && len(password) > 0 && len(s.PasswordName) > 0 {
		user[s.PasswordName] = password
	}
	colNumber := 0
	var values []interface{}
	table := s.UserTable
	querySet := make([]string, 0)
	for colName, v2 := range user {
		values = append(values, v2)
		querySet = append(querySet, fmt.Sprintf("%v= ?", colName))
		colNumber++
	}
	queryWhere := fmt.Sprintf(" %s = ? ",
		s.UserIdName,
		// s.StatusName,
	)
	values = append(values, id)
	// values = append(values, from)
	query := fmt.Sprintf("update %v set %v where %v", table, strings.Join(querySet, ","), queryWhere)
	er1 := session.Query(query, values...).Exec()
	if er1 != nil {
		return false, er1
	}
	return true, nil
}

func BuildInsert(tableName string, user map[string]interface{}) (string, []interface{}) {
	var cols []string
	var values []interface{}
	for col, v := range user {
		cols = append(cols, col)
		values = append(values, v)
	}
	column := fmt.Sprintf("(%v)", strings.Join(cols, ","))
	numCol := len(cols)
	var arrValue []string
	for i := 0; i < numCol; i++ {
		arrValue = append(arrValue, buildParam(i))
	}
	value := fmt.Sprintf("(%v)", strings.Join(arrValue, ","))
	return fmt.Sprintf("INSERT INTO %v %v VALUES %v", tableName, column, value), values
}
func buildParam(i int) string {
	return "?"
}
