# dedicated SAP OData v2.0 Golang SDK

A generic, production-ready Go SDK for consuming SAP OData v2.0 services. This library handles the complexities of SAP OData communication, including CSRF token management, cookie handling, and query string building, allowing you to focus on your business logic.

## 🚀 Features

- **Generic & Reusable**: Not tied to any specific SAP module. Works with any compliant OData v2 service.
- **CSRF Token Management**: Automatically handles `X-CSRF-Token` fetching, caching, and retries for `POST`, `PUT`, `PATCH`, and `DELETE` operations.
- **Fluent Query Builder**: Easily construct complex OData queries with `$filter`, `$select`, `$expand`, `$top`, `$skip`, and `$orderby`.
- **Navigation Properties**: Traverse related entity sets via OData navigation properties (e.g., `EntitySet('key')/NavProperty`).
- **Type-Safe**: Uses Go generics for strict typing of response entities.
- **Configurable**: Supports configuration via environment variables or `.env` file (using Viper).
- **Resilient**: Built on top of [go-resty](https://github.com/go-resty/resty) for robust HTTP communication.

## 📦 Installation

```bash
go get github.com/Willias7788/go-odata-v2-sdk
```

## 🛠️ Configuration

The SDK supports loading configuration from environment variables or a `.env` file.

**Environment Variables:**
```env
SAP_HOST=https://your-sap-gateway.com
SAP_USERNAME=your_username
SAP_PASSWORD=your_password
SAP_CLIENT=100  # Optional
```

## 📚 Usage Examples

### 1. Initialize the Client and Service

```go
package main

import (
	"log"

	"github.com/Willias7788/go-odata-v2-sdk/client"
	"github.com/Willias7788/go-odata-v2-sdk/config"
	"github.com/Willias7788/go-odata-v2-sdk/odata"
)

func main() {
	// Load config from .env or environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create the base SAP HTTP client
	sapClient := client.NewSAPClient(cfg.SAPHost, cfg.SAPUsername, cfg.SAPPassword)
	sapClient.SetDebug(true) // Optional: Log authenticated requests

	// Initialize the OData Service Wrapper
	// Point this to your specific service root
	servicePath := "/sap/opu/odata/IWBEP/GWSAMPLE_BASIC/"
	service := odata.NewService(sapClient, servicePath)
}
```

### 2. Define Your Model

Define a struct that matches your OData entity. Use JSON tags to map fields.

```go
type Product struct {
	ID          string  `json:"Id"`
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	Price       float64 `json:"Price,string"` // Handle string-to-number if needed
}
```

### 3. Fetch Entities (GET)

Use the `QueryOptions` builder to filter and select data.

```go
// Build Query: $top=5&$select=Id,Name,Price&$filter=Price gt 20.00
query := odata.NewQueryOptions().
	Top(5).
	Select([]string{"Id", "Name", "Price"}).
	Filter("Price gt 20.00")

// Execute GET
resp, err := odata.GetEntitySet[Product](service, "ProductSet", query)
if err != nil {
	log.Fatal(err)
}

for _, p := range resp.D.Result {
	log.Printf("Product: %s - %s\n", p.Name, p.ID)
}
```

### 4. Create Entity (POST)

The SDK automatically handles the CSRF token exchange required for creation.

```go
newProduct := Product{
	ID:          "HG-9999",
	Name:        "Antigravity Boots",
	Description: "Defy physics",
	Price:       999.99,
}

resp, err := odata.CreateEntity[Product](service, "ProductSet", newProduct)
if err != nil {
	log.Fatal("Create failed:", err)
}
log.Printf("Created: %s", resp.D.Result.ID)
```

### 5. Navigate to Related Entities (Navigation Property)

Use `GetNavigationSet` to traverse OData navigation properties. This builds a URL like `EntitySet('key')/NavigationProperty`.

```go
// Define the related entity struct
type TransferOrder struct {
	DO_NUM  string `json:"DO_NUM"`
	DO_DATE string `json:"DO_DATE"`
	TO_NUM  string `json:"TO_NUM"`
	TO_DATE string `json:"TO_DATE"`
	SO_NUM  string `json:"SO_NUM"`
	SO_DATE string `json:"SO_DATE"`
}

// Fetch transfer orders linked to a delivery order
// URL: DeliveryOrderSet('8120010348')/toTransferOrder
resp, err := odata.GetNavigationSet[TransferOrder](
	service,
	"DeliveryOrderSet",
	"('8120010348')",
	"toTransferOrder",
	nil,
)
if err != nil {
	log.Fatal(err)
}

for _, to := range resp.D.Result {
	log.Printf("DO: %s → TO: %s → SO: %s\n", to.DO_NUM, to.TO_NUM, to.SO_NUM)
}
```

### 6. Update Entity (PUT/PATCH)

```go
newProduct.Name = "Antigravity Boots V2"

// Entity Key often requires special formatting, e.g., "('HG-9999')"
key := "('HG-9999')" 

err := odata.UpdateEntity(service, "ProductSet", key, newProduct)
if err != nil {
	log.Fatal("Update failed:", err)
}
```

## 📂 Project Structure

```text
├── client/           # Core HTTP client, Auth, & CSRF Logic
├── config/           # Configuration management
├── models/           # Generic OData wrapper structs
├── odata/            # High-level OData service & Query builder
└── examples/         # Runnable usage examples
```

## 🤝 Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

[MIT](https://choosealicense.com/licenses/mit/)
