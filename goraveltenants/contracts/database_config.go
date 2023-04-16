package contracts

type DatabaseConfig interface {
    GetName() string
    GetUsername() string
    GetPassword() string
    MakeCredentials()
    GetTemplateConnectionName() string
    Connection() map[string]interface{}
    TenantConfig() map[string]interface{}
    Manager() TenantDatabaseManager
    GeneratePasswordUsing(func (args ...interface{}) (string, error))
    GenerateDatabaseNameUsing(func (args ...interface{}) (string, error))
    GenerateUsernameUsing(func (args ...interface{}) (string, error))
}
