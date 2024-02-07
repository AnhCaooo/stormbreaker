Stormbreaker: a service fetches electric price in Finland and does handlings behind the scene which is written in [Go](https://go.dev/)

# Basic logic: 
Fetch data from oomi.fi  once per day and store data to database. Then return this value to client
# Advanced logic: 
Before fetch data from oomi.fi, check from database if the query exists or not. If exists, get from db, otherwise call to Oomi.fi

**Purpose:** this will prevent someone tries to use this service spam Oomi.fi