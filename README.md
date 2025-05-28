# chirpy
    Chirpy is a lightweight Go-based HTTPS server that provides a RESTful API inspired by Twitter. 

## Motivation
a practical project to help my understanding of HTTP request handling, REST API design, and secure server implementation in Go.

## Goal
It allows users to 
*post short messages
*retrieve short messages
*interact with short messages
*unique users
*include authentication


## API Documentation



## Endpoints
### /api/users
#### POST - used to create a user
        expects a json body in the following format:
        ```{
            password:<users password>
            Email:<users email>
        }```
##### Responses
returns a json response body representing the information of the user in the following format:
    ```{
        id:<uuid of the user>
        created_at:<datetime when the user was created>
        updated_at:<datetime when the user record was last updated>
        email:<string of the users email>
        is_chirpy_red:<booliean value of if the user is an active subscriber>
    }```
##### Status Codes
can return the following codes:
* 405 - if the method used is not POST.
* 400 - if the body is not formated correctly.
* 500 - if there is an issue with hashing a users password or an issue with creating the user in the database.
* 201 - user created succefully

#### PUT - used to update a user
        expects a json body and an authentication token in the header.
        the json body should be in the following format:
        ```{
            password:<users password>
            email:<users email>
        }```
##### Responses

##### Status Codes


 /api/chirps

### Responses

### Status Codes



 /api/login

### Responses

### Status Codes


 /api/refresh

### Responses

### Status Codes


 /api/revoke
    


### Responses

### Status Codes


 /api healthz

### Responses

### Status Codes



