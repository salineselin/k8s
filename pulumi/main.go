package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// create a container registry
		reg, err := container.NewRegistry(ctx, "make-registry", &container.RegistryArgs{
			Location: pulumi.String("US"),
		})
		if err != nil {
			return err
		}

		// log the container registrys url
		ctx.Export("registry: ", reg.BucketSelfLink)

		// get the underlying bucket backing the container registry
		regBucket, err := storage.GetBucket(ctx, "get-registry-bucket", reg.ID().ToIDOutput())

		// Create an admin serviceaccount
		saName := "salinesel-in-bucket-read"
		sa, err := serviceaccount.NewAccount(ctx, saName, &serviceaccount.AccountArgs{
			AccountId:   pulumi.String(saName),
			DisplayName: pulumi.String(saName),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		// bind it to the container registry
		_, err = storage.NewBucketIAMMember(ctx, "give-sa-bucket-permissions", &storage.BucketIAMMemberArgs{
			Bucket: regBucket.Name,
			Role:   pulumi.String("roles/storage.admin"),
			Member: sa.Email.ApplyT(func(Email string) string {
				return "serviceAccount:" + Email
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
