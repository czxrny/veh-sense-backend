userauth_insert_query = """
    INSERT INTO user_auths (email, password, role) 
    VALUES (%s, %s, %s)
    RETURNING id
"""

userinfo_insert_query = """
    INSERT INTO user_infos (user_name, organization_id, total_kilometers) 
    VALUES (%s, %s, %s)
    RETURNING id
"""

vehicle_insert_query = """
    INSERT INTO vehicles (
        owner_id,
        organization_id,
        brand,
        model,
        year,
        engine_capacity,
        engine_power,
        plates,
        expected_fuel
    )
    VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
"""

organization_insert_query = """
    INSERT INTO organizations (
        name, address, city, country, zip_code, country_code, contact_number
    ) VALUES (%s, %s, %s, %s, %s, %s, %s)
    RETURNING id
"""