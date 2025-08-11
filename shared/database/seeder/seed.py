import psycopg2
import os
from dotenv import load_dotenv
import bcrypt

from insert_queries import userauth_insert_query, userinfo_insert_query, vehicle_insert_query, organization_insert_query
from seed_data import users, organizations, vehicles

load_dotenv() 

conn = psycopg2.connect(
    host="localhost",
    port=5432,
    dbname=os.getenv("POSTGRES_DB"),
    user=os.getenv("POSTGRES_USER"),
    password=os.getenv("POSTGRES_PASSWORD")
)

cur = conn.cursor()

cur.execute("SELECT * FROM user_auths;")
rows = cur.fetchall()
if len(rows) == 0:
    for user in users:
        hashed_pw = bcrypt.hashpw(user["password"].encode(), bcrypt.gensalt()).decode()
        cur.execute(userauth_insert_query, (user["email"], hashed_pw, user["role"]))
        cur.execute(userinfo_insert_query, (user["user_name"], user["organization_id"], user["total_kilometers"]))
else:
    print("Users table is not empty.")

cur.execute("SELECT * FROM organizations;")
rows = cur.fetchall()
if len(rows) == 0:
    for organization in organizations:
        cur.execute(organization_insert_query, organization)
else:
    print("Organizations table is not empty.")

cur.execute("SELECT * FROM vehicles;")
rows = cur.fetchall()  # i tutaj
if len(rows) == 0:
    for vehicle in vehicles:
        cur.execute(vehicle_insert_query, vehicle)
else:
    print("Vehicles table is not empty.")

conn.commit()

cur.close()
conn.close()
