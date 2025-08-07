users = [
    # ROOT
    {
        "email": "root@example.com",
        "password": "root123",
        "role": "root",
        "user_name": "System Root",
        "organization_id": None,
        "total_kilometers": 0
    },
    # ADMINS
    {
        "email": "admin1@org1.com",
        "password": "admin123",
        "role": "admin",
        "user_name": "Admin Org1",
        "organization_id": 1,
        "total_kilometers": 0
    },
    {
        "email": "admin2@org2.com",
        "password": "admin123",
        "role": "admin",
        "user_name": "Admin Org2",
        "organization_id": 2,
        "total_kilometers": 0
    },
    # USERS
    {
        "email": "user1@org1.com",
        "password": "user123",
        "role": "user",
        "user_name": "User Org1",
        "organization_id": 1,
        "total_kilometers": 123
    },
    {
        "email": "user2@org2.com",
        "password": "user123",
        "role": "user",
        "user_name": "User Org2",
        "organization_id": 2,
        "total_kilometers": 456
    },
    # one is private
    {
        "email": "user3@example.com",
        "password": "user123",
        "role": "user",
        "user_name": "User Bez Org",
        "organization_id": None,
        "total_kilometers": 789
    }
]

vehicles = [
    # user_id=4 (organization_id=1)
    (4, 1, "Toyota", "Corolla", 2019, 1800, 140, "XYZ1234", 6.5),
    # Organization 1 shared vehicle!
    (None, 1, "Ford", "Focus", 2021, 1600, 150, "XYZ5678", 6.2),

    # user_id=2 (organization_id=2)
    (5, 2, "Skoda", "Octavia", 2018, 2000, 190, "ABC1111", 7.1),
     # Organization 2 shared vehicle!
    (None, 2, "Volkswagen", "Golf", 2020, 1600, 170, "DEF2222", 6.4),

    # user_id=3 (private)
    (6, None, "Mazda", "3", 2017, 1500, 120, "GHI3333", 6.0),
    (6, None, "Hyundai", "i30", 2019, 1400, 110, "JKL4444", 5.8),
]

organizations = [
    ('Acme Logistics', '123 Warehouse St', 'New York', 'USA', '10001', '+1', '+1-555-111-222'),
    ('TransGlobal Freight', '456 Distribution Ave', 'Berlin', 'Germany', '10115', '+49', '+49-30-1234567')
]