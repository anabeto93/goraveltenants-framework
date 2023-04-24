package assets

import (
	"path/filepath"

	"github.com/anabeto93/goraveltenants/database/models"
	"github.com/anabeto93/goraveltenants/contracts"
	"github.com/goravel/framework/facades"
)

func Boot() {}

func init() {
	config := facades.Config
	config.Add("tenancy", map[string]any{
		// tenant model to use
		"tenant_model": models.Tenant{},
		"domain_model": models.Domain{},
		/**
		* The list of domains hosting your central app.
		*
		* Only relevant if you're using the domain or subdomain identification middleware.
		*/
		"central_domains": []string {
			"127.0.0.1",
			"localhost",
		},
		"bootstrappers": []contracts.TenancyBootstrapper{},
		"database": map[string]any{
			"central_connection": config.Env("DB_CONNECTION", "central"),
			/**
			* Connection used as a "template" for the dynamically created tenant database connection.
			* Note: don't name your template connection tenant. That name is reserved by package.
			*/
			"template_tenant_connection": nil,
			/**
			* Tenant database names are created like this:
			* prefix + tenant_id + suffix.
			*/
			"prefix": "tenant",
			"suffix": "",
			"managers": map[string]contracts.TenantDatabaseManager{

			},
		},
		/**
		* Cache tenancy config. Used by CacheTenancyBootstrapper.
		*
		* This works for all Cache facade calls, cache() helper
		* calls and direct calls to injected cache stores.
		*
		* Each key in cache will have a tag applied on it. This tag is used to
		* scope the cache both when writing to it and when reading from it.
		*
		* You can clear cache selectively by specifying the tag.
		*/
		"cache": map[string]string{
			"tag_base": "tenant",
		},
		/**
		* Filesystem tenancy config. Used by FilesystemTenancyBootstrapper.
		* https://tenancyforlaravel.com/docs/v3/tenancy-bootstrappers/#filesystem-tenancy-boostrapper.
		*/
		"filesystem": map[string]any{

		},
		/**
		* Redis tenancy config. Used by RedisTenancyBootstrapper.
		*
		* Note: You need phpredis to use Redis tenancy.
		*
		* Note: You don't need to use this if you're using Redis only for cache.
		* Redis tenancy is only relevant if you're making direct Redis calls,
		* either using the Redis facade or by injecting it as a dependency.
		*/
		"redis": map[string]any{
			"prefix_base": "tenant",
			"prefixed_connections": []string{},
		},
		/**
		* Features are classes that provide additional functionality
		* not needed for tenancy to be bootstrapped. They are run
		* regardless of whether tenancy has been initialized.
		*
		* See the documentation page for each class to
		* understand which ones you want to enable.
		*/
		"features": []string{},
		/**
		* Should tenancy routes be registered.
		*
		* Tenancy routes include tenant asset routes. By default, this route is
		* enabled. But it may be useful to disable them if you use external
		* storage (e.g. S3 / Dropbox) or have a custom asset controller.
		*/
		"routes": true,

		/**
		* Parameters used by the tenants:migrate command.
		*/
		"migration_parameters": map[string]any{
			"--force": true,
			"--path": []string{databasePath("/tenant")},
			"--realpath": true,
		},
		/**
		* Parameters used by the tenants:seed command.
		*/
		"seeder_parameters": map[string]string{
			"--class": "DatabaseSeeder",
			//"--force": true,
		},
	})
}



func databasePath(databaseName string) string {
	return filepath.Join("database/migrations", databaseName)
}