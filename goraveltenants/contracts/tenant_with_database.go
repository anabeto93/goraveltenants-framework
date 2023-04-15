package contracts

type TenantWithDatabase interface {
    Tenant
    Database() DatabaseConfig
    GetInternal(key string) interface{}
}
