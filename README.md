# userm

User management and authentication service.

# Usage

To start the user service, especially for local development, there is a dedicated script `run.sh` to use for that. If the Docker image is not built yet locally, the script will
automatically do that for you. If you want to force the rebuild regardless, you can execute:

```shell
sh run.sh --rebuild
```

This will of course only run userm, which means that you have to provide a Postgres instance for yourself. If you're feeling lazy, you could also just use:

```shell
docker-compose up -d 
```