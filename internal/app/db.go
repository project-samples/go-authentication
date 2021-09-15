package app

import (
	"time"

	"github.com/gocql/gocql"
)

const (
	Keyspace = `authentication`

	CreateKeyspace = `create keyspace if not exists authentication with replication = {'class':'SimpleStrategy', 'replication_factor':1}`

	CreateTableAuthentication = `
					create table if not exists authentication.authentication (
					userid varchar,
					password varchar,
					history list<frozen<historytype>>,
					failtime timestamp,
					lockeduntiltime timestamp,
					successtime timestamp,
					failCount int,
					primary key (userid)
	)`
	CreateTableSignupCode = `
					create table if not exists authentication.signupCode (
					id varchar,
					passcode varchar,
					expiredat timestamp,
					primary key ((id,passcode,expiredat))
	)`
	CreateTableUser = `
					create table if not exists authentication.user (
						userid varchar,
						maxPasswordAge int,
						createdTime timestamp,
						updatedTime timestamp,
						updatedBy varchar,
						version int,
						userName varchar,
						email varchar,
						status varchar,
						signedUpTime timestamp,
						createdBy varchar,
						googleAccount varchar,
						googleActive boolean,
						googleemail varchar,
						displayname varchar,
						familyname	varchar,
						picture varchar,
						givenname varchar,
						primary key (userid)
	)`
	CreatePasswordRestCode = `
					create table if not exists authentication.passwordResetCode (
						id varchar,
						passcode varchar,
						expiredat timestamp,
						primary key ((id,passcode,expiredat))
	)`
	CreateGoogleConfig = `
					create table if not exists authentication.integrationconfiguration (
						clientId varchar,
						sourceType varchar,
						clientSecret varchar,
						primary key (clientId)
	)`
	CreateHistoryType = `CREATE TYPE IF NOT EXiSTS authentication.historytype (password varchar, timestamp timestamp)`
	ConfigGoogle      = `INSERT INTO authentication.integrationconfiguration (clientId, sourceType, clientSecret) VALUES ('590476077434-akbvo7usm6njup6ctiignt38adbecvpa.apps.googleusercontent.com', 'google', '-5gsoTyP8WDT3g2QPz-p5HTM' ) `
)

func Db(root *Root) (*gocql.ClusterConfig, error) {
	cluster := gocql.NewCluster(root.Cassandra.Uri)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra", Password: "cassandra"}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	err = session.Query(CreateKeyspace).Exec()
	if err != nil {
		return nil, err
	}

	// create table
	cluster.Keyspace = Keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	err = session.Query(CreateHistoryType).Exec()
	if err != nil {
		return nil, err
	}
	err = session.Query(CreateTableAuthentication).Exec()
	if err != nil {
		return nil, err
	}
	err = session.Query(CreateTableSignupCode).Exec()
	if err != nil {
		return nil, err

	}
	err = session.Query(CreateTableUser).Exec()
	if err != nil {
		return nil, err

	}
	err = session.Query(CreatePasswordRestCode).Exec()
	if err != nil {
		return nil, err
	}
	err = session.Query(CreateGoogleConfig).Exec()
	if err != nil {
		return nil, err

	}
	err = session.Query(ConfigGoogle).Exec()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return cluster, nil
}
