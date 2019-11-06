# GroupPoop

<br>

## Project Description

The target audience for GroupPoop are primarily the **students and faculty** at the **University of Washington**, though in the future we hope this to be scaled to a wider audience. 

We envision this population to use our application to find the best pit stops on campus. This application could be applied to anyone visiting the University of Washington, but we will be focusing on students and faculty because they are the predominant population on campus.

Our audience seeks to use GroupPoop to achieve the satisfying, entertaining break that they deserve. Bathrooms and their locations on campus will be stored in our database and each will have a rating based on ambience, cleanliness, toilet paper quality, etc. Users will also be able to join chats with others in the process of going about their business. With GroupPoop, users will truly feel like royalty on the porcelain throne.

Most of us, if not all of us, have had some uncomfortable situations when needing to use the restroom on campus. There is an inequality that exists between buildings and their bathrooms, and we want people to have the most comfortable experience when they need to relieve themselves. GroupPoop can make that happen.

<br>

## Technical Specifications

### Architectural Diagram

### User Stories

| #   | Priority | User      | Description |
| --- | -------- | --------- | ----------- |
| 1   | P0       | As a user | I want to get the information about a bathroom on campus |
| 2   | P0       | As a user | I want to add a new bathroom to GroupPoop |
| 3   | P0       | As a user | I want to chat with someone while going about my business |
| 4   | P0       | As a user | I want to create an account and log in |
| 5   | P1       | As a user | I want to rate a bathroom on campus |
| 6   | P1       | As a user | I want to make a list of my favorite bathrooms |
| 7   | P2       | As a user | I want to delete a review |
| 8   | P2       | As a user | I want to like a review |

<br>

| #   | Solution |
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

### Appendix