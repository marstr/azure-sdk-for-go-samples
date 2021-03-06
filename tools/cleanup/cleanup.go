// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/helpers"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
)

func main() {
	var quiet bool
	var resourceGroupNamePrefix string
	flag.BoolVar(&quiet, "quiet", false, "Run quietly")
	flag.StringVar(&resourceGroupNamePrefix, "groupPrefix", helpers.GroupPrefix(), "Specify prefix name of resource group for sample resources.")

	err := helpers.ParseSubscriptionID()
	if err != nil {
		log.Fatalf("Error parsing subscriptionID: %v\n", err)
		os.Exit(1)
	}
	err = helpers.ParseDeviceFlow()
	if err != nil {
		log.Fatalf("Error parsing device flow: %v\n", err)
		log.Fatalf("Using device flow: %v", helpers.DeviceFlow())
	}
	flag.Parse()
	helpers.SetPrefix(resourceGroupNamePrefix)

	if !quiet {
		fmt.Println("Are you sure you want to delete all resource groups in the subscription? (yes | no)")
		var input string
		fmt.Scanln(&input)
		if input != "yes" {
			fmt.Println("Keeping resource groups")
			os.Exit(0)
		}
	}

	futures, groups := resources.DeleteAllGroupsWithPrefix(context.Background(), helpers.GroupPrefix())

	var wg sync.WaitGroup
	resources.WaitForDeleteCompletion(context.Background(), &wg, futures, groups)
	wg.Wait()

	fmt.Println("Done")
}
