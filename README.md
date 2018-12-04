# Izly command line tool

## Installation

~~~
go get github.com/zzOzz/izlyctl
~~~

## Usage

retrieve user cards
~~~
izlyctl get cards "my.user@domain.com"
{
  "rightholder": {
    "identifier": "my.user@domain.com",
    "firstName": "My",
    "lastName": "USER",
    "email": "my.user@domain.com",
    "dueDate": "2049-01-01T01:00:00.000+01",
    "idCompanyRate": 68,
    "idRate": 2,
    "birthDate": "1970-01-01T00:00:00.000+01",
    "rneOrgCode": "0600000A",
    "rneDepCode": "0600000A",
    "cellNumber": "33606060606",
    "createdAt": 1440753777000,
    "updatedAt": 1543913520000,
    "idCrous": 69,
    "pic": 998123456
  },
  "smartcards": [
    {
      "idTransmitter": 65,
      "idMapping": 1,
      "idZdc": 36000000,
      "zdcCreationDate": "2015-09-23T02:00:00.000+02",
      "pixSs": "03",
      "pixNn": "0000",
      "appl": "D1",
      "uid": "04000000000000",
      "rid": "04000000000000",
      "cancelDate": null,
      "revalidateDate": null,
      "deliveryDate": null
    }
  ]
}
~~~

set user rate '-r' and due date '-D'

~~~
izlyctl set -r 2 "my.user@domain.com" -D "2049-01-01"
~~~

debug flag '-d'

~~~
izlyctl set "my.user@domain.com" -d -r 2 -D "2049-01-01"
~~~


## config file

Create .izlyctl.yaml

~~~
api:
  auth:
    login: my-izly-api-login
    password: my-izly-api-password
  url: https://api.lescrous.fr
~~~

you can also pass user/password/url directly through command line 

~~~
izlyctl -u "my-izly-api-login" -p "my-izly-api-password" -s "https://api.lescrous.fr" get cards "my.user@domain.com"
~~~