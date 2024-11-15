# Car Management API

This is a RESTful API for managing cars and their images. The API allows users to:
- Add new cars with associated images.
- Retrieve all cars added by a user.
- Update car details such as title, description, and images.
- Delete cars from the database.

## Requirements

- Go (1.16+)
- MySQL Database
- Environment Variables:
  - `ACCESS_TOKEN_SECRET`: Secret key for JWT token generation

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/car-management-api.git
cd car-management-api
