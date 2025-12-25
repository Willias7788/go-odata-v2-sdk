package main

import (
	"fmt"
	"log"

	"github.com/Willias7788/go-odata-v2-sdk/client"
	"github.com/Willias7788/go-odata-v2-sdk/config"
	"github.com/Willias7788/go-odata-v2-sdk/odata"
)

// Example Material Model
type Material struct {
	Material  string `json:"Material"`
	CreatedOn string `json:"CreatedOn"`
	MatType   string `json:"MatType"`
	MatGrp    string `json:"MatGrp"` // OData often sends numbers as strings or numbers
	UOM       string `json:"UOM"`    // OData often sends numbers as strings or numbers
}

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println(cfg)

	// For demo purposes, we fallback if config is empty (so it runs without env file)
	if cfg.SAPHost == "" {
		fmt.Println("No configuration found, using example values...")
		cfg.SAPHost = "https://sapes5.sapdevcenter.com"
		cfg.SAPUsername = "User" // Replace with real one
		cfg.SAPPassword = "Pass" // Replace with real one
	}

	// 2. Initialize Base Client
	sapClient := client.NewSAPClient(cfg.SAPHost, cfg.SAPUsername, cfg.SAPPassword)
	sapClient.SetDebug(true) // Enable to see request/response logs

	// 3. Initialize OData Service
	// We point to a standard demo service
	servicePath := "/sap/opu/odata/sap/YGW_MM_001_SRV/"
	service := odata.NewService(sapClient, servicePath)

	fmt.Println("--- SDK Initialized ---")

	// 4. Perform Operations

	// A. GET Entity Set with Filter and Select
	fmt.Println("\n[GET] Fetching Products...")
	query := odata.NewQueryOptions().
		Top(5).
		Select([]string{"Material", "MatType", "CreatedOn"}) //.
		// Filter("Price gt 20.00") // OData V2 filter syntax

	productsResp, err := odata.GetEntitySet[Material](service, "MaterialMainSet", query)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
	} else {
		// Access results: resp.D.Result (due to our generic wrapper)
		for _, p := range productsResp.D.Result {
			fmt.Printf("Product: %s - %s ($%.2f)\n", p.Material, p.CreatedOn, p.MatType)
		}
	}

	// // B. GET Entity By Key
	// // Assuming we found a product ID from above, or hardcode one
	// targetID := "HT-1000"
	// fmt.Printf("\n[GET] Fetching Product By Key: %s...\n", targetID)

	// productResp, err := odata.GetEntityByKey[Material](service, "ProductSet", fmt.Sprintf("'%s'", targetID), nil)
	// if err != nil {
	// 	log.Printf("Error fetching single product: %v", err)
	// } else {
	// 	p := productResp.D.Result
	// 	fmt.Printf("Found: %s - %s\n", p.Material, p.CreatedOn)
	// }

	// // C. Create Entity
	// fmt.Println("\n[POST] Creating new Product...")
	// newProduct := Material{
	// 	Material:  "HT-1001",
	// 	CreatedOn: "2025-12-25",
	// 	MatType:   "HT",
	// 	MatGrp:    "HT",
	// 	UOM:       "EA",
	// }

	// createResp, err := odata.CreateEntity[Material](service, "MaterialMainSet", newProduct)
	// if err != nil {
	// 	// This might fail on the public demo server due to permissions, but the logic holds
	// 	log.Printf("Create failed (expected on public demo): %v", err)
	// } else {
	// 	fmt.Printf("Created Product: %s\n", createResp.D.Result.Material)
	// }

	// // D. Update Entity
	// fmt.Println("\n[PUT] Updating Product...")
	// newProduct.MatGrp = "HT-1001"
	// err = odata.UpdateEntity(service, "MaterialMainSet", fmt.Sprintf("'%s'", newProduct.Material), newProduct)
	// if err != nil {
	// 	log.Printf("Update failed: %v", err)
	// } else {
	// 	fmt.Println("Product Updated Successfully")
	// }

	// // E. Delete Entity
	// fmt.Println("\n[DELETE] Deleting Product...")
	// err = odata.DeleteEntity(service, "ProductSet", fmt.Sprintf("'%s'", newProduct.Material))
	// if err != nil {
	// 	log.Printf("Delete failed: %v", err)
	// } else {
	// 	fmt.Println("Product Deleted Successfully")
	// }
}
