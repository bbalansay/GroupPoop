# GroupPoop

## Project Description

The target audience for GroupPoop are primarily the **students and faculty** at the **University of Washington**, though in the future we hope this to be scaled to a wider audience. 

We envision this population to use our application to find the best pit stops on campus. This application could be applied to anyone visiting the University of Washington, but we will be focusing on students and faculty because they are the predominant population on campus.

Our audience seeks to use GroupPoop to achieve the satisfying, entertaining break that they deserve. Bathrooms and their locations on campus will be stored in our database and each will have a rating based on ambience, cleanliness, toilet paper quality, etc. Users will also be able to join chats with others in the process of going about their business. With GroupPoop, users will truly feel like royalty on the porcelain throne.

Most of us, if not all of us, have had some uncomfortable situations when needing to use the restroom on campus. There is an inequality that exists between buildings and their bathrooms, and we want people to have the most comfortable experience when they need to relieve themselves. GroupPoop can make that happen.

<br>

## Technical Specifications

### Architectural Diagram
The system we create will implement a microservices architecture. All requests from users get handled by the Gateway layer server, which then creates needs and puts them onto the RabbitMQ request queue. Microservices will be subscribed to the request queue and if they can fulfill a need, they will, and then will return the fulfilled need onto the reply queue, which the gateway layer then receives and processes. The Redis store and MySQL store will be accessible via microservices.

![Architecture Diagram](img/architecture_diagram.png)

### User Stories

| #   | Priority | User      | Issue |
| --- | -------- | --------- | ----- |
| 1   | P0       | As a user | I want to get the information about a bathroom on campus |
| 2   | P0       | As a user | I want to add a new bathroom to GroupPoop |
| 3   | P0       | As a user | I want to chat with someone while going about my business |
| 4   | P0       | As a user | I want to create an account and log in |
| 5   | P1       | As a user | I want to rate a bathroom on campus |
| 6   | P1       | As a user | I want to make a list of my favorite bathrooms |
| 7   | P2       | As a user | I want to delete a review |
| 8   | P2       | As a user | I want to like a review |

<br>

| #   | Solution to Issue |
| --- | -------- |
| 1   | To get information about a bathroom on campus, make a **GET request** at `/bathrooms/{id}`. Upon receiving the request, the server will attempt to fetch data from the **MySQL database** using a **SELECT statement** and display the information if successful. |
| 2   | To add a new bathroom to GroupPoop, a **POST request** will be made at `/bathrooms` containing relevant information about the bathroom to add. The server will then contact the **MySQL database** and execute an **INSERT statement**, returning the object if created successfully. |
| 3   |  |
| 4   |  |
| 5   |  |
| 6   |  |
| 7   |  |
| 8   |  |

### Endpoints

### Models

We will be using MySql as our persistent data store.

**User** \
`User`: Keeps track of user information. \
create table if not exists User ( \
  user_id int not null auto_increment primary key, \
  email varchar(512) not null, \
  pass_hash varchar(128) not null, \
  user_name varchar(256) not null, \
  first_name varchar(128) not null, \
  last_name varchar(128) not null, \
  photo_url varchar(128) not null, \
  index (email, user_name) \
)

**Chat** \
`Chat`: Keeps track of a conversation between two users. \
create table if not exists Chat ( \
  chat_id int not null auto_increment primary key, \
  start_time datetime not null, \
  end_time datetime not null \
)

**Message** \
`Message`: Keeps track of individual message sent from one user. \
create table if not exists Message ( \
  message_id int not null auto_increment primary key, \
  FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE not null, \
  FOREIGN KEY (chat_id) REFERENCES Chat(chat_id) ON DELETE CASCADE not null, \
  content varchar(512) not null \
)

**Bathroom** \
`Bathroom`: Keeps track of information relating to a bathroom. \
create table if not exists Bathroom ( \
  bathroom_id int not null auto_increment primary key, \
  name varchar(128) not null, \
  description varchar(512) not null, \
  location varchar(128) not null, \
  num_sinks int not null, \
  num_toilets int not null, \
  num_urinals int not null, \
  num_trash_cans int not null, \
  num_hand_dryers int not null, \
  num_towel_dispenser int not null \
)

**Review** \
`Review`: Keeps track of review a user makes for a bathroom. \
create table if not exists Review ( \
  review_id int not null auto_increment primary key, \
  FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE not null, \
  FOREIGN KEY (bathroom_id) REFERENCES Bathroom(bathroom_id) ON DELETE CASCADE not null, \
  content varchar(512) not null, \
  time datetime not null \
)


