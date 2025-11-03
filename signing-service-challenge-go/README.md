# Signature Service - Coding Challenge
## New version after changes

### Build and Run
```bash
# Install dependencies
make install

# Build the application
make build

# Run tests
make test

# Build and run
make run
```

### Available Make Commands
```bash
make help              # Show all available commands
make build             # Build the application
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make test-domain       # Run domain tests only
make test-api          # Run API tests only
make test-crypto       # Run crypto tests only
make test-persistence  # Run persistence tests only
make run               # Build and run the application
make clean             # Clean build artifacts
make check             # Format, vet and test
```

## Implemented Features

### ✅ Core Functionality
- **RESTful API** using Gin framework
- **Signature Devices**: Create and manage RSA/ECDSA signing devices
- **Transaction Signing**: Sign data with monotonically increasing counter
- **Signature Chaining**: Each signature includes the previous signature (blockchain-like)
- **Thread-Safe Operations**: Concurrent-safe counter increment with mutex
- **In-Memory Storage**: Thread-safe repository with CRUD operations

### 🔐 Security Features
- **RSA Signing**: RSA-PSS with SHA-256
- **ECDSA Signing**: ECDSA with P-384 curve and SHA-256
- **Signature Counter**: Strictly monotonically increasing, gap-free
- **Signature Format**: `<counter>_<data>_<last_signature_base64>`

### 📡 API Endpoints
```
POST   /api/v0/devices          - Create signature device (RSA or ECDSA)
GET    /api/v0/devices          - List all devices
GET    /api/v0/devices/:id      - Get device by ID
POST   /api/v0/devices/:id/sign - Sign transaction data
GET    /api/v0/health           - Health check
```

### 🧪 Testing
- **27 test cases** across 5 test files
- **100% passing** tests including concurrency tests
- **Table-driven tests** following Go best practices
- **Coverage**: Domain methods, API endpoints, crypto, persistence

### 🏗️ Architecture
```
domain/          - Business logic and device model
api/             - HTTP handlers with Gin
crypto/          - RSA/ECDSA signers and key generation
persistence/     - In-memory repository (ready for DB migration)
```

## AI Tools Usage

### GitHub Copilot
This project was developed with the assistance of **GitHub Copilot** as an AI pair programming tool. Below is a detailed breakdown of how AI was utilized throughout the development process:

#### Areas Where AI Was Used:

**1. Domain Model Implementation (`domain/device.go`)**
- The `GetSecuredDataToSign()` method logic was developed with AI assistance, ensuring correct format: `<counter>_<data>_<last_signature>`
- AI suggested the pattern for base64 encoding the device ID for the initial signature case

**2. Unit Tests (`*_test.go` files)**
- **All test files** were created with GitHub Copilot assistance using table-driven test patterns
- AI helped structure 27 test cases across 5 test files covering:
  - Domain methods (device_methods_test.go) - 7 test functions
  - Crypto key generation (generation_test.go) - 2 test functions
  - In-memory persistence (inmemory_test.go) - 4 test functions
  - API endpoints (device_test.go, server_test.go) - 5 test functions
- The concurrency test (`TestIncrementCounter_Concurrency`) was designed with AI to verify thread-safety with 100 concurrent goroutines

**3. Thread-Safe Counter Implementation**
- GitHub Copilot provided guidance on implementing the monotonically increasing signature counter
- AI recommended using `sync.Mutex` to prevent race conditions in concurrent environments
- The `IncrementCounter()` method implementation with proper locking was developed with AI assistance

**4. Documentation**
- This README structure and content organization was created with GitHub Copilot
- The Makefile commands documentation was formatted with AI assistance

#### Design Decisions Made by Developer:

While AI assisted with implementation, the following architectural decisions were made independently:
- Choice of Gin framework for the HTTP API layer
- Repository pattern for persistence layer (preparing for future database migration)
- Separation of concerns: domain, API, crypto, and persistence packages
- RESTful API endpoint design and routing structure
- Error handling strategies and HTTP status code selections

#### Development Approach:

