package persistence

import (
	"log"

	"github.com/vice-registry/vice-util/models"
	r "gopkg.in/gorethink/gorethink.v3"
)

// InitDatabase initializes an empty database instance (e.g. creates admin account)
func InitDatabase() {
	createDatabase(connectionProperties.Database)
	createTable("users")
	createIndex("users", "Username")
	createTable("environments")
	createIndex("environments", "Userid")
	createTable("images")
	createIndex("images", "Userid")
	createTable("deployments")
	createIndex("deployments", "Userid")
	createAdminUser()
}

func createDatabase(name string) {
	_, err := r.DBCreate(name).RunWrite(connectionProperties.Session)
	if err != nil {
		log.Printf("Database %s not created: %s", name, err)
	}
}

func createTable(name string) {
	log.Printf("Create (if not exist) RethinkDB table %s ...", name)
	_, err := r.DB(connectionProperties.Database).TableCreate(name, r.TableCreateOpts{
		PrimaryKey: "ID",
	}).RunWrite(connectionProperties.Session)
	if err != nil {
		log.Printf("Table %s not created: %s", name, err)
	}
}

func createIndex(table string, field string) {
	_, err := r.DB(connectionProperties.Database).Table(table).IndexCreate(field).RunWrite(connectionProperties.Session)
	if err != nil {
		log.Printf("Index for %s on Table %s not created: %s", field, table, err)
	}
}

func createAdminUser() {
	//var admin models.User
	admin, _ := GetUser("admin")
	if admin == nil || admin.ID == "" {
		// create admin user
		var admin models.User
		admin.ID = "admin"
		admin.Username = "admin"
		admin.Password = "admin"
		admin.Email = "admin@vice-registry.org"
		admin.Fullname = "Admin User"
		err := createItem(admin, "users")
		if err != nil {
			log.Printf("Admin user not created: %s", err)
		}
	}
}
