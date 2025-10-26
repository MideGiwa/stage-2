# Country Data API

A RESTful API built with Go and Gin that fetches country data from external APIs, stores it in a MySQL database, and provides CRUD operations with additional functionalities like image generation and status tracking.

## Table of Contents

- [Features](#features)
- [Country Fields](#country-fields)
- [Validation Rules](#validation-rules)
- [External APIs Used](#external-apis-used)
- [Setup Instructions](#setup-instructions)
  - [Prerequisites](#prerequisites)
  - [Environment Variables](#environment-variables)
  - [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
  - [POST /countries/refresh](#post-countriesrefresh)
  - [GET /countries](#get-countries)
  - [GET /countries/:name](#get-countriesname)
  - [DELETE /countries/:name](#delete-countriesname)
  - [GET /status](#get-status)
  - [GET /countries/image](#get-countriesimage)
- [Error Handling](#error-handling)
- [Image Generation](#image-generation)

## Features

- **Fetch and Cache Data**: Retrieves country information from `restcountries.com` and exchange rates from `open.er-api.com`.
- **Database Storage**: Stores and updates data in a MySQL database.
- **Computed Fields**: Calculates `estimated_gdp` based on population, a random multiplier, and exchange rates.
- **CRUD Operations**: Provides endpoints for fetching all countries, fetching by name, and deleting by name.
- **Filtering and Sorting**: Supports filtering countries by `region` and `currency`, and sorting by `gdp_desc`, `name_asc`, etc.
- **Status Endpoint**: Shows the total number of cached countries and the timestamp of the last refresh.
- **Summary Image Generation**: Creates a `summary.png` image with total countries, top 5 by GDP, and refresh timestamp.
- **Consistent Error Handling**: Returns standardized JSON error responses.

## Country Fields

| Field Name        | Type     | Description                                                          | Constraints             |
| :---------------- | :------- | :------------------------------------------------------------------- | :---------------------- |
| `id`              | `uint`   | Auto-generated primary key                                           | Primary Key             |
| `name`            | `string` | Country name                                                         | Required, Unique        |
| `capital`         | `string` | Capital city                                                         | Optional                |
| `region`          | `string` | Region (e.g., Africa, Europe)                                        | Optional                |
| `population`      | `uint64` | Total population                                                     | Required                |
| `currency_code`   | `string` | ISO currency code (e.g., USD, NGN)                                   | Required                |
| `exchange_rate`   | `float64`| Exchange rate against USD                                            | Optional                |
| `estimated_gdp`   | `float64`| Computed as `population × random(1000–2000) ÷ exchange_rate`         | Optional                |
| `flag_url`        | `string` | URL to the country's flag image                                      | Optional                |
| `last_refreshed_at` | `time.Time` | Timestamp of the last update for this country                       | Auto-updated            |

## Validation Rules

- `name`, `population`, and `currency_code` are required.
- Invalid or missing data will result in a `400 Bad Request` with a JSON error body.

Example:
```json
{
  "error": "Validation failed",
  "details": {
    "currency_code": "is required"
  }
}
```

## External APIs Used

- **Countries Data**: `https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies`
- **Exchange Rates**: `https://open.er-api.com/v6/latest/USD`

## Setup Instructions

### Prerequisites

Before running the application, ensure you have the following installed:

- **Go**: Version 1.16 or higher.
- **MySQL Database**: A running MySQL instance.
- **Git**: For cloning the repository.

### Environment Variables

Create a `.env` file in the root directory of the project based on the `.env.example` file.

```
/Users/mide/Documents/Projects/Hng-13/Backend/stage-2/.env.example#L1-6
PORT=8080
DB_USER=root
DB_PASSWORD=password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=country_db
```

Replace the placeholder values with your actual database credentials and desired port.

### Running the Application

1. **Clone the repository:**
   ```bash
   git clone <repository_url>
   cd <repository_directory>/Backend/stage-2
   ```

2. **Install Go modules:**
   ```bash
   go mod tidy
   ```

3. **Ensure font for image generation:**
   The image generation uses a font. By default, it looks for `assets/Roboto-Bold.ttf`. If this file is not present, it will use a default sans font. For optimal appearance, you can download `Roboto-Bold.ttf` from Google Fonts and place it in an `assets` directory at the project root.
   ```bash
   mkdir -p assets
   # Place Roboto-Bold.ttf inside the assets directory
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```
   The API will start listening on the port specified in your `.env` file (default: `8080`).

## API Endpoints

All responses are in JSON format.

### `POST /countries/refresh`

Fetches all countries and exchange rates from external APIs, then caches them in the database.
This endpoint also triggers the generation of the `cache/summary.png` image.

- **URL**: `/countries/refresh`
- **Method**: `POST`
- **Response**:
  ```json
  {
    "message": "Countries refreshed successfully",
    "total_countries": 250,
    "last_refreshed_at": "2025-10-22T18:00:00Z"
  }
  ```
- **Error Response (External API failure)**:
  ```json
  {
    "error": "External data source unavailable",
    "details": "Could not fetch data from restcountries.com"
  }
  ```

### `GET /countries`

Retrieves all cached countries from the database. Supports filtering and sorting.

- **URL**: `/countries`
- **Method**: `GET`
- **Query Parameters**:
  - `?region=[region_name]`: Filter by country region (e.g., `?region=Africa`). Case-insensitive.
  - `?currency=[currency_code]`: Filter by currency code (e.g., `?currency=NGN`). Case-insensitive.
  - `?sort=[field_order]`: Sort results.
    - `gdp_desc`: Sort by estimated GDP in descending order.
    - `gdp_asc`: Sort by estimated GDP in ascending order.
    - `name_desc`: Sort by name in descending order.
    - `name_asc`: Sort by name in ascending order (default).
    - `population_desc`: Sort by population in descending order.
    - `population_asc`: Sort by population in ascending order.
- **Example Response (`GET /countries?region=Africa`):**
  ```json
  [
    {
      "id": 1,
      "name": "Nigeria",
      "capital": "Abuja",
      "region": "Africa",
      "population": 206139589,
      "currency_code": "NGN",
      "exchange_rate": 1600.23,
      "estimated_gdp": 25767448125.2,
      "flag_url": "https://flagcdn.com/ng.svg",
      "last_refreshed_at": "2025-10-22T18:00:00Z"
    },
    {
      "id": 2,
      "name": "Ghana",
      "capital": "Accra",
      "region": "Africa",
      "population": 31072940,
      "currency_code": "GHS",
      "exchange_rate": 15.34,
      "estimated_gdp": 3029834520.6,
      "flag_url": "https://flagcdn.com/gh.svg",
      "last_refreshed_at": "2025-10-22T18:00:00Z"
    }
  ]
  ```

### `GET /countries/:name`

Retrieves a single country by its name.

- **URL**: `/countries/{country_name}` (e.g., `/countries/Nigeria`)
- **Method**: `GET`
- **Example Response:**
  ```json
  {
    "id": 1,
    "name": "Nigeria",
    "capital": "Abuja",
    "region": "Africa",
    "population": 206139589,
    "currency_code": "NGN",
    "exchange_rate": 1600.23,
    "estimated_gdp": 25767448125.2,
    "flag_url": "https://flagcdn.com/ng.svg",
    "last_refreshed_at": "2025-10-22T18:00:00Z"
  }
  ```
- **Error Response (Country not found)**:
  ```json
  {
    "error": "Country not found"
  }
  ```

### `DELETE /countries/:name`

Deletes a country record by its name.

- **URL**: `/countries/{country_name}` (e.g., `/countries/Nigeria`)
- **Method**: `DELETE`
- **Response**:
  ```json
  {
    "message": "Country deleted successfully"
  }
  ```
- **Error Response (Country not found)**:
  ```json
  {
    "error": "Country not found"
  }
  ```

### `GET /status`

Retrieves the overall status of the cached data.

- **URL**: `/status`
- **Method**: `GET`
- **Example Response:**
  ```json
  {
    "total_countries": 250,
    "last_refreshed_at": "2025-10-22T18:00:00Z"
  }
  ```

### `GET /countries/image`

Serves the generated summary image.

- **URL**: `/countries/image`
- **Method**: `GET`
- **Response**: Serves the `cache/summary.png` image directly.
- **Error Response (Image not found)**:
  ```json
  {
    "error": "Summary image not found"
  }
  ```

## Error Handling

The API returns consistent JSON error responses:

- `400 Bad Request`: `{ "error": "Validation failed", "details": { "field": "is required" } }`
- `404 Not Found`: `{ "error": "Resource not found" }` (e.g., `Country not found`)
- `500 Internal Server Error`: `{ "error": "Internal server error" }`
- `503 Service Unavailable`: `{ "error": "External data source unavailable", "details": "Could not fetch data from [API name]" }`

## Image Generation

Upon a successful `POST /countries/refresh` request, the API generates an image named `summary.png` in the `cache/` directory. This image includes:
- The total number of countries cached.
- The top 5 countries by estimated GDP.
- The timestamp of the last data refresh.

This image can then be accessed via the `GET /countries/image` endpoint.