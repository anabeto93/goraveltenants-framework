package contracts

type ManagesDatabaseUsers interface {
    TenantDatabaseManager
    CreateUser(databaseConfig DatabaseConfig) bool
    DeleteUser(databaseConfig DatabaseConfig) bool
    UserExists(username string) bool
}
