# Monorepo

## Commandes

Docker build :
- `docker build --build-arg DIRECTORY=server --file server/Dockerfile --progress=plain . -t server`
- `docker build --build-arg DIRECTORY=bookHandler --file bookHandler/Dockerfile --progress=plain . -t bookhandler`

Récupérer le schema de la DB :
- `dc exec database pg_dump -U $POSTGRES_USER novel_database --schema-only`
- psql databasename < data_base_dump (c'est quoi déjà ?)
  
Générer les structure de sqlboiler :
- `sqlboiler -c schemas/sqlboiler.toml psql`