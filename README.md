# Monorepo

## Commandes

Docker build :
- `docker build --build-arg DIRECTORY=server --file server/Dockerfile --progress=plain . -t server`
- `docker build --build-arg DIRECTORY=book-handler --file book-handler/Dockerfile --progress=plain . -t book-handler`

Récupérer le schema de la DB :
- `dc exec database pg_dump -U $POSTGRES_USER novel --schema-only`
- psql databasename < data_base_dump (c'est quoi déjà ?)
  
Générer les structure de sqlboiler :
- `sqlboiler -c ../schemas/sqlboiler.toml --wipe psql`

Se connecter sur le DB : 
- `dc exec database psql -h database -U $POSTGRES_USER novel`
- `docker exec -it database psql -h database -U $POSTGRES_USER novel`
- `psql --host=host.docker.internal --user=root_user --password novel`

pg_dump -U dbusername dbname > dbexport.pgsql
dc exec database pg_dump -U $POSTGRES_USER novel > dbexport.pgsql
dc exec database pg_dump -U $POSTGRES_USER postgres --schema-only
