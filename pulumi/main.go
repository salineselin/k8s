package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/container"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := container.NewRegistry(ctx, "registry", &container.RegistryArgs{
			Location: pulumi.String("US"),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
