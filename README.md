# tweetstorm
A service that indexes tweets

# Running the app
1. Install docker
1. Update configuration as desired - e.g., select terms to send to filter API
1. Set up a `.env` file in the repository root with the following values from your Twitter API account (or provide them in the `docker-compose up` command below)
    * TWITTER_API_KEY
    * TWITTER_API_SECRET
    * TWITTER_ACCESS_TOKEN
    * TWITTER_ACCESS_TOKEN_SECRET
1. Run `docker-compose up` with optional `-d`

If using visual studio code, [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension and the included configuration to open the editor in container. You can then provide the environment variables in a local `.env` file and run the app in the editor.

# Viewing data
## Elasticsearch
* The docker-compose file will create a cerebro container and map the port at `localhost:9000`. *Use `es01` as the host, not `localhost`.*
* Kibana will also be running at `localhost:5601`.

You can also consume elasticsearch API directly or with another tool, through the localhost mapped port (`localhost:9200`).

## MongoDB
* Connect any MongoDB client to `localhost:27017` if you have a preferred tool (e.g. Compass)
* Run `docker exec -it tweetstorm_mongo_1 mongo` (possibly replacing the container name with a value from `docker ps`) to open a mongodb CLI without the need to install anything.
