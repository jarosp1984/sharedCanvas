# List of major prompts used in this project

1. create html with a canvas and javascript code that reacts to mause events and allow painting on the canvas
2. now add go application that host this file
3. plan REST API for drawing line from X1, Y1, to X2, Y2
4. Add Open API file describing current API, put shi in /doc/openapi/ folder
5. Plan new API method for receiveing stored lines in order.
6. Now we need to modify script in index.html that old lines are loaded and drawn on page load
7. Plan the changes to enable auto update of data: all other clients should get newly added lines every 1 second.
8. Plan changes to add clearing canvas feature, the button is already present in html
9. update @openapi.yaml to cover clear methods
10. Currently server supports only one session without any identification, I want to add sessions using UUID as identifier, this UUID will be part of URL, so users can send URL of the sessin to someone to participate. Please plan the solution.