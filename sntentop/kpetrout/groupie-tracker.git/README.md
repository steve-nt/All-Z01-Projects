# Groupie Tracker 

## Description 
Groupie Tracker, is an online web application that is taking data from a given API and presents them in the order given on the port `:8080` of the `localhost`.

## API content and app features
The data of the API are basically informations about some famous music bands or solo music artists, that the user can see in the app with ease.
Some basic information examples are :
- **Locations :** some of the locations that the band or artist has played before 
- **Dates :** some of the dates of their music concert
- **Relations :** a combination of the locations and dates that are matching 

but you can see a lot more.

## Usage
To run this program you need to install **golang** and after :
1. Clone the repository :
- Open a terminal, and type the following command :
    ```bash
    git clone https://platform.zone01.gr/git/kpetrout/groupie-tracker.git
2. Set the API keys :
- Navigate into the project with `cd` command and run the command:
    ```bash
    . set_keys.sh
3. Run the program :
- Once the keys are setted run the command:
    ```go
    go run .
- If the program run succesfully you will see a message :
    ```
    Server is running on http://localhost:8080
4. Open in browser :
- Now you can (ctrl+click) on the link of the message or just click here http://localhost:8080

## Authors
Creators and Primary Developers :
- Konstantinos Petroutsos
- Christos Markos
- Christos Gkaldanidis