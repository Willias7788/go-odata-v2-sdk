package main

import (
	"fmt"
	"log"

	"github.com/Willias7788/go-odata-v2-sdk/client"
	"github.com/Willias7788/go-odata-v2-sdk/config"
	"github.com/Willias7788/go-odata-v2-sdk/odata"
)

// Example Product Model
type Product struct {
	ID          string  `json:"Id"`
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	Price       float64 `json:"Price,string"` // OData often sends numbers as strings or numbers
}

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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
	servicePath := "/sap/opu/odata/IWBEP/GWSAMPLE_BASIC/"
	service := odata.NewService(sapClient, servicePath)

	fmt.Println("--- SDK Initialized ---")

	// 4. Perform Operations

	// A. GET Entity Set with Filter and Select
	fmt.Println("\n[GET] Fetching Products...")
	query := odata.NewQueryOptions().
		Top(5).
		Select([]string{"Id", "Name", "Price"}).
		Filter("Price gt 20.00") // OData V2 filter syntax

	productsResp, err := odata.GetEntitySet[Product](service, "ProductSet", query)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
	} else {
		// Access results: resp.D.Result (due to our generic wrapper)
		for _, p := range productsResp.D.Result {
			fmt.Printf("Product: %s - %s ($%.2f)\n", p.ID, p.Name, p.Price)
		}
	}

	// B. GET Entity By Key
	// Assuming we found a product ID from above, or hardcode one
	targetID := "HT-1000"
	fmt.Printf("\n[GET] Fetching Product By Key: %s...\n", targetID)
	
	productResp, err := odata.GetEntityByKey[Product](service, "ProductSet", fmt.Sprintf("'%s'", targetID), nil)
	if err != nil {
		log.Printf("Error fetching single product: %v", err)
	} else {
		p := productResp.D.Result
		fmt.Printf("Found: %s - %s\n", p.Name, p.Description)
	}

	// C. Create Entity
	fmt.Println("\n[POST] Creating new Product...")
	newProduct := Product{
		ID:   "HG-9999",
		Name: "Antigravity Boots",
		Description: "Defy physics",
		Price: 999.99,
	}
	
	createResp, err := odata.CreateEntity[Product](service, "ProductSet", newProduct)
	if err != nil {
		// This might fail on the public demo server due to permissions, but the logic holds
		log.Printf("Create failed (expected on public demo): %v", err)
	} else {
		fmt.Printf("Created Product: %s\n", createResp.D.Result.ID)
	}

	// D. Update Entity
	fmt.Println("\n[PUT] Updating Product...")
	newProduct.Name = "Antigravity Boots V2"
	err = odata.UpdateEntity(service, "ProductSet", fmt.Sprintf("'%s'", newProduct.ID), newProduct)
	if err != nil {
		log.Printf("Update failed: %v", err)
	} else {
		fmt.Println("Product Updated Successfully")
	}

	// E. Delete Entity
	fmt.Println("\n[DELETE] Deleting Product...")
	err = odata.DeleteEntity(service, "ProductSet", fmt.Sprintf("'%s'", newProduct.ID))
	if err != nil {
		log.Printf("Delete failed: %v", err)
	} else {
		fmt.Println("Product Deleted Successfully")
	}
}
