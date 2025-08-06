# VehSense BACKEND

This project aims to enable the communication between `VehSense` applications with the drivers database. 

## Overlook

`VehSense BACKEND` goal is to gather the data sent in batches from the `VehSense Android App` (that is connected to the esp32), sum and gather all of the info, and create raports. Each user can only have one active raport that gathers the data - the way it is planned to  work is it checks if there is a raport that has and end timestamp smaller than 10 mins from now. if yes, use the data to add the info that was sent to backend and update the asset. This approach ensures that a sudden stop of the app won't erase all the data that was meant for the report. That is also a better approach since the backend restart, or multiple backend instances could potentailly mean that the data stored in ram would disappear. 

`VehSense BACKEND` also will provide REST API for all of the crud operations on the application database.

> It is also planned to provide a websocket communication for live drivers informations.

There will be couple of roles: 
- driver - provided by special endpoint, where the user provides the active jwt and rest api returns the jwt with driver role - then he can provide the report info. 
- user - check out your own raports.
- admin - check out raports of all of the users in your organization.
- app admin - to add organizations and their admins

### Used chi + GORM + PostgreSQL + Docker.

# Current stage of the development - TODO

- [ ] Query params for all GET requests
- [ ] Decide what endpoints are really required and needed (root needs to be able to delete vehicles?)
- [ ] Validate tokens when the owner was deleted or changed the password

## Then...

- [ ] Unify the status codes
- [ ] Develop the Raport service with sessions etc.
- [ ] Prepare the tests of each endpoint

## And then...

- [ ] Implement changing the password by email
- [ ] HTTPS
- [ ] WebSocket for the admin panel LIVE view
- [ ] Microservices
