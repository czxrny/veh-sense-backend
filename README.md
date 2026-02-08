# VehSense Backend

`VehSense Backend` enables communication between **VehSense** applications through a server.  

## Overview

The project aims to implement endpoints that provide:  
- user authorization,  
- retrieval of requested data,  
- analysis of ride data submitted via the [VehSense Mobile App](https://github.com/czxrny/veh-sense-app),  
- report generation based on that data.  

The project consists of two services: **REST API** and **Batch Receiver**:  
- **REST API** – handles authorization and CRUD operations for the client.  
- **Batch Receiver** – receives data batches, generates reports, and stores data for frontend visualization.  

The project uses: `chi`, `GORM`, `PostgreSQL`, and `Docker`.

## User Roles

The backend supports different types of users:  
- **user** – [VehSense Mobile App](https://github.com/czxrny/veh-sense-app) user, either a *fleet driver* or *private*.  
- **admin** – [VehSense Web App](https://github.com/czxrny/veh-sense-panel) user with access to all fleet-related data.  
- **root** – system user who can modify most assets in the server database.  

Endpoint documentation is available in the `/docs` folder.  

## Getting Started

To run VehSense Backend:  

1. Clone the repository:  
```bash
git clone https://github.com/czxrny/veh-sense-backend.git
```

2. Navigate to the main project directory:
```
cd ./veh-sense-backend
```

3. Configure environment variables in .env.example:
```
nano ./.env.example
```

4. Copy .env.example to .env:
```
cp ./.env.example ./.env
```

5. Run the project:
```
docker compose up -d --build
```

> Optionally, you can create a docker-compose.override.yml to configure a local PostgreSQL database.

## TODO

- [ ] Unify error handling

- [ ] Standardize status codes

- [ ] Update REST API Swagger documentation

- [ ] Refactor User Service code

- [ ] Allow admins to edit vehicle status

- [ ] Implement the rest of the tests

## Future Improvements

- [ ] Validate tokens if the owner is deleted or changes their password

- [ ] Enable HTTPS

- [ ] WebSocket for live admin panel view

- [ ] Microservices architecture

- [ ] Implement password reset via email