The development followed a **collaborative approach** where:
1. Core architecture and design patterns were defined by the developer
2. GitHub Copilot assisted with boilerplate code, test structures, and implementation details
3. All AI-generated code was reviewed, tested, and adjusted to meet the challenge requirements
4. The developer retained full understanding and ownership of the codebase

This approach allowed for **faster development** while maintaining **code quality** and ensuring complete understanding of all implemented functionality, which enables thorough discussion during the interview process.

---

# Signature Service - Coding Challenge - Original README
## Instructions

This challenge is part of the software engineering interview process at fiskaly.

If you see this challenge, you've passed the first round of interviews and are now at the second and last stage.

We would like you to attempt the challenge below. You will then be able to discuss your solution in the skill-fit interview with two of our colleagues from the development department.

The quality of your code is more important to us than the quantity.

### Project Setup

For the challenge, we provide you with:

- Go project containing the setup
- Basic API structure and functionality
- Encoding / decoding of different key types (only needed to serialize keys to a persistent storage)
- Key generation algorithms (ECC, RSA)
- Library to generate UUIDs, included in `go.mod`

You can use these things as a foundation, but you're also free to modify them as you see fit.

### Prerequisites & Tooling

- Golang (v1.20+)

### The Challenge

The goal is to implement an API service that allows customers to create `signature devices` with which they can sign arbitrary transaction data.

#### Domain Description

The `signature service` can manage multiple `signature devices`. Such a device is identified by a unique identifier (e.g. UUID). For now you can pretend there is only one user / organization using the system (e.g. a dedicated node for them), therefore you do not need to think about user management at all.

When creating the `signature device`, the client of the API has to choose the signature algorithm that the device will be using to sign transaction data. During the creation process, a new key pair (`public key` & `private key`) has to be generated and assigned to the device.

The `signature device` should also have a `label` that can be used to display it in the UI and a `signature_counter` that tracks how many signatures have been created with this device. The `label` is provided by the user. The `signature_counter` shall only be modified internally.

##### Signature Creation

For the signature creation, the client will have to provide `data_to_be_signed` through the API. In order to increase the security of the system, we will extend this raw data with the current `signature_counter` and the `last_signature`.

The resulting string (`secured_data_to_be_signed`) should follow this format: `<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>`

In the base case there is no `last_signature` (= `signature_counter == 0`). Use the `base64`-encoded device ID (`last_signature = base64(device.id)`) instead of the `last_signature`.

This special string will be signed (`Signer.sign(secured_data_to_be_signed)`) and the resulting signature (`base64` encoded) will be returned to the client. The signature response could look like this:

```json
{ 
    "signature": <signature_base64_encoded>,
    "signed_data": "<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>"
}
```

After the signature has been created, the signature counter's value has to be incremented (`signature_counter += 1`).

#### API

For now we need to provide two main operations to our customers:

- `CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse`
- `SignTransaction(deviceId: string, data: string): SignatureResponse`

Think of how to expose these operations through a RESTful HTTP-based API.

In addition, `list / retrieval operations` for the resources generated in the previous operations should be made available to the customers.

#### QA / Testing

As we are in the business of compliance technology, we need to make sure that our implementation is verifiably correct. Think of an automatable way to assure the correctness (in this challenge: adherence to the specifications) of the system.

#### Technical Constraints & Considerations

- The system will be used by many concurrent clients accessing the same resources.
- The `signature_counter` has to be strictly monotonically increasing and ideally without any gaps.
- The system currently only supports `RSA` and `ECDSA` as signature algorithms. Try to design the signing mechanism in a way that allows easy extension to other algorithms without changing the core domain logic.
- For now it is enough to store signature devices in memory. Efficiency is not a priority for this. In the future we might want to scale out. As you design your storage logic, keep in mind that we may later want to switch to a relational database.

### AI Tools

The use of AI tools to aid completing the challenge is permitted, but you will need to be able to reason about the design and implementation choices made when you reach the interview stage. Furthermore, if you used any AI tools, you need to clearly state which tools were used for different parts of the challenge. Ensure that you document this inside the `README` for your repository, so that it is visible to the reviewers.

### Credits

This challenge is heavily influenced by the regulations for `KassenSichV` (Germany) as well as the `RKSV` (Austria) and our solutions for them.
