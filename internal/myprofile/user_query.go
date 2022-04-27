package myprofile

import (
	"fmt"
	q "github.com/core-go/sql"
	"github.com/lib/pq"
	"strings"
)

func BuildQuery(filter interface{}) (query string, params []interface{}) {
	query = `select * from users`
	s := filter.(*UserFilter)
	var where []string

	i := 1
	if s.Interests != nil && len(s.Interests) > 0 {
		where = append(where, fmt.Sprintf(`interests && %s`, q.BuildDollarParam(i)))
		params = append(params, pq.Array(s.Interests))
		i++
	}
	if s.Skills != nil && len(s.Skills) > 0 {
		var skills []string
		for _, value := range s.Skills {
			skills = append(skills, fmt.Sprintf(`%s <@ ANY(skills)`, q.BuildDollarParam(i)))
			params = append(params, value)
			i++
		}
		where = append(where, fmt.Sprintf(`(%s)`, strings.Join(skills, " or ")))
	}
	if s.DateOfBirth != nil {
		if s.DateOfBirth.Min != nil {
			where = append(where, fmt.Sprintf(`date_of_birth >= %s`, q.BuildDollarParam(i)))
			params = append(params, s.DateOfBirth.Min)
			i++
		}
		if s.DateOfBirth.Max != nil {
			where = append(where, fmt.Sprintf(`date_of_birth <= %s`, q.BuildDollarParam(i)))
			params = append(params, s.DateOfBirth.Max)
			i++
		}
	}
	if len(s.Id) > 0 {
		where = append(where, fmt.Sprintf(`id = %s`, q.BuildDollarParam(i)))
		params = append(params, s.Id)
		i++
	}
	if len(s.Username) > 0 {
		where = append(where, fmt.Sprintf(`username ilike %s`, q.BuildDollarParam(i)))
		params = append(params, "%"+s.Username+"%")
		i++
	}
	if len(s.Email) > 0 {
		where = append(where, fmt.Sprintf(`email ilike %s`, q.BuildDollarParam(i)))
		params = append(params, s.Email+"%")
		i++
	}
	if len(s.Phone) > 0 {
		where = append(where, fmt.Sprintf(`phone ilike %s`, q.BuildDollarParam(i)))
		params = append(params, "%"+s.Phone+"%")
		i++
	}
	if s.Settings != nil {
		params = append(params, s.Settings)
		where = append(where, fmt.Sprintf(`settings && %s`, q.BuildDollarParam(i)))
		i++
	}
	if s.Achievements != nil && len(s.Achievements) > 0 {
		var achievements []string
		for _, value := range s.Achievements {
			params = append(params, value)
			achievements = append(achievements, fmt.Sprintf(`%s <@ ANY(achievements)`, q.BuildDollarParam(i)))
			i++
		}
		where = append(where, fmt.Sprintf(`(%s)`, strings.Join(achievements, " or ")))
	}
	if len(where) > 0 {
		query = query + ` where ` + strings.Join(where, " and ")
	}
	return
}
