package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/golang-jwt/jwt"
	"github.com/knadh/stuffbin"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/api_types"
	"github.com/wapikit/wapikit/interfaces"
	"golang.org/x/crypto/bcrypt"
)

func applyMigrations(db *sql.DB, fs stuffbin.FileSystem) error {
	fmt.Println("applying migrations.....")
	// now because we will using this function from the binary we have to use the path in reference with the binary stuffed static files
	migrationFilesPathPatterns := "/migrations/*.sql"
	migrationFilePaths, err := fs.Glob(migrationFilesPathPatterns)

	if err != nil {
		fmt.Println("error reading migration files with pattern: ", err)
		return fmt.Errorf("error reading migration files with pattern: %w", err)
	}
	// Sort migration files by filename
	sort.Slice(migrationFilePaths, func(i, j int) bool {
		return migrationFilePaths[i] < migrationFilePaths[j]
	})
	// Regex to match migration file names (YYYYMMDDHHmmSS.sql)
	migrationFilePattern := regexp.MustCompile(`^/migrations/\d{14}\.sql$`)

	for _, filePath := range migrationFilePaths {
		isMatch := migrationFilePattern.MatchString(filePath)
		if isMatch {
			// Read SQL from file
			sqlBytes, err := fs.Read(filePath)
			if err != nil {
				return fmt.Errorf("error reading migration file %s: %w", filePath, err)
			}

			// execute the migration
			_, err = db.Exec(string(sqlBytes))
			if err != nil {
				return fmt.Errorf("error executing migration %s: %w", filePath, err)
			}
		}
	}

	return nil
}

func installApp(db *sql.DB, fs stuffbin.FileSystem, prompt, idempotent bool) {
	if !idempotent {
		fmt.Println("** first time installation **")
		fmt.Println("** IMPORTANT: This will wipe existing wapikit tables and types in the DB")
	} else {
		fmt.Println("** first time (idempotent) installation **")
	}

	if prompt {
		var ok string
		fmt.Print("continue (y/N)?  ")
		if _, err := fmt.Scanf("%s", &ok); err != nil {
			logger.Error("error reading value from terminal: %v", err.Error(), nil)
		}
		if strings.ToLower(ok) != "y" {
			fmt.Println("install cancelled.")
			return
		}
	}

	// If idempotence is on, check if the DB is already setup.
	if idempotent {
		// fetch a user from db and check if this is a default user
		userQuery := SELECT(table.User.AllColumns).FROM(table.User)
		var users []model.User
		err := userQuery.Query(db, &users)
		if err != nil {
			logger.Error("error checking existing DB schema: %v", err)
			logger.Error("db is not initialized yet: %v", err)
		} else {
			logger.Error("skipping install as database seems to be already ready with schema migrations")
			os.Exit(0)
		}
	}

	// Migrate the tables.
	if err := applyMigrations(db, fs); err != nil {
		logger.Error("error migrating DB schema: %v", err)
		os.Exit(1)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(koa.String("app.default_user_password")), bcrypt.DefaultCost)
	fmt.Println("hashed password: ", string(hashedPassword))
	if err != nil {
		panic(err)
	}

	password := string(hashedPassword)

	defaultUser := model.User{
		Name:      koa.String("app.default_user_name"),
		Email:     koa.String("app.default_user_email"),
		Username:  koa.String("app.default_user_username"),
		Password:  &password,
		Status:    model.UserAccountStatusEnum_Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	defaultOrganization := model.Organization{
		Name:      "Default Organization",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var insertedUser model.User
	var insertedOrg model.Organization
	var insertedMember model.OrganizationMember

	insertDefaultUserQuery := table.User.INSERT(table.User.MutableColumns).
		MODEL(defaultUser).RETURNING(table.User.UniqueId)
	insertDefaultOrganizationQuery := table.Organization.INSERT(table.Organization.MutableColumns).MODEL(defaultOrganization).RETURNING(table.Organization.UniqueId)
	err = insertDefaultUserQuery.Query(db, &insertedUser)

	if err != nil {
		panic(err)
	}

	err = insertDefaultOrganizationQuery.Query(db, &insertedOrg)

	if err != nil {
		panic(err)
	}

	logger.Info("inserted default user: %v", insertedUser, insertedOrg)

	defaultOrgMember := model.OrganizationMember{
		AccessLevel:    model.UserPermissionLevelEnum_Owner,
		OrganizationId: insertedOrg.UniqueId,
		UserId:         insertedUser.UniqueId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	insertDefaultOrganizationMemberQuery := table.OrganizationMember.INSERT(table.OrganizationMember.MutableColumns).MODEL(defaultOrgMember).
		RETURNING(table.OrganizationMember.UniqueId)
	err = insertDefaultOrganizationMemberQuery.Query(db, &insertedMember)

	if err != nil {
		panic(err)
	}

	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username:       insertedUser.Username,
			Email:          insertedUser.Email,
			Role:           api_types.UserPermissionLevelEnum(insertedMember.AccessLevel),
			UniqueId:       insertedUser.UniqueId.String(),
			OrganizationId: insertedOrg.UniqueId.String(),
			Name:           insertedOrg.Name,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 2).Unix(), // 60-day expiration
			Issuer:    "wapikit",
		},
	}

	//Create the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(koa.String("app.jwt_secret")))
	if err != nil {
		panic(err)
	}

	defaultUserApiKey := model.ApiKey{
		MemberId:       insertedMember.UniqueId,
		OrganizationId: insertedOrg.UniqueId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Key:            token,
	}

	insertDefaultUserApiKeyQuery := table.ApiKey.INSERT(table.ApiKey.MutableColumns).MODEL(defaultUserApiKey).RETURNING(table.ApiKey.UniqueId)

	var insertedApiKey model.ApiKey

	err = insertDefaultUserApiKeyQuery.Query(db, &insertedApiKey)

	if err != nil {
		panic(err)
	}

	logger.Info("setup complete")
	logger.Info(`run the program and access the dashboard at %s`, koa.MustString("app.address"))
}
