# Cab API

Cab Data Researcher is a company that provides insights on the open
data about NY cab trips

In order to make it more useful we want to wrap the data in a public
API.

The API provide the number of trips cabs have done in a given day. Cabs
are identified by their medallion id.

Example API call:

```
    http://localhost:8000/trip/2013-12-01?medallion=00377E15077848677B32CE184CE7E871,01C905F5CF4CD4D366979BBCF281A15E,2B1A06E9228B7278227621EF1B879A1D,56DD5BCBDD7DE1F36A455AA7FAEA7C88
```

Return:

```
    [
        {
            "medallion": "00377E15077848677B32CE184CE7E871",
            "date": "2013-12-1",
            "total": 3
        },
        {
            "medallion": "01C905F5CF4CD4D366979BBCF281A15E",
            "date": "2013-12-1",
            "total": 1
        },
        {
            "medallion": "2B1A06E9228B7278227621EF1B879A1D",
            "date": "2013-12-1",
            "total": 4
        },
        {
            "medallion": "56DD5BCBDD7DE1F36A455AA7FAEA7C88",
            "date": "2013-12-1",
            "total": 3
        }
    ]
```

## Running

Download:

```
    $ go get github.com/gambarini/cabapi
```

You must have the following installed to run the API server:

- golang 1.9
- docker
- mysql cmdline client

From the cabapi folder execute the start script:

```
    $ start.sh
```

It will take a while (be patient!) to initialize the following:

- Redis docker container.
- MySQL docker container.
- Import the data from ny_cab_data_cab_trip_data_full.sql to MySQL db.
- Start the api server

If it all work, the API should be running on localhost:8000

Just Ctrl+C to terminate the server when you done. The script will cleanup the
containers and data folders.

## Endpoints

##### /trip/{yyyy-mm-dd}?medalion=id1,id2&cache=true|false

Return all cabs with the number of total trips for the given date.


Query Params:

- medallion (optional): comma separated list with the medalliond ids to search.
When "medallion" param is added it's a cached search by default.

- cache (optional): true|false. If false the trips cache is disabled. The default
is true. But there is only cache for requests with the "medallion" query param

