package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/knadh/stuffbin"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
	"golang.org/x/crypto/bcrypt"
)

func applyMigrations(db *sql.DB, fs stuffbin.FileSystem) error {
	// now because we will using this function from the binary we have to use the path in reference with the binary stuffed static files
	migrationFilesPathPatterns := "/migrations/*.sql"
	migrationFilePaths, err := fs.Glob(migrationFilesPathPatterns)
	if err != nil {
		return fmt.Errorf("error reading migration files with pattern: %w", err)
	}
	// Sort migration files by filename
	sort.Slice(migrationFilePaths, func(i, j int) bool {
		return migrationFilePaths[i] < migrationFilePaths[j]
	})
	// Regex to match migration file names (YYYYMMDDHHmmSS.sql)
	migrationFilePattern := regexp.MustCompile(`^\d{14}\.sql$`)
	for _, filePath := range migrationFilePaths {
		if migrationFilePattern.MatchString(filePath) {
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

func installApp(lastVer string, db *sql.DB, fs stuffbin.FileSystem, prompt, idempotent bool) {
	if !idempotent {
		fmt.Println("** first time installation **")
		fmt.Printf("** IMPORTANT: This will wipe existing listmonk tables and types in the DB '%s' **",
			koa.String("db.database"))
	} else {
		fmt.Println("** first time (idempotent) installation **")
	}

	if prompt {
		var ok string
		fmt.Print("continue (y/N)?  ")
		if _, err := fmt.Scanf("%s", &ok); err != nil {
			logger.Error("error reading value from terminal: %v", err)
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
		logger.Info("checking existing DB schema", users)
		if len(users) == 0 {
			logger.Error("db is not initialized yet: %v", err)

		} else {
			logger.Error("skipping install as database seems to be already ready with schema migrations")
			os.Exit(0)
		}
	}

	// Migrate the tables.
	if err := applyMigrations(db, fs); err != nil {
		logger.Error("error migrating DB schema: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(koa.String("app.default_user_password")), bcrypt.DefaultCost)
	fmt.Println("hashed password: ", string(hashedPassword))
	if err != nil {
		panic(err)
	}

	defaultUser := model.User{
		Name:     "Default User",
		Email:    koa.String("app.default_user_email"),
		Username: koa.String("app.default_user_username"),
		Password: string(hashedPassword),
		Status:   model.UserAccountStatusEnum_Active,
	}

	defaultOrganization := model.Organization{
		Name: "Default Organization",
	}

	var insertedUser []model.User
	var insertedOrg []model.Organization
	var insertedMember []model.OrganizationMember

	insertDefaultUserQuery := table.User.INSERT().
		MODEL(defaultUser).RETURNING(table.User.UniqueId)
	insertDefaultOrganizationQuery := table.Organization.INSERT().MODEL(defaultOrganization).RETURNING(table.Organization.UniqueId)
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
		AccessLevel:    model.UserPermissionLevel_Owner,
		OrganizationId: insertedOrg[0].UniqueId,
		UserId:         insertedUser[0].UniqueId,
	}

	insertDefaultOrganizationMemberQuery := table.OrganizationMember.INSERT().MODEL(defaultOrgMember)
	err = insertDefaultOrganizationMemberQuery.Query(db, &insertedMember)
	if err != nil {
		panic(err)
	}

	logger.Info("setup complete")
	logger.Info(`run the program and access the dashboard at %s`, koa.MustString("app.address"))
}
