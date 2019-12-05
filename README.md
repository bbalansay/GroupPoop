# GroupPoop

## Project Description

The target audience for GroupPoop are primarily the **students and faculty** at the **University of Washington**, though in the future we hope this to be scaled to a wider audience. 

We envision this population to use our application to find the best pit stops on campus. This application could be applied to anyone visiting the University of Washington, but we will be focusing on students and faculty because they are the predominant population on campus.

Our audience seeks to use GroupPoop to achieve the satisfying, entertaining break that they deserve. Bathrooms and their locations on campus will be stored in our database and each will have a rating based on ambience, cleanliness, toilet paper quality, etc. Users will also be able to join chats with others in the process of going about their business. With GroupPoop, users will truly feel like royalty on the porcelain throne.

Most of us, if not all of us, have had some uncomfortable situations when needing to use the restroom on campus. There is an inequality that exists between buildings and their bathrooms, and we want people to have the most comfortable experience when they need to relieve themselves. GroupPoop can make that happen.

<br>

## Technical Specifications

### Initial Architectural Diagram
The system we create will implement a microservices architecture. All requests from users get handled by the Gateway layer server, which then creates needs and puts them onto the RabbitMQ request queue. Microservices will be subscribed to the request queue and if they can fulfill a need, they will, and then will return the fulfilled need onto the reply queue, which the gateway layer then receives and processes. The Redis store and MySQL store will be accessible via microservices.

![Architecture Diagram](img/architecture_diagram.png)

### Final Architectural Diagram
Here is our final architecture diagram. As you can see we got rid of RabbitMQ handling microservices and our API gateway now functions as a reverse proxy.

![Final Architecture Diagram](img/final_architecture_diagram.png)

### User Stories

| #   | Priority | User      | Issue |
| --- | -------- | --------- | ----- |
| 1   | P0       | As a user | I want to get the information about a bathroom on campus |
| 2   | P0       | As a user | I want to chat with someone while going about my business |
| 3   | P0       | As a user | I want to create an account and log in |
| 4   | P1       | As a user | I want to rate a bathroom on campus |
| 5   | P1       | As a user | I want to make a list of my favorite bathrooms |
| 6   | P2       | As a user | I want to delete a review |

<br>

| #   | Solution to Issue |
| --- | -------- |
| 1   | To get information about a bathroom on campus, make a **GET request** at `/bathrooms/{id}`. Upon receiving the request, the server will attempt to fetch data from the **MySQL database** using a **SELECT statement** and display the information if successful. |
| 2   | To chat with someone, a user must utilize a websocket connection to connect with other users who are logged in. |
| 3   | To create an account, a user must make a **POST request** at `/users/{id}`. Upon receiving the request, add a new user to the **MySQL database** using an **INSERT statement** and the provided credentials. |
| 4   | To review a bathroom on campus, make a **POST request** at `/bathrooms/{id}`. Upon receiving the request, the server will create a new **INSERT statement** using the information prvided to add to the **MySQL database**.|
| 5   | To make a list of favorite bathrooms, make a **PATCH request** at `/user/{id}`. Upon receiving the request, the server will update the user information in the **MySQL database** to include a list of bathrooms. |
| 6   | To delete a review of a bathroom make a **DELETE request** at `/review/{id}`. Upon receiving the request, the server will delete a review from the **MySQL database** that matches the given information. |

### Endpoints
`/user/login`:

- `POST`: `application/json`: Log in user and returns session token.
	- `200`: `application/json`: Successfully logs in user; returns session token in `Authorization` header.
  - `401`: Cannot authenticate provided credentials.
  - `415`: Cannot decode body / received unsupported body.
  - `500`: Internal server error.

- `DELETE`: Log out a user.
  - `200`: Successfully logs out user. 
  - `401`: Cannot verify session token or no session token. 
  - `500`: Internal server error.

<br>

`/user`:

- `GET`: Get user information, including reviews.
	- `200`; `application/json`: Succesfully retrieves user information, returns encoded user model in body.
	- `401`: Cannot verify session token or no session token.
	- `500`: Internal server error.
- `POST`: `application/json`: Create a new user.
	- `201`; `application/json`: Successfully creates a new user, returns encoded user model in body. 
	- `401`: Cannot verify session token or no session token.  
	- `415`: Cannot decode body / received unsupported body. 
	- `500`: Internal server error. 
- `PATCH`: `application/json`: Update password for user.
	- `200`; `application/json`: Successfully updates password for user. 
	- `401`: Cannot verify session token or no session token. 
	- `415`: Cannot decode body / received unsupported body. 
	- `500`: Internal server error. 
- `DELETE`: Delete a user.
	- `200`: Successfully deletes user. 
	- `401`: Cannot verify session token or no session token. 
	- `500`: Internal server error. 

<br>

`/review`: 

- `GET`: Get review information
	- `200`; `application/json`: Succesfully retrieves review information, returns encoded review model in body. 
	- `401`: Cannot verify session token or no session token. 
	- `500`: Internal server error. 
- `POST`: `application/json`: Create a new review.
	- `201`; `application/json`: Successfully creates a new review, returns encoded review model in body. 
	- `401`: Cannot verify session token or no session token. 
	- `415`: Cannot decode body / received unsupported body. 
	- `500`: Internal server error. 
- `PATCH`: `application/json`: Update review.
	- `200`; `application/json`: Successfully updates review. 
	- `401`: Cannot verify session token or no session token. 
	- `415`: Cannot decode body / received unsupported body. 
	- `500`: Internal server error. 
- `DELETE`: Delete a review.
	- `200`: Successfully deletes review. 
	- `401`: Cannot verify session token or no session token. 
	- `500`: Internal server error. 

<br>

`/bathroom`: 

- `GET`: Get bathroom information
	- `200`; `application/json`: Succesfully retrieves bathroom information, returns encoded review model in body. 
	- `401`: Cannot verify session token or no session token. 
	- `500`: Internal server error. 

<br>

`/chat`:
- Websocket connection for users to chat with each other.
- User is required to connect with session token otherwise they are not logged in.




### Models

We will be using MySql as our persistent data store.

`User`: Keeps track of user information. 

```
create table if not exists User ( 
  user_id int not null auto_increment primary key, 
  email varchar(512) not null, 
  pass_hash varchar(128) not null, 
  user_name varchar(256) not null, 
  first_name varchar(128) not null, 
  last_name varchar(128) not null, 
  photo_url varchar(128) not null, 
  index (email, user_name) 
)
```

<br>

`Chat`: Keeps track of a conversation between two users. 
```
create table if not exists Chat ( 
  chat_id int not null auto_increment primary key, 
  start_time datetime not null, 
  end_time datetime not null 
)
```

<br>

`Message`: Keeps track of individual message sent from one user. 
```
create table if not exists Message ( 
  message_id int not null auto_increment primary key, 
  FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE not null, 
  FOREIGN KEY (chat_id) REFERENCES Chat(chat_id) ON DELETE CASCADE not null, 
  content varchar(512) not null 
)
```

<br>

`Bathroom`: Keeps track of information relating to a bathroom. 
```
create table if not exists Bathroom ( 
  bathroom_id int not null auto_increment primary key, 
  name varchar(128) not null, 
  description varchar(512) not null, 
  location varchar(128) not null, 
  num_sinks int not null, 
  num_toilets int not null, 
  num_urinals int not null, 
  num_trash_cans int not null, 
  num_hand_dryers int not null, 
  num_towel_dispenser int not null 
)
```

<br>

`Review`: Keeps track of review a user makes for a bathroom. 
```
create table if not exists Review ( 
  review_id int not null auto_increment primary key, 
  FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE not null, 
  FOREIGN KEY (bathroom_id) REFERENCES Bathroom(bathroom_id) ON DELETE CASCADE not null, 
  content varchar(512) not null, 
  time datetime not null 
)
```

