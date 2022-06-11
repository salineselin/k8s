package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// add google beta provider
		projectName := config.Get(ctx, "gcp:project")
		google_beta, err := gcp.NewProvider(ctx, "google-beta", &gcp.ProviderArgs{
			Project: pulumi.String(projectName),
		})
		if err != nil {
			return err
		}

		// enable artifact registry api
		api, err := projects.NewService(ctx, "enable-artifacts-registry-api", &projects.ServiceArgs{
			DisableDependentServices: pulumi.Bool(true),
			Service:                  pulumi.String("artifactregistry.googleapis.com"),
		})
		if err != nil {
			return err
		}

		// create an artifact registry
		reg, err := artifactregistry.NewRepository(ctx, "saline-selin-artifacts-repo", &artifactregistry.RepositoryArgs{
			Location:     pulumi.String("us-west3"),
			RepositoryId: pulumi.String("salinesel.in"),
			Description:  pulumi.String("container registry for gcr.salinesel.in"),
			Format:       pulumi.String("DOCKER"),
		}, pulumi.Provider(google_beta), pulumi.DependsOn([]pulumi.Resource{api}))
		if err != nil {
			return err
		}

		// Create an admin serviceaccount
		saName := "saline-selin-artifacts-admin"
		sa, err := serviceaccount.NewAccount(ctx, saName, &serviceaccount.AccountArgs{
			AccountId:   pulumi.String(saName),
			DisplayName: pulumi.String(saName),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		// bind it to the artifact registry
		_, err = artifactregistry.NewRepositoryIamMember(ctx, "saline-selin-artifacts-binding", &artifactregistry.RepositoryIamMemberArgs{
			Repository: reg.Name,
			Role:       pulumi.String("roles/artifactregistry.admin"),
			Member: sa.Email.ApplyT(func(Email string) string {
				return "serviceAccount:" + Email
			}).(pulumi.StringOutput),
		}, pulumi.Provider(google_beta))
		if err != nil {
			return err
		}
		return nil
	})
}
