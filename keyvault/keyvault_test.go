// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package keyvault

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
)

var (
	vaultName = "vault-sample-go-" + internal.GetRandomLetterSequence(5)
)

func TestMain(m *testing.M) {
	flag.StringVar(&vaultName, "vaultName", vaultName, "Specify name of vault to create.")

	err := internal.ParseArgs()
	if err != nil {
		log.Fatalln("failed to parse args")
	}
	os.Exit(m.Run())
}

func ExampleSetVaultPermissions() {
	internal.SetResourceGroupName("SetVaultPermissions")
	ctx := context.Background()
	defer resources.Cleanup(ctx)
	_, err := resources.CreateGroup(ctx, internal.ResourceGroupName())
	if err != nil {
		internal.PrintAndLog(err.Error())
	}

	_, err = CreateVault(ctx, vaultName)
	if err != nil {
		internal.PrintAndLog(err.Error())
	}
	internal.PrintAndLog("vault created")

	_, err = SetVaultPermissions(ctx, vaultName)
	if err != nil {
		internal.PrintAndLog(err.Error())
	}
	internal.PrintAndLog("set vault permissions")

	// Output:
	// vault created
	// set vault permissions
}
